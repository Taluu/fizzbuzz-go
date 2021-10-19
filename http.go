package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
)

type fizzbuzzRequest struct {
	Int1  uint   `json:"int1"`
	Int2  uint   `json:"int2"`
	Limit uint   `json:"limit"`
	Str1  string `json:"str1"`
	Str2  string `json:"str2"`
}

var stats = struct {
	sync.Mutex
	counter map[fizzbuzzRequest]int
}{
	counter: make(map[fizzbuzzRequest]int),
}

func fizzbuzzHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Println()

	if r.Method != http.MethodPost {
		jsonError(w, "Only POST requests on this endpoint", http.StatusMethodNotAllowed)
		fmt.Printf("HTTP %s /fizzbuzz : 405 (Method Not Allowed)\n", r.Method)

		return
	}

	// The original fizzbuzz is printing "fizz" for each multiples of 3,
	// "buzz" for each multiples of 5, "fizzbuzz" for each multiples of 3 AND 5
	// and the number otherwise from 1 to 100.
	request := fizzbuzzRequest{
		Int1:  3,
		Int2:  5,
		Limit: 100,
		Str1:  "fizz",
		Str2:  "buzz",
	}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&request)

	if err != nil && err != io.EOF {
		jsonError(w, "Could not correctly unmarshall json", http.StatusBadRequest)
		fmt.Printf("HTTP POST /fizzbuzz : 400 (%s)\n", err.Error())

		return
	}

	fizzbuzz := Fizzbuzz(request.Int1, request.Int2, request.Limit, request.Str1, request.Str2)

	result := struct {
		Fizzbuzz string `json:"fizzbuzz"`
	}{
		Fizzbuzz: fizzbuzz,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)

	stats.Lock()
	stats.counter[request]++
	stats.Unlock()

	fmt.Printf("HTTP POST /fizzbuzz : 200\n")
}

func jsonError(w http.ResponseWriter, error string, code int) {
	HTTPError := struct {
		Error string `json:"error"`
	}{
		Error: error,
	}

	result, _ := json.Marshal(HTTPError)

	w.WriteHeader(code)
	w.Write(result)
}
