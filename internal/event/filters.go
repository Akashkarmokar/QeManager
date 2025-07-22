package event

type Filter struct {
	filteredMessages chan Event
}

func NewFilter() *Filter {
	return &Filter{
		filteredMessages: make(chan Event, Batch_size),
	}
}

func (f *Filter) CheckValidMessage(event Event) {
	// Primarily no filtering logic is implemented here.
	// This is just a placeholder for future filtering logic.
	f.filteredMessages <- event
}
