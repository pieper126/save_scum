package internal_test

import (
	"testing"

	"eu4_save_scum.com/saver/internal"
	"github.com/stretchr/testify/assert"
)

func TestAbleToCreateEu4Filter(t *testing.T) {
	filter := internal.BuildEu4Filter()
	assert.NotNil(t, filter)
}

func TestAbleToFilterTmpFiles(t *testing.T) {
	filter := internal.BuildEu4Filter()
	assert.NotNil(t, filter)

	in := make(chan internal.Event, 10)
	out := filter.Filter(in)

	in <- internal.Event{EventType: internal.CREATE, FileName: "tmp1.tmp"}
	in <- internal.Event{EventType: internal.CREATE, FileName: "tmp2.tmp"}
	in <- internal.Event{EventType: internal.CREATE, FileName: "tmp3.tmp"}
	in <- internal.Event{EventType: internal.CREATE, FileName: "real.eu4"}
	close(in)

	events := []internal.Event{}
	for event := range out {
		events = append(events, event)
	}

	assert.Len(t, events, 1)
}

func TestFiltersBackUps(t *testing.T) {
	filter := internal.BuildEu4Filter()
	assert.NotNil(t, filter)

	in := make(chan internal.Event, 10)
	out := filter.Filter(in)

	in <- internal.Event{EventType: internal.CREATE, FileName: "bengal_test_for_Code_Backup.eu4"}
	in <- internal.Event{EventType: internal.CREATE, FileName: "bengal_test_for_Code_Backup.eu4"}
	in <- internal.Event{EventType: internal.CREATE, FileName: "bengal_test_for_Code_Backup.eu4"}
	in <- internal.Event{EventType: internal.CREATE, FileName: "bengal_test_for_Code_Backup.eu4"}
	in <- internal.Event{EventType: internal.CREATE, FileName: "real.eu4"}
	close(in)

	events := []internal.Event{}
	for event := range out {
		events = append(events, event)
	}

	assert.Len(t, events, 1)
}
