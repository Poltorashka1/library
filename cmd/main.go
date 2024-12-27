package main

import (
	"book/internal/app"
	_ "github.com/mattn/go-sqlite3"
)

const cfgFileName = ".env"

func main() {
	a := app.New()
	a.Start(cfgFileName)
}
