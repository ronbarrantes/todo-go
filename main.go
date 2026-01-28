package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

func fib(n int) int {
	cache := make(map[int]int)

	cache[0] = 0
	cache[1] = 1

	return fibHelper(n, cache)
}

func fibHelper(n int, c map[int]int) int {
	if n < 0 {
		return 0
	}

	if val, exists := c[n]; exists {
		return val
	}

	c[n] = fibHelper(n-1, c) + fibHelper(n-2, c)

	return c[n]
}

func main() {
	numStr := os.Args[1]
	num, err := strconv.Atoi(numStr)

	if err != nil || num < 0 {
		fmt.Println("Please provide a valid non-negative integer")
		return
	}

	s := time.Now()

	defer func() {
		duration := time.Since(s)
		fmt.Printf("This program took %v to complete\n", duration)
	}()

	fmt.Printf("fib of %d is %d\n", num, fib(num))
}
