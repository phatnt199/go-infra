---
sidebar_position: 1
---

# Deployment to Production

Learn how to deploy your go-infra application to production environments.

## Pre-Deployment Checklist

Before deploying to production, ensure:

- [ ] All tests pass
- [ ] Security vulnerabilities checked
- [ ] Environment variables configured
- [ ] Database migrations ready
- [ ] Monitoring and logging configured
- [ ] Backup strategy in place
- [ ] Health check endpoints working
- [ ] Performance tested

## Building for Production

### Optimize Binary Size

```bash
# Build with optimizations
go build -ldflags="-s -w" -o app ./cmd/server

# Further compress with upx (optional)
upx --best --lzma app
```

### Multi-Stage Docker Build

```dockerfile
# Dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build application
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o main ./cmd/server

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy binary from builder
COPY --from=builder /app/main .

# Copy configuration
COPY config/config.production.json ./config/

EXPOSE 8080

CMD ["./main"]
```

### Build Script

```bash
#!/bin/bash
# build.sh

set -e

VERSION=${1:-latest}
IMAGE_NAME="myapp"

echo "Building Docker image: $IMAGE_NAME:$VERSION"

docker build -t $IMAGE_NAME:$VERSION .

echo "Build completed successfully"
```

## Environment Configuration

### Production Config

```json
// config/config.production.json
{
	"environment": "production",
	"server": {
		"host": "0.0.0.0",
		"port": 8080,
		"read_timeout": 30,
		"write_timeout": 30
	},
	"database": {
		"host": "${DB_HOST}",
		"port": 5432,
		"username": "${DB_USER}",
		"password": "${DB_PASSWORD}",
		"database": "${DB_NAME}",
		"ssl_mode": "require",
		"max_open_conns": 100,
		"max_idle_conns": 10,
		"conn_max_lifetime": 3600
	},
	"jwt": {
		"secret": "${JWT_SECRET}",
		"expiration": 1
	},
	"logging": {
		"level": "info",
		"format": "json"
	}
}
```

### Environment Variables

```bash
# .env.production (never commit this!)
APP_ENV=production

# Server
SERVER_HOST=0.0.0.0
SERVER_PORT=8080

# Database
DB_HOST=your-db-host.example.com
DB_PORT=5432
DB_USER=myapp_user
DB_PASSWORD=super-secret-password
DB_NAME=myapp_production
DB_SSL_MODE=require

# JWT
JWT_SECRET=your-super-secret-jwt-key-min-32-characters

# Logging
LOG_LEVEL=info
LOG_FORMAT=json
```

## Docker Deployment

### Docker Compose Production

```yaml
# docker-compose.production.yml
version: "3.8"

services:
  app:
    image: myapp:latest
    restart: always
    ports:
      - "8080:8080"
    environment:
      - APP_ENV=production
      - DB_HOST=${DB_HOST}
      - DB_USER=${DB_USER}
      - DB_PASSWORD=${DB_PASSWORD}
      - DB_NAME=${DB_NAME}
      - JWT_SECRET=${JWT_SECRET}
    depends_on:
      - postgres
    healthcheck:
      test:
        [
          "CMD",
          "wget",
          "--quiet",
          "--tries=1",
          "--spider",
          "http://localhost:8080/health",
        ]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 40s

  postgres:
    image: postgres:15-alpine
    restart: always
    environment:
      POSTGRES_USER: ${DB_USER}
      POSTGRES_PASSWORD: ${DB_PASSWORD}
      POSTGRES_DB: ${DB_NAME}
    volumes:
      - postgres_data:/var/lib/postgresql/data
    ports:
      - "5432:5432"

  nginx:
    image: nginx:alpine
    restart: always
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf:ro
      - ./ssl:/etc/nginx/ssl:ro
    depends_on:
      - app

volumes:
  postgres_data:
    driver: local
```

### Deploy with Docker

```bash
# Load environment variables
export $(cat .env.production | xargs)

# Build and start services
docker-compose -f docker-compose.production.yml up -d

# View logs
docker-compose -f docker-compose.production.yml logs -f

# Scale application
docker-compose -f docker-compose.production.yml up -d --scale app=3
```

## Kubernetes Deployment

### Deployment Manifest

```yaml
# k8s/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: myapp
  labels:
    app: myapp
spec:
  replicas: 3
  selector:
    matchLabels:
      app: myapp
  template:
    metadata:
      labels:
        app: myapp
    spec:
      containers:
        - name: myapp
          image: myapp:latest
          ports:
            - containerPort: 8080
          env:
            - name: APP_ENV
              value: "production"
            - name: DB_HOST
              valueFrom:
                secretKeyRef:
                  name: myapp-secrets
                  key: db-host
            - name: DB_PASSWORD
              valueFrom:
                secretKeyRef:
                  name: myapp-secrets
                  key: db-password
            - name: JWT_SECRET
              valueFrom:
                secretKeyRef:
                  name: myapp-secrets
                  key: jwt-secret
          livenessProbe:
            httpGet:
              path: /health
              port: 8080
            initialDelaySeconds: 30
            periodSeconds: 10
          readinessProbe:
            httpGet:
              path: /health
              port: 8080
            initialDelaySeconds: 5
            periodSeconds: 5
          resources:
            requests:
              memory: "128Mi"
              cpu: "100m"
            limits:
              memory: "512Mi"
              cpu: "500m"
---
apiVersion: v1
kind: Service
metadata:
  name: myapp-service
spec:
  selector:
    app: myapp
  ports:
    - protocol: TCP
      port: 80
      targetPort: 8080
  type: LoadBalancer
```

### Secrets

```yaml
# k8s/secrets.yaml
apiVersion: v1
kind: Secret
metadata:
  name: myapp-secrets
type: Opaque
stringData:
  db-host: "postgres.example.com"
  db-password: "super-secret-password"
  jwt-secret: "your-super-secret-jwt-key"
```

### Deploy to Kubernetes

```bash
# Create secrets
kubectl apply -f k8s/secrets.yaml

# Deploy application
kubectl apply -f k8s/deployment.yaml

# Check status
kubectl get pods
kubectl get services

# View logs
kubectl logs -f deployment/myapp

# Scale deployment
kubectl scale deployment myapp --replicas=5
```

## Cloud Deployments

### AWS (ECS)

```json
// task-definition.json
{
	"family": "myapp",
	"networkMode": "awsvpc",
	"requiresCompatibilities": ["FARGATE"],
	"cpu": "256",
	"memory": "512",
	"containerDefinitions": [
		{
			"name": "myapp",
			"image": "123456789.dkr.ecr.us-east-1.amazonaws.com/myapp:latest",
			"portMappings": [
				{
					"containerPort": 8080,
					"protocol": "tcp"
				}
			],
			"environment": [
				{
					"name": "APP_ENV",
					"value": "production"
				}
			],
			"secrets": [
				{
					"name": "DB_PASSWORD",
					"valueFrom": "arn:aws:secretsmanager:region:account:secret:db-password"
				}
			],
			"logConfiguration": {
				"logDriver": "awslogs",
				"options": {
					"awslogs-group": "/ecs/myapp",
					"awslogs-region": "us-east-1",
					"awslogs-stream-prefix": "ecs"
				}
			}
		}
	]
}
```

### Google Cloud Run

```yaml
# cloudrun.yaml
apiVersion: serving.knative.dev/v1
kind: Service
metadata:
  name: myapp
spec:
  template:
    spec:
      containers:
        - image: gcr.io/project-id/myapp:latest
          ports:
            - containerPort: 8080
          env:
            - name: APP_ENV
              value: production
            - name: DB_PASSWORD
              valueFrom:
                secretKeyRef:
                  name: db-password
                  key: latest
          resources:
            limits:
              memory: 512Mi
              cpu: "1"
```

Deploy:

```bash
gcloud run deploy myapp \
  --image gcr.io/project-id/myapp:latest \
  --platform managed \
  --region us-central1 \
  --allow-unauthenticated
```

### Heroku

```yaml
# heroku.yml
build:
  docker:
    web: Dockerfile

run:
  web: ./main
```

```bash
# Deploy to Heroku
heroku container:push web
heroku container:release web
```

## Nginx Reverse Proxy

```nginx
# nginx.conf
events {
    worker_connections 1024;
}

http {
    upstream app {
        least_conn;
        server app:8080 max_fails=3 fail_timeout=30s;
    }

    server {
        listen 80;
        server_name example.com;

        # Redirect to HTTPS
        return 301 https://$server_name$request_uri;
    }

    server {
        listen 443 ssl http2;
        server_name example.com;

        ssl_certificate /etc/nginx/ssl/cert.pem;
        ssl_certificate_key /etc/nginx/ssl/key.pem;

        # Security headers
        add_header X-Frame-Options "SAMEORIGIN" always;
        add_header X-Content-Type-Options "nosniff" always;
        add_header X-XSS-Protection "1; mode=block" always;

        location / {
            proxy_pass http://app;
            proxy_http_version 1.1;
            proxy_set_header Upgrade $http_upgrade;
            proxy_set_header Connection 'upgrade';
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto $scheme;
            proxy_cache_bypass $http_upgrade;
        }

        location /health {
            access_log off;
            proxy_pass http://app;
        }
    }
}
```

## Monitoring

### Health Check Endpoint

```go
// Add to your application
router.Get("/health", func(c *fiber.Ctx) error {
    // Check database connection
    sqlDB, _ := db.DB()
    if err := sqlDB.Ping(); err != nil {
        return c.Status(503).JSON(fiber.Map{
            "status": "unhealthy",
            "database": "down",
        })
    }

    return c.JSON(fiber.Map{
        "status": "healthy",
        "version": "1.0.0",
        "uptime": time.Since(startTime).Seconds(),
    })
})
```

### Prometheus Metrics

```go
import (
    "github.com/gofiber/fiber/v2/middleware/monitor"
)

// Add metrics endpoint
router.Get("/metrics", monitor.New())
```

## Database Migration in Production

```bash
# Run migrations before deploying
migrate -path migrations \
  -database "postgres://user:pass@prod-db:5432/mydb?sslmode=require" \
  up

# Or use init container in Kubernetes
# See k8s/migration-job.yaml
```

## Logging

```go
import (
    "github.com/phatnt199/go-infra/pkg/logger"
    zaplogger "github.com/phatnt199/go-infra/pkg/logger/zap"
    "go.uber.org/fx"
)

// Production logger with Fx
func main() {
    app := fx.New(
        zaplogger.Module,
        fx.Invoke(runApp),
    )
    app.Run()
}

func runApp(log logger.Logger) {
    // Structured logging
    log.Infow("Server started", logger.Fields{
        "port": cfg.Server.Port,
        "environment": cfg.Environment,
    })
}
```

## Security Checklist

- [ ] Use HTTPS/TLS in production
- [ ] Store secrets in secret management service
- [ ] Enable CORS properly
- [ ] Set security headers
- [ ] Rate limit API endpoints
- [ ] Use strong JWT secrets (min 32 characters)
- [ ] Enable database SSL mode
- [ ] Regular security updates
- [ ] Principle of least privilege

## Backup Strategy

```bash
# Automated PostgreSQL backup
#!/bin/bash
# backup.sh

BACKUP_DIR="/backups"
DB_NAME="myapp_production"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)

pg_dump -h localhost -U postgres $DB_NAME > $BACKUP_DIR/$DB_NAME_$TIMESTAMP.sql

# Keep only last 7 days
find $BACKUP_DIR -name "*.sql" -mtime +7 -delete
```

## Next Steps

- Set up [Monitoring and Observability](./monitoring)
- Configure [CI/CD Pipeline](./cicd)
- Learn about [Scaling Strategies](./scaling)
