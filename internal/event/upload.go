package event

import (
	"context"
	"errors"
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

			// Create a new context with timeout for the upload operation
			uploadCtx, cancel := context.WithTimeout(ctx, 30*time.Second) // Set reasonable timeout
			defer cancel()

			// Execute the upload with context support
			err := uploadToThirdParty(uploadCtx, msg)
			if err != nil {
				if errors.Is(err, context.Canceled) {
					fmt.Println("Upload canceled:", err)
					return
				}
				if errors.Is(err, context.DeadlineExceeded) {
					fmt.Println("Upload timeout:", err)
					continue // or return depending on your requirements
				}
				fmt.Println("Upload failed:", err)
				continue
			}

			fmt.Println("Order Updated successfully", msg, "at index", idx)
		}
	}
}

// Example third-party upload function with context support
func uploadToThirdParty(ctx context.Context, msg Event) error {
	// Simulate upload work that respects context
	done := make(chan struct{})
	var err error
	time.Sleep(3 * time.Second)
	go func() {
		// Actual upload logic here
		// This could be an HTTP request, SDK call, etc.
		// err = thirdPartyClient.Upload(ctx, msg.Message)
		close(done)
	}()

	select {
	case <-done:
		return err
	case <-ctx.Done():
		// Cleanup resources if possible
		return ctx.Err()
	}
}
