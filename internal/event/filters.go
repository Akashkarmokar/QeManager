package event

import (
	"context"
	"fmt"
)

type Filter struct {
	rowMessages      chan Event
	filteredMessages chan Event
}

func NewFilter(bufferSize int) *Filter {
	return &Filter{
		rowMessages:      make(chan Event, bufferSize),
		filteredMessages: make(chan Event, bufferSize),
	}
}

func (f *Filter) CheckValidMessage(ctx context.Context, event Event) error {
	select {
	case f.rowMessages <- event:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (f *Filter) StartFiltering(ctx context.Context) {
	defer close(f.filteredMessages)
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Filtering canceled by context")
			return
		case msg, ok := <-f.rowMessages:
			if !ok {
				return
			}
			select {
			case f.filteredMessages <- msg:
			case <-ctx.Done():
				return
			}
		}
	}
}
