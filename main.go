package main

import (
	"flag"
	"fmt"
	"net/http"
	"sync"
	"time"
)

var (
	url         string
	requests    int
	concurrency int
)

func init() {
	flag.StringVar(&url, "url", "", "URL do serviço a ser testado.")
	flag.IntVar(&requests, "requests", 100, "Número total de requests.")
	flag.IntVar(&concurrency, "concurrency", 10, "Número de chamadas simultâneas.")
}

func main() {
	flag.Parse()

	if url == "" {
		fmt.Println("URL do serviço é obrigatória.")
		return
	}

	start := time.Now()
	var wg sync.WaitGroup
	requestsPerGoroutine := requests / concurrency
	results := make(chan int, requests)

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < requestsPerGoroutine; j++ {
				resp, err := http.Get(url)
				if err != nil {
					results <- 0
					continue
				}
				results <- resp.StatusCode
				resp.Body.Close()
			}
		}()
	}

	wg.Wait()
	close(results)

	statusCount := make(map[int]int)
	for status := range results {
		statusCount[status]++
	}

	totalTime := time.Since(start)

	fmt.Printf("Tempo total gasto: %v\n", totalTime)
	fmt.Printf("Quantidade total de requests realizados: %d\n", requests)
	fmt.Printf("Quantidade de requests com status HTTP 200: %d\n", statusCount[http.StatusOK])
	fmt.Println("Distribuição de outros códigos de status HTTP:")
	for status, count := range statusCount {
		if status != http.StatusOK {
			fmt.Printf("HTTP %d: %d\n", status, count)
		}
	}
}
