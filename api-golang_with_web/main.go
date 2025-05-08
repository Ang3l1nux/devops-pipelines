package main

import (
    "log"
    "net/http"
    "github.com/gorilla/mux"
)

type Pessoa struct {
    ID        int    `json:"id"`
    Nome      string `json:"nome"`
    CPF       string `json:"cpf"`
    Nascimento string `json:"nascimento"`
}

func main() {
    InitDB()
    router := mux.NewRouter()

    router.HandleFunc("/api/pessoas", GetPessoas).Methods("GET")
    router.HandleFunc("/api/pessoas", CreatePessoa).Methods("POST")
    router.HandleFunc("/api/pessoas/{cpf}", UpdatePessoa).Methods("PUT")
    router.HandleFunc("/api/pessoas/{cpf}", DeletePessoa).Methods("DELETE")
    router.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("./static"))))

    log.Println("Servidor rodando em http://localhost:8080")
    log.Fatal(http.ListenAndServe(":8080", router))
}
