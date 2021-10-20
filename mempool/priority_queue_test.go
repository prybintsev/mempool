package mempool

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

func TestPriorityQueue_PriorityOrdering(t *testing.T) {
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
			priorityQueue := NewPriorityQueue(100)
			for _, tx := range tc.inputTxs {
				priorityQueue.Push(tx)
			}

			output := []testPriorityItem{}
			for i := 0; i < len(tc.inputTxs); i++ {
				item := priorityQueue.Pop()
				testItem := item.(testPriorityItem)
				output = append(output, testItem)
			}

			require.Equal(t, tc.expectedOutput, output)
		})
	}
}


func TestPriorityQueue_Capacity(t *testing.T) {
	priorityQueue := NewPriorityQueue(10)
	for i := 0; i <  15; i++ {
		priorityQueue.Push(testPriorityItem{priority: float64(i)})
	}

	var expectedOutput []PriorityItem
	for i := 14; i >= 5; i-- {
		expectedOutput = append(expectedOutput, testPriorityItem{priority: float64(i)})
	}

	var output []PriorityItem
	for {
		if priorityQueue.Len() == 0 {
			break
		}
		item := priorityQueue.Pop()
		output = append(output, item)
	}

	require.Equal(t, expectedOutput, output)
}

func TestPriorityQueue_PoppingFromEmpty(t *testing.T) {
	priorityQueue := NewPriorityQueue(100)
	priorityQueue.Push(testPriorityItem{priority: float64(1)})
	priorityQueue.Push(testPriorityItem{priority: float64(2)})
	priorityQueue.Push(testPriorityItem{priority: float64(3)})

	require.Equal(t, 3, priorityQueue.Len())

	require.Equal(t, float64(3), priorityQueue.Pop().Priority())
	require.Equal(t, 2, priorityQueue.Len())
	require.Equal(t, float64(2), priorityQueue.Pop().Priority())
	require.Equal(t, 1, priorityQueue.Len())
	require.Equal(t, float64(1), priorityQueue.Pop().Priority())
	require.Equal(t, 0, priorityQueue.Len())

	require.Panics(t, func() {
		priorityQueue.Pop()
	})
}