package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

type Pool struct {
	handler     Handler
	numWorkers  int
	messageChan chan []byte
	errChan     chan *worker
	wg          *sync.WaitGroup
	cancelChan  chan context.CancelFunc
}

func NewPool(ctx context.Context, handler Handler, numWorkers int) *Pool {
	messageChan := make(chan []byte, numWorkers)
	errChan := make(chan *worker, numWorkers)
	pool := &Pool{
		handler:     handler,
		numWorkers:  numWorkers,
		messageChan: messageChan,
		errChan:     errChan,
		wg:          &sync.WaitGroup{},
		cancelChan:  make(chan context.CancelFunc),
	}
	pool.init(ctx)

	fmt.Printf("starting worker pool with %d workers", numWorkers)
	return pool
}

func (p *Pool) Consume(message []byte) {
	p.messageChan <- message
}

func (p *Pool) init(ctx context.Context) {
	for i := 1; i <= p.numWorkers; i++ {
		worker := &worker{
			id:          i,
			handler:     p.handler,
			messageChan: p.messageChan,
			errChan:     p.errChan,
			wg:          p.wg,
		}

		worker.wg.Add(1)
		ctxC, cancel := context.WithCancel(ctx)
		p.cancelChan <- cancel
		go worker.Work(ctxC)
	}

	go func() {
		for w := range p.errChan {
			log.Error().Msgf("worker %d encounter a panic", w.id)
			time.Sleep(2 * time.Second)
			log.Info().Msgf("worker %d restarting", w.id)
			ctxC, cancel := context.WithCancel(ctx)
			p.cancelChan <- cancel
			go w.Work(ctxC)
		}
	}()
}

func (p *Pool) RemoveWorker() {
	fn := <-p.cancelChan
	fn()
}
