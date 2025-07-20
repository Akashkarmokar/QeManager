package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

const (
	batch_size     = 5 // Maximum No Of Message will be available at a time
	max_time_limit = 5 // Maximum Timelimit to run function
)

type Event struct {
	message string
}

func EvenHandler(ctx context.Context, wg *sync.WaitGroup, eventData []Event) {
	defer wg.Done()

	// If there's no data, return immediately
	if len(eventData) == 0 {
		fmt.Println("No event data to process.")
		return
	}

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

func main() {
	var messages []Event
	for i := range batch_size {
		messages = append(messages, Event{
			message: "Message " + fmt.Sprintf("%d", i),
		})
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*(max_time_limit-2))
	var wg sync.WaitGroup

	wg.Add(1)
	go EvenHandler(ctx, &wg, messages)

	time.Sleep(time.Second * (max_time_limit))
	cancel()
	wg.Wait()
}
