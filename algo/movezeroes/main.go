package main

import "fmt"

// Задав целочисленный массив nums, переместите все 0 в его конец, сохранив относительный порядок ненулевых элементов.
// Обратите внимание, что вы должны сделать это на месте, не создавая копию массива.

// **Example 1:**
// **Input:** nums = [0,1,0,3,12]
// **Output:** [1,3,12,0,0]

// **Example 2:**
// **Input:** nums = [0]
// **Output:** [0]

func moveZeroes(nums []int) {
	if len(nums) < 2 {
		return
	}

	l := -1 // указатель на самый левый ноль
	r := 0  // идет по слайсу

	for r < len(nums) {
		if nums[r] > 0 {
			if l != -1 {
				swap(nums, r, l)
				l++
			}
		} else {
			if l == -1 {
				l = r
			}
		}

		r++
	}
}

func swap(nums []int, i, j int) {
	nums[i], nums[j] = nums[j], nums[i]
}

func main() {
	nums1 := []int{0, 1, 0, 3, 12}
	moveZeroes(nums1)
	fmt.Println(nums1)

	nums2 := []int{0}
	moveZeroes(nums2)
	fmt.Println(nums2)

	nums3 := []int{0, 5, 9, 7, 0, 0, 0, 0, 12}
	moveZeroes(nums3)
	fmt.Println(nums3)
}
