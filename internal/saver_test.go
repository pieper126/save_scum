package internal_test

import (
	"os"
	"testing"
	"time"

	"eu4_save_scum.com/saver/internal"
	"github.com/stretchr/testify/assert"
)

func TestAbleToSaveOneSave(t *testing.T) {
	cfg := internal.Config{
		SaveTo: "./save_to",
	}
	defer os.RemoveAll(cfg.SaveTo)

	saver, err := internal.BuildSaver(cfg)
	assert.Nil(t, err)
	assert.NotNil(t, saver)

	events := make(chan internal.Event, 100)
	saver.StartBackUp(events)

	file_name := "blas.eu4"
	fs, err := os.Create(file_name)
	assert.Nil(t, err)

	fs.Write([]byte("something new is happening"))
	fs.Close()
	defer os.Remove(file_name)

	events <- internal.Event{EventType: internal.CREATE, FileName: file_name}

	time.Sleep(10 * time.Millisecond)

	entries, err := os.ReadDir(cfg.SaveTo)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(entries))

	saved_file_name := "blas_0.eu4"
	saved_path := cfg.SaveTo + "/" + saved_file_name
	stat, err := os.Stat(saved_path)
	assert.Nil(t, err, err)
	assert.NotNil(t, stat)

	assert.Equal(t, saved_file_name, stat.Name())

	saved, err := os.ReadFile(saved_path)
	assert.Nil(t, err)
	og, err := os.ReadFile(file_name)
	assert.Nil(t, err)

	assert.Equal(t, og, saved)
}

func TestOnlyListensToCreates(t *testing.T) {
	cfg := internal.Config{
		SaveTo: "./save_to",
	}
	defer os.RemoveAll(cfg.SaveTo)

	saver, err := internal.BuildSaver(cfg)
	assert.Nil(t, err)
	assert.NotNil(t, saver)

	events := make(chan internal.Event, 100)
	saver.StartBackUp(events)

	file_name := "blas.eu4"
	fs, err := os.Create(file_name)
	assert.Nil(t, err)

	fs.Write([]byte("something new is happening"))
	fs.Close()
	defer os.Remove(file_name)

	events <- internal.Event{EventType: internal.CREATE, FileName: file_name}
	events <- internal.Event{EventType: internal.DELETE, FileName: file_name}
	events <- internal.Event{EventType: internal.UPDATE, FileName: file_name}
	events <- internal.Event{EventType: internal.DELETE, FileName: file_name}
	events <- internal.Event{EventType: internal.UPDATE, FileName: file_name}
	events <- internal.Event{EventType: internal.CREATE, FileName: file_name}

	time.Sleep(10 * time.Millisecond)

	entries, err := os.ReadDir(cfg.SaveTo)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(entries))

	saved_file_name := "blas_0.eu4"
	saved_path := cfg.SaveTo + "/" + saved_file_name
	stat, err := os.Stat(saved_path)
	assert.Nil(t, err, err)
	assert.NotNil(t, stat)

	assert.Equal(t, saved_file_name, stat.Name())

	saved, err := os.ReadFile(saved_path)
	assert.Nil(t, err)
	og, err := os.ReadFile(file_name)
	assert.Nil(t, err)

	assert.Equal(t, og, saved)
}

func TestEveryGameHasItsOwnVersioning(t *testing.T) {
	cfg := internal.Config{
		SaveTo: "./save_to",
	}
	defer os.RemoveAll(cfg.SaveTo)

	saver, err := internal.BuildSaver(cfg)
	assert.Nil(t, err)
	assert.NotNil(t, saver)

	events := make(chan internal.Event, 100)
	saver.StartBackUp(events)

	file_name := "blas.eu4"
	fs, err := os.Create(file_name)
	assert.Nil(t, err)

	fs.Write([]byte("something new is happening"))
	fs.Close()
	defer os.Remove(file_name)

	events <- internal.Event{EventType: internal.CREATE, FileName: file_name}
	events <- internal.Event{EventType: internal.DELETE, FileName: file_name}
	events <- internal.Event{EventType: internal.UPDATE, FileName: file_name}
	events <- internal.Event{EventType: internal.DELETE, FileName: file_name}
	events <- internal.Event{EventType: internal.UPDATE, FileName: file_name}
	events <- internal.Event{EventType: internal.CREATE, FileName: file_name}

	time.Sleep(10 * time.Millisecond)

	entries, err := os.ReadDir(cfg.SaveTo)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(entries))

	saved_file_name := "blas_0.eu4"
	saved_path := cfg.SaveTo + "/" + saved_file_name
	stat, err := os.Stat(saved_path)
	assert.Nil(t, err, err)
	assert.NotNil(t, stat)

	assert.Equal(t, saved_file_name, stat.Name())

	saved, err := os.ReadFile(saved_path)
	assert.Nil(t, err)
	og, err := os.ReadFile(file_name)
	assert.Nil(t, err)

	assert.Equal(t, og, saved)
}

func TestAbleToRetrieveLatest(t *testing.T) {
	cfg := internal.Config{
		SaveTo: "./save_to",
	}
	defer os.RemoveAll(cfg.SaveTo)

	saver, err := internal.BuildSaver(cfg)
	assert.Nil(t, err)
	assert.NotNil(t, saver)

	events := make(chan internal.Event, 100)
	saver.StartBackUp(events)

	file_name := "blas.eu4"
	fs, err := os.Create(file_name)
	assert.Nil(t, err)

	fs.Write([]byte("1"))
	fs.Close()
	defer os.Remove(file_name)

	events <- internal.Event{EventType: internal.CREATE, FileName: file_name}

	fs1, err := os.OpenFile(file_name, os.O_APPEND|os.O_RDWR, os.ModeAppend)
	assert.Nil(t, err)

	fs1.Write([]byte("2"))
	fs1.Close()
	events <- internal.Event{EventType: internal.CREATE, FileName: file_name}

	time.Sleep(10 * time.Millisecond)

	save, err := saver.Latest(file_name)
	assert.Nil(t, err)
	assert.NotNil(t, save)

	blas := string(save.Bytes)
	assert.Equal(t, "12", blas)
}

func TestAbleToRetrieveLatestFor2Saves(t *testing.T) {
	cfg := internal.Config{
		SaveTo: "./save_to",
	}
	defer os.RemoveAll(cfg.SaveTo)

	saver, err := internal.BuildSaver(cfg)
	assert.Nil(t, err)
	assert.NotNil(t, saver)

	events := make(chan internal.Event, 100)
	saver.StartBackUp(events)

	file_name := "blas.eu4"
	second_file_name := "blas2.eu4"

	fs, err := os.Create(file_name)
	assert.Nil(t, err)
	fs.Write([]byte("1"))
	fs.Close()
	defer os.Remove(file_name)

	events <- internal.Event{EventType: internal.CREATE, FileName: file_name}

	second, err := os.Create(second_file_name)
	assert.Nil(t, err)
	second.Write([]byte("3"))
	second.Close()
	defer os.Remove(second_file_name)
	events <- internal.Event{EventType: internal.CREATE, FileName: second_file_name}

	fs1, err := os.OpenFile(file_name, os.O_APPEND|os.O_RDWR, os.ModeAppend)
	assert.Nil(t, err)
	fs1.Write([]byte("2"))
	fs1.Close()
	events <- internal.Event{EventType: internal.CREATE, FileName: file_name}

	time.Sleep(10 * time.Millisecond)

	save, err := saver.Latest(file_name)
	assert.Nil(t, err)
	assert.NotNil(t, save)

	content := string(save.Bytes)
	assert.Equal(t, "12", content)

	save2, err := saver.Latest(second_file_name)
	assert.Nil(t, err)
	assert.NotNil(t, save)

	content2 := string(save2.Bytes)
	assert.Equal(t, "3", content2)
}

func TestLatestWhichIsUknown(t *testing.T) {
	cfg := internal.Config{
		SaveTo: "./save_to",
	}
	defer os.RemoveAll(cfg.SaveTo)

	saver, err := internal.BuildSaver(cfg)
	assert.Nil(t, err)
	assert.NotNil(t, saver)

	events := make(chan internal.Event, 100)
	saver.StartBackUp(events)

	_, err = saver.Latest("unkown")
	assert.NotNil(t, err)
}

func TestAbleRetrieveOneBeforeLatest(t *testing.T) {
	cfg := internal.Config{
		SaveTo: "./save_to",
	}
	defer os.RemoveAll(cfg.SaveTo)

	saver, err := internal.BuildSaver(cfg)
	assert.Nil(t, err)
	assert.NotNil(t, saver)

	events := make(chan internal.Event, 100)
	saver.StartBackUp(events)

	file_name := "blas.eu4"
	second_file_name := "blas2.eu4"

	fs, err := os.Create(file_name)
	assert.Nil(t, err)
	fs.Write([]byte("1"))
	fs.Close()
	defer os.Remove(file_name)

	events <- internal.Event{EventType: internal.CREATE, FileName: file_name}

	second, err := os.Create(second_file_name)
	assert.Nil(t, err)
	second.Write([]byte("3"))
	second.Close()
	defer os.Remove(second_file_name)
	events <- internal.Event{EventType: internal.CREATE, FileName: second_file_name}

	time.Sleep(10 * time.Millisecond)

	fs1, err := os.OpenFile(file_name, os.O_APPEND|os.O_RDWR, os.ModeAppend)
	assert.Nil(t, err)
	fs1.Write([]byte("2"))
	fs1.Close()
	events <- internal.Event{EventType: internal.CREATE, FileName: file_name}

	time.Sleep(10 * time.Millisecond)

	save, err := saver.Retrieve(file_name, 1)
	assert.Nil(t, err)
	assert.NotNil(t, save)

	content := string(save.Bytes)
	assert.Equal(t, "1", content)
}

func TestAbleToListKnownSaveGames(t *testing.T) {
	cfg := internal.Config{
		SaveTo: "./save_to",
	}
	defer os.RemoveAll(cfg.SaveTo)

	saver, err := internal.BuildSaver(cfg)
	assert.Nil(t, err)
	assert.NotNil(t, saver)

	events := make(chan internal.Event, 100)
	saver.StartBackUp(events)

	file_name := "blas_to_make_it_harder.eu4"
	second_file_name := "blas2.eu4"

	fs, err := os.Create(file_name)
	assert.Nil(t, err)
	fs.Write([]byte("1"))
	fs.Close()
	defer os.Remove(file_name)

	events <- internal.Event{EventType: internal.CREATE, FileName: file_name}

	second, err := os.Create(second_file_name)
	assert.Nil(t, err)
	second.Write([]byte("3"))
	second.Close()
	defer os.Remove(second_file_name)
	events <- internal.Event{EventType: internal.CREATE, FileName: second_file_name}

	time.Sleep(10 * time.Millisecond)

	fs1, err := os.OpenFile(file_name, os.O_APPEND|os.O_RDWR, os.ModeAppend)
	assert.Nil(t, err)
	fs1.Write([]byte("2"))
	fs1.Close()
	events <- internal.Event{EventType: internal.CREATE, FileName: file_name}

	time.Sleep(10 * time.Millisecond)

	saveNames := saver.ListSaveNames()
	assert.NotNil(t, saveNames)

	assert.Len(t, saveNames, 2)
	assert.ElementsMatch(t, []string{file_name, second_file_name}, saveNames)
}

func TestDealsWithRestartCorrectly(t *testing.T) {
	cfg := internal.Config{
		SaveTo: "./save_to",
	}
	defer os.RemoveAll(cfg.SaveTo)

	saver, err := internal.BuildSaver(cfg)
	assert.Nil(t, err)
	assert.NotNil(t, saver)

	events := make(chan internal.Event, 100)
	saver.StartBackUp(events)

	file_name := "bengal_test_for_code.eu4"

	fs, err := os.Create(file_name)
	assert.Nil(t, err)
	fs.Write([]byte("1"))
	fs.Close()
	defer os.Remove(file_name)

	events <- internal.Event{EventType: internal.CREATE, FileName: file_name}

	time.Sleep(10 * time.Millisecond)

	saver.StopBackUp()

	time.Sleep(10 * time.Millisecond)

	saver, err = internal.BuildSaver(cfg)
	assert.Nil(t, err)
	assert.NotNil(t, saver)

	saver.StartBackUp(events)
	defer saver.StopBackUp()

	fs1, err := os.OpenFile(file_name, os.O_APPEND|os.O_RDWR, os.ModeAppend)
	assert.Nil(t, err)
	fs1.Write([]byte("2"))
	fs1.Close()
	events <- internal.Event{EventType: internal.CREATE, FileName: file_name}

	time.Sleep(10 * time.Millisecond)

	s, err := saver.Retrieve(file_name, 1)
	assert.Nil(t, err)
	assert.Equal(t, "1", string(s.Bytes))

	s, err = saver.Retrieve(file_name, 0)
	assert.Nil(t, err)
	assert.Equal(t, "12", string(s.Bytes))

	assert.ElementsMatch(t, saver.ListSaveNames(), []string{file_name})
}

func TestBackUp(t *testing.T) {
	cfg := internal.Config{
		SaveTo: "./save_to",
	}
	defer os.RemoveAll(cfg.SaveTo)

	saver, err := internal.BuildSaver(cfg)
	assert.Nil(t, err)
	assert.NotNil(t, saver)

	file_name := "bengal_test_for_code.eu4"

	fs, err := os.Create(file_name)
	assert.Nil(t, err)
	fs.Write([]byte("1"))
	fs.Close()
	defer os.Remove(file_name)

	err = saver.BackUp(file_name)
	assert.Nil(t, err)

	err = os.WriteFile(file_name, []byte("2"), 0644)
	assert.Nil(t, err)

	err = saver.BackUp(file_name)
	assert.Nil(t, err)

	save, err := saver.Retrieve(file_name, 0)
	assert.Nil(t, err)
	assert.Equal(t, "2", string(save.Bytes))

	save, err = saver.Retrieve(file_name, 1)
	assert.Nil(t, err)
	assert.Equal(t, "1", string(save.Bytes))
}
