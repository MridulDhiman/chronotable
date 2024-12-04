package main

import (
	
	"os"
	"path/filepath"

	"github.com/MridulDhiman/chronotable/chronotable"
	"github.com/MridulDhiman/chronotable/config"
)

func init() {
	newpath := filepath.Join(".", config.CHRONO_MAIN_DIR)
	 os.MkdirAll(newpath, os.FileMode(0755))
}

func main() {
	table := chronotable.New(&chronotable.Options{
		EnableAOF:      true,
		AOFPath:        config.MAIN_AOF_FILE,
		EnableSnapshot: true,
	})
	table.Put("key1", 23)
	table.Put("key2", "hello")
	table.Put("key3", "yo")
	table.Commit()
	table.Put("key4", 324)
	table.Commit()
	table.Put("key5", "namaste")
	table.Timetravel(2)
	
}
