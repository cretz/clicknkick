package main

import (
	"log"

	"github.com/cretz/clicknkick/game"
)

func main() {
	g, err := game.New(21*64, 10*64)
	if err == nil {
		err = g.Run()
	}
	if err != nil {
		log.Fatal(err)
	}
}
