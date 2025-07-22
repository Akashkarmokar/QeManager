package event

import (
	"fmt"
	"sync"
	"time"
)

func UploadFile(idx int, parentWg *sync.WaitGroup, validatedMessages chan Event) {
	defer parentWg.Done()
	for msg := range validatedMessages {
		fmt.Println("Uploading validated message:", msg, "at index:", idx)
		time.Sleep(1 * time.Second)
	}
	fmt.Println("No more messages to upload.", "at index:", idx)
}
