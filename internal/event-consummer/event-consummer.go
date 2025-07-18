package eventconsummer

import (
	"drillCore/internal/events/telegram"
	"drillCore/internal/session"
	"errors"
	"log"
	"time"
)

type Consumer struct {
	fetcher   session.Fetcher
	processor session.Processor
	batchSize int
}

func New(fetcher session.Fetcher, processor session.Processor, batchSize int) Consumer {
	return Consumer{
		fetcher:   fetcher,
		processor: processor,
		batchSize: batchSize,
	}
}

func (c *Consumer) Start() error {
	for {
		gotEvents, err := c.fetcher.Fetch(c.batchSize)
		if err != nil {
			if errors.Is(err, telegram.ErrNoUpdatesFound) {
				time.Sleep(1 * time.Second)

				continue
			}

			log.Printf("Error fetching events: %v", err)

			continue
		}

		if err := c.handleEvents(gotEvents); err != nil {
			log.Println("Error handling events: ", err)

			continue
		}
	}
}

func (c *Consumer) handleEvents(event []session.Event) error {
	for _, e := range event {
		log.Printf("got new event: %s", e.Text)

		if err := c.processor.Process(e); err != nil {
			log.Printf("failed to handle event: %v", err.Error())

			continue
		}
	}

	return nil
}
