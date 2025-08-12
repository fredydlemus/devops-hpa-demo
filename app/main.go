package main

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"time"
)

func burnCPU(ms int){
	end := time.Now().Add(time.Duration(ms) * time.Millisecond)
	x := 0.0001
	for time.Now().Before(end) {
		x += math.Sqrt(x) * 1.00001
		if x > 1e9 {
			x = 0.0001
		}
	}
}

func handler(w http.ResponseWriter, r *http.Request){
	burn := 0
	if v := r.URL.Query().Get("burnMs"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			burn = n
		}
	}

	if burn == 0{
		if env := os.Getenv("BURN_MS"); env != "" {
			if n, err := strconv.Atoi(env); err == nil {
				burn = n
			}
		}
	}

	if burn < 0{
		burn = 0
	}

	start := time.Now()
	burnCPU(burn)
	fmt.Fprintf(w, "ok | burned=%dms | took=%s\n", burn, time.Since(start))
}

func health(w http.ResponseWriter, r *http.Request){
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}

func main(){
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)
	mux.HandleFunc("/healthz", health)
	addr := ":8080"
	log.Printf("listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}