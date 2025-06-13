# 📋 API Reference - Gestar Bem

Documentação completa dos endpoints da API REST da rede social Gestar Bem.

## 🔐 Autenticação

A API utiliza **JWT (JSON Web Tokens)** para autenticação. 

```mermaid
sequenceDiagram
    participant C as Cliente
    participant A as API
    participant D as Database
    
    C->>A: POST /login (email, senha)
    A->>D: Verificar credenciais
    D-->>A: Usuário válido
    A->>A: Gerar JWT Token
    A-->>C: Token JWT
    
    Note over C: Armazenar token
    
    C->>A: GET /usuarios (Authorization: Bearer token)
    A->>A: Validar JWT
    A->>D: Buscar dados
    D-->>A: Dados do usuário
    A-->>C: Resposta com dados
```

### Como obter um token:
```bash
POST /login
Content-Type: application/json

{
  "email": "usuario@exemplo.com",
  "senha": "minhasenha"
}
```

### Como usar o token:
```bash
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

## 👥 Usuários

### Criar Usuário
```http
POST /usuarios
Content-Type: application/json
```

**Body:**
```json
{
  "nome": "João Silva",
  "nick": "joaosilva",
  "email": "joao@exemplo.com",
  "senha": "minhasenha123"
}
```

**Resposta (201):**
```json
{
  "id": 1,
  "nome": "João Silva",
  "nick": "joaosilva",
  "email": "joao@exemplo.com",
  "criadoem": "2025-06-11T22:30:00Z"
}
```

### Fazer Login
```http
POST /login
Content-Type: application/json
```

**Body:**
```json
{
  "email": "joao@exemplo.com",
  "senha": "minhasenha123"
}
```

**Resposta (200):**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### Listar Usuários
```http
GET /usuarios
Authorization: Bearer {token}
```

**Query Parameters:**
- `usuario` (opcional) - Filtrar por nome ou nick

**Resposta (200):**
```json
[
  {
    "id": 1,
    "nome": "João Silva",
    "nick": "joaosilva",
    "email": "joao@exemplo.com",
    "foto_perfil": "perfil_1_1749691234.jpg",
    "criadoem": "2025-06-11T22:30:00Z"
  }
]
```

### Buscar Usuário por ID
```http
GET /usuarios/{id}
Authorization: Bearer {token}
```

**Resposta (200):**
```json
{
  "id": 1,
  "nome": "João Silva",
  "nick": "joaosilva",
  "email": "joao@exemplo.com",
  "foto_perfil": "perfil_1_1749691234.jpg",
  "criadoem": "2025-06-11T22:30:00Z"
}
```

### Atualizar Usuário
```http
PUT /usuarios/{id}
Authorization: Bearer {token}
Content-Type: application/json
```

**Body:**
```json
{
  "nome": "João Silva Santos",
  "nick": "joaosilva",
  "email": "joao.santos@exemplo.com"
}
```

### Deletar Usuário
```http
DELETE /usuarios/{id}
Authorization: Bearer {token}
```

### Upload de Foto de Perfil
```http
POST /usuarios/{id}/foto-perfil
Authorization: Bearer {token}
Content-Type: multipart/form-data
```

**Form Data:**
- `foto` - Arquivo de imagem (JPG, PNG, GIF)

**Validações:**
- Tamanho máximo: 10MB
- Tipos aceitos: image/jpeg, image/png, image/gif
- Apenas o próprio usuário pode alterar sua foto

**Resposta (200):**
```json
{
  "mensagem": "Foto de perfil atualizada com sucesso",
  "foto_perfil": "perfil_1_1749691234.jpg"
}
```

### Seguir Usuário
```http
POST /usuarios/{id}/seguir
Authorization: Bearer {token}
```

### Parar de Seguir
```http
DELETE /usuarios/{id}/parar-de-seguir
Authorization: Bearer {token}
```

### Listar Seguidores
```http
GET /usuarios/{id}/seguidores
Authorization: Bearer {token}
```

### Listar Seguindo
```http
GET /usuarios/{id}/seguindo
Authorization: Bearer {token}
```

### Estatísticas do Usuário
```http
GET /usuarios/{id}/estatisticas
Authorization: Bearer {token}
```

**Resposta (200):**
```json
{
  "seguidores": 150,
  "seguindo": 89,
  "publicacoes": 42
}
```

## 📝 Publicações

### Criar Publicação
```http
POST /publicacoes
Authorization: Bearer {token}
Content-Type: application/json
```

**Body:**
```json
{
  "titulo": "Minha primeira publicação",
  "conteudo": "Este é o conteúdo da minha publicação..."
}
```

### Listar Publicações
```http
GET /publicacoes
Authorization: Bearer {token}
```

### Buscar Publicação
```http
GET /publicacoes/{id}
Authorization: Bearer {token}
```

### Atualizar Publicação
```http
PUT /publicacoes/{id}
Authorization: Bearer {token}
```

### Deletar Publicação
```http
DELETE /publicacoes/{id}
Authorization: Bearer {token}
```

### Listar Publicações do Usuário
```http
GET /usuarios/{id}/publicacoes
Authorization: Bearer {token}
```

### Curtir Publicação
```http
POST /publicacoes/{id}/curtir
Authorization: Bearer {token}
```

### Descurtir Publicação
```http
DELETE /publicacoes/{id}/descurtir
Authorization: Bearer {token}
```

## 💬 Comentários

### Criar Comentário
```http
POST /publicacoes/{id}/comentarios
Authorization: Bearer {token}
Content-Type: application/json
```

**Body:**
```json
{
  "conteudo": "Ótima publicação!"
}
```

### Listar Comentários
```http
GET /publicacoes/{id}/comentarios
Authorization: Bearer {token}
```

### Atualizar Comentário
```http
PUT /comentarios/{id}
Authorization: Bearer {token}
```

### Deletar Comentário
```http
DELETE /comentarios/{id}
Authorization: Bearer {token}
```

### Curtir Comentário
```http
POST /comentarios/{id}/curtir
Authorization: Bearer {token}
```

### Descurtir Comentário
```http
DELETE /comentarios/{id}/descurtir
Authorization: Bearer {token}
```

## 📁 Arquivos

### Servir Foto de Perfil
```http
GET /uploads/perfil/{arquivo}
```

**Exemplo:**
```
GET /uploads/perfil/perfil_1_1749691234.jpg
```

## ❌ Códigos de Erro

| Código | Descrição |
|--------|-----------|
| `400` | Bad Request - Dados inválidos |
| `401` | Unauthorized - Token inválido/ausente |
| `403` | Forbidden - Sem permissão |
| `404` | Not Found - Recurso não encontrado |
| `409` | Conflict - Conflito (ex: email já existe) |
| `422` | Unprocessable Entity - Validação falhou |
| `500` | Internal Server Error - Erro interno |

### Formato de Erro
```json
{
  "erro": "Descrição do erro"
}
```

## 🔧 Headers Importantes

### Request Headers
```http
Content-Type: application/json
Authorization: Bearer {token}
```

### Response Headers
```http
Content-Type: application/json
Access-Control-Allow-Origin: *
Access-Control-Allow-Methods: GET, POST, PUT, DELETE, OPTIONS
Access-Control-Allow-Headers: Content-Type, Authorization
```

## 📊 Rate Limiting

- **Limite:** 1000 requests por hora por IP
- **Header de resposta:** `X-RateLimit-Remaining`

## 🧪 Exemplos com cURL

### Criar usuário e fazer upload de foto:
```bash
# 1. Criar usuário
curl -X POST http://localhost:5000/usuarios \
  -H "Content-Type: application/json" \
  -d '{"nome":"João","nick":"joao","email":"joao@teste.com","senha":"123456"}'

# 2. Fazer login
TOKEN=$(curl -s -X POST http://localhost:5000/login \
  -H "Content-Type: application/json" \
  -d '{"email":"joao@teste.com","senha":"123456"}' | jq -r '.token')

# 3. Upload de foto
curl -X POST http://localhost:5000/usuarios/1/foto-perfil \
  -H "Authorization: Bearer $TOKEN" \
  -F "foto=@minha_foto.jpg"
```

## 🔍 Filtros e Paginação

### Usuários
```http
GET /usuarios?usuario=joao&page=1&limit=10
```

### Publicações
```http
GET /publicacoes?page=1&limit=20&order=desc
```

---

**Base URL:** `http://localhost:5000`  
**Versão:** v1.0  
**Última atualização:** 11/06/2025 