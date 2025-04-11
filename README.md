# 🔗 TinyLiar

Um encurtador de links simples, funcional e levemente sarcástico — feito com Go, PostgreSQL e uma pitada de desespero.

## 🚀 Funcionalidades

- Criar links curtos (`POST /?shorten=...`)
- Redirecionar para o link original (`GET /{short_url}`)
- Contar cliques porque estatísticas importam
- Verificar expiração dos links
- 🔒 Suporte a encriptação (opcional, mas chique)

## 🛠️ Tecnologias

- **Go** (com Chi router)
- **PostgreSQL**
- **Docker & Docker Compose**
- `sqlx` para acesso ao banco
- `slog` pra fingir que fazemos logging profissional

## 📦 Instalação

Clone o projeto:

```bash
git clone https://github.com/JeanGrijp/tinyliar.git
cd tinyliar
```
