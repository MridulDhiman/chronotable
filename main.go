package main

import (
	"github.com/MridulDhiman/chronotable/chronotable"
)

func main() {
	table := chronotable.New(&chronotable.Options{
		EnableAOF: true,
		AOFPath: "./chrono.aof",
	})
	table.Put("key1", 23)
	table.Put("key2", "hello")
	table.Put("key3", "yo")
}
