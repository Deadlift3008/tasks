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
