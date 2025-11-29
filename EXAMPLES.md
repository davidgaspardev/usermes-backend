# API Usage Examples

Este documento contém exemplos práticos de como usar a API do UserMes.

## 📋 Índice

- [Configuração](#configuração)
- [Autenticação](#autenticação)
- [Operações de Usuário](#operações-de-usuário)
- [Exemplos com cURL](#exemplos-com-curl)
- [Exemplos com JavaScript/Fetch](#exemplos-com-javascriptfetch)
- [Exemplos com Go](#exemplos-com-go)
- [Códigos de Status](#códigos-de-status)
- [Tratamento de Erros](#tratamento-de-erros)

## Configuração

Base URL: `http://localhost:3000`

Todas as rotas da API começam com `/api`

## Autenticação

A API usa JWT (JSON Web Tokens) para autenticação. Após fazer login, você receberá um token que deve ser incluído no header `Authorization` de todas as requisições protegidas.

Formato: `Authorization: Bearer <seu-token>`

## Operações de Usuário

### 1. Registrar Novo Usuário

**Endpoint:** `POST /api/users/register`

**Headers:**
```
Content-Type: application/json
```

**Body:**
```json
{
  "email": "joao@example.com",
  "password": "SenhaSegura123",
  "name": "João Silva"
}
```

**Resposta de Sucesso (201):**
```json
{
  "message": "User registered successfully",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "joao@example.com",
    "name": "João Silva",
    "is_active": true,
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z"
  }
}
```

### 2. Login

**Endpoint:** `POST /api/users/login`

**Headers:**
```
Content-Type: application/json
```

**Body:**
```json
{
  "email": "joao@example.com",
  "password": "SenhaSegura123"
}
```

**Resposta de Sucesso (200):**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "joao@example.com",
    "name": "João Silva",
    "is_active": true,
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z",
    "last_login_at": "2024-01-15T11:00:00Z"
  }
}
```

### 3. Obter Perfil Atual

**Endpoint:** `GET /api/users/me`

**Headers:**
```
Authorization: Bearer <seu-token>
```

**Resposta de Sucesso (200):**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "email": "joao@example.com",
  "name": "João Silva",
  "is_active": true,
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T10:30:00Z",
  "last_login_at": "2024-01-15T11:00:00Z"
}
```

### 4. Obter Usuário por ID

**Endpoint:** `GET /api/users/:id`

**Headers:**
```
Authorization: Bearer <seu-token>
```

**Resposta de Sucesso (200):**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "email": "joao@example.com",
  "name": "João Silva",
  "is_active": true,
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T10:30:00Z",
  "last_login_at": "2024-01-15T11:00:00Z"
}
```

### 5. Atualizar Usuário

**Endpoint:** `PUT /api/users/:id`

**Headers:**
```
Content-Type: application/json
Authorization: Bearer <seu-token>
```

**Body:**
```json
{
  "name": "João Pedro Silva"
}
```

**Resposta de Sucesso (200):**
```json
{
  "message": "User updated successfully",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "joao@example.com",
    "name": "João Pedro Silva",
    "is_active": true,
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T12:00:00Z",
    "last_login_at": "2024-01-15T11:00:00Z"
  }
}
```

### 6. Alterar Senha

**Endpoint:** `POST /api/users/:id/change-password`

**Headers:**
```
Content-Type: application/json
Authorization: Bearer <seu-token>
```

**Body:**
```json
{
  "old_password": "SenhaSegura123",
  "new_password": "NovaSenhaSegura456"
}
```

**Resposta de Sucesso (200):**
```json
{
  "message": "Password changed successfully"
}
```

### 7. Desativar Conta

**Endpoint:** `POST /api/users/:id/deactivate`

**Headers:**
```
Authorization: Bearer <seu-token>
```

**Resposta de Sucesso (200):**
```json
{
  "message": "User deactivated successfully"
}
```

### 8. Ativar Conta

**Endpoint:** `POST /api/users/:id/activate`

**Headers:**
```
Authorization: Bearer <seu-token>
```

**Resposta de Sucesso (200):**
```json
{
  "message": "User activated successfully"
}
```

### 9. Health Check

**Endpoint:** `GET /health`

**Resposta de Sucesso (200):**
```json
{
  "status": "ok",
  "timestamp": 1705318800
}
```

## Exemplos com cURL

### Registrar Usuário

```bash
curl -X POST http://localhost:3000/api/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "maria@example.com",
    "password": "SenhaForte123",
    "name": "Maria Santos"
  }'
```

### Login

```bash
curl -X POST http://localhost:3000/api/users/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "maria@example.com",
    "password": "SenhaForte123"
  }'
```

### Salvar token em variável (Linux/Mac)

```bash
TOKEN=$(curl -s -X POST http://localhost:3000/api/users/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "maria@example.com",
    "password": "SenhaForte123"
  }' | jq -r '.token')

echo $TOKEN
```

### Obter Perfil (usando token)

```bash
curl -X GET http://localhost:3000/api/users/me \
  -H "Authorization: Bearer $TOKEN"
```

### Atualizar Nome

```bash
curl -X PUT http://localhost:3000/api/users/550e8400-e29b-41d4-a716-446655440000 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "name": "Maria Santos Silva"
  }'
```

### Alterar Senha

```bash
curl -X POST http://localhost:3000/api/users/550e8400-e29b-41d4-a716-446655440000/change-password \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "old_password": "SenhaForte123",
    "new_password": "NovaSenhaForte456"
  }'
```

## Exemplos com JavaScript/Fetch

### Registrar Usuário

```javascript
async function registerUser(email, password, name) {
  try {
    const response = await fetch('http://localhost:3000/api/users/register', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ email, password, name }),
    });

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.message);
    }

    const data = await response.json();
    console.log('Usuário registrado:', data);
    return data;
  } catch (error) {
    console.error('Erro ao registrar usuário:', error);
    throw error;
  }
}

// Uso
registerUser('pedro@example.com', 'Senha123', 'Pedro Costa');
```

### Login

```javascript
async function login(email, password) {
  try {
    const response = await fetch('http://localhost:3000/api/users/login', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ email, password }),
    });

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.message);
    }

    const data = await response.json();
    
    // Salvar token no localStorage
    localStorage.setItem('token', data.token);
    localStorage.setItem('user', JSON.stringify(data.user));
    
    console.log('Login realizado com sucesso');
    return data;
  } catch (error) {
    console.error('Erro ao fazer login:', error);
    throw error;
  }
}

// Uso
login('pedro@example.com', 'Senha123');
```

### Obter Perfil

```javascript
async function getProfile() {
  try {
    const token = localStorage.getItem('token');
    
    if (!token) {
      throw new Error('Token não encontrado. Faça login primeiro.');
    }

    const response = await fetch('http://localhost:3000/api/users/me', {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${token}`,
      },
    });

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.message);
    }

    const user = await response.json();
    console.log('Perfil do usuário:', user);
    return user;
  } catch (error) {
    console.error('Erro ao obter perfil:', error);
    throw error;
  }
}

// Uso
getProfile();
```

### Atualizar Nome

```javascript
async function updateUserName(userId, newName) {
  try {
    const token = localStorage.getItem('token');
    
    const response = await fetch(`http://localhost:3000/api/users/${userId}`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}`,
      },
      body: JSON.stringify({ name: newName }),
    });

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.message);
    }

    const data = await response.json();
    console.log('Nome atualizado:', data);
    return data;
  } catch (error) {
    console.error('Erro ao atualizar nome:', error);
    throw error;
  }
}

// Uso
updateUserName('550e8400-e29b-41d4-a716-446655440000', 'Novo Nome');
```

### Cliente API Completo

```javascript
class UserMesAPI {
  constructor(baseURL = 'http://localhost:3000') {
    this.baseURL = baseURL;
    this.token = localStorage.getItem('token');
  }

  setToken(token) {
    this.token = token;
    localStorage.setItem('token', token);
  }

  clearToken() {
    this.token = null;
    localStorage.removeItem('token');
  }

  async request(endpoint, options = {}) {
    const url = `${this.baseURL}${endpoint}`;
    const headers = {
      'Content-Type': 'application/json',
      ...options.headers,
    };

    if (this.token && !options.skipAuth) {
      headers['Authorization'] = `Bearer ${this.token}`;
    }

    const config = {
      ...options,
      headers,
    };

    try {
      const response = await fetch(url, config);
      const data = await response.json();

      if (!response.ok) {
        throw new Error(data.message || 'Erro na requisição');
      }

      return data;
    } catch (error) {
      console.error('API Error:', error);
      throw error;
    }
  }

  async register(email, password, name) {
    const data = await this.request('/api/users/register', {
      method: 'POST',
      body: JSON.stringify({ email, password, name }),
      skipAuth: true,
    });
    return data;
  }

  async login(email, password) {
    const data = await this.request('/api/users/login', {
      method: 'POST',
      body: JSON.stringify({ email, password }),
      skipAuth: true,
    });
    
    this.setToken(data.token);
    return data;
  }

  async logout() {
    this.clearToken();
  }

  async getProfile() {
    return await this.request('/api/users/me');
  }

  async getUserById(userId) {
    return await this.request(`/api/users/${userId}`);
  }

  async updateUser(userId, name) {
    return await this.request(`/api/users/${userId}`, {
      method: 'PUT',
      body: JSON.stringify({ name }),
    });
  }

  async changePassword(userId, oldPassword, newPassword) {
    return await this.request(`/api/users/${userId}/change-password`, {
      method: 'POST',
      body: JSON.stringify({ old_password: oldPassword, new_password: newPassword }),
    });
  }

  async deactivateUser(userId) {
    return await this.request(`/api/users/${userId}/deactivate`, {
      method: 'POST',
    });
  }

  async activateUser(userId) {
    return await this.request(`/api/users/${userId}/activate`, {
      method: 'POST',
    });
  }

  async healthCheck() {
    return await this.request('/health', { skipAuth: true });
  }
}

// Uso
const api = new UserMesAPI();

// Registrar
await api.register('ana@example.com', 'Senha123', 'Ana Silva');

// Login
const loginData = await api.login('ana@example.com', 'Senha123');

// Obter perfil
const profile = await api.getProfile();

// Atualizar nome
await api.updateUser(profile.id, 'Ana Paula Silva');

// Logout
api.logout();
```

## Exemplos com Go

### Cliente API

```go
package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
)

const baseURL = "http://localhost:3000"

type Client struct {
    httpClient *http.Client
    token      string
}

func NewClient() *Client {
    return &Client{
        httpClient: &http.Client{},
    }
}

func (c *Client) SetToken(token string) {
    c.token = token
}

func (c *Client) request(method, endpoint string, body interface{}, needsAuth bool) ([]byte, error) {
    var reqBody io.Reader
    if body != nil {
        jsonData, err := json.Marshal(body)
        if err != nil {
            return nil, err
        }
        reqBody = bytes.NewBuffer(jsonData)
    }

    req, err := http.NewRequest(method, baseURL+endpoint, reqBody)
    if err != nil {
        return nil, err
    }

    req.Header.Set("Content-Type", "application/json")
    if needsAuth && c.token != "" {
        req.Header.Set("Authorization", "Bearer "+c.token)
    }

    resp, err := c.httpClient.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    respBody, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }

    if resp.StatusCode >= 400 {
        return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(respBody))
    }

    return respBody, nil
}

type RegisterRequest struct {
    Email    string `json:"email"`
    Password string `json:"password"`
    Name     string `json:"name"`
}

type LoginRequest struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}

type LoginResponse struct {
    Token string      `json:"token"`
    User  UserResponse `json:"user"`
}

type UserResponse struct {
    ID        string `json:"id"`
    Email     string `json:"email"`
    Name      string `json:"name"`
    IsActive  bool   `json:"is_active"`
    CreatedAt string `json:"created_at"`
    UpdatedAt string `json:"updated_at"`
}

func (c *Client) Register(email, password, name string) (*UserResponse, error) {
    req := RegisterRequest{
        Email:    email,
        Password: password,
        Name:     name,
    }

    respBody, err := c.request("POST", "/api/users/register", req, false)
    if err != nil {
        return nil, err
    }

    var response struct {
        Message string       `json:"message"`
        Data    UserResponse `json:"data"`
    }

    if err := json.Unmarshal(respBody, &response); err != nil {
        return nil, err
    }

    return &response.Data, nil
}

func (c *Client) Login(email, password string) (*LoginResponse, error) {
    req := LoginRequest{
        Email:    email,
        Password: password,
    }

    respBody, err := c.request("POST", "/api/users/login", req, false)
    if err != nil {
        return nil, err
    }

    var response LoginResponse
    if err := json.Unmarshal(respBody, &response); err != nil {
        return nil, err
    }

    c.SetToken(response.Token)
    return &response, nil
}

func (c *Client) GetProfile() (*UserResponse, error) {
    respBody, err := c.request("GET", "/api/users/me", nil, true)
    if err != nil {
        return nil, err
    }

    var user UserResponse
    if err := json.Unmarshal(respBody, &user); err != nil {
        return nil, err
    }

    return &user, nil
}

func main() {
    client := NewClient()

    // Registrar
    user, err := client.Register("carlos@example.com", "Senha123", "Carlos Oliveira")
    if err != nil {
        fmt.Printf("Erro ao registrar: %v\n", err)
        return
    }
    fmt.Printf("Usuário registrado: %+v\n", user)

    // Login
    loginResp, err := client.Login("carlos@example.com", "Senha123")
    if err != nil {
        fmt.Printf("Erro ao fazer login: %v\n", err)
        return
    }
    fmt.Printf("Login realizado. Token: %s\n", loginResp.Token)

    // Obter perfil
    profile, err := client.GetProfile()
    if err != nil {
        fmt.Printf("Erro ao obter perfil: %v\n", err)
        return
    }
    fmt.Printf("Perfil: %+v\n", profile)
}
```

## Códigos de Status

| Código | Significado | Descrição |
|--------|-------------|-----------|
| 200 | OK | Requisição bem-sucedida |
| 201 | Created | Recurso criado com sucesso |
| 400 | Bad Request | Requisição inválida ou malformada |
| 401 | Unauthorized | Não autenticado ou token inválido |
| 403 | Forbidden | Sem permissão para acessar o recurso |
| 404 | Not Found | Recurso não encontrado |
| 409 | Conflict | Conflito (ex: email já existe) |
| 500 | Internal Server Error | Erro interno do servidor |

## Tratamento de Erros

Todas as respostas de erro seguem o mesmo formato:

```json
{
  "error": "error_code",
  "message": "Descrição legível do erro"
}
```

### Exemplos de Erros

#### Email já existe (409)
```json
{
  "error": "conflict",
  "message": "email already exists"
}
```

#### Credenciais inválidas (401)
```json
{
  "error": "unauthorized",
  "message": "Invalid email or password"
}
```

#### Token inválido (401)
```json
{
  "error": "unauthorized",
  "message": "Invalid or expired token"
}
```

#### Senha muito curta (400)
```json
{
  "error": "invalid_password",
  "message": "password must be at least 8 characters long"
}
```

#### Usuário não encontrado (404)
```json
{
  "error": "not_found",
  "message": "user not found"
}
```

#### Validação de campo (400)
```json
{
  "error": "validation_error",
  "message": "email is required"
}
```

## Dicas de Segurança

1. **Nunca exponha o token JWT** em logs ou console em produção
2. **Use HTTPS** em produção para proteger o token em trânsito
3. **Armazene o token de forma segura** (httpOnly cookies ou memória, não localStorage em produção)
4. **Implemente refresh tokens** para sessões de longa duração
5. **Valide todos os inputs** no lado do cliente antes de enviar
6. **Implemente rate limiting** para prevenir ataques de força bruta
7. **Use senhas fortes** com pelo menos 8 caracteres, letras e números
8. **Implemente logout adequado** limpando tokens e sessões

## Fluxo Completo de Autenticação

```javascript
// 1. Usuário se registra
const registerData = await api.register('user@example.com', 'Pass123', 'User Name');

// 2. Usuário faz login
const loginData = await api.login('user@example.com', 'Pass123');
// Token é salvo automaticamente

// 3. Fazer requisições autenticadas
const profile = await api.getProfile();
const userData = await api.getUserById(profile.id);
await api.updateUser(profile.id, 'New Name');

// 4. Logout
api.logout();
// Token é removido
```

## Testes com Postman

1. Importe a collection no Postman
2. Configure a variável `baseUrl` como `http://localhost:3000`
3. Após o login, salve o token retornado em uma variável de ambiente
4. Use `{{token}}` nos headers de autorização das outras requisições

## Troubleshooting

### Erro "Authorization header is required"
- Verifique se está incluindo o header `Authorization: Bearer <token>`
- Confirme que o token está sendo enviado corretamente

### Erro "Invalid or expired token"
- Faça login novamente para obter um novo token
- Tokens expiram após 24 horas

### Erro "email already exists"
- Use um email diferente para registro
- Ou faça login com o email existente

### Erro CORS
- Certifique-se de que o servidor está rodando
- Verifique se está fazendo requisições para a URL correta
- Em desenvolvimento, o CORS está configurado para aceitar todas as origens