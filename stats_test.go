package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestFizzbuzzStats(t *testing.T) {
	requests, _ := getResponse()

	if stats.counter[requests[0]] != 2 && stats.counter[requests[1]] != 1 {
		t.Fatalf("Not counted right : had %d and %d, expected 2 and 1", stats.counter[requests[0]], stats.counter[requests[1]])
	}
}

func TestStatsEndpoint(t *testing.T) {
	requests, resp := getResponse()

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

func getResponse() ([]fizzbuzzRequest, *http.Response) {
	stats.Clear()

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

	return requests, resp
}
