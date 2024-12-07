package aof

import (
	"bufio"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/MridulDhiman/chronotable/config"
)

type AOF struct {
	MainPath   string
	File   *os.File
	Writer *bufio.Writer
	VersionToPath map[int]string
	SeekCurrent int64
}

func New(_path string, initialized bool) *AOF {
	// open file in append, write only mode and create as well, if not created.
	// File has user permissions set: 6(rw-)4(r--)4(r--)
	_path = filepath.Join("./", config.CHRONO_MAIN_DIR, _path)
	flags := os.O_WRONLY|os.O_APPEND;
	fmt.Println("initialized", initialized)
	if !initialized {
		flags = os.O_APPEND|os.O_CREATE|os.O_WRONLY
	}
	file, err := os.OpenFile(_path,flags, fs.FileMode(0777))
	if err != nil {
		log.Fatal("Could not open file: ", err)
	}
	return &AOF{
		File:   file,
		MainPath:   _path,
		Writer: bufio.NewWriter(file),
		VersionToPath: make(map[int]string),
	}
}

func (aof *AOF) Log(operation string) error {
	fmt.Println("operation: ", operation)

	defer func(aof* AOF) {
		ptr, _ := aof.File.Seek(0, io.SeekCurrent)
		aof.SeekCurrent = ptr
	}(aof)

	if _, err := aof.Writer.WriteString(operation + "\n"); err != nil {
		fmt.Println("could not write to file")
		return err
	}
	if err := aof.Writer.Flush(); err != nil {
		return err
	}

	return aof.File.Sync()
}

func (aof *AOF) Replay(start, end int64) error  {
	file, err := os.OpenFile(filepath.Join("./", config.CHRONO_MAIN_DIR, config.MAIN_AOF_FILE), os.O_RDONLY, fs.FileMode(0644))
	if err != nil {
		log.Fatal("Could not open file: ", err)
	}
	defer file.Close()
// move the file pointer to start
	if _, err:= file.Seek(start, io.SeekStart); err != nil {
		log.Fatal("(error) could not rewind file pointer: ", err)
	}
	scanner:= bufio.NewScanner(file)
 	var delta int64 = end - start
	// why does file pointer location keep shifting by few bytes
	var error_rate int64 = 10
	for scanner.Scan() {
		if delta <= error_rate {
			break;
		}
		line := scanner.Text()
		delta -= int64(len(line))
		fmt.Println("Line: ", line);
	}
	
return nil
}

func Format(key string, value any) string {
	return fmt.Sprintf("Key: %s, Value: %v, Timestamp: %v", key, value, time.Now().UTC().Format(time.RFC3339Nano)) //ISO format with nanoseconds for AOF logs
}


