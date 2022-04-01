package rx

func BufferCount[T any](source Observable[T], bufferSize int, options ...BufferCountOption) Observable[[]T] {
	return Func(func(subscriber Writer[[]T]) (err error) {
		opts := buildOptions(options)

		if opts.startBufferEvery == 0 {
			opts.startBufferEvery = bufferSize
		}

		var buffers [][]T
		count := 0

		writer, reader := Pipe[T](subscriber)
		source.Subscribe(writer)

		defer func() {
			for _, buffer := range buffers {
				if !subscriber.Write(buffer) {
					return
				}
			}
		}()

		for {
			var ok bool
			var value T
			var toEmit [][]T

			if value, ok = reader.Read(); !ok {
				return
			}

			if count%opts.startBufferEvery == 0 {
				buffers = append(buffers, nil)
			}
			count++

			var newBuffers [][]T
			for _, buffer := range buffers {
				buffer = append(buffer, value)
				if bufferSize <= len(buffer) {
					toEmit = append(toEmit, buffer)
				} else {
					newBuffers = append(newBuffers, buffer)
				}
			}
			buffers = newBuffers

			for _, buffer := range toEmit {
				if !subscriber.Write(buffer) {
					return
				}
			}
		}
	})
}

type BufferCountOption func(*bufferCountOptions)

type bufferCountOptions struct {
	startBufferEvery int
}
