package eventconsummer

import (
	"drillCore/internal/events"
	"log"
	"time"
)

type Consumer struct {
	fetcher   events.Fetcher
	processor events.Processor
	batchSize int
}

func New(fetcher events.Fetcher, processor events.Processor, batchSize int) Consumer {
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
			log.Printf("Error fetching events: %v", err)

			continue
		}

		if len(gotEvents) == 0 {
			time.Sleep(1 * time.Second)

			continue
		}

		if err := c.handleEvents(gotEvents); err != nil {
			log.Println("Error handling events: ", err)

			continue
		}
	}
}

func (c *Consumer) handleEvents(event []events.Event) error {
	for _, e := range event {
		log.Printf("got new event: %s", e.Text)

		if err := c.processor.Process(e); err != nil {
			log.Printf("failed to handle event: %v", err.Error())

			continue
		}
	}

	return nil
}
