package main

import (
	"github.com/MridulDhiman/chronotable/pkg/chronotable"
)

func main() {
	table := chronotable.New()
	table.Put("key1", 23)
	table.Put("key2", "hello")
}
