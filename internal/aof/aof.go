package aof

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/MridulDhiman/chronotable/config"
)

type AOF struct {
	MainPath      string
	File          *os.File
	Writer        *bufio.Writer
	VersionToPath map[int]string
}

func New(_path string, initialized bool) *AOF {
	// open file in append, write only mode and create as well, if not created.
	// File has user permissions set: 6(rw-)4(r--)4(r--)
	_path = filepath.Join("./", config.CHRONO_MAIN_DIR, _path)
	flags := os.O_WRONLY | os.O_APPEND
	fmt.Println("initialized", initialized)
	if !initialized {
		flags = os.O_APPEND | os.O_CREATE | os.O_WRONLY
	}
	file, err := os.OpenFile(_path, flags, fs.FileMode(0777))
	if err != nil {
		log.Fatal("Could not open file: ", err)
	}
	return &AOF{
		File:          file,
		MainPath:      _path,
		Writer:        bufio.NewWriter(file),
		VersionToPath: make(map[int]string),
	}
}

func (aof *AOF) Log(operation string) error {
	fmt.Println("operation: ", operation)

	if _, err := aof.Writer.WriteString(operation + "\n"); err != nil {
		fmt.Println("could not write to file")
		return err
	}
	if err := aof.Writer.Flush(); err != nil {
		return err
	}

	return aof.File.Sync()
}

// Scan AOF line by line
// json unmarshal line to *Operation
func (aof *AOF) MustReplay() []*Operation {
	// TODO: add file descripter as dependency
	file, err := os.OpenFile(filepath.Join("./", config.CHRONO_MAIN_DIR, config.MAIN_AOF_FILE), os.O_RDONLY, fs.FileMode(0644))
	if err != nil {
		log.Fatal("Could not open file: ", err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	var operations []*Operation
	for scanner.Scan() {
		line := scanner.Text()
		// TODO: add decoder as an dependency
		operation := new(Operation)
		if err := json.Unmarshal([]byte(line), operation); err != nil {
			panic(err)
		}
		operations = append(operations, operation)
	}
	return operations
}

// clears the AOF log upon commit
func (aof *AOF) Clear() error {
	if err := aof.File.Close(); err != nil {
		return err
	}
	if err := os.Truncate(aof.MainPath, 0); err != nil {
		return err
	}
	f, err := os.OpenFile(aof.MainPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		log.Fatalf("Unable to open log file '%v'", aof.MainPath)
		return err
	}
	aof.File = f
	return nil
}

func MustFormat(key string, value any, operationType OperationType) string {
	operation := Operation{
		Key:           key,
		Value:         value,
		OperationType: operationType,
		Timestamp:     time.Now().UTC().Format(time.RFC3339Nano), //ISO format with nanoseconds for AOF logs
	}

	// TODO: use encoder as an dependency
	data, err := json.Marshal(operation)
	if err != nil {
		panic(err)
	}

	return string(data)
}
