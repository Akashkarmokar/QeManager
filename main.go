// Example usage in main.go or any handler
// package main

// import (
// 	"context"
// 	"fmt"
// 	"time"

// 	"github.com/Akashkarmokar/QeManager/internal/aws"
// )

// func main() {
// 	ctx := context.Background()
// 	s3util, err := aws.NewS3Util(ctx, "teststreamtoupload", "ap-south-1")
// 	if err != nil {
// 		fmt.Println("S3 init error:", err)
// 		return
// 	}

// 	key := "4782376-uhd_3840_2160_30fps.mp4"
// 	expiry := 15 * time.Minute

// 	getURL, err := s3util.PresignGetURL(ctx, key, expiry)
// 	if err != nil {
// 		fmt.Println("Failed to get presigned GET URL:", err)
// 		return
// 	}
// 	fmt.Println("Presigned GET URL:", getURL)

// 	putURL, err := s3util.PresignPutURL(ctx, key, expiry)
// 	if err != nil {
// 		fmt.Println("Failed to get presigned PUT URL:", err)
// 		return
// 	}
// 	fmt.Println("Presigned PUT URL:", putURL)
// }

// package main

// import (
// 	"context"
// 	"fmt"
// 	"io"
// 	"os"

// 	"github.com/Akashkarmokar/QeManager/internal/aws"
// )

// func main() {
// 	ctx := context.Background()
// 	s3util, err := aws.NewS3Util(ctx, "teststreamtoupload", "ap-south-1")
// 	if err != nil {
// 		fmt.Println("S3 init error:", err)
// 		return
// 	}

// 	key := "4782376-uhd_3840_2160_30fps.mp4"      // S3 object key
// 	localPath := "./downloaded-file.mp4" // Local file path

// 	// Stream from S3
// 	reader, err := s3util.StreamObject(ctx, key)
// 	if err != nil {
// 		fmt.Println("Failed to stream object:", err)
// 		return
// 	}
// 	defer reader.Close()

// 	// Create local file
// 	outFile, err := os.Create(localPath)
// 	if err != nil {
// 		fmt.Println("Failed to create local file:", err)
// 		return
// 	}
// 	defer outFile.Close()

// 	// Copy S3 stream to local file
// 	written, err := io.Copy(outFile, reader)
// 	if err != nil {
// 		fmt.Println("Failed to copy data:", err)
// 		return
// 	}

// 	fmt.Printf("Downloaded %d bytes from S3 to %s\n", written, localPath)
// }

// package main

// import (
// 	"context"
// 	"fmt"
// 	"io"
// 	"os"
// 	"time"

// 	"github.com/Akashkarmokar/QeManager/internal/aws"
// )

// func main() {
// 	ctx := context.Background()
// 	s3util, err := aws.NewS3Util(ctx, "teststreamtoupload", "ap-south-1")
// 	if err != nil {
// 		fmt.Println("S3 init error:", err)
// 		return
// 	}

// 	key := "4782376-uhd_3840_2160_30fps.mp4" // S3 object key
// 	localPath := "./downloaded-file.mp4"     // Local file path

// 	start := time.Now() // Start timing

// 	// Stream from S3
// 	reader, err := s3util.StreamObject(ctx, key)
// 	if err != nil {
// 		fmt.Println("Failed to stream object:", err)
// 		return
// 	}
// 	defer reader.Close()

// 	// Create local file
// 	outFile, err := os.Create(localPath)
// 	if err != nil {
// 		fmt.Println("Failed to create local file:", err)
// 		return
// 	}
// 	defer outFile.Close()

// 	// Copy S3 stream to local file
// 	written, err := io.Copy(outFile, reader)
// 	if err != nil {
// 		fmt.Println("Failed to copy data:", err)
// 		return
// 	}

// 	elapsed := time.Since(start) // End timing

//		fmt.Printf("Downloaded %d bytes from S3 to %s\n", written, localPath)
//		fmt.Printf("Time taken: %s\n", elapsed)
//	}
package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/Akashkarmokar/QeManager/internal/event"
)

func main() {
	var messages []event.Event
	for i := 0; i < event.Batch_size; i++ {
		messages = append(messages, event.Event{
			Message: "Message " + fmt.Sprintf("%d", i),
		})
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*event.Max_time_limit)
	var wg sync.WaitGroup

	wg.Add(1)
	go event.EventHandler(ctx, &wg, messages)

	time.Sleep(time.Minute * (event.Max_time_limit + 2))
	cancel()
	wg.Wait()
}
