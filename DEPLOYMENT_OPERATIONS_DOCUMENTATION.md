# Deployment and Operations Documentation

## Table of Contents
1. [Overview](#overview)
2. [Prerequisites](#prerequisites)
3. [Local Development Setup](#local-development-setup)
4. [Docker Deployment](#docker-deployment)
5. [Production Deployment](#production-deployment)
6. [Configuration Management](#configuration-management)
7. [Database Management](#database-management)
8. [Monitoring and Logging](#monitoring-and-logging)
9. [Security](#security)
10. [Backup and Recovery](#backup-and-recovery)
11. [Troubleshooting](#troubleshooting)
12. [Maintenance](#maintenance)

## Overview

This document provides comprehensive guidance for deploying, operating, and maintaining the Go CRM microservices application in various environments.

### Architecture Components

- **Frontend**: Nginx serving static files and acting as reverse proxy
- **API Gateway**: Go service handling rate limiting and IP filtering
- **API Service**: Go service providing CRM business logic
- **Database**: PostgreSQL for data persistence
- **Containerization**: Docker and Docker Compose

## Prerequisites

### System Requirements

**Minimum Requirements:**
- CPU: 2 cores
- RAM: 4GB
- Storage: 20GB
- OS: Linux (Ubuntu 20.04+ recommended)

**Recommended Requirements:**
- CPU: 4+ cores
- RAM: 8GB+
- Storage: 50GB+ SSD
- OS: Ubuntu 22.04 LTS

### Software Dependencies

**Required Software:**
- Docker 20.10+
- Docker Compose 2.0+
- Make (for build automation)
- Git (for version control)

**Optional Software:**
- Docker Desktop (for local development)
- kubectl (for Kubernetes deployment)
- helm (for Kubernetes package management)

### Network Requirements

**Ports:**
- 3000: Frontend (HTTP)
- 8080: API Gateway
- 8085: API Service
- 5432: PostgreSQL

**Firewall Configuration:**
```bash
# Allow required ports
sudo ufw allow 3000/tcp
sudo ufw allow 8080/tcp
sudo ufw allow 8085/tcp
sudo ufw allow 5432/tcp
```

## Local Development Setup

### Quick Start

1. **Clone Repository:**
   ```bash
   git clone <repository-url>
   cd go-crm-microservices
   ```

2. **Start Services:**
   ```bash
   make docker-up
   ```

3. **Access Application:**
   - URL: http://localhost:3000
   - Email: `demo@example.com`
   - Password: `demo123`

4. **Stop Services:**
   ```bash
   make docker-down
   ```

### Development Workflow

1. **Code Changes:**
   ```bash
   # Make changes to Go code
   # Rebuild and restart
   make docker-up
   ```

2. **Database Changes:**
   ```bash
   # Add new migration files to migrations/
   # Restart services to apply migrations
   make docker-down
   make docker-up
   ```

3. **Frontend Changes:**
   ```bash
   # Edit files in public/
   # Changes are immediately available (no rebuild needed)
   ```

### Development Commands

```bash
# View logs
make logs

# View specific service logs
docker-compose logs api
docker-compose logs gateway
docker-compose logs postgres

# Access database
make db-connect

# Run tests
make test

# Clean up
make clean
```

## Docker Deployment

### Docker Compose Configuration

**docker-compose.yml:**
```yaml
version: '3.8'

services:
  nginx:
    image: nginx:alpine
    ports:
      - "3000:80"
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf
      - ./public:/usr/share/nginx/html
    depends_on:
      - gateway

  gateway:
    build:
      context: .
      dockerfile: Dockerfile
      target: gateway
    ports:
      - "8080:8080"
    environment:
      - PORT=8080

  api:
    build:
      context: .
      dockerfile: Dockerfile
      target: api
    ports:
      - "8085:8085"
    environment:
      - DB_HOST=postgres
      - DB_PORT=5432
      - DB_USER=user
      - DB_PASSWORD=pass
      - DB_NAME=users
      - PORT=8085
    depends_on:
      - postgres

  postgres:
    image: postgres:15-alpine
    environment:
      - POSTGRES_USER=user
      - POSTGRES_PASSWORD=pass
      - POSTGRES_DB=users
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./migrations:/docker-entrypoint-initdb.d
    ports:
      - "5432:5432"

volumes:
  postgres_data:
```

### Dockerfile

**Multi-stage Dockerfile:**
```dockerfile
# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o gateway ./cmd/gateway
RUN go build -o api ./cmd/api

# Gateway stage
FROM alpine:latest AS gateway
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/gateway .
EXPOSE 8080
CMD ["./gateway"]

# API stage
FROM alpine:latest AS api
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/api .
EXPOSE 8085
CMD ["./api"]
```

### Environment Variables

**Development (.env):**
```env
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=user
DB_PASSWORD=pass
DB_NAME=users

# Services
API_PORT=8085
GATEWAY_PORT=8080
NGINX_PORT=3000

# JWT
JWT_SECRET=supersecret

# Rate Limiting
RATE_LIMIT_REQUESTS=10
RATE_LIMIT_WINDOW=60
```

**Production (.env.prod):**
```env
# Database
DB_HOST=postgres
DB_PORT=5432
DB_USER=prod_user
DB_PASSWORD=secure_password_here
DB_NAME=crm_production

# Services
API_PORT=8085
GATEWAY_PORT=8080
NGINX_PORT=3000

# JWT
JWT_SECRET=your_very_secure_jwt_secret_here

# Rate Limiting
RATE_LIMIT_REQUESTS=100
RATE_LIMIT_WINDOW=60
```

### Build and Deploy

```bash
# Build images
docker-compose build

# Start services
docker-compose up -d

# View status
docker-compose ps

# View logs
docker-compose logs -f

# Stop services
docker-compose down

# Clean up volumes
docker-compose down -v
```

## Production Deployment

### Production Architecture

```mermaid
graph TD
    subgraph "Load Balancer"
        LB[Load Balancer / Nginx]
    end

    subgraph "Application Layer"
        Gateway1[API Gateway 1]
        Gateway2[API Gateway 2]
        API1[API Service 1]
        API2[API Service 2]
    end

    subgraph "Database Layer"
        DB1[(PostgreSQL Primary)]
        DB2[(PostgreSQL Replica)]
    end

    subgraph "Monitoring"
        Prometheus[Prometheus]
        Grafana[Grafana]
        ELK[ELK Stack]
    end

    LB --> Gateway1
    LB --> Gateway2
    Gateway1 --> API1
    Gateway1 --> API2
    Gateway2 --> API1
    Gateway2 --> API2
    API1 --> DB1
    API2 --> DB1
    DB1 --> DB2
```

### Kubernetes Deployment

**namespace.yaml:**
```yaml
apiVersion: v1
kind: Namespace
metadata:
  name: crm-system
```

**configmap.yaml:**
```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: crm-config
  namespace: crm-system
data:
  DB_HOST: "postgres-service"
  DB_PORT: "5432"
  DB_NAME: "crm_production"
  API_PORT: "8085"
  GATEWAY_PORT: "8080"
```

**secret.yaml:**
```yaml
apiVersion: v1
kind: Secret
metadata:
  name: crm-secrets
  namespace: crm-system
type: Opaque
data:
  DB_USER: cHJvZF91c2Vy
  DB_PASSWORD: c2VjdXJlX3Bhc3N3b3Jk
  JWT_SECRET: eW91cl92ZXJ5X3NlY3VyZV9qd3Rfc2VjcmV0
```

**deployment.yaml:**
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: api-service
  namespace: crm-system
spec:
  replicas: 3
  selector:
    matchLabels:
      app: api-service
  template:
    metadata:
      labels:
        app: api-service
    spec:
      containers:
      - name: api
        image: crm-api:latest
        ports:
        - containerPort: 8085
        env:
        - name: DB_HOST
          valueFrom:
            configMapKeyRef:
              name: crm-config
              key: DB_HOST
        - name: DB_USER
          valueFrom:
            secretKeyRef:
              name: crm-secrets
              key: DB_USER
        - name: DB_PASSWORD
          valueFrom:
            secretKeyRef:
              name: crm-secrets
              key: DB_PASSWORD
        - name: JWT_SECRET
          valueFrom:
            secretKeyRef:
              name: crm-secrets
              key: JWT_SECRET
        resources:
          requests:
            memory: "256Mi"
            cpu: "250m"
          limits:
            memory: "512Mi"
            cpu: "500m"
        livenessProbe:
          httpGet:
            path: /healthz
            port: 8085
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /healthz
            port: 8085
          initialDelaySeconds: 5
          periodSeconds: 5
```

**service.yaml:**
```yaml
apiVersion: v1
kind: Service
metadata:
  name: api-service
  namespace: crm-system
spec:
  selector:
    app: api-service
  ports:
  - protocol: TCP
    port: 8085
    targetPort: 8085
  type: ClusterIP
```

**ingress.yaml:**
```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: crm-ingress
  namespace: crm-system
  annotations:
    nginx.ingress.kubernetes.io/rewrite-target: /
    cert-manager.io/cluster-issuer: letsencrypt-prod
spec:
  tls:
  - hosts:
    - crm.example.com
    secretName: crm-tls
  rules:
  - host: crm.example.com
    http:
      paths:
      - path: /api
        pathType: Prefix
        backend:
          service:
            name: api-service
            port:
              number: 8085
      - path: /
        pathType: Prefix
        backend:
          service:
            name: frontend-service
            port:
              number: 80
```

### Deployment Commands

```bash
# Create namespace
kubectl apply -f namespace.yaml

# Apply configurations
kubectl apply -f configmap.yaml
kubectl apply -f secret.yaml

# Deploy services
kubectl apply -f deployment.yaml
kubectl apply -f service.yaml
kubectl apply -f ingress.yaml

# Check deployment status
kubectl get pods -n crm-system
kubectl get services -n crm-system
kubectl get ingress -n crm-system

# View logs
kubectl logs -f deployment/api-service -n crm-system

# Scale deployment
kubectl scale deployment api-service --replicas=5 -n crm-system
```

## Configuration Management

### Environment-Specific Configurations

**Development:**
```bash
# Use development configuration
export ENV=development
docker-compose -f docker-compose.yml -f docker-compose.dev.yml up
```

**Staging:**
```bash
# Use staging configuration
export ENV=staging
docker-compose -f docker-compose.yml -f docker-compose.staging.yml up
```

**Production:**
```bash
# Use production configuration
export ENV=production
docker-compose -f docker-compose.yml -f docker-compose.prod.yml up
```

### Configuration Validation

```bash
# Validate Docker Compose configuration
docker-compose config

# Validate Kubernetes manifests
kubectl apply --dry-run=client -f deployment.yaml

# Validate environment variables
./scripts/validate-env.sh
```

### Configuration Updates

```bash
# Update configuration without downtime
docker-compose up -d --no-deps --build api

# Rolling update in Kubernetes
kubectl set image deployment/api-service api=crm-api:new-version -n crm-system

# Update secrets
kubectl create secret generic crm-secrets --from-literal=JWT_SECRET=new-secret --dry-run=client -o yaml | kubectl apply -f -
```

## Database Management

### Database Setup

**Initial Setup:**
```bash
# Create database
createdb crm_production

# Run migrations
migrate -path migrations -database "postgres://user:pass@localhost:5432/crm_production?sslmode=disable" up

# Seed initial data
go run scripts/seed.go
```

**Migration Management:**
```bash
# Create new migration
migrate create -ext sql -dir migrations -seq add_new_table

# Apply migrations
migrate -path migrations -database "postgres://user:pass@localhost:5432/crm_production?sslmode=disable" up

# Rollback migrations
migrate -path migrations -database "postgres://user:pass@localhost:5432/crm_production?sslmode=disable" down 1

# Check migration status
migrate -path migrations -database "postgres://user:pass@localhost:5432/crm_production?sslmode=disable" version
```

### Database Backup

**Automated Backup Script:**
```bash
#!/bin/bash
# backup.sh

DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_DIR="/backups"
DB_NAME="crm_production"
DB_USER="user"
DB_HOST="localhost"

# Create backup
pg_dump -h $DB_HOST -U $DB_USER -d $DB_NAME > $BACKUP_DIR/backup_$DATE.sql

# Compress backup
gzip $BACKUP_DIR/backup_$DATE.sql

# Keep only last 7 days of backups
find $BACKUP_DIR -name "backup_*.sql.gz" -mtime +7 -delete

echo "Backup completed: backup_$DATE.sql.gz"
```

**Cron Job Setup:**
```bash
# Add to crontab
0 2 * * * /path/to/backup.sh >> /var/log/backup.log 2>&1
```

### Database Monitoring

**Connection Pool Monitoring:**
```sql
-- Check active connections
SELECT count(*) as active_connections 
FROM pg_stat_activity 
WHERE state = 'active';

-- Check connection pool usage
SELECT 
    max_conn,
    used,
    res_for_super,
    max_conn - used - res_for_super AS remaining
FROM 
    (SELECT setting::int AS max_conn FROM pg_settings WHERE name = 'max_connections') AS max_connections,
    (SELECT count(*) AS used FROM pg_stat_activity) AS active_connections,
    (SELECT setting::int AS res_for_super FROM pg_settings WHERE name = 'superuser_reserved_connections') AS reserved_connections;
```

**Performance Monitoring:**
```sql
-- Slow queries
SELECT 
    query,
    calls,
    total_time,
    mean_time,
    rows
FROM pg_stat_statements 
ORDER BY mean_time DESC 
LIMIT 10;

-- Table sizes
SELECT 
    schemaname,
    tablename,
    pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) as size
FROM pg_tables 
WHERE schemaname = 'public'
ORDER BY pg_total_relation_size(schemaname||'.'||tablename) DESC;
```

## Monitoring and Logging

### Application Monitoring

**Health Check Endpoints:**
```bash
# API Service Health
curl http://localhost:8085/healthz

# Gateway Health
curl http://localhost:8080/healthz

# Database Health
pg_isready -h localhost -p 5432
```

**Prometheus Metrics:**
```go
// Add to API service
import (
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
    httpRequestsTotal = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Total number of HTTP requests",
        },
        []string{"method", "endpoint", "status"},
    )
    
    httpRequestDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name: "http_request_duration_seconds",
            Help: "HTTP request duration in seconds",
        },
        []string{"method", "endpoint"},
    )
)

func init() {
    prometheus.MustRegister(httpRequestsTotal)
    prometheus.MustRegister(httpRequestDuration)
}

// Add metrics endpoint
http.Handle("/metrics", promhttp.Handler())
```

### Logging Configuration

**Structured Logging:**
```go
import (
    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
)

func setupLogger() *zap.Logger {
    config := zap.NewProductionConfig()
    config.EncoderConfig.TimeKey = "timestamp"
    config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
    
    logger, err := config.Build()
    if err != nil {
        log.Fatal("Failed to initialize logger:", err)
    }
    
    return logger
}

// Usage
logger.Info("Request processed",
    zap.String("method", r.Method),
    zap.String("path", r.URL.Path),
    zap.Int("status", statusCode),
    zap.Duration("duration", duration),
)
```

**Log Aggregation (ELK Stack):**
```yaml
# logstash.conf
input {
  beats {
    port => 5044
  }
}

filter {
  if [fields][service] == "api" {
    grok {
      match => { "message" => "%{TIMESTAMP_ISO8601:timestamp} %{LOGLEVEL:level} %{GREEDYDATA:message}" }
    }
  }
}

output {
  elasticsearch {
    hosts => ["elasticsearch:9200"]
    index => "crm-logs-%{+YYYY.MM.dd}"
  }
}
```

### Alerting

**Prometheus Alert Rules:**
```yaml
# alerts.yml
groups:
- name: crm-alerts
  rules:
  - alert: HighErrorRate
    expr: rate(http_requests_total{status=~"5.."}[5m]) > 0.1
    for: 2m
    labels:
      severity: warning
    annotations:
      summary: "High error rate detected"
      description: "Error rate is {{ $value }} errors per second"

  - alert: ServiceDown
    expr: up == 0
    for: 1m
    labels:
      severity: critical
    annotations:
      summary: "Service is down"
      description: "Service {{ $labels.instance }} is down"

  - alert: HighResponseTime
    expr: histogram_quantile(0.95, http_request_duration_seconds) > 1
    for: 5m
    labels:
      severity: warning
    annotations:
      summary: "High response time"
      description: "95th percentile response time is {{ $value }}s"
```

**Alertmanager Configuration:**
```yaml
# alertmanager.yml
global:
  smtp_smarthost: 'smtp.example.com:587'
  smtp_from: 'alerts@example.com'
  smtp_auth_username: 'alerts@example.com'
  smtp_auth_password: 'password'

route:
  group_by: ['alertname']
  group_wait: 10s
  group_interval: 10s
  repeat_interval: 1h
  receiver: 'team-crm'

receivers:
- name: 'team-crm'
  email_configs:
  - to: 'team@example.com'
    send_resolved: true
```

## Security

### Security Hardening

**Docker Security:**
```dockerfile
# Use non-root user
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup
USER appuser

# Remove unnecessary packages
RUN apk del --purge curl wget

# Scan for vulnerabilities
RUN trivy filesystem --exit-code 1 --severity HIGH,CRITICAL /
```

**Network Security:**
```bash
# Configure firewall
sudo ufw default deny incoming
sudo ufw default allow outgoing
sudo ufw allow ssh
sudo ufw allow 3000/tcp
sudo ufw enable

# Use VPN for database access
sudo ufw allow from 10.0.0.0/8 to any port 5432
```

**SSL/TLS Configuration:**
```nginx
# nginx.conf
server {
    listen 443 ssl http2;
    server_name crm.example.com;
    
    ssl_certificate /etc/ssl/certs/crm.crt;
    ssl_certificate_key /etc/ssl/private/crm.key;
    
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers ECDHE-RSA-AES256-GCM-SHA512:DHE-RSA-AES256-GCM-SHA512;
    ssl_prefer_server_ciphers off;
    
    # Security headers
    add_header Strict-Transport-Security "max-age=31536000; includeSubDomains" always;
    add_header X-Frame-Options DENY always;
    add_header X-Content-Type-Options nosniff always;
    add_header X-XSS-Protection "1; mode=block" always;
}
```

### Access Control

**API Rate Limiting:**
```go
// Enhanced rate limiting
type RateLimiter struct {
    requests map[string][]time.Time
    mu       sync.RWMutex
    limit    int
    window   time.Duration
}

func (rl *RateLimiter) Allow(ip string) bool {
    rl.mu.Lock()
    defer rl.mu.Unlock()
    
    now := time.Now()
    windowStart := now.Add(-rl.window)
    
    // Clean old requests
    if requests, exists := rl.requests[ip]; exists {
        var valid []time.Time
        for _, req := range requests {
            if req.After(windowStart) {
                valid = append(valid, req)
            }
        }
        rl.requests[ip] = valid
    }
    
    // Check limit
    if len(rl.requests[ip]) >= rl.limit {
        return false
    }
    
    // Add current request
    rl.requests[ip] = append(rl.requests[ip], now)
    return true
}
```

**JWT Security:**
```go
// Secure JWT configuration
var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func createToken(email string) (string, error) {
    claims := jwt.MapClaims{
        "email": email,
        "exp":   time.Now().Add(24 * time.Hour).Unix(),
        "iat":   time.Now().Unix(),
        "iss":   "crm-system",
        "aud":   "crm-users",
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(jwtSecret)
}

func validateToken(tokenString string) (jwt.MapClaims, error) {
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return jwtSecret, nil
    })
    
    if err != nil {
        return nil, err
    }
    
    if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
        return claims, nil
    }
    
    return nil, fmt.Errorf("invalid token")
}
```

## Backup and Recovery

### Backup Strategy

**Full Backup Schedule:**
- Daily: Full database backup
- Weekly: Full system backup
- Monthly: Archive backup

**Backup Scripts:**
```bash
#!/bin/bash
# full-backup.sh

DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_DIR="/backups"
S3_BUCKET="crm-backups"

# Database backup
pg_dump -h localhost -U user -d crm_production | gzip > $BACKUP_DIR/db_backup_$DATE.sql.gz

# Configuration backup
tar -czf $BACKUP_DIR/config_backup_$DATE.tar.gz /etc/crm/

# Application backup
tar -czf $BACKUP_DIR/app_backup_$DATE.tar.gz /opt/crm/

# Upload to S3
aws s3 cp $BACKUP_DIR/db_backup_$DATE.sql.gz s3://$S3_BUCKET/database/
aws s3 cp $BACKUP_DIR/config_backup_$DATE.tar.gz s3://$S3_BUCKET/config/
aws s3 cp $BACKUP_DIR/app_backup_$DATE.tar.gz s3://$S3_BUCKET/application/

# Cleanup old backups
find $BACKUP_DIR -name "*.gz" -mtime +30 -delete
```

### Recovery Procedures

**Database Recovery:**
```bash
#!/bin/bash
# db-recovery.sh

BACKUP_FILE=$1
DB_NAME="crm_production"

if [ -z "$BACKUP_FILE" ]; then
    echo "Usage: $0 <backup_file>"
    exit 1
fi

# Stop application
docker-compose stop api

# Drop and recreate database
dropdb $DB_NAME
createdb $DB_NAME

# Restore from backup
gunzip -c $BACKUP_FILE | psql -d $DB_NAME

# Start application
docker-compose start api

echo "Database recovery completed"
```

**Full System Recovery:**
```bash
#!/bin/bash
# system-recovery.sh

BACKUP_DATE=$1
S3_BUCKET="crm-backups"

if [ -z "$BACKUP_DATE" ]; then
    echo "Usage: $0 <backup_date>"
    exit 1
fi

# Download backups from S3
aws s3 cp s3://$S3_BUCKET/database/db_backup_$BACKUP_DATE.sql.gz /tmp/
aws s3 cp s3://$S3_BUCKET/config/config_backup_$BACKUP_DATE.tar.gz /tmp/
aws s3 cp s3://$S3_BUCKET/application/app_backup_$BACKUP_DATE.tar.gz /tmp/

# Stop all services
docker-compose down

# Restore configuration
tar -xzf /tmp/config_backup_$BACKUP_DATE.tar.gz -C /

# Restore application
tar -xzf /tmp/app_backup_$BACKUP_DATE.tar.gz -C /

# Restore database
gunzip -c /tmp/db_backup_$BACKUP_DATE.sql.gz | psql -d crm_production

# Start services
docker-compose up -d

echo "System recovery completed"
```

### Disaster Recovery

**RTO (Recovery Time Objective):** 4 hours
**RPO (Recovery Point Objective):** 1 hour

**Recovery Checklist:**
1. Verify backup integrity
2. Restore database
3. Restore configuration
4. Restore application files
5. Verify service connectivity
6. Run health checks
7. Verify data integrity
8. Update DNS if necessary
9. Notify stakeholders

## Troubleshooting

### Common Issues

**Service Won't Start:**
```bash
# Check logs
docker-compose logs api

# Check port conflicts
netstat -tulpn | grep :8085

# Check disk space
df -h

# Check memory
free -h
```

**Database Connection Issues:**
```bash
# Test database connectivity
pg_isready -h localhost -p 5432

# Check database logs
docker-compose logs postgres

# Check connection pool
psql -d crm_production -c "SELECT count(*) FROM pg_stat_activity;"
```

**High Memory Usage:**
```bash
# Check memory usage
docker stats

# Check for memory leaks
docker exec -it <container_id> top

# Restart service
docker-compose restart api
```

**Slow Performance:**
```bash
# Check CPU usage
htop

# Check disk I/O
iotop

# Check network
iftop

# Check database performance
psql -d crm_production -c "SELECT * FROM pg_stat_statements ORDER BY mean_time DESC LIMIT 10;"
```

### Debugging Tools

**Application Debugging:**
```go
// Add debug logging
log.SetLevel(log.DebugLevel)

// Add profiling
import _ "net/http/pprof"

func main() {
    go func() {
        log.Println(http.ListenAndServe("localhost:6060", nil))
    }()
}
```

**Network Debugging:**
```bash
# Check network connectivity
telnet localhost 8085

# Check DNS resolution
nslookup crm.example.com

# Check SSL certificate
openssl s_client -connect crm.example.com:443 -servername crm.example.com
```

**Database Debugging:**
```sql
-- Enable query logging
ALTER SYSTEM SET log_statement = 'all';
SELECT pg_reload_conf();

-- Check slow queries
SELECT query, calls, total_time, mean_time 
FROM pg_stat_statements 
ORDER BY mean_time DESC 
LIMIT 10;

-- Check locks
SELECT * FROM pg_locks WHERE NOT granted;
```

## Maintenance

### Regular Maintenance Tasks

**Daily:**
- Check service health
- Review error logs
- Monitor resource usage
- Verify backup completion

**Weekly:**
- Update security patches
- Review performance metrics
- Clean up old logs
- Test backup restoration

**Monthly:**
- Update dependencies
- Review access logs
- Performance optimization
- Security audit

**Quarterly:**
- Full system backup
- Disaster recovery test
- Capacity planning
- Documentation review

### Maintenance Scripts

**Health Check Script:**
```bash
#!/bin/bash
# health-check.sh

SERVICES=("api" "gateway" "postgres")
FAILED_SERVICES=()

for service in "${SERVICES[@]}"; do
    if ! docker-compose ps $service | grep -q "Up"; then
        FAILED_SERVICES+=($service)
    fi
done

if [ ${#FAILED_SERVICES[@]} -ne 0 ]; then
    echo "Failed services: ${FAILED_SERVICES[@]}"
    exit 1
else
    echo "All services healthy"
    exit 0
fi
```

**Log Rotation:**
```bash
# /etc/logrotate.d/crm
/var/log/crm/*.log {
    daily
    missingok
    rotate 30
    compress
    delaycompress
    notifempty
    create 644 root root
    postrotate
        docker-compose restart api gateway
    endscript
}
```

**Performance Monitoring:**
```bash
#!/bin/bash
# performance-check.sh

# Check CPU usage
CPU_USAGE=$(top -bn1 | grep "Cpu(s)" | awk '{print $2}' | cut -d'%' -f1)

# Check memory usage
MEMORY_USAGE=$(free | grep Mem | awk '{printf("%.2f", $3/$2 * 100.0)}')

# Check disk usage
DISK_USAGE=$(df / | tail -1 | awk '{print $5}' | cut -d'%' -f1)

echo "CPU Usage: ${CPU_USAGE}%"
echo "Memory Usage: ${MEMORY_USAGE}%"
echo "Disk Usage: ${DISK_USAGE}%"

# Alert if thresholds exceeded
if (( $(echo "$CPU_USAGE > 80" | bc -l) )); then
    echo "High CPU usage detected"
fi

if (( $(echo "$MEMORY_USAGE > 80" | bc -l) )); then
    echo "High memory usage detected"
fi

if [ "$DISK_USAGE" -gt 80 ]; then
    echo "High disk usage detected"
fi
```

### Update Procedures

**Application Updates:**
```bash
#!/bin/bash
# update-application.sh

VERSION=$1

if [ -z "$VERSION" ]; then
    echo "Usage: $0 <version>"
    exit 1
fi

# Backup current version
docker-compose exec api cp /app/api /app/api.backup

# Pull new image
docker-compose pull api

# Update configuration if needed
# ...

# Restart service
docker-compose up -d --no-deps api

# Verify update
sleep 10
if curl -f http://localhost:8085/healthz; then
    echo "Update successful"
else
    echo "Update failed, rolling back"
    docker-compose exec api cp /app/api.backup /app/api
    docker-compose restart api
fi
```

**Database Updates:**
```bash
#!/bin/bash
# update-database.sh

# Create backup before update
pg_dump -h localhost -U user -d crm_production > backup_$(date +%Y%m%d_%H%M%S).sql

# Run migrations
migrate -path migrations -database "postgres://user:pass@localhost:5432/crm_production?sslmode=disable" up

# Verify migration
migrate -path migrations -database "postgres://user:pass@localhost:5432/crm_production?sslmode=disable" version

echo "Database update completed"
```

This comprehensive documentation provides all the information needed to deploy, operate, and maintain the Go CRM microservices application in various environments, from local development to production deployment.