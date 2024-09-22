package list

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestListStorage(t *testing.T) {
	type testingType struct {
		name string
	}

	tests := []testingType{
		{
			name: "Data 1",
		},
		{
			name: "Data 2",
		},
	}

	list := NewNodeList()
	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list.AddRecord(tt)
			result := list.GetRecords(3)
			firstItem := result[0]

			switch item := firstItem.(type) {
			case testingType:
				require.Equal(t, tt.name, item.name)
			default:
				t.Error("No expect type")
			}
			require.Len(t, result, i+1)
		})
	}
}
