package event

import (
	"context"
	"fmt"
	"sync"
)

const (
	Batch_size     = 5 // Maximum No Of Message will be available at a time
	Max_time_limit = 5 // Maximum Timelimit to run function
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

	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	newFilter := NewFilter()

	wg.Add(1)
	go newFilter.StartFiltering(ctx, &wg)

	wg.Add(1)
	go func() {
		defer func() {
			close(newFilter.rowMessages)
			wg.Done()
		}()
		for _, ev := range eventData {
			select {
			case <-ctx.Done():
				fmt.Println("Context Done is called during event processing!")
				return
			default:
				fmt.Println("Data from event:", ev)
				newFilter.CheckValidMessage(ev)
			}
		}
	}()

	workerCount := Batch_size
	for i := 0; i < workerCount+3; i++ {
		wg.Add(1)
		go UploadFile(i, &wg, newFilter.filteredMessages)
	}

	wg.Wait()
	fmt.Println("All events processed successfully.")
}
