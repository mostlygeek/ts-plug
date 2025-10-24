#!/usr/bin/env node

import http from "http";
import { parseArgs } from "util";

// Parse command line arguments
const { values } = parseArgs({
  options: {
    listen: {
      type: "string",
      short: "l",
      default: "localhost:8080",
    },
  },
});

// Parse listen address
const listenAddr = values.listen;
let host = "";
let port = 8080;

if (listenAddr.startsWith(":")) {
  port = parseInt(listenAddr.slice(1), 10);
} else if (listenAddr.includes(":")) {
  [host, port] = listenAddr.split(":");
  port = parseInt(port, 10);
} else {
  port = parseInt(listenAddr, 10);
}

function getClientIP(req) {
  // Try to get IP from X-Forwarded-For header (for proxies)
  const forwarded = req.headers["x-forwarded-for"];
  if (forwarded) {
    const ips = forwarded.split(",");
    if (ips.length > 0) {
      return ips[0].trim();
    }
  }

  // Try to get IP from X-Real-IP header
  const realIP = req.headers["x-real-ip"];
  if (realIP) {
    return realIP;
  }

  // Fall back to remote address
  const remoteAddr = req.socket.remoteAddress;
  return remoteAddr || "unknown";
}

const server = http.createServer((req, res) => {
  // Check if required headers exist
  const loginName = req.headers["tailscale-user-login"];
  const displayName = req.headers["tailscale-user-name"];
  const profilePicURL = req.headers["tailscale-user-profile-pic"];

  res.setHeader("Cache-Control", "no-cache");

  if (loginName && displayName && profilePicURL) {
    // Serve HTML page with user information
    res.setHeader("Content-Type", "text/html");
    res.writeHead(200);
    res.end(`
<!DOCTYPE html>
<html>
<head>
    <title>Hello from node.js!</title>
</head>
<body>
    <h1>(NodeJS) Tailscale User Information</h1>
    <p><strong>Login Name:</strong> ${loginName}</p>
    <p><strong>Name:</strong> ${displayName}</p>
    <p><strong>Profile Picture:</strong></p>
    <img src="${profilePicURL}" alt="Profile Picture" style="max-width: 200px;">
</body>
</html>
`);
  } else {
    // Print anonymous message with IP
    const ip = getClientIP(req);
    res.setHeader("Content-Type", "text/plain");
    res.writeHead(200);
    res.end(`Hello anonymous from ${ip}\n`);
  }
});

server.listen(port, host, () => {
  const addr = host || "localhost";
  console.log(`Starting server on ${addr}:${port}`);
});
