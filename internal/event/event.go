package event

import (
	"context"
	"fmt"
	"sync"
)

type Event struct {
	Message string
}

func EventHandler(ctx context.Context, wg *sync.WaitGroup, eventData []Event) {
	defer wg.Done()

	// If there's no data, return immediately
	if len(eventData) == 0 {
		fmt.Println("No event data to process.")
		return
	}
	// newFilter, err := filter.NewFilter()
	for len(eventData) > 0 {
		select {
		case <-ctx.Done():
			fmt.Println("Context Done is called during event processing!")
			return
		default:
			fmt.Println("Data from event:", eventData[0])
			eventData = eventData[1:]
			// Optional: simulate processing delay
			// time.Sleep(1 * time.Second)
		}
	}

	fmt.Println("All events processed successfully.")
}
