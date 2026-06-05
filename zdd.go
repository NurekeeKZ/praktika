package main

import "fmt"

func calc(nums []int) (int, int, int) {
	sum := 0
	min := nums[0]
	max := nums[0]

	for _, n := range nums {
		sum += n

		if n < min {
			min = n
		}
		if n > max {
			max = n
		}
	}

	return sum, min, max
}

func main() {
	nums := []int{5, 2, 9, 1, 7}

	sum, min, max := calc(nums)

	fmt.Println("Sum:", sum)
	fmt.Println("Min:", min)
	fmt.Println("Max:", max)
}
