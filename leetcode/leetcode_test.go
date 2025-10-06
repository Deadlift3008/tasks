package leetcode

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBuyerDissatisfaction(t *testing.T) {
	t.Run("should return sum of buyer dissatisfaction", func(t *testing.T) {
		buyerNeeds := []int{1, 5, 10}
		goods := []int{1, 4, 7, 11}

		result := BuyerDissatisfaction(goods, buyerNeeds)

		require.Equal(t, 2, result)
	})

	t.Run("should return sum of buyer dissatisfaction on limit of slices", func(t *testing.T) {
		buyerNeeds := []int{1, 25, 7}
		goods := []int{8, 6, 15, 9}

		result := BuyerDissatisfaction(goods, buyerNeeds)

		require.Equal(t, 16, result)
	})
}

func TestGetChampions(t *testing.T) {
	t.Run("with only one user with all days records", func(t *testing.T) {
		testData := [][]Stat{{{
			userId: 1,
			steps:  100,
		}, {
			userId: 2,
			steps:  200,
		}}, {{
			userId: 1,
			steps:  100,
		}}, {{
			userId: 1,
			steps:  500,
		}}}

		result := GetChampions(testData)

		require.Equal(t, Champions{userIds: []int{1}, steps: 700}, result)
	})

	t.Run("if 2 users with all days records", func(t *testing.T) {
		testData := [][]Stat{{{
			userId: 1,
			steps:  100,
		}, {
			userId: 2,
			steps:  200,
		}}, {{
			userId: 1,
			steps:  100,
		}, {
			userId: 2,
			steps:  300,
		}}, {{
			userId: 1,
			steps:  500,
		}, {
			userId: 2,
			steps:  300,
		}}}

		result := GetChampions(testData)

		require.Equal(t, Champions{userIds: []int{2}, steps: 800}, result)
	})

	t.Run("if 2 users with all days records and equal steps", func(t *testing.T) {
		testData := [][]Stat{{{
			userId: 1,
			steps:  100,
		}, {
			userId: 2,
			steps:  200,
		}}, {{
			userId: 1,
			steps:  100,
		}, {
			userId: 2,
			steps:  300,
		}}, {{
			userId: 1,
			steps:  500,
		}, {
			userId: 2,
			steps:  200,
		}}}

		result := GetChampions(testData)

		require.Equal(t, Champions{userIds: []int{1, 2}, steps: 700}, result)
	})
}

func TestTopKLargest(t *testing.T) {
	t.Run("simple case", func(t *testing.T) {
		testData := []int{100, 20, 10, 500, 1}

		result := TopKLargest(testData, 2)

		require.Equal(t, []int{500, 100}, result)
	})
}
