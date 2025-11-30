# UserMes Backend - Hexagonal Architecture with Modular Monolith

[![CI](https://github.com/YOUR_USERNAME/usermes-backend/actions/workflows/ci.yml/badge.svg)](https://github.com/YOUR_USERNAME/usermes-backend/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/YOUR_USERNAME/usermes-backend/branch/main/graph/badge.svg)](https://codecov.io/gh/YOUR_USERNAME/usermes-backend)
[![Go Report Card](https://goreportcard.com/badge/github.com/YOUR_USERNAME/usermes-backend)](https://goreportcard.com/report/github.com/YOUR_USERNAME/usermes-backend)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A RESTful API built in Go following **Hexagonal Architecture** (Ports and Adapters) principles with **Modular Monolith**.

## 🏗️ Architecture

This project implements hexagonal architecture (also known as Ports and Adapters), which promotes:

- **Separation of concerns**: Domain isolated from infrastructure
- **Testability**: Facilitates unit and integration testing
- **Flexibility**: Easy to swap adapters (database, frameworks, etc)
- **Scalability**: Clear path to migrate to microservices

### Directory Structure

```
usermes-backend/
├── cmd/                                    # Application entry points
│   └── main.go                            # Main application
├── internal/                              # Private application code
│   ├── modules/                           # Monolith modules
│   │   └── user/                          # User module
│   │       ├── domain/                    # Domain Layer (Core)
│   │       │   ├── entity/               # Domain entities
│   │       │   │   └── user.go
│   │       │   ├── valueobject/          # Value Objects
│   │       │   │   ├── email.go
│   │       │   │   └── password.go
│   │       │   └── errors/               # Domain errors
│   │       │       └── errors.go
│   │       ├── application/               # Application Layer (Use Cases)
│   │       │   ├── port/                 # Ports (Interfaces)
│   │       │   │   ├── input/           # Input ports (Use Cases)
│   │       │   │   │   └── user_service.go
│   │       │   │   └── output/          # Output ports (Repositories, etc)
│   │       │   │       ├── user_repository.go
│   │       │   │       └── token_generator.go
│   │       │   └── usecase/              # Use case implementations
│   │       │       └── user_service_impl.go
│   │       └── infrastructure/            # Infrastructure Layer (Adapters)
│   │           ├── adapter/
│   │           │   ├── input/
│   │           │   │   └── http/         # HTTP Adapter (Controllers)
│   │           │   │       ├── user_handler.go
│   │           │   │       ├── middleware.go
│   │           │   │       └── routes.go
│   │           │   └── output/
│   │           │       └── persistence/  # Persistence adapter
│   │           │           └── memory_user_repository.go
│   │           └── dto/                   # Data Transfer Objects
│   │               └── user_dto.go
│   └── shared/                            # Code shared between modules
│       └── infrastructure/
│           ├── http/
│           │   └── server/
│           │       └── fiber_server.go
│           └── security/
│               └── jwt_token_generator.go
└── pkg/                                   # Public reusable packages
```

## 📦 Hexagonal Architecture Layers

### 1. **Domain Layer** (Core)
The heart of the application, containing pure business logic:

- **Entities**: Objects with unique identity (`User`)
- **Value Objects**: Immutable objects without identity (`Email`, `Password`)
- **Domain Errors**: Domain-specific errors
- **Business Rules**: Domain validations and behaviors

**Characteristics:**
- Does not depend on any other layer
- Does not know about frameworks or external libraries
- Contains only pure business logic

### 2. **Application Layer** (Use Cases)
Orchestrates data flow and coordinates operations:

- **Input Ports**: Interfaces that define use cases (what the application does)
- **Output Ports**: Interfaces that define external dependencies (repositories, services)
- **Use Cases**: Implementation of use cases using domain entities

**Characteristics:**
- Depends only on the domain layer
- Defines interfaces (ports) that will be implemented by the infrastructure layer
- Contains application logic (orchestration)

### 3. **Infrastructure Layer** (Adapters)
Implements technical details and connects with the external world:

- **Input Adapters**: HTTP handlers, CLI, gRPC, etc.
- **Output Adapters**: Repository implementations, API clients, etc.
- **DTOs**: Objects for data transfer between layers

**Characteristics:**
- Implements interfaces (ports) defined in the application layer
- Contains framework and library-specific code
- Can be easily replaced without affecting the domain

## 🎯 User Module

The User module is the first module of the monolith, responsible for:

- ✅ User registration
- ✅ Authentication (Login)
- ✅ Profile management
- ✅ Password change
- ✅ Account activation/deactivation

### Available Endpoints

#### Public (no authentication required)

```http
POST /api/users/register
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "SecurePass123",
  "name": "John Doe"
}
```

```http
POST /api/users/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "SecurePass123"
}
```

#### Protected (authentication required)

```http
GET /api/users/me
Authorization: Bearer <token>
```

```http
GET /api/users/:id
Authorization: Bearer <token>
```

```http
PUT /api/users/:id
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": "Jane Doe"
}
```

```http
POST /api/users/:id/change-password
Authorization: Bearer <token>
Content-Type: application/json

{
  "old_password": "OldPass123",
  "new_password": "NewSecurePass123"
}
```

```http
POST /api/users/:id/deactivate
Authorization: Bearer <token>
```

```http
POST /api/users/:id/activate
Authorization: Bearer <token>
```

## 🚀 How to Run

### Prerequisites

- Go 1.22.1 or higher
- Make (optional)

### Installation

```bash
# Clone the repository
git clone <repository-url>
cd usermes-backend

# Download dependencies
go mod download

# Run the application
go run cmd/main.go
```

### Using Make

```bash
# Run the application
make run

# Build the application
make build

# Run tests
make test

# Run tests with coverage (90% minimum)
make coverage

# Generate HTML coverage report
make coverage-html
```

## 🧪 Testing the API

### 1. Register a new user

```bash
curl -X POST http://localhost:3000/api/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "SecurePass123",
    "name": "Test User"
  }'
```

### 2. Login

```bash
curl -X POST http://localhost:3000/api/users/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "SecurePass123"
  }'
```

Copy the returned token to use in the next requests.

### 3. Get user profile

```bash
curl -X GET http://localhost:3000/api/users/me \
  -H "Authorization: Bearer <your-token>"
```

### 4. Health Check

```bash
curl http://localhost:3000/health
```

## 🔒 Security

- **Password**: Hashed with bcrypt (cost factor 12)
- **JWT**: Tokens with 24-hour expiration
- **Validation**: Email and password validated with business rules
- **CORS**: Configured to allow requests from any origin (adjust in production)

### Password Rules

- Minimum 8 characters
- Maximum 72 characters (bcrypt limitation)
- Must contain at least one letter
- Must contain at least one number

## 🧪 Testing

This project maintains **90% minimum code coverage** enforced by CI/CD pipeline.

### Running Tests

```bash
# Run all tests
make test

# Run tests with coverage check (90% threshold)
make coverage

# Generate HTML coverage report
make coverage-html
open coverage.html

# Show coverage by function
make coverage-func
```

### Coverage Status

- **Domain Layer**: 90%+ coverage (entities, value objects)
- **Application Layer**: 84%+ coverage (use cases)
- **Infrastructure Layer**: 100% coverage (repositories)

### Writing Tests

We follow these testing principles:
- **Test behavior, not implementation**
- **Use table-driven tests** for multiple scenarios
- **Mock external dependencies** using interfaces
- **Test edge cases** and error scenarios

See [CI/CD Documentation](.github/workflows/README.md) for more details.

## 📝 Implemented Best Practices

1. **Dependency Inversion**: Upper layers don't depend on concrete implementations
2. **Single Responsibility**: Each component has a single responsibility
3. **Open/Closed**: Open for extension, closed for modification
4. **Interface Segregation**: Small and focused interfaces
5. **Domain-Driven Design**: Rich domain modeling
6. **Value Objects**: Validation and value encapsulation
7. **Repository Pattern**: Persistence abstraction
8. **Use Case Pattern**: Explicit and testable use cases

## 🔄 Next Steps

### Future Implementations

- [ ] Integration with database (PostgreSQL/MySQL)
- [ ] Redis for cache and sessions
- [ ] Refresh tokens
- [ ] Rate limiting
- [ ] Structured logging
- [ ] Metrics and observability
- [x] Unit and integration tests with 90% coverage
- [x] CI/CD pipeline with GitHub Actions
- [ ] Docker and Docker Compose
- [ ] Database migrations
- [ ] Swagger/OpenAPI documentation
- [ ] New modules (Tasks, Projects, etc)

### Adding a New Module

To add a new module to the monolith:

1. Create the directory structure in `internal/modules/[module-name]`
2. Implement the layers: domain → application → infrastructure
3. Register routes in `main.go`
4. Maintain independence between modules

Example:
```
internal/modules/task/
├── domain/
├── application/
└── infrastructure/
```

## 🤝 Contributing

1. Fork the project
2. Create a branch for your feature (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License.

## 📚 References

- [Hexagonal Architecture](https://alistair.cockburn.us/hexagonal-architecture/)
- [Domain-Driven Design](https://martinfowler.com/bliki/DomainDrivenDesign.html)
- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Modular Monolith](https://www.kamilgrzybek.com/design/modular-monolith-primer/)