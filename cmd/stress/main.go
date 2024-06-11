package main

import (
	"fmt"
	"github.com/magicx-ai/groq-go/cmd/stress/stress"
	"os"
	"sync"
	"time"
)

// StressTestConfig holds the configuration for the stress test
type StressTestConfig struct {
	Name              string
	RequestsPerSecond int
	Duration          time.Duration
	ExecFunc          func() error
}

// StressTestResult holds the result data of the stress test
type StressTestResult struct {
	Name              string
	RequestsPerSecond int
	TotalRequests     int
	ErrorCount        int
	SuccessCount      int
	ErrorRate         float64
	SuccessRate       float64
}

// StressTest runs a stress test based on the provided configuration
func StressTest(config StressTestConfig) StressTestResult {
	var wg sync.WaitGroup
	var mu sync.Mutex
	var errorCount, successCount int

	// Calculate the delay between each request
	delay := time.Second / time.Duration(config.RequestsPerSecond)

	// Create a ticker that will trigger at the specified rate
	ticker := time.NewTicker(delay)
	defer ticker.Stop()

	// Create a timeout to stop the stress test after the specified duration
	timeout := time.After(config.Duration)

	// Run the stress test
	for {
		select {
		case <-timeout:
			// Timeout reached, stop the stress test
			wg.Wait() // Wait for all goroutines to finish
			totalRequests := errorCount + successCount
			errorRate := float64(errorCount) / float64(totalRequests) * 100
			successRate := float64(successCount) / float64(totalRequests) * 100
			return StressTestResult{
				Name:              config.Name,
				RequestsPerSecond: config.RequestsPerSecond,
				TotalRequests:     totalRequests,
				ErrorCount:        errorCount,
				SuccessCount:      successCount,
				ErrorRate:         errorRate,
				SuccessRate:       successRate,
			}
		case <-ticker.C:
			// Launch a new request
			wg.Add(1)
			go func() {
				defer wg.Done()
				err := config.ExecFunc()
				mu.Lock()
				defer mu.Unlock()
				if err != nil {
					errorCount++
				} else {
					successCount++
				}
			}()
		}
	}
}

func main() {
	apiKey := os.Getenv("GROQ_API_KEY")
	if apiKey == "" {
		panic("GROQ_API_KEY not set")
	}

	const (
		targetRPS = 10
		duration  = 5 * time.Second
	)

	// Configure the stress test
	testcases := []StressTestConfig{
		{
			Name:              "Chat completion",
			RequestsPerSecond: targetRPS,
			Duration:          duration,
			ExecFunc:          stress.CompletionFunc(apiKey),
		},
		{
			Name:              "Chat stream completion",
			RequestsPerSecond: targetRPS,
			Duration:          duration,
			ExecFunc:          stress.CompletionStreamFunc(apiKey),
		},
		{
			Name:              "List Models",
			RequestsPerSecond: targetRPS,
			Duration:          duration,
			ExecFunc:          stress.ListModelsFunc(apiKey),
		},
		{
			Name:              "Retrieve Model",
			RequestsPerSecond: targetRPS,
			Duration:          duration,
			ExecFunc:          stress.RetrieveModelFunc(apiKey),
		},
	}

	var wg sync.WaitGroup
	for _, tc := range testcases {
		tc := tc
		wg.Add(1)
		go func() {
			defer wg.Done()
			result := StressTest(tc)
			printTestResult(result)
		}()
	}
	wg.Wait()
}

func printTestResult(r StressTestResult) {
	// Print the result
	fmt.Printf("Stress test: %s completed\n", r.Name)
	fmt.Printf("Requests per second: %d\n", r.RequestsPerSecond)
	fmt.Printf("Total requests: %d\n", r.TotalRequests)
	fmt.Printf("Error count: %d\n", r.ErrorCount)
	fmt.Printf("Success count: %d\n", r.SuccessCount)
	fmt.Printf("Error rate: %.2f%%\n", r.ErrorRate)
	fmt.Printf("Success rate: %.2f%%\n", r.SuccessRate)
}
