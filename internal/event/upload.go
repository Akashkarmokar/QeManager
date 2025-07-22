package event

import (
	"fmt"
	"sync"
)

func UploadFile(idx int, parentWg *sync.WaitGroup, validatedMessages chan Event) {
	defer parentWg.Done()
	for msg := range validatedMessages {
		fmt.Println("Uploading validated message:", msg, "at index:", idx)
	}
	fmt.Println("No more messages to upload.", "at index:", idx)
}
