package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
)

type Game struct {
	serial   string
	name     string
	longname string
	category string
}

func main() {
	gameDB, err := os.Open("games.csv")
	if err != nil {
		log.Fatal(err)
		log.Fatal("(you can get games.csv from https://github.com/Lynnesbian/psp-screenshot-manager/blob/master/games.csv)")
	}
	reader := csv.NewReader(bufio.NewReader(gameDB))
	var games []Game
	_, _ = reader.Read() //discard the first line, it's just the header

	for {
		line, err := reader.Read()
		if err == io.EOF {
			//we've reached the end of the file
			break
		} else if err != nil {
			log.Fatal(err)
		}

		games = append(games, Game{
			serial:   line[0],
			name:     line[1],
			longname: line[2],
			category: line[3],
		})
	} //end csv reader loop

}
