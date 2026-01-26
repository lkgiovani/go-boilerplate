# 🚀 Boilerplate Go 2026

[![Go](https://img.shields.io/badge/Go-1.25-00ADD8.svg)](https://golang.org/)
[![Fiber](https://img.shields.io/badge/Fiber-v2-00ACD7.svg)](https://gofiber.io/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

Um boilerplate moderno e completo para aplicações **Go** com autenticação JWT, OAuth2, Rate Limiting, monitoramento de segurança e arquitetura limpa.

---

## 📋 Índice

- [✨ Features](#-features)
- [🛠️ Stack Tecnológica](#️-stack-tecnológica)
- [📁 Arquitetura](#-arquitetura)
- [🚀 Quick Start](#-quick-start)
- [⚙️ Configuração](#️-configuração)
- [🔐 Autenticação](#-autenticação)
- [🛡️ Segurança](#️-segurança)
- [📊 API Endpoints](#-api-endpoints)
- [🧪 Testes](#-testes)
- [🌍 Internacionalização](#-internacionalização)
- [📚 Documentação API](#-documentação-api)

---

## ✨ Features

### 🔐 Autenticação & Autorização

- **JWT (JSON Web Tokens)** com access tokens de curta duração
- **Refresh Token Rotation** com detecção de reutilização (família de tokens)
- **OAuth2** com Google (Web e Mobile)
- **Múltiplos dispositivos** - controle de sessões por device
- **Logout** individual ou de todos os dispositivos

### 👤 Gestão de Usuários

- CRUD completo de usuários
- Controle de roles (USER/ADMIN)
- Ativação/desativação de contas
- Atualização de senha
- Busca paginada com filtros
- Reset de senha via email

### 🛡️ Segurança Avançada

- **Rate Limiting** por IP com middleware customizado
  - Limites por método HTTP (GET, POST, PUT, DELETE)
  - Limites específicos por endpoint
  - Whitelist de IPs
- **Detecção de Atividades Suspeitas**
  - Detecção de bots via User-Agent
  - Monitoramento de taxa de requisições anormal
  - Detecção de criação em massa
- **Auto-bloqueio de Usuários**
  - Bloqueio automático baseado em severidade
  - Bloqueio temporário ou permanente
  - Desbloqueio por admin

### 🌍 Internacionalização (i18n)

- Suporte a múltiplos idiomas
- Português (Brasil) e Inglês incluídos
- Mensagens de erro traduzidas
- Detecção automática via `Accept-Language`

### 📊 Observabilidade

- **Health Check** endpoints
- **Metrics** com Prometheus (opcional)
- **Structured Logging** com Zap
- **Request ID** tracking

---

## 🛠️ Stack Tecnológica

### Core

| Tecnologia  | Versão | Descrição                         |
| ----------- | ------ | --------------------------------- |
| **Go**      | 1.25+  | Linguagem principal               |
| **Fiber**   | v2     | Framework web de alta performance |
| **Uber Fx** | 1.24+  | Dependency injection              |
| **Air**     | -      | Hot reload para desenvolvimento   |

### Segurança & Auth

| Tecnologia              | Versão | Descrição                  |
| ----------------------- | ------ | -------------------------- |
| **golang-jwt/jwt**      | v3     | Tokens JWT                 |
| **golang.org/x/crypto** | latest | Bcrypt para hash de senhas |
| **Google OAuth2**       | -      | Login social               |
| **Rate Limiter**        | custom | Rate limiting middleware   |

### Persistência

| Tecnologia         | Versão | Descrição                |
| ------------------ | ------ | ------------------------ |
| **PostgreSQL**     | 16+    | Banco de dados principal |
| **GORM**           | v1.31+ | ORM                      |
| **golang-migrate** | -      | Migrações de banco       |

### Ferramentas

| Tecnologia   | Versão | Descrição                              |
| ------------ | ------ | -------------------------------------- |
| **Zap**      | 1.27+  | Structured logging                     |
| **godotenv** | 1.5+   | Gerenciamento de variáveis de ambiente |
| **UUID**     | 1.6+   | Geração de UUIDs                       |
| **Swagger**  | -      | Documentação API (opcional)            |

### Testes

| Tecnologia     | Versão | Descrição                  |
| -------------- | ------ | -------------------------- |
| **testing**    | stdlib | Framework de testes nativo |
| **testify**    | 1.11+  | Assertions e mocks         |
| **dockertest** | -      | Containers para testes     |

---

## 📁 Arquitetura

O projeto segue os princípios de **Clean Architecture** com camadas bem definidas:

```
go-boilerplate/
├── 📂 cmd/                     # Aplicação principal
│   └── 📂 startServer/         # Entry point
├── 📂 internal/                # Código privado da aplicação
│   ├── 📂 contracts/           # Interfaces de domínio
│   ├── 📂 delivery/            # Handlers HTTP (controllers)
│   ├── 📂 domain/              # Lógica de negócio
│   │   ├── 📂 auth/            # Serviços de autenticação
│   │   └── 📂 user/            # Serviços de usuário
│   ├── 📂 entities/            # Modelos de dados (GORM)
│   ├── 📂 dto/                 # Data Transfer Objects
│   │   ├── 📂 request/         # DTOs de requisição
│   │   └── 📂 response/        # DTOs de resposta
│   └── 📂 errors/              # Erros customizados
├── 📂 infra/                   # Infraestrutura
│   ├── 📂 config/              # Configurações
│   ├── 📂 database/            # Conexão com banco
│   ├── 📂 fx/                  # Dependency injection (Uber Fx)
│   ├── 📂 middleware/          # Middlewares
│   └── 📂 security/            # Segurança
├── 📂 pkg/                     # Pacotes reutilizáveis
│   ├── 📂 encrypt/             # Bcrypt helpers
│   ├── 📂 jwt/                 # JWT utilities
│   ├── 📂 logger/              # Logger (Zap)
│   ├── 📂 utils/               # Utilitários gerais
│   └── 📂 uuid/                # UUID helpers
├── 📂 resources/               # Recursos estáticos
│   ├── 📂 db/
│   │   └── 📂 migrations/      # Migrações SQL
│   └── 📂 i18n/                # Arquivos de tradução
├── 📂 tests/                   # Testes
│   ├── 📂 integration/         # Testes de integração
│   ├── 📂 unit/                # Testes unitários
│   └── 📂 helpers/             # Helpers de teste
├── .env                        # Variáveis de ambiente
├── .air.toml                   # Configuração do Air
├── docker-compose.yml          # Docker Compose
├── Dockerfile                  # Dockerfile
└── go.mod                      # Dependências Go
```

### Fluxo de Dados

```
HTTP Request → Middleware → Handler (Delivery) → Service (Domain) → Repository → Database
      ↑                                                                              ↓
      └──────────────────────────────────────────────────────────────────────────────┘
```

---

## 🚀 Quick Start

### Pré-requisitos

- **Go 1.25+** ([Download](https://golang.org/dl/))
- **Docker** (para PostgreSQL local)
- **Air** (opcional, para hot reload)

### 1️⃣ Clone o repositório

```bash
git clone https://github.com/lkgiovani/go-boilerplate.git
cd go-boilerplate
```

### 2️⃣ Instale as dependências

```bash
go mod download
```

### 3️⃣ Configure as variáveis de ambiente

```bash
cp .env.example .env
# Edite o .env com suas configurações
```

### 4️⃣ Inicie o PostgreSQL

```bash
docker-compose up -d postgres
```

### 5️⃣ Execute as migrações

```bash
# Usando golang-migrate
migrate -path resources/db/migrations -database "postgresql://postgres:postgres@localhost:5432/boilerplate?sslmode=disable" up
```

### 6️⃣ Execute a aplicação

**Com Air (hot reload):**

```bash
air
```

**Sem Air:**

```bash
go run cmd/startServer/main.go
```

### 7️⃣ Acesse a aplicação

- **API**: http://localhost:8080
- **Health Check**: http://localhost:8080/health

---

## ⚙️ Configuração

### Variáveis de Ambiente

| Variável                        | Descrição                           | Padrão  |
| ------------------------------- | ----------------------------------- | ------- |
| `SERVER_PORT`                   | Porta do servidor                   | `8080`  |
| `JWT_SECRET`                    | Chave secreta JWT (mín. 32 chars)   | -       |
| `JWT_EXPIRATION_MINUTES`        | Expiração do access token (minutos) | `15`    |
| `REFRESH_TOKEN_EXPIRATION_DAYS` | Expiração do refresh token (dias)   | `30`    |
| `COOKIE_DOMAIN`                 | Domínio dos cookies                 | -       |
| `COOKIE_SECURE`                 | Cookies apenas HTTPS                | `false` |
| `GOOGLE_CLIENT_ID`              | OAuth2 Google Client ID             | -       |
| `GOOGLE_CLIENT_SECRET`          | OAuth2 Google Client Secret         | -       |
| `GOOGLE_ANDROID_CLIENT_ID`      | OAuth2 Google Android Client ID     | -       |
| `GOOGLE_IOS_CLIENT_ID`          | OAuth2 Google iOS Client ID         | -       |

### Database

| Variável      | Descrição           | Padrão        |
| ------------- | ------------------- | ------------- |
| `DB_HOST`     | Host do PostgreSQL  | `localhost`   |
| `DB_PORT`     | Porta do PostgreSQL | `5432`        |
| `DB_USER`     | Usuário do banco    | `postgres`    |
| `DB_PASSWORD` | Senha do banco      | `postgres`    |
| `DB_NAME`     | Nome do banco       | `boilerplate` |
| `DB_SSL_MODE` | SSL mode            | `disable`     |

### Rate Limiting

| Variável             | Descrição            | Padrão |
| -------------------- | -------------------- | ------ |
| `RATE_LIMIT_ENABLED` | Ativar rate limiting | `true` |
| `RATE_LIMIT_GLOBAL`  | Limite global/minuto | `300`  |
| `RATE_LIMIT_GET`     | Limite GET/minuto    | `600`  |
| `RATE_LIMIT_POST`    | Limite POST/minuto   | `60`   |
| `RATE_LIMIT_LOGIN`   | Limite login/minuto  | `5`    |
| `RATE_LIMIT_SIGNUP`  | Limite signup/minuto | `3`    |

### Segurança (Auto-bloqueio)

| Variável                             | Descrição                       | Padrão |
| ------------------------------------ | ------------------------------- | ------ |
| `SECURITY_AUTO_BLOCK_CRITICAL`       | Atividades críticas p/ bloqueio | `3`    |
| `SECURITY_AUTO_BLOCK_HIGH`           | Atividades alta severidade      | `10`   |
| `SECURITY_AUTO_BLOCK_TOTAL`          | Total de atividades             | `20`   |
| `SECURITY_AUTO_BLOCK_WINDOW_HOURS`   | Janela de análise (horas)       | `24`   |
| `SECURITY_AUTO_BLOCK_DURATION_HOURS` | Duração do bloqueio (horas)     | `72`   |

### Ambientes

| Ambiente      | Descrição                                |
| ------------- | ---------------------------------------- |
| `development` | Desenvolvimento local (logs detalhados)  |
| `staging`     | Staging (similar produção)               |
| `production`  | Produção (otimizações, segurança máxima) |

---

## 🔐 Autenticação

### Fluxo de Autenticação (Web)

```mermaid
sequenceDiagram
    participant C as Client
    participant S as Server
    participant DB as Database

    C->>S: POST /v1/auth/login (email, password)
    S->>DB: Validar credenciais
    DB-->>S: User data
    S-->>C: Access Token (body) + Refresh Token (HttpOnly Cookie)

    Note over C,S: Access Token expira em 15 min

    C->>S: POST /v1/auth/refresh (com Cookie)
    S->>DB: Validar e rotacionar refresh token
    DB-->>S: Novo token
    S-->>C: Novo Access Token + Novo Refresh Token (Cookie)
```

### Endpoints de Autenticação

#### Login

```http
POST /v1/auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "securePassword123"
}
```

**Resposta:**

```json
{
  "userId": "550e8400-e29b-41d4-a716-446655440000",
  "email": "user@example.com",
  "accessToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "expiresIn": 900
}
```

#### Signup

```http
POST /v1/auth/signup
Content-Type: application/json

{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "securePassword123"
}
```

#### Refresh Token

```http
POST /v1/auth/refresh
Cookie: refreshToken=...
```

#### Logout

```http
POST /v1/auth/logout
Cookie: refreshToken=...
Authorization: Bearer <access_token>
```

#### Logout de Todos os Dispositivos

```http
POST /v1/auth/logout-all
Authorization: Bearer <access_token>
```

### OAuth2 (Google)

#### Web (Browser)

Redirecione para:

```
GET /oauth2/authorization/google
```

#### Mobile

```http
POST /v1/auth/mobile/oauth2/google
Content-Type: application/json

{
  "idToken": "eyJhbGciOiJSUzI1NiIsInR5cC...",
  "deviceId": "device-uuid-123"
}
```

---

## 🛡️ Segurança

### Rate Limiting

O sistema implementa **rate limiting multinível**:

1. **Global**: 300 req/min por IP
2. **Por Método HTTP**:
   - GET: 600/min
   - POST: 60/min
   - PUT/PATCH: 120/min
   - DELETE: 30/min
3. **Por Endpoint Crítico**:
   - Login: 5/min
   - Signup: 3/min
   - Forgot Password: 3/min

**Headers de Resposta:**

```
X-RateLimit-Limit: 5
X-RateLimit-Remaining: 3
X-RateLimit-Reset: 1640000000
Retry-After: 60
```

### Detecção de Atividades Suspeitas

| Tipo de Atividade     | Severidade | Descrição                           |
| --------------------- | ---------- | ----------------------------------- |
| `RATE_LIMIT_EXCEEDED` | HIGH       | Excedeu rate limit                  |
| `MASS_CREATION`       | HIGH       | Criação em massa de recursos        |
| `PATTERN_ABUSE`       | CRITICAL   | Padrão de abuso detectado           |
| `AUTOMATED_BEHAVIOR`  | HIGH       | Comportamento automatizado (bot)    |
| `UNAUTHORIZED_ACCESS` | MEDIUM     | Tentativas de acesso não autorizado |

### Auto-bloqueio

Usuários são automaticamente bloqueados quando:

- 3+ atividades **CRITICAL** em 24h
- 10+ atividades **HIGH** em 24h
- 20+ atividades totais em 24h

**Duração do bloqueio**: 72 horas (configurável)

---

## 📊 API Endpoints

### Auth Controller (`/v1/auth`)

| Método | Endpoint           | Descrição                       | Auth |
| ------ | ------------------ | ------------------------------- | ---- |
| POST   | `/login`           | Login com email/senha           | ❌   |
| POST   | `/signup`          | Criar nova conta                | ❌   |
| POST   | `/refresh`         | Renovar access token            | ❌   |
| POST   | `/logout`          | Logout do dispositivo atual     | ✅   |
| POST   | `/logout-all`      | Logout de todos os dispositivos | ✅   |
| POST   | `/forgot-password` | Solicitar reset de senha        | ❌   |
| POST   | `/reset-password`  | Resetar senha com token         | ❌   |

### Mobile Auth (`/v1/auth/mobile`)

| Método | Endpoint         | Descrição                 | Auth |
| ------ | ---------------- | ------------------------- | ---- |
| POST   | `/oauth2/google` | Login com Google (mobile) | ❌   |
| POST   | `/refresh`       | Renovar token (mobile)    | ❌   |

### User Controller (`/v1/users`)

| Método | Endpoint        | Descrição                  | Auth | Role  |
| ------ | --------------- | -------------------------- | ---- | ----- |
| GET    | `/me`           | Dados do usuário atual     | ✅   | USER  |
| PUT    | `/`             | Atualizar perfil           | ✅   | USER  |
| PATCH  | `/password`     | Alterar senha              | ✅   | USER  |
| GET    | `/`             | Listar usuários (paginado) | ✅   | ADMIN |
| GET    | `/:id`          | Buscar por ID              | ✅   | ADMIN |
| GET    | `/email/:email` | Buscar por email           | ✅   | ADMIN |
| POST   | `/`             | Criar usuário              | ✅   | ADMIN |
| DELETE | `/:id`          | Deletar usuário            | ✅   | ADMIN |
| DELETE | `/bulk`         | Deletar múltiplos          | ✅   | ADMIN |
| PATCH  | `/:id/status`   | Ativar/desativar           | ✅   | ADMIN |

### Health & Monitoring

| Método | Endpoint   | Descrição           | Auth |
| ------ | ---------- | ------------------- | ---- |
| GET    | `/health`  | Health check        | ❌   |
| GET    | `/metrics` | Métricas Prometheus | ❌   |

---

## 🧪 Testes

### Executar Todos os Testes

```bash
go test ./... -v
```

### Cobertura de Testes

```bash
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### Testes de Integração

Os testes usam **dockertest** para subir um PostgreSQL real:

```go
func TestAuthService(t *testing.T) {
    // Setup PostgreSQL container
    pool, resource := setupTestDB(t)
    defer pool.Purge(resource)

    // Testes reais com banco PostgreSQL
}
```

### Estrutura de Testes

```
tests/
├── integration/
│   ├── auth_test.go
│   └── user_test.go
├── unit/
│   ├── service_test.go
│   └── repository_test.go
└── helpers/
    ├── testdb.go
    └── fixtures.go
```

---

## 🌍 Internacionalização

### Idiomas Suportados

- 🇺🇸 **Inglês** (padrão)
- 🇧🇷 **Português (Brasil)**

### Como Usar

Envie o header `Accept-Language`:

```http
GET /v1/users/me
Accept-Language: pt-BR
Authorization: Bearer <token>
```

**Resposta de erro em português:**

```json
{
  "error": "Credenciais inválidas",
  "message": "Email ou senha incorretos",
  "status": 401
}
```

### Adicionar Novo Idioma

1. Crie `resources/i18n/XX.json`
2. Traduza as chaves de `en.json`

---

## 📚 Documentação API

### Swagger UI (Opcional)

Se configurado, acesse:

```
http://localhost:8080/swagger/index.html
```

### Health Check

```
http://localhost:8080/health
```

**Resposta:**

```json
{
  "status": "UP",
  "database": "UP",
  "timestamp": "2026-01-26T12:00:00Z"
}
```

---

## 📝 Migrações de Banco

O projeto usa **golang-migrate** para versionamento do schema:

| Versão | Descrição                                                |
| ------ | -------------------------------------------------------- |
| V1     | Tabela `users`                                           |
| V2     | Tabela `refresh_tokens`                                  |
| V3     | Tabelas `suspicious_activities` e `user_security_blocks` |
| V4     | Tabela `oauth2_accounts`                                 |
| V5     | Tabela `password_reset_tokens`                           |

### Estrutura do Banco

```sql
-- users: Usuários do sistema
-- refresh_tokens: Tokens de refresh (família, rotação)
-- suspicious_activities: Registro de atividades suspeitas
-- user_security_blocks: Bloqueios de segurança
-- oauth2_accounts: Contas OAuth2 vinculadas
-- password_reset_tokens: Tokens de reset de senha
```

## 🤝 Contribuição

1. Fork o projeto
2. Crie sua branch (`git checkout -b feature/AmazingFeature`)
3. Commit suas mudanças (`git commit -m 'Add AmazingFeature'`)
4. Push para a branch (`git push origin feature/AmazingFeature`)
5. Abra um Pull Request

---

## 📄 Licença

Este projeto está sob a licença MIT. Veja o arquivo [LICENSE](LICENSE) para mais detalhes.

---

## 🙏 Agradecimentos

- [Fiber](https://gofiber.io/)
- [GORM](https://gorm.io/)
- [Uber Fx](https://uber-go.github.io/fx/)
- [Zap](https://github.com/uber-go/zap)

---

<p align="center">
  Feito com ❤️ para a comunidade Go
</p>
