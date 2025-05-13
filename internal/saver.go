package internal

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

type Saver struct {
	saveTo string
	latest map[string][]string
	done   chan struct{}
}

type Save struct {
	Bytes []byte
}

func BuildSaver(cfg Config) (*Saver, error) {
	if _, err := os.Stat(cfg.SaveTo); err != nil {
		err := os.Mkdir(cfg.SaveTo, 0755)
		if err != nil {
			return nil, err
		}
	}

	saver := &Saver{saveTo: cfg.SaveTo, latest: map[string][]string{}, done: make(chan struct{})}

	err := saver.resetLatest()
	if err != nil {
		return nil, err
	}

	return saver, nil
}

func (s *Saver) StartBackUp(events <-chan Event) {
	go func() {
		for {
			select {
			case event := <-events:
				if event.EventType != CREATE {
					continue
				}

				s.BackUp(event.FileName)
			case <-s.done:
				return
			}
		}
	}()
}

func (s *Saver) BackUp(save string) error {
	copied_name := s.updateLatest(save)

	fs, err := os.Create(copied_name)
	if err != nil {
		fmt.Printf("err creating copied  for %s, %s", fs.Name(), err)
		return err
	}
	defer fs.Close()

	og, err := os.ReadFile(save)
	if err != nil {
		fmt.Printf("err opening original file for %s, %s", fs.Name(), err)
		return err
	}

	_, err = fs.Write(og)
	if err != nil {
		fmt.Printf("err copying file for %s, %s", fs.Name(), err)
		return err
	}

	return nil
}

func (s *Saver) Latest(saveName string) (Save, error) {
	return s.Retrieve(saveName, 0)
}

func (s *Saver) Retrieve(saveName string, offset int) (Save, error) {
	fileNames, ok := s.latest[saveName]
	if !ok {
		return Save{}, fmt.Errorf("unknown save name: %s", saveName)
	}

	fileName := fileNames[len(fileNames)-1-offset]

	bytes, err := os.ReadFile(fileName)
	if err != nil {
		return Save{}, fmt.Errorf("unable to read latest: %s", fileName)
	}

	return Save{Bytes: bytes}, nil
}

func (s *Saver) ListSaveNames() []string {
	keys := make([]string, 0, len(s.latest))
	for k := range s.latest {
		keys = append(keys, k)
	}
	return keys
}

func (s *Saver) StopBackUp() {
	s.done <- struct{}{}
}

func (s *Saver) resetLatest() error {
	s.latest = map[string][]string{}

	entries, err := os.ReadDir(s.saveTo)
	if err != nil {
		return err
	}

	sort.Slice(entries, func(x, y int) bool {
		return entries[x].Name() < entries[y].Name()
	})

	for _, entry := range entries {
		dirs := strings.Split(entry.Name(), "/")
		file_name := dirs[len(dirs)-1]
		name_and_suffix := strings.Split(file_name, ".")
		name_with_number := name_and_suffix[0]

		for i := len(name_with_number) - 1; i > 0; i-- {
			if name_with_number[i] == '_' {
				name := fmt.Sprintf("%s.%s", name_with_number[:i], name_and_suffix[1])
				s.updateLatest(name)
				break
			}
		}

	}

	return nil
}

func (s *Saver) updateLatest(fileName string) string {
	fileNames, ok := s.latest[fileName]
	if !ok {
		fileNames = []string{}
	}

	splitted := strings.Split(fileName, ".")
	copied_name := fmt.Sprintf("%s/%s_%d.%s", s.saveTo, splitted[0], len(fileNames), splitted[1])
	s.latest[fileName] = append(fileNames, copied_name)

	return copied_name
}
