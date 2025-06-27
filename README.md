# 🚀 Afrikpay Gateway

> **Microservices-based crypto gateway enabling cryptocurrency purchases and Mobile Money wallet deposits**

[![Go Version](https://img.shields.io/badge/Go-1.21+-blue.svg)](https://golang.org)
[![Docker](https://img.shields.io/badge/Docker-Compose-blue.svg)](https://docker.com)
[![Temporal](https://img.shields.io/badge/Temporal-Workflow-green.svg)](https://temporal.io)
[![TDD](https://img.shields.io/badge/Development-TDD-red.svg)](https://en.wikipedia.org/wiki/Test-driven_development)

## 🎯 Overview

Afrikpay Gateway is a robust microservices architecture that allows users to:

- **Buy cryptocurrencies** (USDT/BTC) via exchange APIs (Binance, Bitget)
- **Deposit to Afrikpay wallet** via Mobile Money (MTN, Orange)
- **Secure transactions** with JWT authentication and Temporal workflows
- **Scalable architecture** with independent microservices

## 🏗️ Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Auth Service  │    │   CRUD Service  │    │ Temporal Service│
│    (JWT)        │    │   (MongoDB)     │    │  (PostgreSQL)   │
│   Port: 8001    │    │   Port: 8002    │    │   Port: 8003    │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         └───────────────────────┼───────────────────────┘
                                 │
                    ┌─────────────────┐
                    │ Client Service  │
                    │ (External APIs) │
                    │   (Library)     │
                    └─────────────────┘
```

### 🔧 Services

| Service | Port | Purpose | Technology |
|---------|------|---------|------------|
| **Auth** | 8001 | JWT Authentication | Go + JWT + RSA |
| **CRUD** | 8002 | Data Operations | Go + MongoDB |
| **Temporal** | 8003 | Workflow Orchestration | Go + Temporal + Saga |
| **Client** | - | External API Connections | Go + Circuit Breaker |

## 🚀 Quick Start

### Prerequisites

- **Go 1.21+**
- **Docker & Docker Compose**
- **Make** (optional but recommended)

### 1. Clone & Setup

```bash
git clone <repository-url>
cd afrikpay-gateway
cp .env.example .env
```

### 2. Start Development Environment

```bash
# Using Make (recommended)
make dev

# Or using Docker Compose directly
docker-compose up --build
```

### 3. Verify Services

```bash
# Check all services health
make health

# Or manually
curl http://localhost:8001/health  # Auth Service
curl http://localhost:8002/health  # CRUD Service
curl http://localhost:8003/health  # Temporal Service
```

## 🛠️ Development

### Available Commands

```bash
# Development
make dev              # Start all services
make dev-detached     # Start in background
make stop             # Stop all services
make restart          # Restart services
make logs             # View all logs

# Building
make build            # Build all services
make build-auth       # Build auth service only
make build-crud       # Build crud service only

# Testing
make test             # Run all tests
make test-auth        # Test auth service only
make test-coverage    # Run with coverage

# Utilities
make clean            # Clean build artifacts
make deps             # Manage dependencies
make fmt              # Format code
make lint             # Run linter
```

### 🧪 Test-Driven Development (TDD)

This project follows **strict TDD approach**:

1. **RED** → Write failing test first
2. **GREEN** → Write minimal code to pass
3. **REFACTOR** → Improve code quality

**Coverage Requirements:**
- Business Logic: **90%+**
- Handlers: **85%+**
- Overall: **85%+**

### 📁 Project Structure

```
afrikpay-gateway/
├── services/           # Microservices
│   ├── auth/          # JWT Authentication
│   ├── crud/          # Data Operations
│   ├── temporal/      # Workflow Engine
│   └── client/        # External APIs
├── shared/            # Shared libraries
├── config/            # Configuration files
├── tests/             # Integration tests
├── docs/              # Documentation
├── scripts/           # Utility scripts
└── deployments/       # Deployment configs
```

## 🔐 Security

- **JWT Authentication** with RSA key pairs
- **Environment-based configuration**
- **CORS protection**
- **Rate limiting**
- **Input validation**

## 📊 Monitoring

Access monitoring dashboards:

- **Temporal UI**: http://localhost:8080
- **Prometheus**: http://localhost:9090
- **Grafana**: http://localhost:3000 (admin/admin123)

## 🔌 API Endpoints

### Auth Service (8001)
- `POST /auth/login` - User authentication
- `GET /auth/verify` - Token verification
- `POST /auth/refresh` - Token refresh

### CRUD Service (8002)
- `GET /users` - List users
- `POST /users` - Create user
- `GET /wallets` - List wallets
- `POST /transactions` - Create transaction

### Temporal Service (8003)
- `POST /workflows/crypto-purchase` - Start crypto purchase
- `POST /workflows/wallet-deposit` - Start wallet deposit
- `GET /workflows/{id}` - Get workflow status

## 🧪 Testing

```bash
# Run all tests
make test

# Run specific service tests
make test-auth
make test-crud
make test-temporal

# Run with coverage
make test-coverage
```

## 📚 Documentation

- [API Documentation](./docs/api/) - OpenAPI/Swagger specs
- [Architecture Guide](./docs/architecture/) - System design
- [Development Guide](./docs/guides/) - Development workflow
- [Deployment Guide](./deployments/) - Production deployment

## 🤝 Contributing

1. **Fork** the repository
2. **Create** feature branch (`git checkout -b feature/amazing-feature`)
3. **Follow TDD** approach (RED → GREEN → REFACTOR)
4. **Write tests** first, then implementation
5. **Ensure** 85%+ test coverage
6. **Commit** changes (`git commit -m 'Add amazing feature'`)
7. **Push** to branch (`git push origin feature/amazing-feature`)
8. **Open** Pull Request

### Code Standards

- **Go Effective Go** guidelines
- **TDD mandatory** for all features
- **90%+ coverage** for business logic
- **Comprehensive error handling**
- **Interface segregation**

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🆘 Support

- **Issues**: [GitHub Issues](https://github.com/afrikpay/gateway/issues)
- **Discussions**: [GitHub Discussions](https://github.com/afrikpay/gateway/discussions)
- **Documentation**: [Project Wiki](https://github.com/afrikpay/gateway/wiki)

---

**Built with ❤️ using Go, Temporal, and Test-Driven Development**
