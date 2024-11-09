package main

import (
	"github.com/MridulDhiman/chronotable/chronotable"
)

func main() {
	table := chronotable.New(&chronotable.Options{
		EnableAOF: true,
		AOFPath:   "./chrono.aof",
		EnableSnapshot: true,
	})
	table.Put("key1", 23)
	table.Put("key2", "hello")
	table.Put("key3", "yo")
	table.Commit()
	table.Put("key4", 324)
	table.Commit()
}
