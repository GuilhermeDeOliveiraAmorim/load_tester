package main

import (
	"flag"
	"fmt"
	"net/http"
	"sync"
	"time"
)

const (
	maxRetries     = 5
	initialBackoff = 100 * time.Millisecond
	maxBackoff     = 10 * time.Second
)

type result struct {
	status int
	err    error
}

func main() {
	url := flag.String("url", "", "URL do serviço a ser testado")
	totalRequests := flag.Int("requests", 100, "Número total de requests")
	concurrency := flag.Int("concurrency", 10, "Número de chamadas simultâneas")
	flag.Parse()

	if *url == "" {
		fmt.Println("A URL do serviço é obrigatória.")
		flag.Usage()
		return
	}

	results := make(chan result, *totalRequests)
	var wg sync.WaitGroup

	startTime := time.Now()

	for i := 0; i < *concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < *totalRequests / *concurrency; j++ {
				res, err := getWithRetry(*url)
				results <- result{status: res, err: err}
			}
		}()
	}

	wg.Wait()
	close(results)

	totalTime := time.Since(startTime)

	report(totalTime, results)
}

func getWithRetry(url string) (int, error) {
	var status int
	var err error
	backoff := initialBackoff

	for i := 0; i < maxRetries; i++ {
		resp, err := http.Get(url)
		if err != nil {
			status = 0
		} else {
			status = resp.StatusCode
			resp.Body.Close()
		}

		if err == nil && status != http.StatusTooManyRequests {
			return status, nil
		}

		time.Sleep(backoff)
		if backoff < maxBackoff {
			backoff *= 2
			if backoff > maxBackoff {
				backoff = maxBackoff
			}
		}
	}

	return status, err
}

func report(totalTime time.Duration, results chan result) {
	totalRequests := 0
	successfulRequests := 0
	statusCodes := make(map[int]int)

	for res := range results {
		totalRequests++
		if res.err == nil {
			statusCodes[res.status]++
			if res.status == http.StatusOK {
				successfulRequests++
			}
		}
	}

	fmt.Printf("Tempo total gasto: %v\n", totalTime)
	fmt.Printf("Quantidade total de requests realizados: %d\n", totalRequests)
	fmt.Printf("Quantidade de requests com status HTTP 200: %d\n", successfulRequests)
	fmt.Println("Distribuição dos códigos de status HTTP:")
	for code, count := range statusCodes {
		fmt.Printf("  %d: %d\n", code, count)
	}
}
