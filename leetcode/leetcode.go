package leetcode

import (
	"math"
	"slices"
	"sort"
)

func findClosestIndex(target int, areaIndex int, goods []int) int {
	indicesToCheck := make([]int, 0, 3)

	if areaIndex != 0 {
		indicesToCheck = append(indicesToCheck, areaIndex-1)
	}

	if areaIndex >= len(goods)-1 {
		indicesToCheck = append(indicesToCheck, len(goods)-1)
	} else {
		indicesToCheck = append(indicesToCheck, areaIndex, areaIndex+1)
	}

	minIndex := indicesToCheck[0]
	minDiff := int(math.Abs(float64(goods[minIndex] - target)))

	for _, indexToCheck := range indicesToCheck {
		currentDiff := int(math.Abs(float64(goods[indexToCheck] - target)))

		if currentDiff < minDiff {
			minDiff = currentDiff
			minIndex = indexToCheck
		}
	}

	return minIndex
}

func BuyerDissatisfaction(goods []int, buyerNeeds []int) int {
	sort.Ints(goods)

	result := 0

	for _, need := range buyerNeeds {
		idx, found := slices.BinarySearch(goods, need)

		if found {
			continue
		}

		closestIndex := findClosestIndex(need, idx, goods)
		result += int(math.Abs(float64(goods[closestIndex] - need)))
	}

	return result
}
