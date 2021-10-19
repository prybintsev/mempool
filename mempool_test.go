package memepool

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type testPriorityItem struct {
	priority float64
}

func (t testPriorityItem) Priority() float64  {
	return t.priority
}

func TestMemPool_PriorityOrdering(t *testing.T) {
	tests := map[string]struct {
		inputTxs []testPriorityItem
		expectedOutput []testPriorityItem
	} {
		"permutation 1": {
			inputTxs: []testPriorityItem{{priority: 1},{priority: 2},{priority: 3}},
			expectedOutput: []testPriorityItem{{priority: 3},{priority: 2},{priority: 1}},
		},
		"permutation 2": {
			inputTxs: []testPriorityItem{{priority: 1},{priority: 3},{priority: 2}},
			expectedOutput: []testPriorityItem{{priority: 3},{priority: 2},{priority: 1}},
		},
		"permutation 3": {
			inputTxs: []testPriorityItem{{priority: 2},{priority: 1},{priority: 3}},
			expectedOutput: []testPriorityItem{{priority: 3},{priority: 2},{priority: 1}},
		},
		"permutation 4": {
			inputTxs: []testPriorityItem{{priority: 2},{priority: 3},{priority: 1}},
			expectedOutput: []testPriorityItem{{priority: 3},{priority: 2},{priority: 1}},
		},
		"permutation 5": {
			inputTxs: []testPriorityItem{{priority: 3},{priority: 1},{priority: 2}},
			expectedOutput: []testPriorityItem{{priority: 3},{priority: 2},{priority: 1}},
		},
		"permutation 6": {
			inputTxs: []testPriorityItem{{priority: 3},{priority: 2},{priority: 1}},
			expectedOutput: []testPriorityItem{{priority: 3},{priority: 2},{priority: 1}},
		},
		"empty": {
			inputTxs: []testPriorityItem{},
			expectedOutput: []testPriorityItem{},
		},
		"equal items": {
			inputTxs: []testPriorityItem{{priority: 1},{priority: 1},{priority: 3}},
			expectedOutput: []testPriorityItem{{priority: 3},{priority: 1},{priority: 1}},
		},
	}

	for tName, tc := range tests {
		tc := tc
		t.Run(tName, func(t *testing.T) {
			memPool := NewMemPool()
			for _, tx := range tc.inputTxs {
				memPool.Push(tx)
			}

			output := []testPriorityItem{}
			for i := 0; i < len(tc.inputTxs); i++ {
				item := memPool.Pop()
				testItem := item.(testPriorityItem)
				output = append(output, testItem)
			}

			require.Equal(t, tc.expectedOutput, output)
		})
	}
}


func TestMemPool_OverCapacity(t *testing.T) {
	memPool := NewMemPool()
	for i := 0; i <  5500; i++ {
		memPool.Push(testPriorityItem{priority: float64(i)})
	}

	var expectedOutput []PriorityItem
	for i := 5499; i >= 500; i-- {
		expectedOutput = append(expectedOutput, testPriorityItem{priority: float64(i)})
	}

	var output []PriorityItem
	for {
		if memPool.Len() == 0 {
			break
		}
		item := memPool.Pop()
		output = append(output, item)
	}

	require.Equal(t, expectedOutput, output)
}

func TestMemPool_PoppingFromEmptyMemPool(t *testing.T) {
	memPool := NewMemPool()
	memPool.Push(testPriorityItem{priority: float64(1)})
	memPool.Push(testPriorityItem{priority: float64(2)})
	memPool.Push(testPriorityItem{priority: float64(3)})

	require.Equal(t, 3, memPool.Len())

	require.Equal(t, float64(3), memPool.Pop().Priority())
	require.Equal(t, 2, memPool.Len())
	require.Equal(t, float64(2), memPool.Pop().Priority())
	require.Equal(t, 1, memPool.Len())
	require.Equal(t, float64(1), memPool.Pop().Priority())
	require.Equal(t, 0, memPool.Len())

	require.Panics(t, func() {
		memPool.Pop()
	})
}