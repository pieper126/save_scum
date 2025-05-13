package internal_test

import (
	"os"
	"testing"

	"eu4_save_scum.com/saver/internal"
	"github.com/stretchr/testify/assert"
)

func TestLoadNotExistingConfig(t *testing.T) {
	_, err := internal.LoadConfig("./")
	assert.Error(t, err)
}

func TestLoadExistingConfig(t *testing.T) {
	file_name := "bla_test"
	os.Create(file_name)

	test_config := `{
	  "readFrom": "./",
	  "saveTo": "./"
	}`

	os.WriteFile(file_name, []byte(test_config), os.ModeAppend)

	_, err := internal.LoadConfig(file_name)
	assert.Nil(t, err)

	os.Remove(file_name)
}

func TestMalFormedConfig(t *testing.T) {
	file_name := "bla_test"
	os.Create(file_name)

	test_config := `{
	  "readFrom": "./",
	  "saveTo": "./"<>
	}`

	os.WriteFile(file_name, []byte(test_config), os.ModeAppend)

	_, err := internal.LoadConfig(file_name)
	assert.Error(t, err)

	os.Remove(file_name)
}

func TestIncorrectReadFrom(t *testing.T) {
	cfg := internal.Config{
		ReadFrom: "incorrect",
		SaveTo:   "./",
	}

	assert.Error(t, cfg.Validate())
}

func TestIncorrectSaveTo(t *testing.T) {
	cfg := internal.Config{
		ReadFrom: "./",
		SaveTo:   "incorrect",
	}

	assert.Error(t, cfg.Validate())
}

func TestCorrectSaveTo(t *testing.T) {
	cfg := internal.Config{
		ReadFrom: "./",
		SaveTo:   "./",
	}

	assert.Nil(t, cfg.Validate())
}
