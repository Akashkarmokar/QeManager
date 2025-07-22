package event

import (
	"context"
	"fmt"
	"time"
)

func UploadFile(ctx context.Context, idx int, validatedMessages <-chan Event) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Order context cancel on upload queue")
			return
		case msg, ok := <-validatedMessages:
			if !ok {
				fmt.Println("No more messages to upload at index:", idx)
				return
			}

			fmt.Println("Uploading validated message:", msg, "at index:", idx)

			// Context-aware sleep
			select {
			case <-ctx.Done():
				fmt.Println("Order context canceled during Order process")
				return
			case <-time.After(3 * time.Second):
				fmt.Println("Order Updated successfully", msg, "at index", idx)
			}
		}
	}
}
