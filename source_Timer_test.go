package rx

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSourceTimer(test *testing.T) {
	fmt.Println("Testing source Timer...")
	writer, reader := Pipe[int](nil)
	start := time.Now()
	Timer(time.Second, 100*time.Millisecond).Subscribe(writer)
	v, ok := reader.Read()
	assert.Equal(test, 0, v)
	assert.True(test, ok)
	assert.WithinDuration(test, start.Add(time.Second), time.Now(), 50*time.Millisecond)
	v, ok = reader.Read()
	assert.Equal(test, 1, v)
	assert.True(test, ok)
	assert.WithinDuration(test, start.Add(1100*time.Millisecond), time.Now(), 50*time.Millisecond)
	v, ok = reader.Read()
	assert.Equal(test, 2, v)
	assert.True(test, ok)
	assert.WithinDuration(test, start.Add(1200*time.Millisecond), time.Now(), 50*time.Millisecond)
	v, ok = reader.Read()
	assert.Equal(test, 3, v)
	assert.True(test, ok)
	assert.WithinDuration(test, start.Add(1300*time.Millisecond), time.Now(), 50*time.Millisecond)
	reader.Kill(context.Canceled)
	v, ok = reader.Read()
	assert.Equal(test, 0, v)
	assert.False(test, ok)
	assert.WithinDuration(test, start.Add(1300*time.Millisecond), time.Now(), 50*time.Millisecond)
	err := reader.Wait()
	assert.Equal(test, context.Canceled, err)
}
