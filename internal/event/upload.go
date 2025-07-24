package event

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/Akashkarmokar/QeManager/internal/aws"
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

			fmt.Println("Uploading validated message:", msg, "at index:----->", idx, len(validatedMessages))

			// Create a new context with timeout for the upload operation
			uploadCtx, cancel := context.WithTimeout(ctx, Max_time_limit*time.Minute) // Set reasonable timeout
			defer cancel()

			// Execute the upload with context support
			err := uploadToThirdParty(uploadCtx, msg, idx)
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
// func uploadToThirdParty(ctx context.Context, msg Event, idx int) error {
// 	// Simulate upload work that respects context
// 	// done := make(chan struct{})
// 	// var err error
// 	// time.Sleep(3 * time.Second)
// 	// go func() {
// 	// 	// Actual upload logic here
// 	// 	// This could be an HTTP request, SDK call, etc.
// 	// 	// err = thirdPartyClient.Upload(ctx, msg.Message)
// 	// 	close(done)
// 	// }()

// 	// select {
// 	// case <-done:
// 	// 	return err
// 	// case <-ctx.Done():
// 	// 	// Cleanup resources if possible
// 	// 	return ctx.Err()
// 	nctx, cancel := context.WithCancel(ctx)
// 	defer cancel()
// 	s3util, err := aws.NewS3Util(nctx, "teststreamtoupload", "ap-south-1")
// 	if err != nil {
// 		fmt.Println("S3 init error:", err)
// 		return fmt.Errorf("failed to initialize S3 util ->: %w", err)
// 	}

// 	key := "4782376-uhd_3840_2160_30fps.mp4"                           // S3 object key
// 	localPath := "./downloaded-file" + fmt.Sprintf("%d", idx) + ".mp4" // Local file path

// 	start := time.Now() // Start timing

// 	// Stream from S3
// 	reader, err := s3util.StreamObject(nctx, key)
// 	if err != nil {
// 		fmt.Println("Failed to stream object:", err)
// 		return fmt.Errorf("failed to initialize S3 util **: %w", err)
// 	}
// 	defer reader.Close()

// 	// Create local file
// 	outFile, err := os.Create(localPath)
// 	if err != nil {
// 		fmt.Println("Failed to create local file:", err)
// 		return fmt.Errorf("failed to initialize S3 util &&: %w and file: %s", err, localPath)
// 	}
// 	defer outFile.Close()

// 	// Copy S3 stream to local file
// 	written, err := io.Copy(outFile, reader)
// 	if err != nil {
// 		fmt.Println("Failed to copy data:", err)
// 		return fmt.Errorf("failed to initialize S3 util ##: %w and file: %s", err, localPath)
// 	}

// 	elapsed := time.Since(start) // End timing

// 	fmt.Printf("Downloaded %d bytes from S3 to %s\n", written, localPath)
// 	fmt.Printf("Time taken: %s\n", elapsed)
// 	return nil
// }

func uploadToThirdParty(ctx context.Context, msg Event, idx int) error {
	s3util, err := aws.NewS3Util(ctx, "teststreamtoupload", "ap-south-1")
	if err != nil {
		return fmt.Errorf("failed to initialize S3 util: %w", err)
	}
	fmt.Println("Message to upload:", msg.Message)
	s3Key := "4782376-uhd_3840_2160_30fps.mp4" // Use the message as the S3 key (or set as needed)
	localPath := fmt.Sprintf("./downloaded-file-%d.mp4", idx)

	start := time.Now()

	// Stream from S3
	reader, err := s3util.StreamObject(ctx, s3Key)
	if err != nil {
		return fmt.Errorf("failed to stream object from S3: %w", err)
	}
	defer reader.Close()

	// Create local file
	outFile, err := os.Create(localPath)
	if err != nil {
		return fmt.Errorf("failed to create local file %s: %w", localPath, err)
	}
	defer outFile.Close()

	// Copy S3 stream to local file
	written, err := io.Copy(outFile, reader)
	if err != nil {
		return fmt.Errorf("failed to copy data to %s: %w", localPath, err)
	}

	elapsed := time.Since(start)
	fmt.Printf("Downloaded %d bytes from S3 key %s to %s in %s\n", written, s3Key, localPath, elapsed)
	return nil
}
