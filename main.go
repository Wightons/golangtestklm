package main

import (
	"encoding/json"
	"log"
	"math/big"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type Numbers struct {
	A int `json:"a"`
	B int `json:"b"`
}

type Factorial struct {
	A string `json:"a"`
	B string `json:"b"`
}

func main() {
	router := httprouter.New()
	router.POST("/calculate", CalculateHandler)
	log.Fatal(http.ListenAndServe(":8989", router))
}

func CalculateHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var numbers Numbers
	err := json.NewDecoder(r.Body).Decode(&numbers)
	if err != nil || numbers.A < 0 || numbers.B < 0 {
		http.Error(w, `{"error":"Incorrect input"}`, http.StatusBadRequest)
		return
	}

	chA := make(chan *big.Int)
	chB := make(chan *big.Int)

	go func() {
		chA <- factorial(numbers.A)
	}()

	go func() {
		chB <- factorial(numbers.B)
	}()

	factorial := Factorial{
		A: (<-chA).String(),
		B: (<-chB).String(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(factorial)
}

func factorial(n int) *big.Int {
	result := big.NewInt(1)
	for i := 2; i <= n; i++ {
		result.Mul(result, big.NewInt(int64(i)))
	}
	return result
}
