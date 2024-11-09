package aof

import (
	"bufio"
	"fmt"
	"io/fs"
	"log"
	"os"
	"time"
)

type AOF struct {
	path string
	file *os.File
	writer *bufio.Writer
}

func New(path string) *AOF {
	// open file in append, write only mode and create as well, if not created.
	// File has user permissions set: 6(rw-)4(r--)4(r--)
	file, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, fs.FileMode(0644))
	if err != nil {
		log.Fatal("Could not open file: ",err)
	}
	return &AOF{
		file: file,
		path : path,
		writer: bufio.NewWriter(file),
	}
}

func (aof* AOF) Log(operation string) error {
	fmt.Println("operation: ", operation)
	if _, err := aof.writer.WriteString(operation + "\n"); err != nil {
		fmt.Println("could not write to file")
		return err
	}
	if err := aof.writer.Flush(); err != nil {
		return err
	}

	return aof.file.Sync()
}

func Format(key string, value any) string {
	return fmt.Sprintf("Key: %s, Value: %v, Timestamp: %v", key, value, time.Now().UTC().Format(time.RFC3339Nano)) //ISO format with nanoseconds for AOF logs
}