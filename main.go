package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/Akashkarmokar/QeManager/internal/event"
)

const (
	batch_size     = 5 // Maximum No Of Message will be available at a time
	max_time_limit = 5 // Maximum Timelimit to run function
)

func main() {
	var messages []event.Event
	for i := range batch_size {
		messages = append(messages, event.Event{
			Message: "Message " + fmt.Sprintf("%d", i),
		})
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*(max_time_limit-2))
	var wg sync.WaitGroup

	wg.Add(1)
	go event.EventHandler(ctx, &wg, messages)

	time.Sleep(time.Second * (max_time_limit))
	cancel()
	wg.Wait()
}
