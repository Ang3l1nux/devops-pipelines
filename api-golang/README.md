# **API Golang simples para testes**

<img src="./img/golang.png" width="100" height="100" /><br>
by jose.lima - Maio de 2025

---

## 🚀 Funcionalidades

Api simples para cadastro de usuários com os campos:

- ID
- Nome
- CPF
- DataNascimento

## 🧪 Tecnologias

- [Golang](https://golang.org/)
- [PostgreSQL](https://www.postgresql.org/)
- [Docker](https://www.docker.com/)

---

## ⚙️ Como rodar localmente

## **Step 1 - Install Docker**

```bash
# Update System
sudo apt update
sudo apt upgrade

# Install packages
sudo apt-get install  curl apt-transport-https ca-certificates software-properties-common

# add GPG Key
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -

# add Repo
sudo add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable"

# Update
sudo apt update

# Install Docker
sudo apt install docker-ce

# Verify status Docker
sudo systemctl status docker
```

---

## **Step 2 - Create container Postgres**

Baixar a imagem do postgres:

```bash
docker pull postgres
```

Check seu repositorio local:

```bash
docker images
```

Execute o container do Postgres

```bash
docker run -it -d -p 5432:5432 -e POSTGRES_PASSWORD=1234 postgres
```

Check se esta em execução:

```bash
docker ps
```

### **Alternativa para persistir dados**

Criar um diretório com o nome docker:

```bash
mkdir docker
cd docker
```

Criar o arquivo docker-compose.yaml

```bash
vim docker-compose.yaml
```

```yaml
version: '3.8'

services:
  postgres:
    image: postgres:latest
    container_name: postgres-db
    restart: always
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: 1234
    volumes:
      - pgdata:/var/lib/postgresql/data
    ports:
      - "5432:5432"

volumes:
  pgdata:
```

subir o container a partir do docker-compose:

```bash
docker-compose up -d
docker ps
```


> No Docker, `pgdata` =é um nome padrão para um volume Docker, usado para persistir os dados de um contêiner PostgreSQL=**


**Onde arquivos serão persistidos ?**

>  Em sistemas Linux ou WSL o volume fica em:

```bash
/var/lib/docker/volumes/pgdata/_data/
```

>  Em sistermas Windows o volume fica em:

```bash
/var/lib/docker/volumes/docker_pgdata/_data
```

---

Step 3 - Create Database

Com o DBeaver logar no banco e executar:

```sql
CREATE DATABASE cadastro_pessoas;

\c cadastro_pessoas;

CREATE TABLE pessoas (
    id SERIAL PRIMARY KEY,
    nome VARCHAR(100) NOT NULL,
    cpf VARCHAR(11) NOT NULL UNIQUE,
    data_nascimento DATE NOT NULL
);

```

---

## **Step 4 - Install Go**

💾 Baixe e instale o Go

🔗 [Página oficial de downloads do Go]()

Baixe a versão estável mais recente para seu sistema operacional.

Exemplo para Linux (Go 1.22.0):

```bash
wget https://go.dev/dl/go1.21.6.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.6.linux-amd64.tar.gz
```

📌 Configure as variáveis de ambiente

Adicione ao seu `~/.bashrc`, `~/.zshrc` ou equivalente:

```bash
export PATH=$PATH:/usr/local/go/bin
```

E atualize o terminal:

```bash
source ~/.bashrc
# ou
source ~/.zshrc
```

---

## **Step 5 - Create o código em GO**

🛠️ 1. Crie uma pasta para o projeto e entre nela:

```bash
mkdir cadastro-pessoas
cd cadastro-pessoas
```

📦 2. Inicie o módulo Go:

```bash
go mod init cadastro-pessoas
```

Esse comando cria um arquivo `go.mod`, que é necessário para gerenciar dependências.

🔧 3. Agora você pode instalar as dependências:

```bash
go get github.com/gorilla/mux
go get github.com/lib/pq
```

 🗂️ 4. Código Go:

```go
package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

var db *sql.DB
var err error

type Pessoa struct {
	ID            int    `json:"id"`
	Nome          string `json:"nome"`
	CPF           string `json:"cpf"`
	DataNascimento string `json:"data_nascimento"`
}

func init() {
	// Configuração do banco de dados
	connStr := "user=postgres password=1234 dbname=cadastro_pessoas sslmode=disable"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/pessoas", getPessoas).Methods("GET")
	router.HandleFunc("/pessoas/{id}", getPessoa).Methods("GET")
	router.HandleFunc("/pessoas", criarPessoa).Methods("POST")
	router.HandleFunc("/pessoas/{id}", atualizarPessoa).Methods("PUT")
	router.HandleFunc("/pessoas/{id}", deletarPessoa).Methods("DELETE")

	fmt.Println("Servidor rodando na porta 8080") // <- adicione essa linha
	log.Println("Servidor rodando na porta 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

// Função para listar todas as pessoas
func getPessoas(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, nome, cpf, data_nascimento FROM pessoas")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var pessoas []Pessoa
	for rows.Next() {
		var p Pessoa
		err := rows.Scan(&p.ID, &p.Nome, &p.CPF, &p.DataNascimento)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		pessoas = append(pessoas, p)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pessoas)
}

// Função para buscar pessoa por ID
func getPessoa(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	var p Pessoa
	err := db.QueryRow("SELECT id, nome, cpf, data_nascimento FROM pessoas WHERE id = $1", id).Scan(&p.ID, &p.Nome, &p.CPF, &p.DataNascimento)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

// Função para cadastrar uma pessoa
func criarPessoa(w http.ResponseWriter, r *http.Request) {
	var p Pessoa
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = db.Exec("INSERT INTO pessoas (nome, cpf, data_nascimento) VALUES ($1, $2, $3)", p.Nome, p.CPF, p.DataNascimento)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(p)
}

// Função para atualizar uma pessoa
func atualizarPessoa(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	var p Pessoa
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = db.Exec("UPDATE pessoas SET nome = $1, cpf = $2, data_nascimento = $3 WHERE id = $4", p.Nome, p.CPF, p.DataNascimento, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(p)
}

// Função para deletar uma pessoa
func deletarPessoa(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	_, err := db.Exec("DELETE FROM pessoas WHERE id = $1", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Pessoa deletada com sucesso"))
}
```

▶️ 4. Execute sua aplicação:

```bash
go run main.go &
```

🌐 5. Acessar:

http://localhost:8080/pessoas

---

## **Step 6 - Testes A**

**🧪 Testes com Curl:**

---

## **Step 7 - Testes B**

🧪 Testes com Postman:

---
