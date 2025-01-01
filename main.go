package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"github.com/MridulDhiman/chronotable/chronotable"
	"github.com/MridulDhiman/chronotable/config"
	"github.com/MridulDhiman/chronotable/internal/utils"
)

var initialized bool = false;

func init() {
	newpath := filepath.Join(".", config.CHRONO_MAIN_DIR)
	// check if directory exists or not
	yes, err:= utils.Exists(newpath)
	if err != nil {
		log.Fatalln(err)
	}

	if !yes {
	os.MkdirAll(newpath, os.FileMode(0755))
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
		Initialized: initialized,
		Logger: logger,
		Mode: "local",
	})
	
	if initialized {
		// get current state by replaying logs
		currentVersion,latestVersion, err:= table.ConfigHandler.FetchLatestVersion()
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println("latest version: ", latestVersion)
		fmt.Println("current version: ", currentVersion)
		if currentVersion != 0 {
			if err := table.ReplayOnRestart(currentVersion, latestVersion); err != nil {
				log.Fatal("(error) could not replay writes: ", err)
			}
		}
	}

	table.Put("key1", "hello");
	table.Put("key2", "hello2");
	table.Put("key3", "hello3");
}
