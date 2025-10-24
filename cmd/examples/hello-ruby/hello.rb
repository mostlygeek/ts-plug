#!/usr/bin/env ruby

require 'webrick'
require 'optparse'

# Parse command line arguments
options = { listen: 'localhost:8080' }
OptionParser.new do |opts|
  opts.banner = "Usage: hello.rb [options]"

  opts.on("--listen ADDR", "IP:port to listen on (default: localhost:8080)") do |addr|
    options[:listen] = addr
  end
end.parse!

# Parse listen address
listen_addr = options[:listen]
host = ''
port = 8080

if listen_addr.start_with?(':')
  port = listen_addr[1..-1].to_i
elsif listen_addr.include?(':')
  parts = listen_addr.rpartition(':')
  host = parts[0]
  port = parts[2].to_i
else
  port = listen_addr.to_i
end

# Get client IP from request
def get_client_ip(req)
  # Try to get IP from X-Forwarded-For header (for proxies)
  if req.header['x-forwarded-for']&.any?
    forwarded = req.header['x-forwarded-for'].first
    ips = forwarded.split(',')
    return ips[0].strip if ips.any?
  end

  # Try to get IP from X-Real-IP header
  if req.header['x-real-ip']&.any?
    return req.header['x-real-ip'].first
  end

  # Fall back to remote address
  req.peeraddr[3]
end

# Create HTTP server
server = WEBrick::HTTPServer.new(
  BindAddress: host.empty? ? nil : host,
  Port: port,
  Logger: WEBrick::Log.new($stdout, WEBrick::Log::INFO),
  AccessLog: []
)

# Define request handler
server.mount_proc '/' do |req, res|
  # Check if required headers exist
  login_name = req.header['tailscale-user-login']&.first
  display_name = req.header['tailscale-user-name']&.first
  profile_pic_url = req.header['tailscale-user-profile-pic']&.first

  res['Cache-Control'] = 'no-cache'

  if login_name && display_name && profile_pic_url
    # Serve HTML page with user information
    res['Content-Type'] = 'text/html'
    res.body = <<~HTML
      <!DOCTYPE html>
      <html>
      <head>
          <title>Hello from Ruby!</title>
      </head>
      <body>
          <h1>(Ruby) Tailscale User Information</h1>
          <p><strong>Login Name:</strong> #{login_name}</p>
          <p><strong>Name:</strong> #{display_name}</p>
          <p><strong>Profile Picture:</strong></p>
          <img src="#{profile_pic_url}" alt="Profile Picture" style="max-width: 200px;">
      </body>
      </html>
    HTML
  else
    # Print anonymous message with IP
    ip = get_client_ip(req)
    res['Content-Type'] = 'text/plain'
    res.body = "Hello anonymous from #{ip}\n"
  end
end

# Handle graceful shutdown
trap('INT') do
  puts "\nShutting down server..."
  server.shutdown
end

display_host = host.empty? ? '0.0.0.0' : host
puts "Starting server on #{display_host}:#{port}"

server.start
