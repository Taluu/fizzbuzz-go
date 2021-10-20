package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestFails(t *testing.T) {
	type test struct {
		name            string
		method          string
		body            io.Reader
		expectedCode    int
		expectedMessage string
	}

	tests := []test{
		{name: "Wrong HTTP verb", method: "GET", body: nil, expectedCode: http.StatusMethodNotAllowed, expectedMessage: "Only POST requests on this endpoint"},
		{name: "Wrong type for parameter", method: "POST", body: strings.NewReader("{\"int2\": \"foo\"}"), expectedCode: http.StatusBadRequest, expectedMessage: "Could not correctly unmarshall json"},
	}

	for _, tc := range tests {
		r := httptest.NewRequest(tc.method, "/fizzbuzz", tc.body)
		w := httptest.NewRecorder()
		fizzbuzzHandler(w, r)

		resp := w.Result()

		if resp.Header.Get("Content-type") != "application/json" {
			t.Fatalf("%s :: Expected a application/json content-type, got %s", tc.name, resp.Header.Get("Content-type"))
		}

		if resp.StatusCode != tc.expectedCode {
			t.Fatalf("%s :: Did not expect HTTP %d (%s)", tc.name, resp.StatusCode, resp.Status)
		}

		error := struct {
			Error string `json:"error"`
		}{}

		decoder := json.NewDecoder(resp.Body)
		decoder.Decode(&error)

		if error.Error != tc.expectedMessage {
			t.Fatalf("%s :: Did not expect message %s", tc.name, error.Error)
		}
	}
}

func TestFizzbuzzOK(t *testing.T) {
	type test struct {
		name     string
		body     io.Reader
		expected []string
	}

	tests := []test{
		{name: "With all parameters", body: strings.NewReader("{\"int1\": 2, \"int2\": 3, \"limit\": 10, \"str1\": \"foo\", \"str2\": \"bar\"}"), expected: Fizzbuzz(2, 3, 10, "foo", "bar")},
		{name: "With only one parameter, should take default value", body: strings.NewReader("{\"limit\": 20}"), expected: Fizzbuzz(3, 5, 20, "fizz", "buzz")},
		{name: "Without any parameters, should take default values", body: nil, expected: Fizzbuzz(3, 5, 100, "fizz", "buzz")},
	}

	for _, tc := range tests {
		r := httptest.NewRequest("POST", "/fizzbuzz", tc.body)
		w := httptest.NewRecorder()
		fizzbuzzHandler(w, r)

		resp := w.Result()

		if resp.Header.Get("Content-type") != "application/json" {
			t.Fatalf("%s :: Expected a application/json content-type, got %s", tc.name, resp.Header.Get("Content-type"))
		}

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("%s :: Expected a 200 (OK), got %d (%s)", tc.name, resp.StatusCode, resp.Status)
		}

		result := struct {
			Fizzbuzz []string `json:"fizzbuzz"`
		}{}

		decoder := json.NewDecoder(resp.Body)
		decoder.Decode(&result)

		for k, v := range result.Fizzbuzz {
			if v != tc.expected[k] {
				t.Fatalf("%s :: Expected a correct fizzbuzz return, got %s", tc.name, result.Fizzbuzz)
			}
		}
	}
}
