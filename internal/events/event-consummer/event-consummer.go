package eventconsummer

import (
	"context"
	"errors"
	"time"

	"drillCore/internal/events"
	eventprocessor "drillCore/internal/events/event-processor"

	"go.uber.org/zap"
)

type Consumer struct {
	fetcher   events.Fetcher
	processor events.Processor
	batchSize int

	logger *zap.SugaredLogger
}

func New(fetcher events.Fetcher, processor events.Processor, batchSize int, logger *zap.SugaredLogger) Consumer {
	return Consumer{
		fetcher:   fetcher,
		processor: processor,
		batchSize: batchSize,
		logger:    logger,
	}
}

func (c *Consumer) Start(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			gotEvents, err := c.fetcher.Fetch(ctx, c.batchSize)
			if err != nil {
				if errors.Is(err, eventprocessor.ErrNoUpdatesFound) {
					time.Sleep(1 * time.Second)

					continue
				}

				c.logger.Errorw("failed to fetch events", "error", err)

				continue
			}

			if err := c.handleEvents(ctx, gotEvents); err != nil {
				c.logger.Errorw("failed to handle events", "error", err)

				continue
			}
		}
	}
}

func (c *Consumer) handleEvents(ctx context.Context, event []*events.Event) error {
	for _, e := range event {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			c.logger.Infow("processing event", "event", e)

			if err := c.processor.Process(ctx, e); err != nil {
				c.logger.Errorw("failed to process event", "event", e, "error", err)

				continue
			}
		}
	}

	return nil
}
