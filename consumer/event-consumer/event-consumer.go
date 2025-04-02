package event_consumer

import (
	"link-saver-bot/clients/events"
	"log/slog"
	"time"
)

type Consumer struct {
	fetcher   events.Fetcher
	processor events.Processor
	batchSize int
}

func New(fetcher events.Fetcher, processor events.Processor, batchSize int) Consumer {
	slog.Info("Initializing event consumer", "batchSize", batchSize)
	return Consumer{
		fetcher:   fetcher,
		processor: processor,
		batchSize: batchSize,
	}
}

func (consumer Consumer) Start() error {
	slog.Info("Starting event consumer")

	for {
		eventsFetched, err := consumer.fetcher.Fetch(consumer.batchSize, 0)
		if err != nil {
			slog.Error("Failed to fetch events",
				"batchSize", consumer.batchSize,
				"error", err)
			continue
		}

		if len(eventsFetched) == 0 {
			time.Sleep(1 * time.Second)
			continue
		}

		slog.Info("Fetched new events", "count", len(eventsFetched))

		if err := consumer.handleEvents(eventsFetched); err != nil {
			slog.Error("Failed to handle events", "error", err)
			continue
		}
	}
}

/*
*
*
  - potantial problems here:
  - 1) loosing of events due to network errors
  - solutions: retry, fallback with put event to temporary storage or ram, confirmation from fetcher that event is processed
  - 2) continue to process next event even if problem was not in particular event (for intance, network issue)
    solutons: stop entire batch on first failure, error counter with limit
    3) performance issue when iterating using just single process
    solutions: using sync.WaitGroup{}
*/
func (consumer *Consumer) handleEvents(events []events.Event) error {
	for _, event := range events {
		slog.Info("Processing event",
			"text", event.Text,
			"meta", event.Meta)

		if err := consumer.processor.Process(event); err != nil {
			slog.Error("Failed to process event",
				"text", event.Text,
				"meta", event.Meta,
				"error", err)
			continue
		}

		slog.Debug("Event processed successfully",
			"text", event.Text,
			"meta", event.Meta)
	}

	return nil
}
