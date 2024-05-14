package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Carro struct {
	Nome   string `json:"nome"`
	Modelo string `json:"modelo"`
	Ano    int    `json:"-"`
}

func (c Carro) Andar() {
	fmt.Println("O carro ", c.Nome, " está andando")
}

func (c Carro) Parar() {
	fmt.Println("O carro ", c.Nome, " parou")
}

func main() {

	carro1 := Carro{Nome: "Fusca", Modelo: "VW", Ano: 1970}
	//carro2:= Carro{Nome: "Gol", Modelo: "VW", Ano: 2000}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		//convert struct to json
		json.NewEncoder(w).Encode(carro1)
	})

	http.ListenAndServe(":8080", nil)
}
