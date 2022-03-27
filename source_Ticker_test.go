package rx

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSourceTicker(test *testing.T) {
	fmt.Println("Testing source Ticker...")
	writer, reader := Pipe[int](nil)
	start := time.Now()
	Ticker(100 * time.Millisecond).Subscribe(writer)
	v, ok := reader.Read()
	assert.Equal(test, 0, v)
	assert.True(test, ok)
	assert.WithinDuration(test, start.Add(100*time.Millisecond), time.Now(), 50*time.Millisecond)
	v, ok = reader.Read()
	assert.Equal(test, 1, v)
	assert.True(test, ok)
	assert.WithinDuration(test, start.Add(200*time.Millisecond), time.Now(), 50*time.Millisecond)
	v, ok = reader.Read()
	assert.Equal(test, 2, v)
	assert.True(test, ok)
	assert.WithinDuration(test, start.Add(300*time.Millisecond), time.Now(), 50*time.Millisecond)
	reader.Kill(context.Canceled)
	v, ok = reader.Read()
	assert.Equal(test, 0, v)
	assert.False(test, ok)
	assert.WithinDuration(test, start.Add(300*time.Millisecond), time.Now(), 50*time.Millisecond)
	err := reader.Wait()
	assert.Equal(test, context.Canceled, err)
}
