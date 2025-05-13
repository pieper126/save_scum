package internal_test

import (
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"eu4_save_scum.com/saver/internal"
	"github.com/stretchr/testify/assert"
)

func TestAbleToCreateWatcher(t *testing.T) {
	read_folder := "./read_folder"
	write_folder := "./write_folder"

	os.Create(read_folder)
	os.Create(write_folder)

	cfg := internal.Config{
		ReadFrom: read_folder,
		SaveTo:   write_folder,
	}

	watcher := internal.BuildWatcher(cfg)
	assert.NotNil(t, watcher)

	os.Remove(read_folder)
	os.Remove(write_folder)
}

func TestForONe(t *testing.T) {
	read_folder := "./read_folder/"
	write_folder := "./write_folder/"

	err := os.Mkdir(read_folder, 0755)
	if err != nil {
		assert.Fail(t, fmt.Sprintf("%s", err))
	}
	defer os.RemoveAll(read_folder)
	err = os.Mkdir(write_folder, 0755)
	if err != nil {
		assert.Fail(t, fmt.Sprintf("%s", err))
	}
	defer os.RemoveAll(write_folder)

	cfg := internal.Config{
		ReadFrom: read_folder,
		SaveTo:   write_folder,
	}
	watcher := internal.BuildWatcher(cfg)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		event := <-watcher.Watch()
		assert.NotNil(t, event)
		watcher.Done()
		wg.Done()
	}()

	file_to_be_created := read_folder + "created"

	go func() {
		for i := 0; i < 10; i++ {
			f, err := os.Create(fmt.Sprintf("%s_%d", file_to_be_created, i))
			if err != nil {
				fmt.Print(err)
				assert.NotNil(t, err)
			}
			f.Close()
			time.Sleep(1 * time.Second)
		}
	}()
	wg.Wait()

	os.RemoveAll(read_folder)
	os.RemoveAll(write_folder)
}

func TestAbleToRead(t *testing.T) {
	read_folder := "./read_folder/"
	write_folder := "./write_folder/"

	err := os.Mkdir(read_folder, 0755)
	if err != nil {
		assert.Fail(t, fmt.Sprintf("%s", err))
	}
	err = os.Mkdir(write_folder, 0755)
	if err != nil {
		assert.Fail(t, fmt.Sprintf("%s", err))
	}

	cfg := internal.Config{
		ReadFrom: read_folder,
		SaveTo:   write_folder,
	}
	watcher := internal.BuildWatcher(cfg)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		event := <-watcher.Watch()
		assert.NotNil(t, event)
		watcher.Done()
		wg.Done()
	}()

	file_to_be_created := read_folder + "created"

	go func() {
		for i := 0; i < 10; i++ {
			f, err := os.Create(fmt.Sprintf("%s_%d", file_to_be_created, i))
			if err != nil {
				fmt.Print(err)
				assert.NotNil(t, err)
			}
			f.Close()
			time.Sleep(1 * time.Second)
		}
	}()
	wg.Wait()

	os.RemoveAll(read_folder)
	os.RemoveAll(write_folder)
}

func TestAbleSeeCreate(t *testing.T) {
	read_folder := "read_folder/"
	write_folder := "write_folder/"

	err := os.Mkdir(read_folder, 0755)
	if err != nil {
		assert.Fail(t, fmt.Sprintf("%s", err))
	}
	defer os.RemoveAll(read_folder)
	err = os.Mkdir(write_folder, 0755)
	if err != nil {
		assert.Fail(t, fmt.Sprintf("%s", err))
	}
	defer os.RemoveAll(write_folder)

	cfg := internal.Config{
		ReadFrom: read_folder,
		SaveTo:   write_folder,
	}
	file_to_be_created := read_folder + "created"

	watcher := internal.BuildWatcher(cfg)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		event := <-watcher.Watch()
		assert.Equal(t, internal.Event{EventType: internal.CREATE, FileName: file_to_be_created}, event)
		watcher.Done()
		wg.Done()
	}()

	go func() {
		f, err := os.Create(file_to_be_created)
		if err != nil {
			fmt.Print(err)
			assert.NotNil(t, err)
		}
		f.Close()
	}()
	wg.Wait()

}

func TestAbleToSeeMultiple(t *testing.T) {
	read_folder := "read_folder/"
	write_folder := "write_folder/"

	err := os.Mkdir(read_folder, 0755)
	if err != nil {
		assert.Fail(t, fmt.Sprintf("%s", err))
	}
	defer os.RemoveAll(read_folder)
	err = os.Mkdir(write_folder, 0755)
	if err != nil {
		assert.Fail(t, fmt.Sprintf("%s", err))
	}
	defer os.RemoveAll(write_folder)

	cfg := internal.Config{
		ReadFrom: read_folder,
		SaveTo:   write_folder,
	}
	file_to_be_created := read_folder + "created"

	watcher := internal.BuildWatcher(cfg)
	var wg sync.WaitGroup
	wg.Add(1)

	received := make([]internal.Event, 0)

	go func() {
		for event := range watcher.Watch() {
			received = append(received, event)
		}
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		for i := 0; i < 2; i++ {
			f, err := os.Create(fmt.Sprintf("%s_%d", file_to_be_created, i))
			if err != nil {
				fmt.Print(err)
				assert.NotNil(t, err)
			}
			f.Close()
		}
		time.Sleep(1 * time.Millisecond)
		watcher.Done()
		wg.Done()
	}()
	wg.Wait()

	assert.Equal(t, 2, len(received))
}

func TestAbleToSeeWriteToFile(t *testing.T) {
	read_folder := "read_folder/"
	write_folder := "write_folder/"

	err := os.Mkdir(read_folder, 0755)
	if err != nil {
		assert.Fail(t, fmt.Sprintf("%s", err))
	}
	defer os.RemoveAll(read_folder)
	err = os.Mkdir(write_folder, 0755)
	if err != nil {
		assert.Fail(t, fmt.Sprintf("%s", err))
	}
	defer os.RemoveAll(write_folder)

	cfg := internal.Config{
		ReadFrom: read_folder,
		SaveTo:   write_folder,
	}
	file_to_be_created := read_folder + "created"

	watcher := internal.BuildWatcher(cfg)
	var wg sync.WaitGroup
	wg.Add(1)

	received := make([]internal.Event, 0)

	go func() {
		for event := range watcher.Watch() {
			received = append(received, event)
		}
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		defer watcher.Done()
		defer wg.Done()

		f, err := os.Create(file_to_be_created)
		if err != nil {
			fmt.Print(err)
			assert.NotNil(t, err)
		}
		f.Close()

		again, err := os.OpenFile(file_to_be_created, os.O_RDWR, 0644)
		if err != nil {
			fmt.Print(err)
			assert.NotNil(t, err)
			return
		}

		_, err = again.Write([]byte("something new"))
		if err != nil {
			fmt.Print(err)
			assert.NotNil(t, err)
			return
		}
		again.Close()

		time.Sleep(1 * time.Millisecond)
	}()
	wg.Wait()

	assert.Equal(t, 2, len(received))
	assert.Equal(t, internal.CREATE, received[0].EventType)
	assert.Equal(t, internal.UPDATE, received[1].EventType)
}
