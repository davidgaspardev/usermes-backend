# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- **Modular Architecture**: Reorganized codebase into bounded contexts
  - `iam/user` - Identity and Access Management
  - `organization/plant` - Manufacturing plant management
  - `production/resource` - Production resource management
- **Plant Module**: Complete CRUD implementation for managing manufacturing facilities
  - Plant entity with code, name, location (latitude/longitude)
  - Plant ownership tied to user accounts
  - Activation/deactivation lifecycle
  - 100% domain test coverage, 91.5% application coverage
- **Multi-tenancy Support**: Resources are now scoped to plants
  - Resources reference plants via `plantCode` (weak reference pattern)
  - Data isolation at application level
  - Hierarchical API routes: `/v1/plants/:plant_code/production/resources`
- **Enhanced Documentation**:
  - Mermaid diagrams for architecture visualization
  - Detailed module descriptions with responsibilities and boundaries
  - System flow diagrams
  - Golden rules for architecture patterns

### Changed
- **Resource Module**: Updated to support multi-tenancy
  - Added `plantCode` field to Resource entity
  - Updated all use cases to include plant context
  - Resources are now plant-scoped instead of global
- **API Routes**: Restructured to reflect domain hierarchy
  - `/v1/users/*` - IAM endpoints
  - `/v1/plants/*` - Organization endpoints
  - `/v1/plants/:plant_code/production/*` - Production endpoints

### Documentation
- Comprehensive README with architecture diagrams
- Module descriptions with clear boundaries
- API documentation with examples
- Quick start guide

## [1.0.0] - 2024-01-15

### Added
- **User Module**: Complete user authentication and management system
  - User registration with email validation (RFC 5322 compliant)
  - Login with JWT token generation (24h expiration)
  - Password hashing using bcrypt (cost factor 10)
  - User profile management (update name, email)
  - Password change functionality
  - User activation/deactivation
  - Protected routes with JWT middleware
  - Domain-driven design with value objects (Email, Password)
  - In-memory repository implementation
  - 95.2% domain coverage, 89.5% application coverage

- **Resource Module**: Complete CRUD implementation for resource management
  - Resource entity with code, type, shiftId, stopFactor, tags
  - Domain layer with business rules and validation
  - RESTful API with 8 endpoints
  - Filtering by type and shift ID
  - Pagination support (limit/offset)
  - Thread-safe in-memory repository
  - 98.3% domain coverage, 92.7% application coverage

### Infrastructure
- **Hexagonal Architecture** (Ports & Adapters)
  - Clear separation between domain, application, and infrastructure
  - Dependency inversion principle enforced
  - High testability with isolated business logic
- **HTTP Server**: Fiber-based REST API
  - JWT authentication middleware
  - Error handling middleware
  - CORS support (configurable)
  - Request logging with emojis for visual clarity
  - Panic recovery
  - Health check endpoint
- **Security**:
  - JWT with configurable expiration
  - bcrypt password hashing
  - Email validation
  - Input sanitization
  - UUID-based IDs to prevent enumeration attacks

### Testing
- Unit tests for domain entities and value objects
- Integration tests for use cases
- Test coverage reporting
- CI/CD pipeline with GitHub Actions
  - Automated testing
  - Linting with golangci-lint
  - Code formatting checks
  - Build verification

### Development Tools
- Go 1.22.1+ support
- Makefile for common tasks (`make build`, `make test`, `make fmt`, `make vet`)
- golangci-lint configuration
- Git workflows configuration

## Architecture Principles

### Hexagonal Architecture
- **Domain Layer**: Pure business logic, framework-independent
- **Application Layer**: Use case orchestration, depends on domain
- **Infrastructure Layer**: Technical details (HTTP, DB), depends on application interfaces

### Domain-Driven Design
- **Bounded Contexts**: Modules represent business domains (IAM, Organization, Production)
- **Entities**: Objects with identity and lifecycle
- **Value Objects**: Immutable objects defined by their attributes
- **Weak References**: Modules communicate via IDs/codes, never direct object references

### Testing Strategy
- **Test Business Logic Only**: 90%+ coverage on domain and application layers
- **No Infrastructure Tests**: HTTP handlers and repositories not tested
- **Fast Tests**: No external dependencies, all in-memory

## Security Notes

- JWT tokens expire after 24 hours (configurable via `TOKEN_DURATION`)
- Passwords hashed with bcrypt, cost factor 10
- Input validation on all endpoints
- UUID-based IDs prevent enumeration attacks
- Multi-tenant data isolation enforced at application level

## Performance

- In-memory storage for development (fast, no external dependencies)
- Thread-safe concurrent access with mutex locks
- Optimized struct memory alignment
- Efficient pagination implementation

## Known Limitations

- In-memory storage means data loss on restart (by design for development)
- No persistent database integration yet (planned)
- CORS configured for all origins (development only - must configure for production)
- No rate limiting (planned)
- No real-time updates yet (planned)

## Future Roadmap

### Planned Features
- PostgreSQL integration for persistent storage
- Shift management module
- Production orders and scheduling
- Stop/downtime tracking with analytics
- Real-time resource status updates
- WebSocket support for real-time data
- Quality management module
- Maintenance tracking
- OEE (Overall Equipment Effectiveness) calculations
- MTBF/MTTR analytics

---

For more information, see the [README.md](README.md) for complete documentation.