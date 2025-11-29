# Quick Start Guide - UserMes Backend

Guia rápido para começar a usar o UserMes Backend.

## 🚀 Instalação Rápida

### Pré-requisitos

- Go 1.22.1 ou superior
- curl (para testes)

### Passos

1. **Clone o repositório**
```bash
git clone <seu-repositório>
cd usermes-backend
```

2. **Instale as dependências**
```bash
go mod download
```

3. **Execute a aplicação**
```bash
go run cmd/main.go
```

A API estará disponível em: `http://localhost:3000`

## ✅ Teste Rápido

Abra outro terminal e execute:

```bash
# Health check
curl http://localhost:3000/health
```

Deve retornar:
```json
{"status":"ok","timestamp":1705318800}
```

## 📝 Seu Primeiro Usuário

### 1. Registrar um usuário

```bash
curl -X POST http://localhost:3000/api/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "seu@email.com",
    "password": "SuaSenha123",
    "name": "Seu Nome"
  }'
```

**Resposta esperada:**
```json
{
  "message": "User registered successfully",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "seu@email.com",
    "name": "Seu Nome",
    "is_active": true,
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z"
  }
}
```

### 2. Fazer login

```bash
curl -X POST http://localhost:3000/api/users/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "seu@email.com",
    "password": "SuaSenha123"
  }'
```

**Resposta esperada:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "seu@email.com",
    "name": "Seu Nome",
    "is_active": true,
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z",
    "last_login_at": "2024-01-15T11:00:00Z"
  }
}
```

**⚠️ IMPORTANTE:** Copie o token retornado! Você precisará dele para as próximas requisições.

### 3. Obter seu perfil

Substitua `<SEU-TOKEN>` pelo token que você recebeu no login:

```bash
curl -X GET http://localhost:3000/api/users/me \
  -H "Authorization: Bearer <SEU-TOKEN>"
```

**Resposta esperada:**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "email": "seu@email.com",
  "name": "Seu Nome",
  "is_active": true,
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T10:30:00Z",
  "last_login_at": "2024-01-15T11:00:00Z"
}
```

## 🎯 Teste Automatizado

Execute todos os testes de uma vez:

```bash
# Terminal 1: Inicie o servidor
go run cmd/main.go

# Terminal 2: Execute os testes
./simple_test.sh
```

## 🛠️ Usando Makefile

```bash
# Executar a aplicação
make run

# Build
make build

# Executar testes
make test

# Ver todos os comandos disponíveis
make help
```

## 📁 Estrutura do Projeto

```
usermes-backend/
├── cmd/                          # Entry point
│   └── main.go
├── internal/
│   ├── modules/                  # Módulos do monolito
│   │   └── user/                 # Módulo de usuário
│   │       ├── domain/           # Entidades e regras de negócio
│   │       ├── application/      # Casos de uso
│   │       └── infrastructure/   # HTTP, Database, etc
│   └── shared/                   # Código compartilhado
│       └── infrastructure/
└── pkg/                          # Pacotes públicos
```

## 📚 Endpoints Disponíveis

### Públicos (sem autenticação)

- `POST /api/users/register` - Registrar novo usuário
- `POST /api/users/login` - Fazer login
- `GET /health` - Health check

### Protegidos (requerem token)

- `GET /api/users/me` - Obter perfil do usuário logado
- `GET /api/users/:id` - Obter usuário por ID
- `PUT /api/users/:id` - Atualizar nome do usuário
- `POST /api/users/:id/change-password` - Alterar senha
- `POST /api/users/:id/deactivate` - Desativar conta
- `POST /api/users/:id/activate` - Ativar conta

## 🔐 Regras de Senha

- Mínimo de 8 caracteres
- Máximo de 72 caracteres
- Deve conter pelo menos uma letra
- Deve conter pelo menos um número

## 💡 Dicas

### Salvar token em variável (Linux/Mac)

```bash
# Fazer login e salvar token
TOKEN=$(curl -s -X POST http://localhost:3000/api/users/login \
  -H "Content-Type: application/json" \
  -d '{"email":"seu@email.com","password":"SuaSenha123"}' \
  | grep -o '"token":"[^"]*"' | sed 's/"token":"//; s/"$//')

# Usar o token
curl -X GET http://localhost:3000/api/users/me \
  -H "Authorization: Bearer $TOKEN"
```

### Salvar token em variável (Windows PowerShell)

```powershell
# Fazer login e salvar token
$response = Invoke-RestMethod -Uri "http://localhost:3000/api/users/login" `
  -Method POST `
  -ContentType "application/json" `
  -Body '{"email":"seu@email.com","password":"SuaSenha123"}'

$token = $response.token

# Usar o token
Invoke-RestMethod -Uri "http://localhost:3000/api/users/me" `
  -Method GET `
  -Headers @{Authorization="Bearer $token"}
```

## 🐛 Troubleshooting

### Porta 3000 já em uso

```bash
# Encontrar processo usando a porta
lsof -i :3000

# Matar o processo
kill -9 <PID>
```

### Erro "Authorization header is required"

Certifique-se de incluir o header de autorização:
```bash
-H "Authorization: Bearer <seu-token>"
```

### Erro "Invalid or expired token"

O token expira após 24 horas. Faça login novamente para obter um novo token.

### Erro "email already exists"

Use um email diferente ou faça login com o email existente.

## 🔄 Resetar Dados

Como está usando repositório em memória, basta reiniciar o servidor:

```bash
# Pressione Ctrl+C para parar o servidor
# Execute novamente
go run cmd/main.go
```

Todos os dados serão perdidos e você poderá começar do zero.

## 📖 Próximos Passos

1. Leia a [Documentação Completa](README.md)
2. Entenda a [Arquitetura](ARCHITECTURE.md)
3. Veja mais [Exemplos de Uso](EXAMPLES.md)
4. Contribua com o projeto!

## 🆘 Precisa de Ajuda?

- Abra uma issue no GitHub
- Consulte a documentação completa
- Verifique os exemplos

## 📝 Notas Importantes

⚠️ **Ambiente de Desenvolvimento:**
- Dados armazenados em memória (não persistem)
- JWT secret está hardcoded (trocar em produção)
- CORS configurado para aceitar qualquer origem

⚠️ **Antes de ir para Produção:**
- Configure variáveis de ambiente
- Use banco de dados real (PostgreSQL, MySQL)
- Configure HTTPS
- Implemente rate limiting
- Configure logging apropriado
- Use JWT secret seguro e único

## ✨ Exemplos Rápidos

### Fluxo completo em um script

```bash
#!/bin/bash

BASE_URL="http://localhost:3000"

# 1. Registrar
curl -X POST $BASE_URL/api/users/register \
  -H "Content-Type: application/json" \
  -d '{"email":"demo@test.com","password":"Demo123","name":"Demo User"}'

# 2. Login e salvar token
TOKEN=$(curl -s -X POST $BASE_URL/api/users/login \
  -H "Content-Type: application/json" \
  -d '{"email":"demo@test.com","password":"Demo123"}' \
  | grep -o '"token":"[^"]*"' | sed 's/"token":"//; s/"$//')

# 3. Ver perfil
curl -X GET $BASE_URL/api/users/me \
  -H "Authorization: Bearer $TOKEN"
```

Salve isso em `test.sh`, dê permissão de execução (`chmod +x test.sh`) e execute (`./test.sh`).

---

**Pronto para começar!** 🚀

Se tudo funcionou, você está pronto para desenvolver novos recursos ou adicionar novos módulos ao monolito.