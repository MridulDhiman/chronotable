package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"syscall"

	"github.com/MridulDhiman/chronotable/chronotable"
	"github.com/MridulDhiman/chronotable/config"
	"github.com/MridulDhiman/chronotable/internal/utils"
)

var initialized bool = false

func init() {
	newpath := filepath.Join(".", config.CHRONO_MAIN_DIR)
	// check if directory exists or not
	yes, err := utils.Exists(newpath)
	if err != nil {
		log.Fatalln(err)
	}

	if !yes {
		os.MkdirAll(newpath, os.FileMode(0755))
		// Set hidden attribute on Windows

		if runtime.GOOS == "windows" {
			ptr, err := syscall.UTF16PtrFromString(newpath)
			if err != nil {
				panic(err)
			}

			var attrs uint32
			attrs, err = syscall.GetFileAttributes(ptr)
			if err != nil {
				panic(err)
			}

			if err := syscall.SetFileAttributes(ptr, attrs|syscall.FILE_ATTRIBUTE_HIDDEN); err != nil {
				panic(err)
			}

		}

	} else {
		initialized = true
	}
}

func main() {
	logger := log.New(os.Stdin, config.DefaultLoggerPrefix, log.LstdFlags)
	table := chronotable.New(&chronotable.Options{
		EnableAOF:      true,
		AOFPath:        config.MAIN_AOF_FILE,
		EnableSnapshot: true,
		Initialized:    initialized,
		Logger:         logger,
	})

	if initialized {
		// get current state by replaying logs
		currentVersion, latestVersion, err := table.ConfigHandler.FetchLatestVersion()
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println("latest version: ", latestVersion)
		fmt.Println("current version: ", currentVersion)
		if err := table.ReplayOnRestart(currentVersion, latestVersion); err != nil {
			log.Fatal("(error) could not replay writes: ", err)
		}
	}

	// table.Put("key1", "hello");
	// table.Put("key2", "hello2");
	// table.Put("key3", "hello3");
	// table.Commit();
	// table.Put("key1", "hello1");
	table.List()
}
