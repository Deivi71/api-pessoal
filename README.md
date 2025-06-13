# 🌟 Gestar Bem - Rede Social

Uma rede social moderna desenvolvida em **Go** com frontend em **HTML/CSS/JavaScript**.

![Status](https://img.shields.io/badge/Status-Em%20Desenvolvimento-yellow)
![Go Version](https://img.shields.io/badge/Go-1.19+-blue)
![MySQL](https://img.shields.io/badge/MySQL-8.0+-orange)
![License](https://img.shields.io/badge/License-MIT-green)

## 🚀 Início Rápido

```bash
# Clone o repositório
git clone <repository-url>
cd APi-gestar-bem

# Configure a base de dados
mysql -u root -p < sql/sql.sql
mysql -u root -p < sql/add_foto_perfil.sql

# Configure as variáveis de ambiente
cp .env.example .env

# Execute o servidor
go run main.go

# Acesse o frontend
open frontend/index.html
```

## ✨ Funcionalidades

- 🔐 **Autenticação JWT** - Sistema seguro de login
- 👥 **Gestão de Usuários** - Perfis completos com fotos
- 📝 **Publicações** - Sistema de posts e feed
- 💬 **Comentários** - Interação nas publicações
- 👥 **Seguidores** - Rede social completa
- 🌙 **Modo Escuro** - Interface adaptável
- 📱 **Responsivo** - Funciona em todos os dispositivos
- 📸 **Upload de Imagens** - Fotos de perfil

## 📚 Documentação Completa

### 🏠 Acesso Principal
**[📖 Documentação Completa](docs/README.md)** - Guia completo do projeto

### 🔗 Links Rápidos
- **[📋 API Reference](docs/api/README.md)** - Endpoints e exemplos
- **[🎨 Frontend Guide](docs/frontend/README.md)** - Interface e componentes  
- **[🗄️ Database Schema](docs/database/README.md)** - Estrutura da base de dados
- **[📸 Upload de Fotos](docs/features/UPLOAD_FOTO_PERFIL.md)** - Funcionalidade completa
- **[📚 Índice Geral](docs/INDEX.md)** - Navegação por toda documentação

## 🏗️ Arquitetura

```
APi-gestar-bem/
├── 📚 docs/                 # Documentação completa
├── 🔧 src/                  # Código fonte Go
├── 🎨 frontend/             # Interface web
├── 🗄️ sql/                  # Scripts de base de dados
├── 📁 uploads/              # Arquivos enviados
└── 📄 README.md             # Este arquivo
```

## 🛠️ Stack Tecnológica

**Backend:**
- Go 1.19+
- Gorilla Mux (HTTP Router)
- JWT (Autenticação)
- MySQL 8.0+
- bcrypt (Hash de senhas)

**Frontend:**
- HTML5 + CSS3 + JavaScript ES6+
- Fetch API
- Design Responsivo
- Modo Escuro/Claro

## 🤝 Contribuição

1. Fork o projeto
2. Crie uma branch (`git checkout -b feature/nova-funcionalidade`)
3. Commit suas mudanças (`git commit -m 'Adiciona nova funcionalidade'`)
4. Push para a branch (`git push origin feature/nova-funcionalidade`)
5. Abra um Pull Request

## 📝 Licença

Este projeto está sob a licença MIT. Veja o arquivo [LICENSE](LICENSE) para detalhes.

## 📞 Suporte

- 📧 **Email**: suporte@gestarbem.com
- 🐛 **Issues**: [GitHub Issues](https://github.com/user/repo/issues)
- 📚 **Documentação**: [docs/](docs/)

---

**Desenvolvido com ❤️ pela equipe Gestar Bem** 