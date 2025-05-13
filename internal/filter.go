package internal

import "strings"

type Eu4Filter struct {
}

func BuildEu4Filter() *Eu4Filter {
	return &Eu4Filter{}
}

func (f *Eu4Filter) Filter(incoming <-chan Event) <-chan Event {
	res := make(chan Event, 10)
	go func() {
		for event := range incoming {
			if strings.Contains(event.FileName, ".tmp") {
				continue
			}

			if strings.Contains(event.FileName, "_Backup.eu4") {
				continue
			}

			res <- event
		}
		close(res)
	}()

	return res
}
