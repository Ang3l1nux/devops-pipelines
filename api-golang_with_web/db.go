package main

import (
    "database/sql"
    _ "github.com/lib/pq"
    "log"
    "net/http"
    "encoding/json"
    "github.com/gorilla/mux"
)

var db *sql.DB

func InitDB() {
    var err error
    connStr := "user=postgres password=1234 dbname=cadastro sslmode=disable"
    db, err = sql.Open("postgres", connStr)
    if err != nil {
        log.Fatal(err)
    }

    _, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS pessoas (
            id SERIAL PRIMARY KEY,
            nome TEXT NOT NULL,
            cpf TEXT UNIQUE NOT NULL,
            nascimento TEXT NOT NULL
        );
    `)
    if err != nil {
        log.Fatal(err)
    }
}

func GetPessoas(w http.ResponseWriter, r *http.Request) {
    rows, err := db.Query("SELECT id, nome, cpf, nascimento FROM pessoas")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var pessoas []Pessoa
    for rows.Next() {
        var p Pessoa
        rows.Scan(&p.ID, &p.Nome, &p.CPF, &p.Nascimento)
        pessoas = append(pessoas, p)
    }

    json.NewEncoder(w).Encode(pessoas)
}

func CreatePessoa(w http.ResponseWriter, r *http.Request) {
    var p Pessoa
    json.NewDecoder(r.Body).Decode(&p)

    _, err := db.Exec("INSERT INTO pessoas (nome, cpf, nascimento) VALUES ($1, $2, $3)", p.Nome, p.CPF, p.Nascimento)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
}

func UpdatePessoa(w http.ResponseWriter, r *http.Request) {
    cpf := mux.Vars(r)["cpf"]
    var p Pessoa
    json.NewDecoder(r.Body).Decode(&p)

    _, err := db.Exec("UPDATE pessoas SET nome=$1, nascimento=$2 WHERE cpf=$3", p.Nome, p.Nascimento, cpf)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
}

func DeletePessoa(w http.ResponseWriter, r *http.Request) {
    cpf := mux.Vars(r)["cpf"]
    _, err := db.Exec("DELETE FROM pessoas WHERE cpf=$1", cpf)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
}
