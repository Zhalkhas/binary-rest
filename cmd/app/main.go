package main

import (
	"github.com/zhalkhas/binary-rest/internal/app"
	"log"
)

func main() {
	a := app.New(app.Config{})
	if err := a.Run(); err != nil {
		log.Fatalln(err)
	}
}
