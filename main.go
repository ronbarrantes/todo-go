package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

func fib(n int, c map[int]int) int {
	if n <= 0 {
		return 0
	}

	if n <= 1 {
		return 1
	}

	if c[n] <= 0 {
		return c[n]
	}

	num := fib(n-1, c) + fib(n-2, c)

	c[n] = num

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

	cache := make(map[int]int)

	defer func() {
		duration := time.Since(s)
		fmt.Printf("This program took %v to complete\n", duration)
	}()

	fmt.Printf("fib of %d is %d\n", num, fib(num, cache))
}
