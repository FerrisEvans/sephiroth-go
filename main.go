package main

import (
	"sephiroth-go/core"
	"sephiroth-go/init"
)

func main() {
	core.Db = init.Database()
}
