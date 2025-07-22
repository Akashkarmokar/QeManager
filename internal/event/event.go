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

	// If there's no data, return immediately
	if len(eventData) == 0 {
		fmt.Println("No event data to process.")
		return
	}
	newFilter := NewFilter()
	for len(eventData) > 0 {
		select {
		case <-ctx.Done():
			fmt.Println("Context Done is called during event processing!")
			return
		default:
			fmt.Println("Data from event:", eventData[0])
			newFilter.CheckValidMessage(eventData[0])
			eventData = eventData[1:]
			// Optional: simulate processing delay
			// time.Sleep(1 * time.Second)
		}
	}
	close(newFilter.filteredMessages)
	var wg sync.WaitGroup

	if len(newFilter.filteredMessages) != 0 {
		fmt.Println("Filtered messages are available for further processing.")
		workerCount := len(newFilter.filteredMessages)

		for i := 0; i < workerCount; i++ {
			wg.Add(1)
			go UploadFile(i, &wg, newFilter.filteredMessages)
		}
	}

	fmt.Println("All events processed successfully.")
	wg.Wait()
}

// func UploadFile(validatedMessages chan Event) {
// 	for {
// 		select {
// 		case msg := <-validatedMessages:
// 			// Process the validated message
// 			// For example, upload it to a server or save it to a database
// 			fmt.Println("Uploading validated message:", msg)
// 		default:
// 			// If no messages are available, break the loop
// 			fmt.Println("No more messages to upload.")
// 			return
// 		}
// 	}
// }
