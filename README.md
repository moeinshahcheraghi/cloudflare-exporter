# Cloudflare Prometheus Exporter

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Docker](https://img.shields.io/badge/docker-ready-brightgreen.svg)](Dockerfile)

A production-ready Prometheus exporter for Cloudflare zone metrics with comprehensive monitoring capabilities. This exporter collects detailed analytics from Cloudflare's GraphQL API and exposes them in Prometheus format for visualization and alerting.

##  Table of Contents

- [Features](#-features)
- [Architecture](#-architecture)
- [Prerequisites](#-prerequisites)
- [Installation](#-installation)
- [Configuration](#-configuration)
- [Usage](#-usage)
- [Metrics](#-metrics)
- [Grafana Dashboards](#-grafana-dashboards)
- [Development](#-development)
- [Troubleshooting](#-troubleshooting)
- [Contributing](#-contributing)
- [License](#-license)

##  Features

### Core Capabilities

- **Comprehensive Metrics Collection**: HTTP requests, bandwidth, cache statistics, performance metrics
- **Security Monitoring**: Firewall events, threats, attack patterns
- **Geographic Analytics**: Country-based request and bandwidth distribution
- **Content Type Analysis**: Detailed breakdown by content type
- **Status Code Tracking**: HTTP status code distribution and aggregation
- **High Performance**: Efficient API usage with configurable scrape intervals
- **Production Ready**: Health checks, proper error handling, structured logging

### Technical Features

- Modular, maintainable codebase following Go best practices
- Docker support with multi-stage builds
- Prometheus metrics with proper naming conventions
- Graceful error handling and recovery
- Configurable via environment variables
- Health check endpoint for orchestration
- Minimal resource footprint

##  Architecture

```
cloudflare-exporter/
├── cmd/
│   └── exporter/           # Main application entry point
│       └── main.go
├── internal/               # Private application code
│   ├── collector/          # Metrics collection logic
│   │   ├── basic.go        # HTTP and performance metrics
│   │   ├── status.go       # Status code metrics
│   │   ├── contenttype.go  # Content type metrics
│   │   └── firewall.go     # Firewall metrics
│   ├── config/             # Configuration management
│   │   └── config.go
│   └── metrics/            # Prometheus metrics definitions
│       └── metrics.go
├── pkg/                    # Public, reusable packages
│   └── cloudflare/         # Cloudflare API client
│       └── client.go
├── Dockerfile              # Multi-stage Docker build
├── docker-compose.yml      # Docker Compose configuration
├── Makefile                # Build automation
└── go.mod                  # Go module definition
```

### Design Principles

1. **Separation of Concerns**: Clear boundaries between API client, metrics collection, and business logic
2. **Testability**: Modular design enables easy unit testing
3. **Extensibility**: Simple to add new metric collectors
4. **Observability**: Structured logging with clear status indicators
5. **Security**: Follows container security best practices

##  Prerequisites

### Required

- **Cloudflare Account**: Pro plan or higher (for firewall metrics)
- **API Token**: With Analytics read permission
- **Zone ID**: Your Cloudflare zone identifier

### Optional

- Go 1.21+ (for local development)
- Docker & Docker Compose (for containerized deployment)
- Prometheus server (for metrics collection)
- Grafana (for visualization)

##  Installation

### Method 1: Docker (Recommended)

1. **Clone the repository**:
   ```bash
   git clone https://github.com/moeinshahcheraghi/cloudflare-exporter.git
   cd cloudflare-exporter
   ```

2. **Create environment file**:
   ```bash
   cp .env.example .env
   # Edit .env with your credentials
   ```

3. **Start the exporter**:
   ```bash
   docker-compose up -d
   ```

4. **Verify it's running**:
   ```bash
   curl http://localhost:9199/health
   curl http://localhost:9199/metrics
   ```

### Method 2: Build with Docker

1. **Build the image**:
   ```bash
   docker build -t cloudflare-exporter:latest .
   ```

2. **Run the container**:
   ```bash
   docker run -d \
     --name cloudflare-exporter \
     -p 9199:9199 \
     -e CLOUDFLARE_API_TOKEN="your_token_here" \
     -e CLOUDFLARE_ZONE_ID="your_zone_id_here" \
     cloudflare-exporter:latest
   ```

3. **Check logs**:
   ```bash
   docker logs -f cloudflare-exporter
   ```

### Method 3: Build from Source

1. **Clone the repository**:
   ```bash
   git clone https://github.com/moeinshahcheraghi/cloudflare-exporter.git
   cd cloudflare-exporter
   ```

2. **Download dependencies**:
   ```bash
   go mod download
   ```

3. **Build the binary**:
   ```bash
   go build -o cloudflare-exporter ./cmd/exporter
   ```

4. **Run**:
   ```bash
   export CLOUDFLARE_API_TOKEN="your_token_here"
   export CLOUDFLARE_ZONE_ID="your_zone_id_here"
   ./cloudflare-exporter
   ```

##  Configuration

### Environment Variables

| Variable | Description | Required | Default |
|----------|-------------|----------|---------|
| `CLOUDFLARE_API_TOKEN` | Cloudflare API token with Analytics:Read permission | Yes | - |
| `CLOUDFLARE_ZONE_ID` | Zone ID to monitor | Yes | - |
| `EXPORTER_PORT` | Port to expose metrics on | No | `9199` |

### Getting Cloudflare Credentials

#### 1. Create API Token

1. Log into [Cloudflare Dashboard](https://dash.cloudflare.com)
2. Go to **My Profile** → **API Tokens**
3. Click **Create Token**
4. Use the **Read Analytics** template or create custom token with:
   - Permissions: `Zone:Analytics:Read`
   - Zone Resources: `Include → Specific zone → [Your Zone]`
5. Copy the token (you won't see it again!)

#### 2. Get Zone ID

1. In Cloudflare Dashboard, select your domain
2. Scroll down in the **Overview** tab
3. Find **Zone ID** in the API section (right sidebar)

### Prometheus Configuration

Add to your `prometheus.yml`:

```yaml
scrape_configs:
  - job_name: 'cloudflare'
    static_configs:
      - targets: ['localhost:9199']
    scrape_interval: 60s
    scrape_timeout: 30s
```

##  Usage

### Local Development

```bash
# Run directly with Go
go run ./cmd/exporter

# Or build first
go build -o cloudflare-exporter ./cmd/exporter
./cloudflare-exporter
```

### Docker Commands

```bash
# Build image
docker build -t cloudflare-exporter:latest .

# Start with docker-compose
docker-compose up -d

# View logs
docker-compose logs -f cloudflare-exporter

# Stop container
docker-compose down

# Restart
docker-compose restart cloudflare-exporter

# Rebuild after changes
docker-compose up -d --build

# Run without docker-compose
docker run -d \
  --name cloudflare-exporter \
  -p 9199:9199 \
  -e CLOUDFLARE_API_TOKEN="your_token" \
  -e CLOUDFLARE_ZONE_ID="your_zone_id" \
  cloudflare-exporter:latest
```

### Go Commands

```bash
# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests for specific package
go test -v ./internal/collector/

# Format code
go fmt ./...

# Tidy dependencies
go mod tidy

# Download dependencies
go mod download

# Build for Linux
GOOS=linux GOARCH=amd64 go build -o cloudflare-exporter ./cmd/exporter

# Build for Windows
GOOS=windows GOARCH=amd64 go build -o cloudflare-exporter.exe ./cmd/exporter

# Build with optimizations
go build -ldflags="-w -s" -o cloudflare-exporter ./cmd/exporter
```

### Kubernetes Deployment

Create a deployment manifest:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: cloudflare-exporter
  namespace: monitoring
spec:
  replicas: 1
  selector:
    matchLabels:
      app: cloudflare-exporter
  template:
    metadata:
      labels:
        app: cloudflare-exporter
    spec:
      containers:
      - name: exporter
        image: cloudflare-exporter:latest
        ports:
        - containerPort: 9199
          name: metrics
        env:
        - name: CLOUDFLARE_API_TOKEN
          valueFrom:
            secretKeyRef:
              name: cloudflare-credentials
              key: api-token
        - name: CLOUDFLARE_ZONE_ID
          valueFrom:
            configMapKeyRef:
              name: cloudflare-config
              key: zone-id
        livenessProbe:
          httpGet:
            path: /health
            port: 9199
          initialDelaySeconds: 10
          periodSeconds: 30
        readinessProbe:
          httpGet:
            path: /health
            port: 9199
          initialDelaySeconds: 5
          periodSeconds: 10
---
apiVersion: v1
kind: Service
metadata:
  name: cloudflare-exporter
  namespace: monitoring
  labels:
    app: cloudflare-exporter
spec:
  type: ClusterIP
  ports:
  - port: 9199
    targetPort: 9199
    name: metrics
  selector:
    app: cloudflare-exporter
```

## 📈 Metrics

### HTTP Metrics

| Metric | Type | Description |
|--------|------|-------------|
| `cloudflare_zone_requests_total` | Gauge | Total number of requests |
| `cloudflare_zone_requests_cached` | Gauge | Number of cached requests |
| `cloudflare_zone_requests_uncached` | Gauge | Number of uncached requests |
| `cloudflare_zone_requests_encrypted` | Gauge | Number of HTTPS requests |
| `cloudflare_zone_pageviews_total` | Gauge | Total page views |
| `cloudflare_zone_cache_hit_rate_percent` | Gauge | Cache hit rate percentage |
| `cloudflare_zone_encryption_rate_percent` | Gauge | HTTPS usage percentage |

### Bandwidth Metrics

| Metric | Type | Description |
|--------|------|-------------|
| `cloudflare_zone_bandwidth_total_bytes` | Gauge | Total bandwidth in bytes |
| `cloudflare_zone_bandwidth_cached_bytes` | Gauge | Cached bandwidth in bytes |
| `cloudflare_zone_bandwidth_uncached_bytes` | Gauge | Uncached bandwidth in bytes |
| `cloudflare_zone_bandwidth_encrypted_bytes` | Gauge | Encrypted bandwidth in bytes |

### Geographic Metrics

| Metric | Type | Labels | Description |
|--------|------|--------|-------------|
| `cloudflare_zone_requests_country` | Gauge | `country` | Requests by country |
| `cloudflare_zone_bandwidth_country_bytes` | Gauge | `country` | Bandwidth by country |

### Status Code Metrics

| Metric | Type | Labels | Description |
|--------|------|--------|-------------|
| `cloudflare_zone_edge_response_status` | Gauge | `status` | Requests by specific status code |
| `cloudflare_zone_status_2xx_total` | Gauge | - | Total 2xx responses |
| `cloudflare_zone_status_3xx_total` | Gauge | - | Total 3xx responses |
| `cloudflare_zone_status_4xx_total` | Gauge | - | Total 4xx responses |
| `cloudflare_zone_status_5xx_total` | Gauge | - | Total 5xx responses |

### Content Type Metrics

| Metric | Type | Labels | Description |
|--------|------|--------|-------------|
| `cloudflare_zone_requests_content_type` | Gauge | `content_type` | Requests by content type |
| `cloudflare_zone_bandwidth_content_type_bytes` | Gauge | `content_type` | Bandwidth by content type |

### Security Metrics

| Metric | Type | Labels | Description |
|--------|------|--------|-------------|
| `cloudflare_zone_threats_total` | Gauge | - | Total threats detected |
| `cloudflare_zone_firewall_events_total` | Gauge | - | Total firewall events |
| `cloudflare_zone_firewall_action` | Gauge | `action` | Firewall events by action |
| `cloudflare_zone_firewall_source` | Gauge | `source` | Firewall events by source |
| `cloudflare_zone_firewall_country` | Gauge | `country` | Firewall events by country |
| `cloudflare_zone_firewall_ip` | Gauge | `ip` | Top 100 attacking IPs |

### Performance Metrics

| Metric | Type | Description |
|--------|------|-------------|
| `cloudflare_zone_client_wait_time_total_ms` | Gauge | Total client wait time (ms) |
| `cloudflare_zone_client_wait_time_avg_ms` | Gauge | Average wait time per request (ms) |

## 📊 Grafana Dashboards

### Sample PromQL Queries

**Request Rate**:
```promql
rate(cloudflare_zone_requests_total[5m])
```

**Cache Hit Rate Over Time**:
```promql
cloudflare_zone_cache_hit_rate_percent
```

**Top Countries by Traffic**:
```promql
topk(10, cloudflare_zone_requests_country)
```

**Error Rate**:
```promql
(cloudflare_zone_status_4xx_total + cloudflare_zone_status_5xx_total) / cloudflare_zone_requests_total * 100
```

**Bandwidth by Content Type**:
```promql
sum by (content_type) (cloudflare_zone_bandwidth_content_type_bytes)
```

**Firewall Block Rate**:
```promql
cloudflare_zone_firewall_action{action="block"}
```

### Dashboard Panels

Create panels in Grafana for:

1. **Overview**
   - Total requests (stat)
   - Cache hit rate (gauge)
   - Bandwidth usage (graph)
   - Request rate (graph)

2. **Performance**
   - Response time (graph)
   - Status code distribution (pie chart)
   - Top content types (bar gauge)

3. **Security**
   - Threats detected (stat)
   - Firewall events (graph)
   - Top attacking countries (geo map)
   - Top attacking IPs (table)

4. **Geographic**
   - Requests by country (world map)
   - Bandwidth by region (bar chart)

## 🛠 Development

### Project Structure

```
cloudflare-exporter/
├── cmd/exporter/           # Application entry point
├── internal/               # Private packages
│   ├── collector/          # Metric collectors (basic, status, content, firewall)
│   ├── config/             # Configuration loading
│   └── metrics/            # Prometheus metric definitions
├── pkg/cloudflare/         # Cloudflare API client (reusable)
├── Dockerfile              # Container image
├── docker-compose.yml      # Local development
└── go.mod                  # Dependencies
```

### Adding New Metrics

1. **Define metric in `internal/metrics/metrics.go`**:
   ```go
   NewMetric: prometheus.NewGaugeVec(
       prometheus.GaugeOpts{
           Name: "cloudflare_zone_new_metric",
           Help: "Description of new metric",
       },
       []string{"zone_id", "label"},
   ),
   ```

2. **Register in `Register()` method**

3. **Create collector in `internal/collector/`**:
   ```go
   func (c *Collector) CollectNewMetric() error {
       // Query Cloudflare API
       // Process data
       // Set metrics
   }
   ```

4. **Call from `CollectAll()` in `collector/basic.go`**

### Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with verbose output
go test -v ./...

# Run specific package tests
go test -v ./internal/collector/

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Code Quality

```bash
# Format code
go fmt ./...

# Vet code for issues
go vet ./...

# Tidy dependencies
go mod tidy

# Update dependencies
go get -u ./...

# Run static analysis (if golangci-lint installed)
golangci-lint run
```

### Building

```bash
# Build for current platform
go build -o cloudflare-exporter ./cmd/exporter

# Build for Linux
GOOS=linux GOARCH=amd64 go build -o cloudflare-exporter ./cmd/exporter

# Build for Windows
GOOS=windows GOARCH=amd64 go build -o cloudflare-exporter.exe ./cmd/exporter

# Build with size optimization
go build -ldflags="-w -s" -o cloudflare-exporter ./cmd/exporter

# Build Docker image
docker build -t cloudflare-exporter:latest .

# Build with specific tag
docker build -t cloudflare-exporter:v1.0.0 .
```

## 🔧 Troubleshooting

### Common Issues

#### 1. "No zones found"

**Cause**: Invalid Zone ID or API token lacks permissions

**Solution**:
```bash
# Verify Zone ID in Cloudflare dashboard
# Check API token has Zone:Analytics:Read permission
# Test API token:
curl -X GET "https://api.cloudflare.com/client/v4/zones/YOUR_ZONE_ID" \
  -H "Authorization: Bearer YOUR_API_TOKEN"
```

#### 2. "Firewall metrics not available"

**Cause**: Firewall metrics require Pro plan or higher

**Solution**: This is expected on Free plans. The exporter will log this as info and continue.

#### 3. High memory usage

**Cause**: Too many unique label values (IPs, user agents, etc.)

**Solution**: The exporter limits high-cardinality metrics (top 100 IPs, top 20 user agents). Adjust in `collector/firewall.go` if needed.

#### 4. Connection timeout

**Cause**: Cloudflare API rate limiting or network issues

**Solution**:
- Increase scrape interval (default: 60s)
- Check API rate limits in Cloudflare dashboard
- Verify network connectivity

### Debug Mode

Enable verbose logging:

```bash
# Add to main.go temporarily
log.SetFlags(log.LstdFlags | log.Lshortfile)
```

### Health Check

```bash
# Check exporter health
curl http://localhost:9199/health

# Check metrics endpoint
curl http://localhost:9199/metrics | grep cloudflare_zone
```

##  Best Practices

### 1. Security

- **Never commit API tokens** to version control
- Use environment variables or secrets management
- Rotate API tokens regularly
- Use minimal required permissions
- Enable HTTPS if exposing publicly

### 2. Performance

- Use appropriate scrape intervals (60s recommended)
- Monitor exporter resource usage
- Limit high-cardinality labels
- Consider caching for frequently accessed zones

### 3. Monitoring

- Set up alerts for exporter health
- Monitor API quota usage
- Track metric collection errors
- Use Grafana for visualization

### 4. Production Deployment

- Use container orchestration (Kubernetes, Docker Swarm)
- Implement proper logging aggregation
- Set resource limits
- Use health checks in load balancers
- Implement graceful shutdown

##  Contributing

Contributions are welcome! Please follow these steps:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Development Guidelines

- Follow Go best practices and idioms
- Add tests for new functionality
- Update documentation
- Maintain backward compatibility
- Use conventional commit messages

##  License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

##  Acknowledgments

- [Prometheus](https://prometheus.io/) - Monitoring system and time series database
- [Cloudflare](https://www.cloudflare.com/) - Web infrastructure and security company
- [Go Prometheus Client](https://github.com/prometheus/client_golang) - Prometheus instrumentation library

##  Support

- **Issues**: [GitHub Issues](https://github.com/moeinshahcheraghi/cloudflare-exporter/issues)
- **Discussions**: [GitHub Discussions](https://github.com/moeinshahcheraghi/cloudflare-exporter/discussions)
- **Documentation**: This README and inline code comments


**Made with  for the DevOps community**
