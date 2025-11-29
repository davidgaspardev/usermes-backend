# Arquitetura do Sistema - UserMes Backend

## 📐 Visão Geral da Arquitetura

Este documento descreve a arquitetura do **UserMes Backend**, construído seguindo os princípios de **Arquitetura Hexagonal** (Ports and Adapters) dentro de um **Monolito Modular**.

## 🎯 Princípios Arquiteturais

### 1. Arquitetura Hexagonal (Ports and Adapters)

A arquitetura hexagonal separa a lógica de negócio (domínio) dos detalhes de implementação (infraestrutura), permitindo que a aplicação seja agnóstica em relação a frameworks, bancos de dados e outros serviços externos.

```
┌─────────────────────────────────────────────────────────────┐
│                    INFRASTRUCTURE LAYER                      │
│  ┌─────────────────────────────────────────────────────┐    │
│  │              INPUT ADAPTERS (Drivers)                │    │
│  │  ┌──────────┐  ┌──────────┐  ┌──────────┐          │    │
│  │  │   HTTP   │  │   CLI    │  │   gRPC   │  ...     │    │
│  │  │ Handlers │  │ Commands │  │ Handlers │          │    │
│  │  └─────┬────┘  └─────┬────┘  └─────┬────┘          │    │
│  └────────┼─────────────┼─────────────┼────────────────┘    │
│           │             │             │                      │
│           └─────────────┼─────────────┘                      │
│                         ▼                                    │
│  ┌─────────────────────────────────────────────────────┐    │
│  │              APPLICATION LAYER                       │    │
│  │  ┌───────────────────────────────────────────────┐  │    │
│  │  │            INPUT PORTS (Interfaces)           │  │    │
│  │  │          (Define Use Cases)                   │  │    │
│  │  └───────────────────┬───────────────────────────┘  │    │
│  │                      │                              │    │
│  │  ┌───────────────────▼───────────────────────────┐  │    │
│  │  │            USE CASES                          │  │    │
│  │  │     (Business Logic Orchestration)           │  │    │
│  │  └───────────────────┬───────────────────────────┘  │    │
│  │                      │                              │    │
│  │  ┌───────────────────▼───────────────────────────┐  │    │
│  │  │           OUTPUT PORTS (Interfaces)           │  │    │
│  │  │      (Define External Dependencies)           │  │    │
│  │  └───────────────────────────────────────────────┘  │    │
│  └────────────────────┬─────────────────────────────────┘    │
│                       │                                      │
│  ┌────────────────────▼─────────────────────────────────┐    │
│  │               DOMAIN LAYER                           │    │
│  │  ┌──────────┐  ┌──────────┐  ┌──────────┐          │    │
│  │  │ Entities │  │  Value   │  │  Domain  │          │    │
│  │  │          │  │ Objects  │  │  Errors  │          │    │
│  │  └──────────┘  └──────────┘  └──────────┘          │    │
│  │                                                      │    │
│  └──────────────────────────────────────────────────────┘    │
│                       ▲                                      │
│  ┌────────────────────┴─────────────────────────────────┐    │
│  │            OUTPUT ADAPTERS (Driven)                  │    │
│  │  ┌──────────┐  ┌──────────┐  ┌──────────┐          │    │
│  │  │ Database │  │   JWT    │  │  Email   │  ...     │    │
│  │  │Repository│  │ Service  │  │ Service  │          │    │
│  │  └──────────┘  └──────────┘  └──────────┘          │    │
│  └──────────────────────────────────────────────────────┘    │
└─────────────────────────────────────────────────────────────┘
```

### 2. Monolito Modular

O sistema é organizado em módulos independentes que podem evoluir para microsserviços no futuro.

```
┌────────────────────────────────────────────────────────┐
│                   MONOLITH                              │
│                                                         │
│  ┌───────────────┐  ┌───────────────┐  ┌────────────┐ │
│  │  User Module  │  │  Task Module  │  │   Other    │ │
│  │               │  │               │  │  Modules   │ │
│  │  ┌─────────┐  │  │  ┌─────────┐  │  │            │ │
│  │  │ Domain  │  │  │  │ Domain  │  │  │    ...     │ │
│  │  ├─────────┤  │  │  ├─────────┤  │  │            │ │
│  │  │   App   │  │  │  │   App   │  │  │            │ │
│  │  ├─────────┤  │  │  ├─────────┤  │  │            │ │
│  │  │  Infra  │  │  │  │  Infra  │  │  │            │ │
│  │  └─────────┘  │  │  └─────────┘  │  │            │ │
│  └───────────────┘  └───────────────┘  └────────────┘ │
│                                                         │
│  ┌──────────────────────────────────────────────────┐  │
│  │          Shared Infrastructure                    │  │
│  │  (HTTP Server, JWT, Database, Logging, etc.)    │  │
│  └──────────────────────────────────────────────────┘  │
└────────────────────────────────────────────────────────┘
```

## 🏗️ Estrutura de Camadas

### Camada 1: Domain (Núcleo)

**Responsabilidade:** Contém a lógica de negócio pura e as regras do domínio.

**Características:**
- ✅ Não depende de nenhuma outra camada
- ✅ Não conhece frameworks ou bibliotecas externas
- ✅ Contém apenas lógica de negócio pura
- ✅ Altamente testável
- ✅ Reutilizável

**Componentes:**
- **Entities:** Objetos com identidade única e ciclo de vida
- **Value Objects:** Objetos imutáveis sem identidade
- **Domain Events:** Eventos que ocorrem no domínio
- **Domain Errors:** Erros específicos do domínio
- **Business Rules:** Regras e invariantes do negócio

**Exemplo - User Module:**
```
domain/
├── entity/
│   └── user.go              # Entidade User com comportamentos
├── valueobject/
│   ├── email.go            # Value Object Email com validação
│   └── password.go         # Value Object Password com hash
└── errors/
    └── errors.go           # Erros do domínio (ErrUserNotFound, etc)
```

**Regras:**
- Nenhuma dependência externa
- Apenas Go stdlib
- Validações no próprio domínio
- Comportamento rico (não anêmico)

### Camada 2: Application (Casos de Uso)

**Responsabilidade:** Orquestra o fluxo de dados entre as camadas e implementa os casos de uso.

**Características:**
- ✅ Depende apenas da camada de domínio
- ✅ Define interfaces (ports) para o mundo externo
- ✅ Implementa a lógica de aplicação
- ✅ Agnóstica de frameworks
- ✅ Fácil de testar

**Componentes:**
- **Input Ports:** Interfaces que definem o que a aplicação faz (use cases)
- **Output Ports:** Interfaces que definem dependências externas (repositories, services)
- **Use Cases:** Implementações dos casos de uso

**Exemplo - User Module:**
```
application/
├── port/
│   ├── input/
│   │   └── user_service.go      # Interface: Register, Login, etc
│   └── output/
│       ├── user_repository.go   # Interface: Save, FindByID, etc
│       └── token_generator.go   # Interface: GenerateToken, ValidateToken
└── usecase/
    └── user_service_impl.go     # Implementação dos casos de uso
```

**Fluxo de Dados:**
```
Input Port (Interface)
    ↓
Use Case Implementation
    ↓
Domain Entities (Business Logic)
    ↓
Output Port (Interface)
```

**Dependency Inversion Principle:**
- A camada de aplicação define as interfaces
- A camada de infraestrutura implementa as interfaces
- O fluxo de controle vai da aplicação → infraestrutura
- A dependência de código vai infraestrutura → aplicação

### Camada 3: Infrastructure (Adaptadores)

**Responsabilidade:** Implementa os detalhes técnicos e conecta com o mundo externo.

**Características:**
- ✅ Implementa as interfaces (ports) da camada de aplicação
- ✅ Contém código específico de frameworks
- ✅ Pode ser facilmente substituída
- ✅ Isolada do domínio

**Componentes:**

**Input Adapters (Drivers):**
- HTTP Handlers (REST API)
- CLI Commands
- gRPC Handlers
- GraphQL Resolvers
- Message Queue Consumers

**Output Adapters (Driven):**
- Database Repositories
- External API Clients
- Message Queue Publishers
- File System
- Cache (Redis)
- Email Service

**Exemplo - User Module:**
```
infrastructure/
├── adapter/
│   ├── input/
│   │   └── http/
│   │       ├── user_handler.go    # HTTP handlers
│   │       ├── middleware.go      # Auth middleware
│   │       └── routes.go          # Route configuration
│   └── output/
│       └── persistence/
│           └── memory_user_repository.go  # In-memory implementation
└── dto/
    └── user_dto.go                # Data Transfer Objects
```

## 🔄 Fluxo de Requisição Completo

```
1. HTTP Request
   │
   ▼
2. HTTP Handler (Input Adapter)
   │ - Parse request
   │ - Validate DTO
   │ - Extract parameters
   │
   ▼
3. Use Case (Application Layer)
   │ - Validate business rules
   │ - Orchestrate domain logic
   │ - Call domain entities
   │
   ▼
4. Domain Entity
   │ - Execute business logic
   │ - Validate invariants
   │ - Apply domain rules
   │
   ▼
5. Repository (Output Adapter)
   │ - Persist changes
   │ - Query data
   │
   ▼
6. Database
   │
   ▼
7. Response
   │ - Map entity to DTO
   │ - Return HTTP response
```

### Exemplo Concreto: Registro de Usuário

```go
// 1. HTTP Request
POST /api/users/register
{
  "email": "user@example.com",
  "password": "Pass123",
  "name": "John Doe"
}

// 2. HTTP Handler (infrastructure/adapter/input/http/user_handler.go)
func (h *UserHandler) Register(c *fiber.Ctx) error {
    var req dto.RegisterRequest
    c.BodyParser(&req)              // Parse request
    
    // 3. Call Use Case (application/usecase/user_service_impl.go)
    user, err := h.userService.Register(ctx, req.Email, req.Password, req.Name)
    
    return c.JSON(dto.ToUserResponse(user))
}

// 3. Use Case Implementation
func (s *UserServiceImpl) Register(ctx, email, password, name string) (*entity.User, error) {
    // Create value objects (domain/valueobject)
    emailVO, _ := valueobject.NewEmail(email)        // Validates email
    passwordVO, _ := valueobject.NewPassword(password) // Hashes password
    
    // Check if user exists (calls output port)
    exists, _ := s.userRepository.ExistsByEmail(ctx, emailVO)
    
    // Create domain entity (domain/entity)
    user := entity.NewUser(emailVO, passwordVO, name)  // Business logic
    
    // Persist (calls output port)
    s.userRepository.Save(ctx, user)
    
    return user, nil
}

// 4. Domain Entity (domain/entity/user.go)
func NewUser(email Email, password Password, name string) *User {
    // Business rules and validations
    return &User{
        id:        uuid.New(),
        email:     email,
        password:  password,
        name:      name,
        isActive:  true,
        createdAt: time.Now(),
    }
}

// 5. Repository Implementation (infrastructure/adapter/output/persistence)
func (r *MemoryUserRepository) Save(ctx context.Context, user *entity.User) error {
    r.users[user.ID()] = user
    return nil
}

// 6. Response
{
  "message": "User registered successfully",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "user@example.com",
    "name": "John Doe",
    "is_active": true
  }
}
```

## 🔌 Ports (Interfaces)

### Input Ports (Use Cases)

Definem **O QUE** a aplicação faz.

```go
// application/port/input/user_service.go
type UserService interface {
    Register(ctx context.Context, email, password, name string) (*entity.User, error)
    Login(ctx context.Context, email, password string) (token string, user *entity.User, err error)
    GetUserByID(ctx context.Context, userID uuid.UUID) (*entity.User, error)
    UpdateUser(ctx context.Context, userID uuid.UUID, name string) (*entity.User, error)
    ChangePassword(ctx context.Context, userID uuid.UUID, oldPass, newPass string) error
}
```

**Quem chama:** Input Adapters (HTTP Handlers, CLI, etc)
**Quem implementa:** Use Cases (application/usecase)

### Output Ports (Dependencies)

Definem **COMO** a aplicação se comunica com o mundo externo.

```go
// application/port/output/user_repository.go
type UserRepository interface {
    Save(ctx context.Context, user *entity.User) error
    Update(ctx context.Context, user *entity.User) error
    FindByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
    FindByEmail(ctx context.Context, email valueobject.Email) (*entity.User, error)
}

// application/port/output/token_generator.go
type TokenGenerator interface {
    GenerateToken(userID uuid.UUID, email string, expiresIn time.Duration) (string, error)
    ValidateToken(token string) (uuid.UUID, error)
}
```

**Quem define:** Application Layer
**Quem implementa:** Output Adapters (infrastructure/adapter/output)
**Quem usa:** Use Cases

## 📦 Módulos

### Estrutura de um Módulo

Cada módulo segue a mesma estrutura de camadas:

```
internal/modules/[module-name]/
├── domain/                    # Camada de Domínio
│   ├── entity/               # Entidades
│   ├── valueobject/          # Value Objects
│   └── errors/               # Erros de domínio
├── application/               # Camada de Aplicação
│   ├── port/
│   │   ├── input/           # Use Cases (interfaces)
│   │   └── output/          # Dependências externas (interfaces)
│   └── usecase/             # Implementação dos casos de uso
└── infrastructure/            # Camada de Infraestrutura
    ├── adapter/
    │   ├── input/           # HTTP, CLI, etc
    │   └── output/          # Repositories, Services
    └── dto/                 # Data Transfer Objects
```

### Comunicação Entre Módulos

**Regras:**
1. Módulos se comunicam apenas através de interfaces bem definidas
2. Não há dependência direta entre módulos
3. Eventos de domínio podem ser usados para comunicação assíncrona
4. Cada módulo tem seu próprio bounded context

**Opções de Comunicação:**

```
┌───────────────┐         ┌───────────────┐
│  User Module  │         │  Task Module  │
└───────┬───────┘         └───────┬───────┘
        │                         │
        │  1. Direct Call         │
        │  (via interface)        │
        ├────────────────────────►│
        │                         │
        │  2. Domain Events       │
        │  (async)                │
        ├────────►│               │
        │         ▼               │
        │    Event Bus            │
        │         │               │
        │         └──────────────►│
        │                         │
        │  3. Shared Service      │
        │  (via output port)      │
        ├────────►│               │
        │         ▼               │
        │   Shared Context        │
        │         │               │
        │         └──────────────►│
        │                         │
```

## 🔐 Segurança

### Autenticação e Autorização

```
┌─────────────────────────────────────────────────────────┐
│                    Authentication Flow                   │
└─────────────────────────────────────────────────────────┘

1. Login Request
   POST /api/users/login
   { "email": "...", "password": "..." }
   │
   ▼
2. Validate Credentials (Use Case)
   - Find user by email
   - Verify password (bcrypt)
   │
   ▼
3. Generate JWT Token (Output Adapter)
   - Create claims (userID, email, exp)
   - Sign with secret key
   │
   ▼
4. Return Token
   { "token": "eyJhbGc...", "user": {...} }

┌─────────────────────────────────────────────────────────┐
│                   Authorization Flow                     │
└─────────────────────────────────────────────────────────┘

1. Protected Request
   GET /api/users/me
   Authorization: Bearer eyJhbGc...
   │
   ▼
2. Auth Middleware (Input Adapter)
   - Extract token from header
   - Validate token signature
   - Check expiration
   │
   ▼
3. Inject User Context
   - Extract userID from claims
   - Store in request context
   │
   ▼
4. Handler Execution
   - Access userID from context
   - Execute business logic
```

### Camadas de Segurança

```
┌──────────────────────────────────────────────────────┐
│  1. Transport Layer (HTTPS)                          │
│     - Encrypt data in transit                        │
└──────────────────────────────────────────────────────┘
                        ▼
┌──────────────────────────────────────────────────────┐
│  2. Authentication (JWT)                             │
│     - Verify user identity                           │
└──────────────────────────────────────────────────────┘
                        ▼
┌──────────────────────────────────────────────────────┐
│  3. Authorization (Middleware)                       │
│     - Check permissions                              │
└──────────────────────────────────────────────────────┘
                        ▼
┌──────────────────────────────────────────────────────┐
│  4. Input Validation (DTO)                           │
│     - Validate request data                          │
└──────────────────────────────────────────────────────┘
                        ▼
┌──────────────────────────────────────────────────────┐
│  5. Business Rules (Domain)                          │
│     - Validate business constraints                  │
└──────────────────────────────────────────────────────┘
                        ▼
┌──────────────────────────────────────────────────────┐
│  6. Data Layer (Repository)                          │
│     - SQL injection prevention                       │
└──────────────────────────────────────────────────────┘
```

## 🧪 Testabilidade

### Estratégia de Testes

```
┌────────────────────────────────────────────────────────┐
│                    Test Pyramid                         │
└────────────────────────────────────────────────────────┘

                    ▲
                   ╱ ╲
                  ╱   ╲
                 ╱ E2E ╲          - End-to-End Tests
                ╱───────╲         - Few, slow, expensive
               ╱         ╲
              ╱───────────╲
             ╱ Integration ╲      - Integration Tests
            ╱───────────────╲     - Medium amount
           ╱                 ╲
          ╱───────────────────╲
         ╱     Unit Tests      ╲  - Unit Tests
        ╱───────────────────────╲ - Many, fast, cheap
       ╱                         ╲
```

### Níveis de Teste

**1. Unit Tests (Domain Layer)**
```go
// domain/entity/user_test.go
func TestUser_VerifyPassword(t *testing.T) {
    password, _ := valueobject.NewPassword("Pass123")
    user := entity.NewUser(email, password, "John")
    
    assert.True(t, user.VerifyPassword("Pass123"))
    assert.False(t, user.VerifyPassword("wrong"))
}
```

**2. Use Case Tests (Application Layer)**
```go
// application/usecase/user_service_test.go
func TestUserService_Register(t *testing.T) {
    // Mock dependencies
    mockRepo := &MockUserRepository{}
    mockTokenGen := &MockTokenGenerator{}
    
    service := NewUserService(mockRepo, mockTokenGen, 24*time.Hour)
    
    user, err := service.Register(ctx, "test@example.com", "Pass123", "Test")
    
    assert.NoError(t, err)
    assert.NotNil(t, user)
}
```

**3. Integration Tests (Infrastructure Layer)**
```go
// infrastructure/adapter/output/persistence/user_repository_test.go
func TestUserRepository_Save(t *testing.T) {
    repo := NewMemoryUserRepository()
    user := createTestUser()
    
    err := repo.Save(context.Background(), user)
    
    assert.NoError(t, err)
    
    found, err := repo.FindByID(context.Background(), user.ID())
    assert.NoError(t, err)
    assert.Equal(t, user.Email(), found.Email())
}
```

**4. E2E Tests**
```go
// e2e/user_test.go
func TestUserRegistrationFlow(t *testing.T) {
    // Start test server
    server := setupTestServer()
    
    // Register
    resp := httptest.NewRequest("POST", "/api/users/register", body)
    assert.Equal(t, 201, resp.StatusCode)
    
    // Login
    resp = httptest.NewRequest("POST", "/api/users/login", body)
    assert.Equal(t, 200, resp.StatusCode)
}
```

## 🚀 Evolução para Microsserviços

A arquitetura hexagonal facilita a migração para microsserviços:

```
┌──────────────────────────────────────────────────────┐
│                    MONOLITH                           │
│  ┌────────────┐  ┌────────────┐  ┌────────────┐     │
│  │User Module │  │Task Module │  │Other Module│     │
│  └────────────┘  └────────────┘  └────────────┘     │
└──────────────────────────────────────────────────────┘
                        │
                        │ Extract Module
                        ▼
┌──────────────────────────────────────────────────────┐
│                MODULAR MONOLITH                       │
│  ┌────────────┐                ┌────────────┐        │
│  │Task Module │                │Other Module│        │
│  └────────────┘                └────────────┘        │
└──────────────────────────────────────────────────────┘
                                              
         ┌────────────────┐                   
         │ User Service   │  ← Extracted      
         │ (Microservice) │                   
         └────────────────┘                   
```

### Passos para Extração:

1. **Identificar Boundaries:** Módulos com baixo acoplamento
2. **Extrair Domain:** Copiar domain layer (já está isolado)
3. **Extrair Application:** Copiar use cases e ports
4. **Criar Novos Adapters:** Implementar HTTP client, gRPC, etc
5. **Deploy Independente:** Cada serviço em seu próprio runtime
6. **Comunicação:** REST, gRPC, Message Queue

## 📈 Escalabilidade

### Vertical Scaling (Monolith)

```
┌────────────────┐
│   Monolith     │
│   (1 instance) │
└────────────────┘
        │
        ▼
┌────────────────┐
│   Monolith     │
│   (Bigger)     │
└────────────────┘
```

### Horizontal Scaling (Monolith)

```
         Load Balancer
              │
    ┌─────────┼─────────┐
    │         │         │
    ▼         ▼         ▼
┌────────┐┌────────┐┌────────┐
│Instance││Instance││Instance│
│   #1   ││   #2   ││   #3   │
└────────┘└────────┘└────────┘
```

### Module-Based Scaling (Future)

```
         API Gateway
              │
    ┌─────────┼─────────────┐
    │         │             │
    ▼         ▼             ▼
┌────────┐┌────────┐  ┌────────┐
│  User  ││  Task  │  │ Other  │
│Service ││Service │  │Service │
│(3 inst)││(5 inst)│  │(2 inst)│
└────────┘└────────┘  └────────┘
```

## 🎯 Boas Práticas

### SOLID Principles

✅ **Single Responsibility:** Cada classe tem uma única responsabilidade
✅ **Open/Closed:** Aberto para extensão, fechado para modificação
✅ **Liskov Substitution:** Subtipos devem ser substituíveis
✅ **Interface Segregation:** Interfaces pequenas e específicas
✅ **Dependency Inversion:** Dependa de abstrações, não de implementações

### Clean Code

✅ Nomes descritivos e significativos
✅ Funções pequenas e focadas
✅ Evitar acoplamento
✅ Separação de responsabilidades
✅ Testes automatizados

### Domain-Driven Design

✅ Ubiquitous Language
✅ Bounded Contexts
✅ Aggregates e Entities
✅ Value Objects
✅ Domain Events
✅ Repository Pattern

## 📚 Referências

- **Hexagonal Architecture:** Alistair Cockburn
- **Clean Architecture:** Robert C. Martin (Uncle Bob)
- **Domain-Driven Design:** Eric Evans
- **Modular Monolith:** Kamil Grzybek
- **Ports and Adapters:** Alistair Cockburn

## 🔄 Diagrama de Dependências

```
┌─────────────────────────────────────────────────────┐
│                  CMD (main.go)                       │
│  - Wires everything together                        │
│  - Dependency injection                             │
└────────────┬────────────────────────────────────────┘
             │
             ▼
┌─────────────────────────────────────────────────────┐
│              INFRASTRUCTURE LAYER                    │
│  ┌──────────────┐  ┌──────────────┐                │
│  │HTTP Handlers │  │  Repositories │                │
│  └──────┬───────┘  └──────┬────────┘                │
└─────────┼──────────────────┼──────────────────────────┘
          │                  │
          │                  │
          ▼                  ▼
┌─────────────────────────────────────────────────────┐
│              APPLICATION LAYER                       │
│  ┌──────────────────────────────────────┐           │
│  │        Input/Output Ports             │           │
│  │          (Interfaces)                 │           │
│  └──────────────┬───────────────────────┘           │
│                 │                                    │
│  ┌──────────────▼───────────────────────┐           │
│  │         Use Case Implementations      │           │
│  └──────────────┬───────────────────────┘           │
└─────────────────┼──────────────────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────────────────────┐
│                DOMAIN LAYER                          │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐          │
│  │ Entities │  │  Value   │  │  Domain  │          │
│  │          │  │ Objects  │  │  Errors  │          │
│  └──────────┘  └──────────┘  └──────────┘          │
└─────────────────────────────────────────────────────┘

Setas indicam direção de dependência (importação)
Domain não depende de nada
Application depende apenas de Domain
Infrastructure depende de Application e Domain
```

## 🎓 Conclusão

Esta arquitetura proporciona:

✅ **Manutenibilidade:** Código organizado e fácil de entender
✅ **Testabilidade:** Cada camada pode ser testada isoladamente
✅ **Flexibilidade:** Fácil trocar implementações (banco, framework)
✅ **Escalabilidade:** Caminho claro para crescimento
✅ **Independência:** Domínio independente de detalhes técnicos
✅ **Evolução:** Pode evoluir para microsserviços quando necessário

A arquitetura hexagonal com monolito modular oferece o melhor dos dois mundos: a simplicidade de um monolito com a modularidade necessária para crescimento futuro.