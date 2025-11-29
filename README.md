# UserMes Backend - Hexagonal Architecture with Modular Monolith

Uma API RESTful construída em Go seguindo os princípios de **Arquitetura Hexagonal** (Ports and Adapters) com **Monolito Modular**.

## 🏗️ Arquitetura

Este projeto implementa uma arquitetura hexagonal (também conhecida como Ports and Adapters), que promove:

- **Separação de responsabilidades**: Domínio isolado da infraestrutura
- **Testabilidade**: Facilita testes unitários e de integração
- **Flexibilidade**: Fácil troca de adaptadores (banco de dados, frameworks, etc)
- **Escalabilidade**: Caminho claro para migração para microsserviços

### Estrutura de Diretórios

```
usermes-backend/
├── cmd/                                    # Entry points da aplicação
│   └── main.go                            # Main application
├── internal/                              # Código privado da aplicação
│   ├── modules/                           # Módulos do monolito
│   │   └── user/                          # Módulo de usuário
│   │       ├── domain/                    # Camada de Domínio (Núcleo)
│   │       │   ├── entity/               # Entidades do domínio
│   │       │   │   └── user.go
│   │       │   ├── valueobject/          # Value Objects
│   │       │   │   ├── email.go
│   │       │   │   └── password.go
│   │       │   └── errors/               # Erros do domínio
│   │       │       └── errors.go
│   │       ├── application/               # Camada de Aplicação (Casos de Uso)
│   │       │   ├── port/                 # Portas (Interfaces)
│   │       │   │   ├── input/           # Portas de entrada (Use Cases)
│   │       │   │   │   └── user_service.go
│   │       │   │   └── output/          # Portas de saída (Repositories, etc)
│   │       │   │       ├── user_repository.go
│   │       │   │       └── token_generator.go
│   │       │   └── usecase/              # Implementação dos casos de uso
│   │       │       └── user_service_impl.go
│   │       └── infrastructure/            # Camada de Infraestrutura (Adaptadores)
│   │           ├── adapter/
│   │           │   ├── input/
│   │           │   │   └── http/         # Adaptador HTTP (Controllers)
│   │           │   │       ├── user_handler.go
│   │           │   │       ├── middleware.go
│   │           │   │       └── routes.go
│   │           │   └── output/
│   │           │       └── persistence/  # Adaptador de persistência
│   │           │           └── memory_user_repository.go
│   │           └── dto/                   # Data Transfer Objects
│   │               └── user_dto.go
│   └── shared/                            # Código compartilhado entre módulos
│       └── infrastructure/
│           ├── http/
│           │   └── server/
│           │       └── fiber_server.go
│           └── security/
│               └── jwt_token_generator.go
└── pkg/                                   # Pacotes públicos reutilizáveis
```

## 📦 Camadas da Arquitetura Hexagonal

### 1. **Domain Layer** (Núcleo)
O coração da aplicação, contendo a lógica de negócio pura:

- **Entities**: Objetos com identidade única (`User`)
- **Value Objects**: Objetos imutáveis sem identidade (`Email`, `Password`)
- **Domain Errors**: Erros específicos do domínio
- **Regras de Negócio**: Validações e comportamentos do domínio

**Características:**
- Não depende de nenhuma outra camada
- Não conhece frameworks ou bibliotecas externas
- Contém apenas lógica de negócio pura

### 2. **Application Layer** (Casos de Uso)
Orquestra o fluxo de dados e coordena as operações:

- **Input Ports**: Interfaces que definem os casos de uso (o que a aplicação faz)
- **Output Ports**: Interfaces que definem as dependências externas (repositórios, serviços)
- **Use Cases**: Implementação dos casos de uso usando as entidades do domínio

**Características:**
- Depende apenas da camada de domínio
- Define interfaces (ports) que serão implementadas pela camada de infraestrutura
- Contém a lógica de aplicação (orquestração)

### 3. **Infrastructure Layer** (Adaptadores)
Implementa os detalhes técnicos e se conecta com o mundo externo:

- **Input Adapters**: HTTP handlers, CLI, gRPC, etc.
- **Output Adapters**: Implementações de repositórios, clientes de APIs, etc.
- **DTOs**: Objetos para transferência de dados entre camadas

**Características:**
- Implementa as interfaces (ports) definidas na camada de aplicação
- Contém código específico de frameworks e bibliotecas
- Pode ser facilmente substituída sem afetar o domínio

## 🎯 Módulo User

O módulo User é o primeiro módulo do monolito, responsável por:

- ✅ Registro de usuários
- ✅ Autenticação (Login)
- ✅ Gerenciamento de perfil
- ✅ Alteração de senha
- ✅ Ativação/Desativação de contas

### Endpoints Disponíveis

#### Públicos (sem autenticação)

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

#### Protegidos (requerem autenticação)

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

## 🚀 Como Executar

### Pré-requisitos

- Go 1.22.1 ou superior
- Make (opcional)

### Instalação

```bash
# Clone o repositório
git clone <repository-url>
cd usermes-backend

# Baixe as dependências
go mod download

# Execute a aplicação
go run cmd/main.go
```

### Usando Make

```bash
# Execute a aplicação
make run

# Build da aplicação
make build

# Execute os testes
make test
```

## 🧪 Testando a API

### 1. Registrar um novo usuário

```bash
curl -X POST http://localhost:3000/api/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "SecurePass123",
    "name": "Test User"
  }'
```

### 2. Fazer login

```bash
curl -X POST http://localhost:3000/api/users/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "SecurePass123"
  }'
```

Copie o token retornado para usar nas próximas requisições.

### 3. Obter perfil do usuário

```bash
curl -X GET http://localhost:3000/api/users/me \
  -H "Authorization: Bearer <seu-token>"
```

### 4. Health Check

```bash
curl http://localhost:3000/health
```

## 🔒 Segurança

- **Senha**: Hash com bcrypt (cost factor 12)
- **JWT**: Tokens com expiração de 24 horas
- **Validação**: Email e senha validados com regras de negócio
- **CORS**: Configurado para permitir requisições de qualquer origem (ajustar em produção)

### Regras de Senha

- Mínimo de 8 caracteres
- Máximo de 72 caracteres (limitação do bcrypt)
- Deve conter pelo menos uma letra
- Deve conter pelo menos um número

## 📝 Boas Práticas Implementadas

1. **Dependency Inversion**: As camadas superiores não dependem de implementações concretas
2. **Single Responsibility**: Cada componente tem uma única responsabilidade
3. **Open/Closed**: Aberto para extensão, fechado para modificação
4. **Interface Segregation**: Interfaces pequenas e focadas
5. **Domain-Driven Design**: Modelagem rica do domínio
6. **Value Objects**: Validação e encapsulamento de valores
7. **Repository Pattern**: Abstração de persistência
8. **Use Case Pattern**: Casos de uso explícitos e testáveis

## 🔄 Próximos Passos

### Implementações Futuras

- [ ] Integração com banco de dados (PostgreSQL/MySQL)
- [ ] Redis para cache e sessions
- [ ] Refresh tokens
- [ ] Rate limiting
- [ ] Logs estruturados
- [ ] Métricas e observabilidade
- [ ] Testes unitários e de integração
- [ ] CI/CD pipeline
- [ ] Docker e Docker Compose
- [ ] Migração de banco de dados
- [ ] Swagger/OpenAPI documentation
- [ ] Novos módulos (Tasks, Projects, etc)

### Adicionando Novo Módulo

Para adicionar um novo módulo ao monolito:

1. Crie a estrutura de diretórios em `internal/modules/[nome-do-modulo]`
2. Implemente as camadas: domain → application → infrastructure
3. Registre as rotas no `main.go`
4. Mantenha a independência entre módulos

Exemplo:
```
internal/modules/task/
├── domain/
├── application/
└── infrastructure/
```

## 🤝 Contribuindo

1. Fork o projeto
2. Crie uma branch para sua feature (`git checkout -b feature/AmazingFeature`)
3. Commit suas mudanças (`git commit -m 'Add some AmazingFeature'`)
4. Push para a branch (`git push origin feature/AmazingFeature`)
5. Abra um Pull Request

## 📄 Licença

Este projeto está sob a licença MIT.

## 📚 Referências

- [Hexagonal Architecture](https://alistair.cockburn.us/hexagonal-architecture/)
- [Domain-Driven Design](https://martinfowler.com/bliki/DomainDrivenDesign.html)
- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Modular Monolith](https://www.kamilgrzybek.com/design/modular-monolith-primer/)