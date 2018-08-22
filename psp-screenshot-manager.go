package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"github.com/jessevdk/go-flags"
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

func loadGames() []Game {
	gameCSV, err := os.Open("games.csv")
	if err != nil {
		log.Print(err)
		log.Fatal("(you can get games.csv from https://github.com/Lynnesbian/psp-screenshot-manager/blob/master/games.csv)")
	}
	reader := csv.NewReader(bufio.NewReader(gameCSV))
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
	return games
}

func main() {

	var opts struct {
		overwrite bool `short:"o" long:"overwrite" description:"Overwrite existing screenshots of the same name. If not provided, will rename files to avoid overwriting."`

		Filepath struct {
			PathToMe string //the first argument will be the name of this file
			Filename string
			Trailing []string
		} `positional-args:"yes" required:"yes"`
	}

	_, err := flags.ParseArgs(&opts, os.Args)
	if err != nil {
		panic(err)
	}

	gameDB := loadGames()

	gameDB = gameDB
	fmt.Println(opts.Filepath.Filename)
}
