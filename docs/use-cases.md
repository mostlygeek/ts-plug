# Use Cases and Patterns

Real-world scenarios for using ts-plug and ts-unplug together.

## Table of Contents

- [Development Workflows](#development-workflows)
- [Testing Scenarios](#testing-scenarios)
- [Deployment Patterns](#deployment-patterns)
- [Team Collaboration](#team-collaboration)
- [Hybrid Cloud Architectures](#hybrid-cloud-architectures)

## Development Workflows

### Full-Stack Development with Remote Database

**Scenario:** You're developing a web app locally but want to use a shared staging database.

**Solution:** Use ts-unplug to bring the remote database to localhost:

```sh
# Terminal 1: Make remote database available locally
ts-unplug -dir ./state-db -port 5432 postgres-staging.tailnet.ts.net:5432

# Terminal 2: Run your app normally
DATABASE_URL=postgresql://localhost:5432/mydb npm run dev

# Terminal 3: Share your dev instance with teammates
ts-plug -hostname dev-yourname -https-port 443:3000 -- npm run dev
```

**Benefits:**
- Real data without database dumps
- Team can test your work instantly
- No VPN or complex networking

### Microservices Development

**Scenario:** You're working on one microservice that depends on several others.

**Solution:** Use ts-unplug for dependencies, ts-plug to share your service:

```sh
# Expose remote auth service locally
ts-unplug -dir ./state-auth -port 8001 auth-service.tailnet.ts.net &

# Expose remote payment service locally
ts-unplug -dir ./state-payment -port 8002 payment-service.tailnet.ts.net &

# Run your service with local env vars
export AUTH_URL=http://localhost:8001
export PAYMENT_URL=http://localhost:8002
go run main.go &

# Share your development service
ts-plug -hostname orders-dev -- go run main.go
```

### Mobile App Development

**Scenario:** Testing a mobile app that needs to hit your local backend.

**Solution:** Use ts-plug to expose your backend with a stable URL:

```sh
# Start backend with ts-plug
ts-plug -hostname mobile-api-dev -- npm run dev

# Configure mobile app to use:
# https://mobile-api-dev.tailnet.ts.net

# Mobile device must be on your Tailnet (install Tailscale app)
```

**Benefits:**
- No need to update URLs constantly
- Works from physical devices
- Proper HTTPS for realistic testing

## Testing Scenarios

### Webhook Testing

**Scenario:** Testing GitHub/Stripe/Twilio webhooks locally.

**Solution:** Use ts-plug with `-public` flag:

```sh
# Start your webhook handler
ts-plug -public -hostname webhook-test -- python webhook_server.py

# Copy the public URL and paste into webhook settings
# https://webhook-test.tailnet.ts.net

# Webhooks will hit your local server
```

**Benefits:**
- No third-party tunneling services
- Automatic TLS certificates
- Tailscale identity in headers

### End-to-End Testing Against Staging

**Scenario:** Running E2E tests against a staging API.

**Solution:** Use ts-unplug to make staging API appear local:

```sh
# Make staging API available at localhost
ts-unplug -dir ./state -port 8080 api-staging.tailnet.ts.net

# Run tests pointing to localhost
API_URL=http://localhost:8080 npm run test:e2e
```

**Benefits:**
- No need to change test configuration
- Faster than VPN
- Can run in CI with Tailscale

### Load Testing Remote Services

**Scenario:** Load test a service on your Tailnet.

**Solution:** Use ts-unplug to proxy, run local load testing tools:

```sh
# Proxy remote service
ts-unplug -dir ./state -port 8080 service-under-test.tailnet.ts.net

# Run load tests against localhost
ab -n 10000 -c 100 http://localhost:8080/api/endpoint
```

## Deployment Patterns

### Sidecar-Free Container Deployment

**Scenario:** Deploy containerized apps without Tailscale sidecar complexity.

**Solution:** Use ts-plug as the container entrypoint:

```dockerfile
FROM node:18
WORKDIR /app
COPY package*.json ./
RUN npm install
COPY . .

# Add ts-plug binary
COPY --from=ghcr.io/yourorg/ts-plug:latest /ts-plug /usr/local/bin/

# Use ts-plug as entrypoint
ENTRYPOINT ["ts-plug", "-hostname", "myapp", "-dir", "/var/lib/tsplug", "--"]
CMD ["npm", "start"]
```

**Docker Compose:**
```yaml
version: '3'
services:
  app:
    build: .
    environment:
      - TS_AUTHKEY=${TS_AUTHKEY}
    volumes:
      - tsplug-state:/var/lib/tsplug

volumes:
  tsplug-state:
```

**Benefits:**
- No separate sidecar container
- Simpler orchestration
- Automatic HTTPS

### Homelab Services

**Scenario:** Expose homelab services securely.

**Solution:** Use ts-plug for each service:

```sh
# Pi-hole
ts-plug -dns -http -hostname pihole -- pihole-FTL

# Jellyfin media server
ts-plug -hostname jellyfin -https-port 443:8096 -- jellyfin

# Home Assistant
ts-plug -hostname homeassistant -https-port 443:8123 -- hass
```

**Benefits:**
- No port forwarding
- Automatic HTTPS
- No dynamic DNS needed

### Temporary Demo Environments

**Scenario:** Share a demo with a client without deployment.

**Solution:** Use ts-plug with `-public`:

```sh
# Start demo environment
ts-plug -public -hostname demo-acme-corp -- npm start

# Share URL with client (no Tailscale required)
# https://demo-acme-corp.tailnet.ts.net
```

## Team Collaboration

### Code Review Testing

**Scenario:** Reviewer wants to test a branch without checking it out.

**Solution:** Developer shares their branch with ts-plug:

```sh
# Developer runs their branch
ts-plug -hostname feature-xyz-alice -- npm run dev

# Reviewer accesses in browser
# https://feature-xyz-alice.tailnet.ts.net

# No git checkout needed!
```

### Designer Preview Environments

**Scenario:** Designers need to preview work-in-progress features.

**Solution:** Each developer exposes their environment:

```sh
# Frontend developer
ts-plug -hostname frontend-bob -https-port 443:3000 -- npm run dev

# Backend developer uses ts-unplug to access Bob's frontend
ts-unplug -dir ./state -port 3000 frontend-bob.tailnet.ts.net:443

# Backend developer exposes API
ts-plug -hostname api-carol -- go run main.go
```

### Pair Programming Across Locations

**Scenario:** Remote pair programming with live server access.

**Solution:** Host shares their development environment:

```sh
# Host runs server
ts-plug -hostname pairing-session -- bundle exec rails server

# Participant accesses same server
# Both see changes in real-time
```

## Hybrid Cloud Architectures

### Local Development, Cloud Database

**Scenario:** Develop locally but use cloud-hosted database.

**Solution:**
```sh
# Put database on Tailscale (could be RDS with TS subnet router)
# Access it locally
ts-unplug -dir ./state -port 5432 rds-proxy.tailnet.ts.net:5432

# Run app locally
DATABASE_URL=postgresql://localhost:5432/prod npm run dev
```

### Multi-Cloud Service Access

**Scenario:** Services spread across AWS, GCP, on-prem.

**Solution:** Use Tailscale subnet routers and ts-unplug:

```sh
# Access AWS service
ts-unplug -dir ./state-aws -port 8001 aws-api.tailnet.ts.net &

# Access GCP service
ts-unplug -dir ./state-gcp -port 8002 gcp-api.tailnet.ts.net &

# Access on-prem service
ts-unplug -dir ./state-onprem -port 8003 onprem-api.tailnet.ts.net &

# Your app sees everything as localhost
export AWS_API=http://localhost:8001
export GCP_API=http://localhost:8002
export ONPREM_API=http://localhost:8003
./run-app.sh
```

### Edge Computing

**Scenario:** Deploy to edge locations (Raspberry Pi, IoT devices).

**Solution:** Run ts-plug on edge devices:

```sh
# On Raspberry Pi
ts-plug -hostname sensor-living-room -- python sensor.py

# On another Pi
ts-plug -hostname sensor-garage -- python sensor.py

# Access all sensors from central dashboard
# No complex networking, no public IPs
```

## Advanced Patterns

### Service Mesh Alternative

Use ts-plug/ts-unplug as a lightweight service mesh:

```sh
# Each service exposes itself with ts-plug
ts-plug -hostname service-a -- ./service-a
ts-plug -hostname service-b -- ./service-b
ts-plug -hostname service-c -- ./service-c

# Services discover each other via Tailscale DNS
curl https://service-b.tailnet.ts.net/api
```

**Benefits:**
- Built-in encryption
- Automatic service discovery
- Identity-based access
- No complex mesh configuration

### Gradual Migration

**Scenario:** Migrating from monolith to microservices.

**Old monolith:**
```sh
ts-plug -hostname legacy-monolith -- ./monolith
```

**New microservices consume monolith:**
```sh
# New service uses ts-unplug to access old monolith
ts-unplug -dir ./state -port 8080 legacy-monolith.tailnet.ts.net

# New service exposes itself
ts-plug -hostname new-service -- ./new-service
```

**Frontend can use both during migration:**
```javascript
// Old endpoints
fetch('https://legacy-monolith.tailnet.ts.net/api/users')

// New endpoints
fetch('https://new-service.tailnet.ts.net/api/users')
```

### Development Environment Orchestration

**Scenario:** Complex development setup script.

**dev-env.sh:**
```bash
#!/bin/bash

# Start remote services locally
ts-unplug -dir ./state-db -port 5432 postgres.tailnet.ts.net:5432 &
ts-unplug -dir ./state-redis -port 6379 redis.tailnet.ts.net:6379 &
ts-unplug -dir ./state-auth -port 8001 auth.tailnet.ts.net &

# Wait for proxies to start
sleep 2

# Start local services and expose them
ts-plug -hostname api-dev -- npm run dev:api &
ts-plug -hostname web-dev -https-port 443:3000 -- npm run dev:web &

echo "Development environment ready!"
echo "API: https://api-dev.tailnet.ts.net"
echo "Web: https://web-dev.tailnet.ts.net"
```

## Security Patterns

### Zero Trust Access

Leverage Tailscale's built-in identity:

```go
// In your service
func handler(w http.ResponseWriter, r *http.Request) {
    user := r.Header.Get("Tailscale-User-Login")
    if user == "" {
        http.Error(w, "Unauthorized", 401)
        return
    }

    // User is authenticated by Tailscale
    // Implement authorization based on user
    if !isAuthorized(user) {
        http.Error(w, "Forbidden", 403)
        return
    }

    // Handle request
}
```

### Temporary Access

Share access temporarily without changing firewall rules:

```sh
# Start service with ts-plug
ts-plug -public -hostname temp-demo -- python -m http.server 8080

# Share URL with external user
# When done, stop ts-plug
# Access automatically revoked
```

## See Also

- [ts-plug Guide](./ts-plug.md)
- [ts-unplug Guide](./ts-unplug.md)
- [Docker Examples](./docker.md)
- [Main README](../README.md)
