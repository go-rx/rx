package rx

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestOperatorConcatMap(test *testing.T) {
	fmt.Println("Testing operator ConcatMap...")
	project := func(x int) Observable[int] {
		return Map(
			Func(func(subscriber Writer[int]) (err error) {
				if !subscriber.Write(10) {
					return
				}
				fmt.Println("inner emitted 10")
				select {
				case <-time.NewTimer(100 * time.Millisecond).C:
				case <-subscriber.Dying():
					return
				}
				if !subscriber.Write(10) {
					return
				}
				fmt.Println("inner emitted 10")
				select {
				case <-time.NewTimer(100 * time.Millisecond).C:
				case <-subscriber.Dying():
					return
				}
				if !subscriber.Write(10) {
					return
				}
				fmt.Println("inner emitted 10")
				return
			}),
			func(i int) (int, error) {
				fmt.Printf("map emitted %d\n", i*x)
				return i * x, nil
			},
		)
	}
	source := Func(func(subscriber Writer[int]) (err error) {
		select {
		case <-time.NewTimer(200 * time.Millisecond).C:
		case <-subscriber.Dying():
			return
		}
		if !subscriber.Write(1) {
			return
		}
		fmt.Println("source emitted 1")
		select {
		case <-time.NewTimer(300 * time.Millisecond).C:
		case <-subscriber.Dying():
			return
		}
		if !subscriber.Write(3) {
			return
		}
		fmt.Println("source emitted 3")
		select {
		case <-time.NewTimer(150 * time.Millisecond).C:
		case <-subscriber.Dying():
			return
		}
		if !subscriber.Write(5) {
			return
		}
		fmt.Println("source emitted 5")
		return
	})
	writer, reader := Pipe[int](nil)
	start := time.Now()
	ConcatMap(source, project).Subscribe(writer)
	v, ok := reader.Read()
	assert.Equal(test, 10, v)
	assert.True(test, ok)
	assert.WithinDuration(test, start.Add(200*time.Millisecond), time.Now(), 50*time.Millisecond)
	v, ok = reader.Read()
	assert.Equal(test, 10, v)
	assert.True(test, ok)
	assert.WithinDuration(test, start.Add(300*time.Millisecond), time.Now(), 50*time.Millisecond)
	v, ok = reader.Read()
	assert.Equal(test, 10, v)
	assert.True(test, ok)
	assert.WithinDuration(test, start.Add(400*time.Millisecond), time.Now(), 50*time.Millisecond)
	v, ok = reader.Read()
	assert.Equal(test, 30, v)
	assert.True(test, ok)
	assert.WithinDuration(test, start.Add(500*time.Millisecond), time.Now(), 50*time.Millisecond)
	v, ok = reader.Read()
	assert.Equal(test, 30, v)
	assert.True(test, ok)
	assert.WithinDuration(test, start.Add(600*time.Millisecond), time.Now(), 50*time.Millisecond)
	v, ok = reader.Read()
	assert.Equal(test, 30, v)
	assert.True(test, ok)
	assert.WithinDuration(test, start.Add(700*time.Millisecond), time.Now(), 50*time.Millisecond)
	v, ok = reader.Read()
	assert.Equal(test, 50, v)
	assert.True(test, ok)
	assert.WithinDuration(test, start.Add(700*time.Millisecond), time.Now(), 80*time.Millisecond)
	v, ok = reader.Read()
	assert.Equal(test, 50, v)
	assert.True(test, ok)
	assert.WithinDuration(test, start.Add(800*time.Millisecond), time.Now(), 80*time.Millisecond)
	v, ok = reader.Read()
	assert.Equal(test, 50, v)
	assert.True(test, ok)
	assert.WithinDuration(test, start.Add(900*time.Millisecond), time.Now(), 80*time.Millisecond)
	err := reader.Wait()
	assert.Nil(test, err)
}
