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

