package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"eu4_save_scum.com/saver/internal"
	"github.com/c-bata/go-prompt"
)

type Commands string

const (
	LIST   Commands = "List"
	EXIT   Commands = "Exit"
	RESET  Commands = "Reset"
	BACKUP Commands = "Backup"
)

func main() {
	fmt.Println("Start saving failed games!")

	cfg, err := internal.LoadConfig("./config.json")
	if err != nil {
		fmt.Println(err)
		fmt.Println("press any key!")
		reader := bufio.NewReader(os.Stdin)
		reader.ReadString('\n')
		return
	}

	err = cfg.Validate()
	if err != nil {
		fmt.Println("%w", err)
		fmt.Println("press any key!")
		reader := bufio.NewReader(os.Stdin)
		reader.ReadString('\n')
		return
	}

	watcher := internal.BuildWatcher(cfg)

	saver, err := internal.BuildSaver(cfg)
	if err != nil {
		fmt.Println("press any key!")
		reader := bufio.NewReader(os.Stdin)
		reader.ReadString('\n')
		fmt.Println(err)
	}

	fmt.Println("Known save games: ")
	for _, v := range saver.ListSaveNames() {
		fmt.Println(v)
	}

	fmt.Printf("start watching for save games in: %s \n", cfg.ReadFrom)
	fmt.Printf("saving to: %s \n", cfg.SaveTo)

	filtersEvents := internal.BuildEu4Filter().Filter(watcher.Watch())

	saver.StartBackUp(filtersEvents)

	for {
		t := prompt.Input(">", completer(saver))
		switch t {
		case string(EXIT):
			fmt.Print("shutting down")
			saver.StopBackUp()
			return
		case string(LIST):
			fmt.Println("Known save games: ")
			for _, save := range saver.ListSaveNames() {
				fmt.Println(save)
			}
		default:
			fmt.Println("input was: ", t)

			if strings.Contains(t, string(RESET)) {
				command := strings.Split(t, " ")
				if len(command) != 3 {
					fmt.Println("incorrectly formatted command")
				}

				offset, _ := strconv.Atoi(command[2])
				fmt.Println("going back: ", offset)

				fileName := command[1]

				save, err := saver.Retrieve(fileName, offset)
				if err != nil {
					fmt.Println("something went wron resetting save game: ", err)
				}

				err = os.WriteFile(fmt.Sprintf("%s/%s", cfg.ReadFrom, fileName), save.Bytes, 0644)
				if err != nil {
					fmt.Println("error writing file: ", err)
				}

				fmt.Println("reset save game: ", fileName)
			}

			if strings.Contains(t, string(BACKUP)) {
				command := strings.Split(t, " ")
				if len(command) != 2 {
					fmt.Println("incorrectly formatted command")
				}

				fileName := command[1]

				fmt.Println("backing up: ", fileName)

				err := saver.BackUp(fileName)
				if err != nil {
					fmt.Println("error backing up file: ", fileName, err)
				}
			}
		}
	}
}

func completer(saver *internal.Saver) func(prompt.Document) []prompt.Suggest {
	return func(d prompt.Document) []prompt.Suggest {
		if strings.Contains(d.GetWordBeforeCursorUntilSeparator("\n"), string(RESET)) {
			s := []prompt.Suggest{}

			if strings.Contains(d.GetWordBeforeCursorUntilSeparator("\n"), ".eu4") {
				s = append(s, prompt.Suggest{Text: "0", Description: "amount of saves you want to go back 0 means latest"})
			} else {
				for _, save := range saver.ListSaveNames() {
					s = append(s, prompt.Suggest{Text: save, Description: save})
				}
			}
			return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
		}

		if strings.Contains(d.GetWordBeforeCursorUntilSeparator("\n"), string(BACKUP)) {
			s := []prompt.Suggest{}

			for _, save := range saver.ListSaveNames() {
				s = append(s, prompt.Suggest{Text: save, Description: save})
			}

			return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
		}

		s := []prompt.Suggest{
			{Text: string(LIST), Description: "List known save games"},
			{Text: string(EXIT), Description: "Exit"},
			{Text: string(RESET), Description: "Reset save game"},
			{Text: string(BACKUP), Description: "Back up savegame"},
		}

		return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
	}
}
