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
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	var wg sync.WaitGroup

	wg.Add(1)
	go event.EventHandler(ctx, &wg, messages)

	time.Sleep(time.Second * 5)
	cancel()
	wg.Wait()
}
