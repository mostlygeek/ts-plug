#!/bin/sh

# Default listen address
LISTEN_ADDR="localhost:8080"

# Parse command line arguments
while [ $# -gt 0 ]; do
  case "$1" in
    --listen)
      LISTEN_ADDR="$2"
      shift 2
      ;;
    *)
      echo "Unknown option: $1"
      echo "Usage: $0 [--listen address:port]"
      exit 1
      ;;
  esac
done

# Parse listen address
HOST=""
PORT=8080

case "$LISTEN_ADDR" in
  :*)
    PORT="${LISTEN_ADDR#:}"
    ;;
  *:*)
    HOST="${LISTEN_ADDR%:*}"
    PORT="${LISTEN_ADDR##*:}"
    ;;
  *)
    PORT="$LISTEN_ADDR"
    ;;
esac

# Get client IP from headers
get_client_ip() {
  x_forwarded_for="$1"
  x_real_ip="$2"

  # Try X-Forwarded-For first
  if [ -n "$x_forwarded_for" ]; then
    echo "$x_forwarded_for" | cut -d',' -f1 | sed 's/^[[:space:]]*//;s/[[:space:]]*$//'
    return
  fi

  # Try X-Real-IP
  if [ -n "$x_real_ip" ]; then
    echo "$x_real_ip"
    return
  fi

  # Fall back to unknown
  echo "unknown"
}

# HTTP request handler - reads from stdin, writes to stdout
handle_request() {
  # Read request line
  read -r method path protocol

  # Initialize header variables
  login_name=""
  display_name=""
  profile_pic_url=""
  x_forwarded_for=""
  x_real_ip=""

  # Read headers line by line
  while IFS= read -r line; do
    # Remove carriage return
    line="${line%$(printf '\r')}"

    # Empty line marks end of headers
    [ -z "$line" ] && break

    # Extract header name and value
    header_name="${line%%:*}"
    header_value="${line#*: }"

    # Convert header name to lowercase for comparison
    header_name_lower=$(echo "$header_name" | tr '[:upper:]' '[:lower:]')

    # Check for specific headers we care about
    case "$header_name_lower" in
      tailscale-user-login)
        login_name="$header_value"
        ;;
      tailscale-user-name)
        display_name="$header_value"
        ;;
      tailscale-user-profile-pic)
        profile_pic_url="$header_value"
        ;;
      x-forwarded-for)
        x_forwarded_for="$header_value"
        ;;
      x-real-ip)
        x_real_ip="$header_value"
        ;;
    esac
  done

  # Generate response based on whether we have Tailscale headers
  if [ -n "$login_name" ] && [ -n "$display_name" ] && [ -n "$profile_pic_url" ]; then
    # Serve HTML page with user information
    body="<!DOCTYPE html>
<html>
<head>
    <title>Hello from Shell!</title>
</head>
<body>
    <h1>(Shell) Tailscale User Information</h1>
    <p><strong>Login Name:</strong> $login_name</p>
    <p><strong>Name:</strong> $display_name</p>
    <p><strong>Profile Picture:</strong></p>
    <img src=\"$profile_pic_url\" alt=\"Profile Picture\" style=\"max-width: 200px;\">
</body>
</html>"
    content_type="text/html"
  else
    # Print anonymous message with IP
    client_ip=$(get_client_ip "$x_forwarded_for" "$x_real_ip")
    body="Hello anonymous from $client_ip"
    content_type="text/plain"
  fi

  # Calculate content length
  content_length=$(printf "%s" "$body" | wc -c | tr -d ' ')

  # Send HTTP response
  printf "HTTP/1.1 200 OK\r\n"
  printf "Content-Type: %s\r\n" "$content_type"
  printf "Content-Length: %s\r\n" "$content_length"
  printf "Cache-Control: no-cache\r\n"
  printf "Connection: close\r\n"
  printf "\r\n"
  printf "%s" "$body"
}

# Start server
start_server() {
  bind_addr="$HOST"
  [ -z "$bind_addr" ] || [ "$bind_addr" = "localhost" ] && bind_addr="127.0.0.1"

  # Create response FIFO
  response_fifo="/tmp/http_response_$$"
  mkfifo "$response_fifo"

  # Cleanup on exit
  trap "rm -f $response_fifo; echo 'Shutting down server...' >&2; exit 0" INT TERM EXIT

  echo "Starting server on ${bind_addr}:${PORT}" >&2

  # Server loop
  while true; do
    # Start nc listening in background, feed it from our response fifo
    nc -l "$bind_addr" "$PORT" < "$response_fifo" | handle_request > "$response_fifo" &
    wait $!
  done
}

start_server
