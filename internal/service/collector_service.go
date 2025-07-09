package service

import (
	"context"
	"sync"
	"time"

	"github.com/dimryb/system-monitor/internal/collector"
	"github.com/dimryb/system-monitor/internal/entity"
	i "github.com/dimryb/system-monitor/internal/interface"
)

type CollectorService struct {
	ctx    context.Context
	log    i.Logger
	params []entity.Parameter
	wg     sync.WaitGroup
}

func NewCollectorService(ctx context.Context, log i.Logger) *CollectorService {
	return &CollectorService{ctx: ctx, log: log}
}

func (s *CollectorService) Start() error {
	errChan := make(chan error, len(s.params))

	for _, p := range s.params {
		coll := collector.NewCollector(time.Second)

		s.wg.Add(1)
		go func(param entity.Parameter, coll i.Collector) {
			defer s.wg.Done()

			ticker := time.NewTicker(time.Second)
			defer ticker.Stop()

			for {
				select {
				case <-s.ctx.Done():
					return
				case <-ticker.C:
					val, err := coll.Collect(s.ctx)
					if err != nil {
						errChan <- err
						continue
					}
					s.log.Infof("value: %s", val)
				}
			}

		}(p, coll)
	}

	go func() {
		s.wg.Wait()
		close(errChan)
	}()

	for err := range errChan {
		if err != nil {
			return err
		}
	}

	return nil
}
