package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
)

var (
	port   = flag.String("port", "53", "Port to listen on")
	domain = flag.String("domain", "tailscale.com", "Domain to use for responses")
)

func main() {
	flag.Parse()

	addr := fmt.Sprintf("127.0.0.1:%s", *port)
	conn, err := net.ListenPacket("udp", addr)
	if err != nil {
		log.Fatalf("Failed to listen on %s: %v", addr, err)
	}
	defer conn.Close()

	log.Printf("Fake DNS resolver listening on %s", addr)
	log.Printf("Resolving all queries to fixed values:")
	log.Printf("  A:     192.0.2.1")
	log.Printf("  AAAA:  2001:db8::1")
	log.Printf("  CNAME: www.%s -> %s", *domain, *domain)
	log.Printf("  TXT:   v=test dns resolver")
	log.Printf("  MX:    10 mx.%s", *domain)

	buffer := make([]byte, 512)
	for {
		n, clientAddr, err := conn.ReadFrom(buffer)
		if err != nil {
			log.Printf("Error reading: %v", err)
			continue
		}

		go handleDNSQuery(conn, buffer[:n], clientAddr)
	}
}

func handleDNSQuery(conn net.PacketConn, query []byte, clientAddr net.Addr) {
	if len(query) < 12 {
		log.Printf("Query too short: %d bytes", len(query))
		return
	}

	// Parse the DNS query
	txID := binary.BigEndian.Uint16(query[0:2])
	flags := binary.BigEndian.Uint16(query[2:4])
	qdCount := binary.BigEndian.Uint16(query[4:6])

	// Check if this is a query (QR bit = 0)
	if flags&0x8000 != 0 {
		log.Printf("Not a query, ignoring")
		return
	}

	// Parse question section
	if qdCount == 0 {
		log.Printf("No questions in query")
		return
	}

	// Extract domain name and query type
	queryDomain, qtype, offset := parseDNSQuestion(query[12:])
	if offset == 0 {
		log.Printf("Failed to parse question")
		return
	}

	log.Printf("Query from %s: %s (type %d)", clientAddr, queryDomain, qtype)

	// Build response
	response := buildDNSResponse(txID, query[12:12+offset], qtype, *domain)

	// Send response
	if _, err := conn.WriteTo(response, clientAddr); err != nil {
		log.Printf("Error sending response: %v", err)
	}
}

// parseDNSQuestion extracts the domain name and query type from a DNS question
func parseDNSQuestion(question []byte) (domain string, qtype uint16, offset int) {
	pos := 0
	labels := []string{}

	for pos < len(question) {
		length := int(question[pos])
		pos++

		if length == 0 {
			break
		}

		if pos+length > len(question) {
			return "", 0, 0
		}

		labels = append(labels, string(question[pos:pos+length]))
		pos += length
	}

	if pos+4 > len(question) {
		return "", 0, 0
	}

	domain = ""
	for i, label := range labels {
		if i > 0 {
			domain += "."
		}
		domain += label
	}

	qtype = binary.BigEndian.Uint16(question[pos : pos+2])
	offset = pos + 4

	return domain, qtype, offset
}

// buildDNSResponse creates a DNS response packet
func buildDNSResponse(txID uint16, question []byte, qtype uint16, domain string) []byte {
	response := make([]byte, 0, 512)

	// DNS Header
	header := make([]byte, 12)
	binary.BigEndian.PutUint16(header[0:2], txID)        // Transaction ID
	binary.BigEndian.PutUint16(header[2:4], 0x8180)      // Flags: Response, Recursion Desired + Available
	binary.BigEndian.PutUint16(header[4:6], 1)           // Questions: 1
	binary.BigEndian.PutUint16(header[6:8], 1)           // Answer RRs: 1
	binary.BigEndian.PutUint16(header[8:10], 0)          // Authority RRs: 0
	binary.BigEndian.PutUint16(header[10:12], 0)         // Additional RRs: 0

	response = append(response, header...)
	response = append(response, question...)

	// Build answer based on query type
	answer := buildAnswer(question, qtype, domain)
	response = append(response, answer...)

	return response
}

// buildAnswer creates the answer section for different DNS record types
func buildAnswer(question []byte, qtype uint16, domain string) []byte {
	answer := make([]byte, 0, 256)

	// Name (pointer to question)
	answer = append(answer, 0xc0, 0x0c)

	// Type
	answer = append(answer, byte(qtype>>8), byte(qtype&0xff))

	// Class (IN)
	answer = append(answer, 0x00, 0x01)

	// TTL (300 seconds)
	answer = append(answer, 0x00, 0x00, 0x01, 0x2c)

	switch qtype {
	case 1: // A record
		// RDLENGTH: 4
		answer = append(answer, 0x00, 0x04)
		// RDATA: 192.0.2.1 (TEST-NET-1)
		answer = append(answer, 192, 0, 2, 1)

	case 28: // AAAA record
		// RDLENGTH: 16
		answer = append(answer, 0x00, 0x10)
		// RDATA: 2001:db8::1
		answer = append(answer, 0x20, 0x01, 0x0d, 0xb8, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01)

	case 5: // CNAME record
		// Build domain name: www.<domain> -> <domain>
		cname := encodeDomainName(domain)
		// RDLENGTH
		answer = append(answer, byte(len(cname)>>8), byte(len(cname)&0xff))
		// RDATA
		answer = append(answer, cname...)

	case 16: // TXT record
		txt := "v=test dns resolver"
		// RDLENGTH: length byte + text
		rdlen := 1 + len(txt)
		answer = append(answer, byte(rdlen>>8), byte(rdlen&0xff))
		// RDATA: length-prefixed string
		answer = append(answer, byte(len(txt)))
		answer = append(answer, []byte(txt)...)

	case 15: // MX record
		mx := encodeDomainName(fmt.Sprintf("mx.%s", domain))
		// RDLENGTH: 2 (preference) + domain length
		rdlen := 2 + len(mx)
		answer = append(answer, byte(rdlen>>8), byte(rdlen&0xff))
		// RDATA: preference (10) + domain
		answer = append(answer, 0x00, 0x0a) // preference: 10
		answer = append(answer, mx...)

	default:
		// Unsupported query type, return minimal answer
		log.Printf("Unsupported query type: %d", qtype)
		// RDLENGTH: 0
		answer = append(answer, 0x00, 0x00)
	}

	return answer
}

// encodeDomainName encodes a domain name in DNS format (length-prefixed labels)
func encodeDomainName(domain string) []byte {
	result := make([]byte, 0, len(domain)+2)

	start := 0
	for i := 0; i <= len(domain); i++ {
		if i == len(domain) || domain[i] == '.' {
			label := domain[start:i]
			result = append(result, byte(len(label)))
			result = append(result, []byte(label)...)
			start = i + 1
		}
	}

	// Null terminator
	result = append(result, 0x00)

	return result
}

func init() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ltime | log.Lmicroseconds)
}
