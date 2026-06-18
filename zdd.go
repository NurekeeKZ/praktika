package main

import (
	"fmt"
)

// =========================
// 1. SUM / MIN / MAX
// =========================
func calcStats(nums []int) (sum, min, max int, err error) {
	if len(nums) == 0 {
		return 0, 0, 0, fmt.Errorf("empty slice")
	}

	sum = 0
	min = nums[0]
	max = nums[0]

	for _, n := range nums {
		sum += n

		if n < min {
			min = n
		}
		if n > max {
			max = n
		}
	}

	return sum, min, max, nil
}

// =========================
// 2. CALCULATOR (switch + error)
// =========================
func calculator(a, b int, op string) (int, error) {
	switch op {
	case "+":
		return a + b, nil
	case "-":
		return a - b, nil
	case "*":
		return a * b, nil
	case "/":
		if b == 0 {
			return 0, fmt.Errorf("division by zero")
		}
		return a / b, nil
	default:
		return 0, fmt.Errorf("unknown operator: %s", op)
	}
}

// =========================
// 3. WORKER (GOROUTINES)
// =========================
func worker(id int, jobs <-chan int, results chan<- int) {
	for j := range jobs {
		fmt.Println("Worker", id, "processing job", j)
		results <- j * j
	}
}

// =========================
// MAIN
// =========================
func main() {

	// -------- TASK 1 --------
	nums := []int{5, 2, 9, 1, 7}

	sum, min, max, err := calcStats(nums)
	if err != nil {
		fmt.Println("Stats error:", err)
	} else {
		fmt.Println("=== STATS ===")
		fmt.Println("Sum:", sum)
		fmt.Println("Min:", min)
		fmt.Println("Max:", max)
	}

	// -------- TASK 2 --------
	fmt.Println("\n=== CALCULATOR ===")

	if res, err := calculator(10, 5, "+"); err == nil {
		fmt.Println("10 + 5 =", res)
	} else {
		fmt.Println("Error:", err)
	}

	if res, err := calculator(10, 0, "/"); err == nil {
		fmt.Println("10 / 0 =", res)
	} else {
		fmt.Println("Error:", err)
	}

	// -------- TASK 3 --------
	fmt.Println("\n=== WORKERS ===")

	jobs := make(chan int)
	results := make(chan int)

	// запускаем 3 воркера
	for i := 1; i <= 3; i++ {
		go worker(i, jobs, results)
	}

	// отправляем задачи
	go func() {
		for i := 1; i <= 5; i++ {
			jobs <- i
		}
		close(jobs)
	}()

	// получаем результаты
	for i := 1; i <= 5; i++ {
		fmt.Println("Result:", <-results)
	}
}
