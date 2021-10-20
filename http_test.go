package main

import (
	"encoding/json"
	"fmt"
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

func TestFizzbuzzStats(t *testing.T) {
	// clear the map counter
	stats.counter = make(map[fizzbuzzRequest]uint)

	// req1 and req3 are requests with the same parameters
	for _, str1 := range []string{"fuzz", "foo", "fuzz"} {
		r := httptest.NewRequest("POST", "/fizzbuzz", strings.NewReader(fmt.Sprintf("{\"str1\": \"%s\", \"str2\": \"stat\"}", str1)))
		fizzbuzzHandler(httptest.NewRecorder(), r)
	}

	requests := []fizzbuzzRequest{
		{Int1: 3, Int2: 5, Limit: 100, Str1: "fuzz", Str2: "stat"},
		{Int1: 3, Int2: 5, Limit: 100, Str1: "foo", Str2: "stat"},
	}

	if stats.counter[requests[0]] != 2 && stats.counter[requests[1]] != 1 {
		t.Fatalf("Not counted right : had %d and %d, expected 2 and 1", stats.counter[requests[0]], stats.counter[requests[1]])
	}
}

func TestStatsEndpoint(t *testing.T) {
	// clear the map counter
	stats.counter = make(map[fizzbuzzRequest]uint)

	// req1 and req3 are requests with the same parameters
	for _, str1 := range []string{"fuzz", "foo", "fuzz"} {
		r := httptest.NewRequest("POST", "/fizzbuzz", strings.NewReader(fmt.Sprintf("{\"str1\": \"%s\", \"str2\": \"stat\"}", str1)))
		fizzbuzzHandler(httptest.NewRecorder(), r)
	}

	requests := []fizzbuzzRequest{
		{Int1: 3, Int2: 5, Limit: 100, Str1: "fuzz", Str2: "stat"},
		{Int1: 3, Int2: 5, Limit: 100, Str1: "foo", Str2: "stat"},
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/stats", nil)
	statsHandler(w, r)

	resp := w.Result()

	if resp.Header.Get("Content-type") != "application/json" {
		t.Fatalf("Stats endpoint :: Expected a application/json content-type, got %s", resp.Header.Get("Content-type"))
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Stats endpoint :: Expected a 200 (OK), got %d (%s)", resp.StatusCode, resp.Status)
	}

	got := make([]statistic, 0)
	expected := make([]statistic, 0, 2)

	expected = append(expected, statistic{Request: requests[0], Counter: 2})
	expected = append(expected, statistic{Request: requests[1], Counter: 1})

	decoder := json.NewDecoder(resp.Body)
	decoder.Decode(&got)

	if len(got) != 2 {
		t.Fatalf("Stats endpoint :: expected an array of length 2, got an array of length %d", len(got))
		return
	}

	for k, stat := range got {
		if stat != expected[k] {
			t.Errorf("Stats endpoint :: expected %v, got %v", expected[k], stat)
		}
	}
}
