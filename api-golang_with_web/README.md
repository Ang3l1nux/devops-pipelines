# **API Golang simples com interface web**

<img src="./img/golang.png" width="100" height="100" /><br>
by jose.lima - Maio de 2025

---

## 🚀 Funcionalidades

- 📄 Cadastro de pessoas com nome, CPF e data de nascimento
- 🔎 Busca por nome ou CPF (com correspondência parcial)
- 🧹 Botão de limpar resultados
- ❌ Exclusão de registros diretamente pela interface web
- 🔄 Atualização e persistência dos dados em banco de dados PostgreSQL
- 📦 Pronto para deploy via Docker e Kubernetes

## 🧪 Tecnologias

- [Golang](https://golang.org/)
- [PostgreSQL](https://www.postgresql.org/)
- [HTML + JS (vanilla)](https://developer.mozilla.org/)
- [Docker](https://www.docker.com/)
- [Kubernetes](https://kubernetes.io/)

---

## ⚙️ Como rodar localmente

1. Clone o projeto:

```bash
git clone https://github.com/seuusuario/cadastro-pessoas-go.git
cd cadastro-pessoas-go
```

2. Suba um banco PostgreSQL (pode usar Docker):

```bash
docker run --name postgres -e POSTGRES_PASSWORD=123456 -e POSTGRES_USER=postgres -e POSTGRES_DB=meubanco -p 5432:5432 -d postgres
```

3. Configure seu `.env` ou defina as variáveis de ambiente no seu sistema:

```bash
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=123456
export DB_NAME=meubanco
```

4. Execute a aplicação:

```bash
go run main.go db.go
```

5. Acesse a interface web:

[
    http://localhost:8080](http://localhost:8080)

---
