package event_consumer

import (
	"link-saver-bot/clients/events"
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

func (consumer Consumer) Start() error {

	for {

		eventsFetched, err := consumer.fetcher.Fetch(consumer.batchSize, 0)

		if err != nil {
			log.Printf("[ERROR] consumer: %s", err)
			continue
			//@todo implement retry logic at fetcher before we will have to skip this here
		}

		if len(eventsFetched) == 0 {
			time.Sleep(1 * time.Second)

			continue
		}

		if err := consumer.handleEvents(eventsFetched); err != nil {
			log.Print(err)
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
		log.Printf("got new event: %s", event.Text)

		if err := consumer.processor.Process(event); err != nil {
			log.Printf("error while handling event: %s", err.Error())
			continue
		}

	}

	return nil
}
