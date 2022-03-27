package rx

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOperatorTake(test *testing.T) {
	fmt.Println("Testing operator Take...")
	writer, reader := Pipe[int](nil)
	Take(List([]int{1, 2, 3}), 2).Subscribe(writer)
	v, ok := reader.Read()
	assert.Equal(test, 1, v)
	assert.True(test, ok)
	v, ok = reader.Read()
	assert.Equal(test, 2, v)
	assert.True(test, ok)
	v, ok = reader.Read()
	assert.Equal(test, 0, v)
	assert.False(test, ok)
	err := reader.Wait()
	assert.Nil(test, err)
}
