package tests

import (
	"reflect"
	"testing"

	"github.com/mmycin/ndc14/libs"
)

func TestReverseArray(t *testing.T) {
	// Test case for reverseArray function
	t.Run("ReverseArray", func(t *testing.T) {
		arr := []int{1, 2, 3, 4, 5}
		libs.ReverseArray(&arr)
		expected := []int{5, 4, 3, 2, 1}
		if !reflect.DeepEqual(arr, expected) {
			t.Errorf("Expected %v but got %v", expected, arr)
		}
	})
}
