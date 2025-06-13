# Frontend - Gestar Bem

Este é o frontend da aplicação **Gestar Bem**, uma interface web moderna e responsiva que consome a API REST desenvolvida em Go.

## 🚀 Funcionalidades

### Autenticação
- ✅ **Login** - Autenticação de usuários
- ✅ **Cadastro** - Registro de novos usuários
- ✅ **Logout** - Encerramento de sessão
- ✅ **Persistência de sessão** - Mantém usuário logado

### Gestão de Perfil
- ✅ **Visualizar perfil** - Exibição dos dados do usuário
- ✅ **Editar perfil** - Atualização de nome, nick e email
- ✅ **Excluir conta** - Remoção completa da conta

### Gestão de Usuários
- ✅ **Listar usuários** - Visualização de todos os usuários
- ✅ **Buscar usuários** - Pesquisa por nome ou nickname
- ✅ **Seguir usuários** - Sistema de seguimento
- ✅ **Deixar de seguir** - Desfazer seguimento

### Interface
- ✅ **Design responsivo** - Adaptável a qualquer dispositivo
- ✅ **Interface moderna** - Design limpo e intuitivo
- ✅ **Notificações** - Sistema de toast para feedback
- ✅ **Loading states** - Indicadores de carregamento
- ✅ **Modais de confirmação** - Para ações importantes

## 🛠️ Tecnologias Utilizadas

- **HTML5** - Estrutura semântica
- **CSS3** - Estilização moderna com gradientes e animações
- **JavaScript ES6+** - Lógica da aplicação (Vanilla JS)
- **Font Awesome** - Ícones
- **LocalStorage** - Persistência de dados no cliente
- **Fetch API** - Comunicação com a API

## 📁 Estrutura de Arquivos

```
frontend/
├── index.html      # Página principal
├── styles.css      # Estilos da aplicação
├── script.js       # Lógica JavaScript
└── README.md       # Documentação
```

## 🔧 Configuração

### 1. Configuração da API

No arquivo `script.js`, ajuste a URL base da API:

```javascript
const API_BASE_URL = 'http://localhost:5000'; // Altere para sua URL
```

### 2. CORS (Cross-Origin Resource Sharing)

Certifique-se de que sua API Go aceita requisições do frontend. Adicione os headers CORS necessários:

```go
// Exemplo de middleware CORS para sua API
func EnableCORS(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
    
    if r.Method == "OPTIONS" {
        w.WriteHeader(http.StatusOK)
        return
    }
}
```

## 🚀 Como Executar

### Opção 1: Servidor HTTP Simples

```bash
# Python 3
python -m http.server 8080

# Python 2
python -m SimpleHTTPServer 8080

# Node.js (se tiver npx instalado)
npx http-server -p 8080

# PHP
php -S localhost:8080
```

### Opção 2: Live Server (VS Code)

1. Instale a extensão "Live Server" no VS Code
2. Clique com o botão direito no `index.html`
3. Selecione "Open with Live Server"

### Opção 3: Servidor Web (Apache/Nginx)

Coloque os arquivos na pasta do servidor web:
- Apache: `/var/www/html/` ou `htdocs/`
- Nginx: `/usr/share/nginx/html/`

## 📱 Funcionalidades Detalhadas

### Dashboard
- **Estatísticas** - Visualização do número total de usuários
- **Cards informativos** - Interface amigável com informações relevantes

### Autenticação
- **Validação em tempo real** - Feedback imediato para o usuário
- **Tokens JWT** - Autenticação segura
- **Persistência de sessão** - Usuário permanece logado

### Gestão de Usuários
- **Cards de usuário** - Interface visual atrativa
- **Busca dinâmica** - Pesquisa em tempo real
- **Ações contextuais** - Botões de seguir/deixar de seguir

### Responsividade
- **Mobile First** - Otimizado para dispositivos móveis
- **Breakpoints** - Adaptação para tablet e desktop
- **Touch Friendly** - Botões e elementos adequados para touch

## 🎨 Design System

### Cores Principais
- **Primary**: `#667eea` (Azul/Roxo)
- **Secondary**: `#764ba2` (Roxo)
- **Success**: `#28a745` (Verde)
- **Danger**: `#e74c3c` (Vermelho)
- **Warning**: `#ffc107` (Amarelo)

### Tipografia
- **Font Family**: Segoe UI, Tahoma, Geneva, Verdana, sans-serif
- **Font Weights**: 400 (normal), 500 (medium), 700 (bold)

### Animações
- **Hover Effects** - Transformações suaves
- **Loading Spinner** - Indicador de carregamento
- **Toast Animations** - Notificações deslizantes

## 🔒 Segurança

### Autenticação
- **JWT Tokens** - Armazenados no localStorage
- **Auto-logout** - Em caso de token inválido
- **Headers Authorization** - Bearer token em requisições

### Validações
- **Client-side** - Validação básica nos formulários
- **Server-side** - Validação definitiva na API
- **Sanitização** - Tratamento de dados de entrada

## 🐛 Tratamento de Erros

### Tipos de Erro
- **Erro de rede** - Problemas de conectividade
- **Erro de autenticação** - Token inválido/expirado
- **Erro de validação** - Dados inválidos
- **Erro do servidor** - Problemas na API

### Feedback Visual
- **Toast notifications** - Mensagens de sucesso/erro
- **Estados de loading** - Indicadores visuais
- **Modais de confirmação** - Para ações destrutivas

## 📋 Endpoints da API Utilizados

| Método | Endpoint | Descrição | Auth |
|--------|----------|-----------|------|
| POST | `/login` | Autenticação | ❌ |
| POST | `/usuarios` | Criar usuário | ❌ |
| GET | `/usuarios` | Listar usuários | ✅ |
| GET | `/usuarios/{id}` | Buscar usuário | ✅ |
| PUT | `/usuarios/{id}` | Atualizar usuário | ✅ |
| DELETE | `/usuarios/{id}` | Deletar usuário | ✅ |
| POST | `/usuarios/{id}/seguir` | Seguir usuário | ✅ |
| POST | `/usuarios/{id}/parar-de-seguir` | Deixar de seguir | ✅ |

## 🔮 Possíveis Melhorias

### Funcionalidades
- [ ] Sistema de posts/publicações
- [ ] Chat em tempo real
- [ ] Notificações push
- [ ] Upload de fotos de perfil
- [ ] Sistema de grupos
- [ ] Feed de atividades

### Técnicas
- [ ] Service Workers (PWA)
- [ ] Lazy loading de imagens
- [ ] Otimização de performance
- [ ] Testes automatizados
- [ ] TypeScript
- [ ] Framework (React/Vue/Angular)

## 📞 Suporte

Para dúvidas ou problemas:
1. Verifique se a API está rodando
2. Confirme as configurações de CORS
3. Verifique o console do navegador
4. Teste os endpoints da API diretamente

## 📄 Licença

Este projeto é parte da API Gestar Bem e segue a mesma licença. 