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

//==================================================================

type Stat struct {
	userId int
	steps  int
}

type Champions struct {
	userIds []int
	steps   int
}

type TransformedRecordMap map[int]int

func transformRecordsToMap(statistics [][]Stat) []TransformedRecordMap {
	result := make([]TransformedRecordMap, 0)

	for _, dayRecord := range statistics {
		localMap := make(map[int]int, 0)

		for _, userRecord := range dayRecord {
			localMap[userRecord.userId] = userRecord.steps
		}

		result = append(result, localMap)
	}

	return result
}

func GetChampions(statistics [][]Stat) Champions {
	candidates := make(map[int]int, 0)

	for _, userRecord := range statistics[0] {
		candidates[userRecord.userId] = 0
	}

	if len(candidates) == 0 {
		return Champions{}
	}

	statMaps := transformRecordsToMap(statistics)

	maxSteps := 0

	for _, dayStat := range statMaps {
		for candidate := range candidates {
			candidateSteps, exist := dayStat[candidate]

			if !exist {
				delete(candidates, candidate)
				continue
			}

			candidates[candidate] += candidateSteps

			if candidates[candidate] > maxSteps {
				maxSteps = candidates[candidate]
			}
		}
	}

	championsIds := make([]int, 0)

	for candidateId, steps := range candidates {
		if steps == maxSteps {
			championsIds = append(championsIds, candidateId)
		}
	}

	return Champions{userIds: championsIds, steps: maxSteps}
}
