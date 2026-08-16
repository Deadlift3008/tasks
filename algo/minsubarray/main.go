package main

import (
	"fmt"
	"sort"
)

// Дан массив положительных целых чисел `nums` и положительное целое число `target`. Необходимо найти минимальную длину подмассива, сумма элементов которого больше или равна `target`. Если такого подмассива не существует, вернуть `0`.

// Пример 1:
// Ввод: target = 7, nums = [2,3,1,2,4,3]
// Вывод: 2
// Объяснение:
// Подмассив `[4,3]` имеет сумму 4+3=7 и минимальную длину среди всех возможных подмассивов, удовлетворяющих условию.

// Пример 2:
// Ввод: target = 4, nums = [1,4,4]
// Вывод: 1
// Объяснение:
// Подмассив `[4]` уже удовлетворяет условию, и его длина равна 11.

// Пример 3:
// Ввод: target = 11, nums = [1,1,1,1,1,1,1,1]

// Вывод: 0

// Объяснение:
// Ни один подмассив не имеет суммы, достаточной для выполнения условия. Ответ: 0.

func minSubArrayLen(target int, nums []int) int {
	sort.Slice(nums, func(i, j int) bool {
		return nums[i] > nums[j]
	})

	count := 0
	current := 0

	for _, num := range nums {
		current += num
		count++

		if current >= target {
			return count
		}
	}

	return 0
}

func main() {
	// Примеры ввода
	nums1 := []int{2, 3, 1, 2, 4, 3}
	target1 := 7

	nums2 := []int{1, 4, 4}
	target2 := 4

	nums3 := []int{1, 1, 1, 1, 1, 1, 1, 1}
	target3 := 11

	fmt.Println(minSubArrayLen(target1, nums1)) // Ожидается: 2
	fmt.Println(minSubArrayLen(target2, nums2)) // Ожидается: 1
	fmt.Println(minSubArrayLen(target3, nums3)) // Ожидается: 0
}
