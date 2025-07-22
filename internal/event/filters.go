package event

import (
    "context"
    "fmt"
    "sync"
)

type Filter struct {
    rowMessages      chan Event
    filteredMessages chan Event
}

func NewFilter() *Filter {
    return &Filter{
        rowMessages:      make(chan Event, Batch_size),
        filteredMessages: make(chan Event, Batch_size),
    }
}

func (f *Filter) CheckValidMessage(event Event) {
    // Placeholder for future filtering logic.
    f.rowMessages <- event
}

func (f *Filter) StartFiltering(ctx context.Context, wg *sync.WaitGroup) {
    defer wg.Done()
    defer close(f.filteredMessages)
    for {
        select {
        case <-ctx.Done():
            return
        case msg, ok := <-f.rowMessages:
            if !ok {
                fmt.Println("rowMessages channel closed, closing filteredMessages")
                return
            }
            f.filteredMessages <- msg
        }
    }
}