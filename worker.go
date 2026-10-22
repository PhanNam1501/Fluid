package main

import (
	"context"
	"fmt"
	"sync"

	"github.com/rs/zerolog/log"
)

type worker struct {
	id          int
	handler     Handler
	messageChan <-chan []byte
	err         error
	errChan     chan<- *worker
	wg          *sync.WaitGroup
	ctx         context.Context
}

func (w *worker) Work(ctx context.Context) {
	defer func() {
		if r := recover(); r != nil {
			if err, ok := r.(error); ok {
				w.err = err
			} else {
				w.err = fmt.Errorf("worker panicked: %v", r)
			}
			w.errChan <- w
		} else {
			log.Info().Msgf("worker %d exits normally", w.id)
			w.wg.Done()
		}
	}()

	for {
		select {
		case <-ctx.Done():
			log.Info().Msgf("worker %d is closing", w.id)
			return
		default:
			msg, ok := <-w.messageChan
			if !ok {
				log.Info().Msgf("worker %d is closing", w.id)
				return
			}
			log.Info().Msgf("worker %d received message and processing it", w.id)

			err := w.handler.Handle(ctx, msg)
			if err != nil {
				log.Error().Err(err).Msgf("worker %d failed to handle message", w.id)
			} else {
				log.Info().Msgf("worker %d handled message", w.id)
			}
		}
	}

}

func (p *Pool) Close() {
	close(p.messageChan)
	p.wg.Wait()
	close(p.errChan)
}
