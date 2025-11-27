# Docker Integration Guide

Using ts-plug to eliminate Tailscale sidecar containers and simplify container networking.

## Overview

Traditional Tailscale container deployment requires either:
- A sidecar container running Tailscale
- Complex network sharing between containers
- Host network mode (which breaks container isolation)

**ts-plug eliminates this complexity** by combining your application and Tailscale connectivity in a single container.

## Basic Pattern

### Traditional Approach (Sidecar)

```yaml
version: '3'
services:
  tailscale:
    image: tailscale/tailscale:latest
    environment:
      - TS_AUTHKEY=${TS_AUTHKEY}
    volumes:
      - tailscale-state:/var/lib/tailscale
    cap_add:
      - NET_ADMIN

  app:
    image: myapp:latest
    network_mode: "service:tailscale"
    depends_on:
      - tailscale

volumes:
  tailscale-state:
```

### ts-plug Approach (No Sidecar)

```yaml
version: '3'
services:
  app:
    image: myapp-with-tsplug:latest
    environment:
      - TS_AUTHKEY=${TS_AUTHKEY}
    volumes:
      - tsplug-state:/var/lib/tsplug

volumes:
  tsplug-state:
```

**Benefits:**
- Single container instead of two
- Simpler orchestration
- No special network modes
- Easier debugging

## Building Images with ts-plug

### Method 1: Multi-Stage Build

```dockerfile
# Build ts-plug
FROM golang:1.21 AS tsplug-builder
WORKDIR /build
RUN git clone https://github.com/tailscale/tsplug.git
WORKDIR /build/tsplug
RUN make ts-plug

# Build your app
FROM node:18 AS app-builder
WORKDIR /app
COPY package*.json ./
RUN npm install
COPY . .
RUN npm run build

# Final image
FROM node:18-slim
WORKDIR /app

# Copy ts-plug binary
COPY --from=tsplug-builder /build/tsplug/build/ts-plug /usr/local/bin/

# Copy your app
COPY --from=app-builder /app/dist ./dist
COPY --from=app-builder /app/node_modules ./node_modules
COPY package*.json ./

# Use ts-plug as entrypoint
ENTRYPOINT ["ts-plug", "-hostname", "myapp", "-dir", "/var/lib/tsplug", "--"]
CMD ["npm", "start"]
```

### Method 2: Copy Pre-Built Binary

```dockerfile
FROM node:18
WORKDIR /app

# Copy pre-built ts-plug (build it separately)
COPY ts-plug /usr/local/bin/
RUN chmod +x /usr/local/bin/ts-plug

# Copy your app
COPY package*.json ./
RUN npm install
COPY . .

ENTRYPOINT ["ts-plug", "-hostname", "myapp", "-dir", "/var/lib/tsplug", "--"]
CMD ["npm", "start"]
```

### Method 3: Base Image

Create a base image with ts-plug:

```dockerfile
# base.Dockerfile
FROM ubuntu:22.04
RUN apt-get update && apt-get install -y ca-certificates
COPY ts-plug /usr/local/bin/
RUN chmod +x /usr/local/bin/ts-plug
```

Then use it:

```dockerfile
# app.Dockerfile
FROM myregistry/tsplug-base:latest
# ... your app setup ...
ENTRYPOINT ["ts-plug", "-hostname", "myapp", "--"]
CMD ["./myapp"]
```

## Real-World Examples

### Pi-hole DNS Server

**Dockerfile:**
```dockerfile
FROM pihole/pihole:latest

# Add ts-plug
COPY ts-plug /usr/local/bin/
RUN chmod +x /usr/local/bin/ts-plug

# Override entrypoint
ENTRYPOINT ["ts-plug", \
    "-hostname", "pihole", \
    "-dir", "/var/lib/tsplug", \
    "-dns", \
    "-http", \
    "--", \
    "/s6-init"]
```

**docker-compose.yml:**
```yaml
version: '3'
services:
  pihole:
    build: .
    environment:
      - TS_AUTHKEY=${TS_AUTHKEY}
      - WEBPASSWORD=admin
    volumes:
      - pihole-config:/etc/pihole
      - pihole-dnsmasq:/etc/dnsmasq.d
      - tsplug-state:/var/lib/tsplug

volumes:
  pihole-config:
  pihole-dnsmasq:
  tsplug-state:
```

**Usage:**
```sh
# Build and run
docker-compose up -d

# Access web interface at:
# https://pihole.tailnet.ts.net

# Configure devices to use DNS:
# pihole.tailnet.ts.net
```

See [docker/pi-hole/](../docker/pi-hole/) for the complete example.

### Node.js Web Application

**Dockerfile:**
```dockerfile
FROM node:18
WORKDIR /app

# Install ts-plug
COPY ts-plug /usr/local/bin/
RUN chmod +x /usr/local/bin/ts-plug

# Install dependencies
COPY package*.json ./
RUN npm ci --only=production

# Copy app
COPY . .

# Expose via ts-plug
ENTRYPOINT ["ts-plug", \
    "-hostname", "webapp", \
    "-dir", "/var/lib/tsplug", \
    "-https-port", "443:3000", \
    "--"]
CMD ["npm", "start"]
```

**Run:**
```sh
docker build -t myapp .
docker run -d \
  -e TS_AUTHKEY=tskey-auth-xxx \
  -v tsplug-state:/var/lib/tsplug \
  myapp
```

### Python Flask API

**Dockerfile:**
```dockerfile
FROM python:3.11-slim
WORKDIR /app

# Install ts-plug
COPY ts-plug /usr/local/bin/
RUN chmod +x /usr/local/bin/ts-plug

# Install dependencies
COPY requirements.txt .
RUN pip install -r requirements.txt

# Copy app
COPY . .

# Run with ts-plug
ENTRYPOINT ["ts-plug", \
    "-hostname", "api", \
    "-dir", "/var/lib/tsplug", \
    "-https-port", "443:5000", \
    "--"]
CMD ["python", "app.py"]
```

### Static Site with nginx

**Dockerfile:**
```dockerfile
FROM nginx:alpine

# Install ts-plug
COPY ts-plug /usr/local/bin/
RUN chmod +x /usr/local/bin/ts-plug

# Copy static files
COPY dist/ /usr/share/nginx/html/

# Expose via ts-plug
ENTRYPOINT ["ts-plug", \
    "-hostname", "website", \
    "-dir", "/var/lib/tsplug", \
    "-public", \
    "--"]
CMD ["nginx", "-g", "daemon off;"]
```

This makes your static site publicly accessible!

## Configuration Patterns

### Environment-Based Hostname

```dockerfile
ENTRYPOINT ["sh", "-c", "exec ts-plug -hostname ${HOSTNAME:-defaultapp} -dir /var/lib/tsplug -- npm start"]
```

```sh
docker run -e HOSTNAME=myapp-staging myimage
```

### Multi-Protocol Support

```dockerfile
# Support HTTP, HTTPS, and DNS
ENTRYPOINT ["ts-plug", \
    "-hostname", "multiservice", \
    "-dir", "/var/lib/tsplug", \
    "-http", \
    "-https", \
    "-dns", \
    "--"]
CMD ["./myserver"]
```

### Public Access Toggle

```dockerfile
ENTRYPOINT ["sh", "-c", \
    "exec ts-plug -hostname ${HOSTNAME:-app} -dir /var/lib/tsplug ${PUBLIC:+-public} -- npm start"]
```

```sh
# Private (default)
docker run myimage

# Public
docker run -e PUBLIC=true myimage
```

## Docker Compose Examples

### Full Stack Application

```yaml
version: '3.8'

services:
  frontend:
    build: ./frontend
    environment:
      - TS_AUTHKEY=${TS_AUTHKEY}
      - HOSTNAME=frontend
    volumes:
      - frontend-state:/var/lib/tsplug

  api:
    build: ./api
    environment:
      - TS_AUTHKEY=${TS_AUTHKEY}
      - HOSTNAME=api
      - DATABASE_URL=postgresql://db:5432/mydb
    volumes:
      - api-state:/var/lib/tsplug
    depends_on:
      - db

  db:
    image: postgres:15
    environment:
      - POSTGRES_PASSWORD=secret
    volumes:
      - db-data:/var/lib/postgresql/data

volumes:
  frontend-state:
  api-state:
  db-data:
```

**Access:**
- Frontend: `https://frontend.tailnet.ts.net`
- API: `https://api.tailnet.ts.net`
- Database: private (only accessible to api container)

### Multiple Environments

```yaml
# docker-compose.staging.yml
version: '3.8'

services:
  app:
    build: .
    environment:
      - TS_AUTHKEY=${TS_AUTHKEY_STAGING}
      - HOSTNAME=app-staging
    volumes:
      - app-staging-state:/var/lib/tsplug

volumes:
  app-staging-state:
```

```yaml
# docker-compose.prod.yml
version: '3.8'

services:
  app:
    build: .
    environment:
      - TS_AUTHKEY=${TS_AUTHKEY_PROD}
      - HOSTNAME=app-prod
      - PUBLIC=true
    volumes:
      - app-prod-state:/var/lib/tsplug

volumes:
  app-prod-state:
```

```sh
# Deploy staging
docker-compose -f docker-compose.staging.yml up -d

# Deploy production
docker-compose -f docker-compose.prod.yml up -d
```

## Kubernetes Integration

### Deployment

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: myapp
spec:
  replicas: 1
  selector:
    matchLabels:
      app: myapp
  template:
    metadata:
      labels:
        app: myapp
    spec:
      containers:
      - name: app
        image: myapp-with-tsplug:latest
        env:
        - name: TS_AUTHKEY
          valueFrom:
            secretKeyRef:
              name: tailscale-auth
              key: authkey
        - name: HOSTNAME
          value: "myapp-k8s"
        volumeMounts:
        - name: tsplug-state
          mountPath: /var/lib/tsplug
      volumes:
      - name: tsplug-state
        emptyDir: {}
```

### StatefulSet for Persistent State

```yaml
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: myapp
spec:
  serviceName: myapp
  replicas: 1
  selector:
    matchLabels:
      app: myapp
  template:
    metadata:
      labels:
        app: myapp
    spec:
      containers:
      - name: app
        image: myapp-with-tsplug:latest
        env:
        - name: TS_AUTHKEY
          valueFrom:
            secretKeyRef:
              name: tailscale-auth
              key: authkey
        volumeMounts:
        - name: tsplug-state
          mountPath: /var/lib/tsplug
  volumeClaimTemplates:
  - metadata:
      name: tsplug-state
    spec:
      accessModes: [ "ReadWriteOnce" ]
      resources:
        requests:
          storage: 1Gi
```

## Troubleshooting

### Container Starts But Can't Connect to Tailscale

**Check auth key:**
```sh
docker logs <container-id>
# Look for authentication errors
```

**Verify auth key is valid:**
```sh
# Generate a new auth key in Tailscale admin console
docker run -e TS_AUTHKEY=tskey-auth-NEW_KEY myimage
```

### State Directory Permissions

```dockerfile
# Ensure state directory is writable
RUN mkdir -p /var/lib/tsplug && chmod 700 /var/lib/tsplug
```

### ts-plug Not Found

```dockerfile
# Verify binary is executable
COPY ts-plug /usr/local/bin/
RUN chmod +x /usr/local/bin/ts-plug
RUN ls -la /usr/local/bin/ts-plug
```

### Application Not Starting

```sh
# Test without ts-plug first
docker run myimage npm start

# Then test with ts-plug
docker run myimage ts-plug -hostname test -dir /tmp -- npm start
```

### Check Logs

```sh
# View logs
docker logs -f <container-id>

# Enable debug logging
docker run -e LOG_LEVEL=debug myimage
```

Update Dockerfile:
```dockerfile
ENTRYPOINT ["ts-plug", \
    "-log", "${LOG_LEVEL:-info}", \
    "-hostname", "myapp", \
    "-dir", "/var/lib/tsplug", \
    "--"]
```

## Best Practices

### 1. Use Auth Keys Properly

- Generate ephemeral auth keys for development
- Use reusable auth keys for production
- Store auth keys in secrets management (never in images)

### 2. Persist State Correctly

- Always mount `/var/lib/tsplug` as a volume
- Use named volumes for easier management
- In K8s, use PersistentVolumeClaims for StatefulSets

### 3. Security

```dockerfile
# Run as non-root when possible
RUN useradd -m -u 1000 appuser
USER appuser

# ts-plug doesn't require root privileges
```

### 4. Health Checks

```dockerfile
HEALTHCHECK --interval=30s --timeout=3s \
  CMD curl -f http://localhost:8080/health || exit 1
```

### 5. Graceful Shutdown

ts-plug handles signals properly, but ensure your app does too:

```dockerfile
# Use exec form to properly receive signals
CMD ["npm", "start"]

# Not:
# CMD npm start  # This creates a shell that doesn't forward signals
```

## See Also

- [ts-plug Guide](./ts-plug.md)
- [Use Cases](./use-cases.md)
- [Docker Examples](../docker/)
- [Main README](../README.md)
