package event

import (
	"context"
	"fmt"
	"sync"
)

const (
	Batch_size     = 5 // Maximum No Of Message will be available at a time
	Max_time_limit = 5 // Maximum Timelimit to run function - In minutes
)

type Event struct {
	Message string
}

func EventHandler(ctx context.Context, parentWg *sync.WaitGroup, eventData []Event) {
	defer parentWg.Done()

	if len(eventData) == 0 {
		fmt.Println("No event data to process.")
		return
	}

	// Create buffered channels to prevent blocking
	newFilter := NewFilter(Batch_size)

	var wg sync.WaitGroup

	// Immediate context check
	select {
	case <-ctx.Done():
		fmt.Println("Order Context Cancelled before processing")
		return
	default:
	}

	// Start filter worker
	wg.Add(1)
	go func() {
		defer wg.Done()
		newFilter.StartFiltering(ctx)
	}()

	total_send := 0
	// Process events
	wg.Add(1)
	go func() {
		defer func() {
			close(newFilter.rowMessages)
			wg.Done()
		}()

		for _, ev := range eventData {
			select {
			case <-ctx.Done():
				fmt.Println("--->Context Done is called during event processing!")
				return
			default:
				fmt.Println("Order Data from event:", ev)
				// Make CheckValidMessage context-aware
				if err := newFilter.CheckValidMessage(ctx, ev); err != nil {
					fmt.Println("Validation error:", err)
					return
				} else {
					total_send++
					fmt.Println("Total messages sent for validation:", total_send)
				}
			}
		}
	}()
	
	// Start upload workers
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			UploadFile(ctx, idx, newFilter.filteredMessages)
		}(i)
	}

	wg.Wait()
	fmt.Println("All events processed successfully.")
}
