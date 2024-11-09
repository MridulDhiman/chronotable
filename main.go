package main

import (
	"fmt"

	"github.com/MridulDhiman/chronotable/chronotable"
	"github.com/MridulDhiman/chronotable/config"
)

func main() {
	table := chronotable.New(&chronotable.Options{
		EnableAOF:      true,
		AOFPath:        config.AOF_PATH,
		EnableSnapshot: true,
	})
	table.Put("key1", 23)
	table.Put("key2", "hello")
	table.Put("key3", "yo")
	table.Commit()
	table.Put("key4", 324)
	table.Commit()
	table.RollbackTo(1)
	if _, ok := table.Get("key4"); !ok {
		fmt.Println("key4 not found")
	}
}
