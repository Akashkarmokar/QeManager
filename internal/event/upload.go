package event

import (
	"fmt"
	"sync"
)

func UploadFile(idx int, parentWg *sync.WaitGroup, validatedMessages chan Event) {
	defer parentWg.Done()

	// for {
	select {
	case msg := <-validatedMessages:
		// Process the validated message
		// For example, upload it to a server or save it to a database
		fmt.Println("Uploading validated message:", msg, "at index:", idx)
	default:
		// If no messages are available, break the loop
		fmt.Println("No more messages to upload.", "at index:", idx)
		return
	}
	// }
}
