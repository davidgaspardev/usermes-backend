# UserMes Backend

[![CI](https://github.com/YOUR_USERNAME/usermes-backend/actions/workflows/ci.yml/badge.svg)](https://github.com/YOUR_USERNAME/usermes-backend/actions/workflows/ci.yml)
[![Coverage](https://img.shields.io/badge/coverage-92.3%25-green.svg)](https://github.com/YOUR_USERNAME/usermes-backend/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/YOUR_USERNAME/usermes-backend)](https://goreportcard.com/report/github.com/YOUR_USERNAME/usermes-backend)

RESTful API built with Go, following **Hexagonal Architecture** and **Modular Monolith** principles.

## 🚀 Quick Start

```bash
# Clone and run
git clone https://github.com/YOUR_USERNAME/usermes-backend.git
cd usermes-backend
go run cmd/main.go

# Server starts at http://localhost:3001
```

## 📚 API Documentation

Each module has complete API documentation with curl examples:

- **[User Module](./internal/modules/user/README.md)** - Registration, authentication, user management
- **[Resource Module](./internal/modules/resource/README.md)** - Resource CRUD operations

## 🏗️ Architecture

Built with **Hexagonal Architecture** + **Modular Monolith**:

```
internal/modules/
├── user/                    # User Management
│   ├── domain/             # 🎯 Business Logic (Tested)
│   ├── application/        # 🎯 Use Cases (Tested)  
│   └── infrastructure/     # 🚫 Adapters (Not Tested)
└── resource/               # Resource Management
    ├── domain/             # 🎯 Business Logic (Tested)
    ├── application/        # 🎯 Use Cases (Tested)
    └── infrastructure/     # 🚫 Adapters (Not Tested)
```

**Focus:** We only test what matters for business logic! **Current Coverage: 92.3%** ✅

## 🧪 Testing

```bash
# Run business logic tests
go test ./internal/modules/user/domain/... ./internal/modules/user/application/... \
         ./internal/modules/resource/domain/... ./internal/modules/resource/application/... -v

# With coverage
go test -coverprofile=coverage.out -covermode=atomic \
        ./internal/modules/user/domain/... ./internal/modules/user/application/... \
        ./internal/modules/resource/domain/... ./internal/modules/resource/application/...

go tool cover -func=coverage.out | grep total
```

## 🔒 Security Features

- **JWT Authentication** with 24h expiration
- **Password Hashing** with bcrypt
- **Email Validation** with business rules
- **Input Sanitization** on all endpoints

## 🎯 Example Usage

### Register & Login
```bash
# Register
curl -X POST http://localhost:3001/api/users/register \
  -H "Content-Type: application/json" \
  -d '{"email": "test@example.com", "password": "SecurePass123!", "name": "John Doe"}'

# Login
curl -X POST http://localhost:3001/api/users/login \
  -H "Content-Type: application/json" \
  -d '{"email": "test@example.com", "password": "SecurePass123!"}'
```

### Create Resource
```bash
curl -X POST http://localhost:3001/api/resources/ \
  -H "Content-Type: application/json" \
  -d '{"code": "RES001", "type": "MACHINE", "stop_factor": 5}'
```

## 🚧 Development

### Prerequisites
- Go 1.22.1+

### Configuration
The server port can be configured via environment variable:
```bash
# Run with custom port (default: 3001)
PORT=8080 go run cmd/main.go

# Or set environment variable
export PORT=8080
go run cmd/main.go
```

### Commands
```bash
# Install dependencies
go mod download

# Run tests
make test

# Build
make build

# Lint
golangci-lint run
```

## 📋 Project Status

| Module | Endpoints | Business Logic Coverage |
|--------|-----------|------------------------|
| **User** | `/api/users/*` | ✅ 91.2% |
| **Resource** | `/api/resources/*` | ✅ 94.9% |

## 🔄 Next Features

- [ ] Database integration (PostgreSQL)
- [ ] Docker setup
- [ ] API versioning
- [ ] Rate limiting
- [ ] Monitoring & metrics

## 🧩 Critérios para Criar um Módulo

### ✅ **DEVE ser módulo separado quando:**

1. **Bounded Context diferente** - Tem linguagem ubíqua própria
2. **Ciclo de vida independente** - Pode ser desenvolvido/deployado separadamente  
3. **Equipes diferentes** - Times diferentes cuidam de cada módulo
4. **Mudanças independentes** - Alterações não afetam outros módulos
5. **Possível microservice** - Pode virar serviço independente no futuro

### ❌ **NÃO deve ser módulo separado quando:**

1. **Acoplamento forte** - Sempre usado junto
2. **Entidades minúsculas** - Muito pequeno, sem lógica própria
3. **Sem negócio próprio** - É só CRUD simples

## 📄 License

MIT License