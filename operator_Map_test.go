package rx

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOperatorMap(test *testing.T) {
	fmt.Println("Testing operator Map...")
	writer, reader := Pipe[string]()
	Map(
		List([]int{1, 2, 3}),
		func(value int) (string, error) {
			return fmt.Sprintf("number %d", value), nil
		},
	).Subscribe(writer)
	v, ok := reader.Read()
	assert.Equal(test, "number 1", v)
	assert.True(test, ok)
	v, ok = reader.Read()
	assert.Equal(test, "number 2", v)
	assert.True(test, ok)
	v, ok = reader.Read()
	assert.Equal(test, "number 3", v)
	assert.True(test, ok)
	v, ok = reader.Read()
	assert.Equal(test, "", v)
	assert.False(test, ok)
	err := reader.Wait()
	assert.Nil(test, err)
}
