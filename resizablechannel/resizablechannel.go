package resizablechannel

import "github.com/PhanNam1501/Fluid/queue"

const INF = int(^uint(0) >> 1)

type ResizableChannel[V any] struct {
	input                    chan *V
	output                   chan *V
	buffer                   *queue.Queue[V]
	length, capacity, resize chan int
	size                     int
}

func NewResizableChannel[V any]() *ResizableChannel[V] {
	c := &ResizableChannel[V]{
		input:    make(chan *V),
		output:   make(chan *V),
		buffer:   queue.NewQueue[V](),
		length:   make(chan int),
		capacity: make(chan int),
		resize:   make(chan int),
		size:     1,
	}
	go c.dynamicChange()
	return c
}

func (c *ResizableChannel[V]) Input() chan<- *V {
	return c.input
}

func (c *ResizableChannel[V]) Output() <-chan *V {
	return c.output
}

func (c *ResizableChannel[V]) Len() int {
	return <-c.length
}

func (c *ResizableChannel[V]) Capacity() int {
	cap, ok := <-c.capacity
	if ok {
		return cap
	} else {
		return c.size
	}
}

func (c *ResizableChannel[V]) Close() {
	close(c.input)
}

func (c *ResizableChannel[V]) Resize(newSize int) {
	if newSize <= 0 && newSize != INF {
		panic("channels: invalid negative size trying to resize channel")
	}
	c.resize <- newSize
}

func (c *ResizableChannel[V]) dynamicChange() {
	var (
		next   *V
		output chan *V
	)
	nextInput := c.input
	input := nextInput
	for input != nil || output != nil {
		select {
		case elem, open := <-input:
			if open {
				c.buffer.Push(elem)
			} else {
				input = nil
				nextInput = nil
			}
		case output <- next:
			c.buffer.Pop()
		case c.size = <-c.resize:
		case c.length <- c.buffer.Length():
		case c.capacity <- c.size:
		}

		if c.buffer.Length() == 0 {
			output = nil
			next = nil
		} else {
			next = c.buffer.Head()
			output = c.output
		}

		if c.size != INF && c.size <= c.buffer.Length() {
			input = nil
		} else {
			input = nextInput
		}
	}

	close(c.length)
	close(c.output)
	close(c.capacity)
}
