package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"sync"
)

type requestCounter struct {
	sync.Mutex
	counter map[fizzbuzzRequest]uint
}

type statistic struct {
	Request fizzbuzzRequest `json:"request"`
	Counter uint            `json:"counter"`
}

func (s *requestCounter) Incr(request fizzbuzzRequest) {
	s.Lock()
	defer s.Unlock()

	s.counter[request]++
}

func (s *requestCounter) MarshalJSON() ([]byte, error) {
	s.Lock()
	defer s.Unlock()

	statistics := make([]statistic, 0, len(s.counter))

	for request, counter := range s.counter {
		statistics = append(statistics, statistic{Request: request, Counter: counter})
	}

	sort.Slice(statistics, func(i, j int) bool {
		if statistics[i].Counter == statistics[j].Counter {
			return i < j
		}

		return statistics[i].Counter > statistics[j].Counter
	})

	return json.Marshal(statistics)
}

func (s *requestCounter) Clear() {
	s.Lock()
	defer s.Unlock()

	s.counter = make(map[fizzbuzzRequest]uint)
}

var stats requestCounter = requestCounter{
	counter: make(map[fizzbuzzRequest]uint),
}

func statsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Println()

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&stats)

	fmt.Printf("HTTP GET /stats : 200\n")
}
