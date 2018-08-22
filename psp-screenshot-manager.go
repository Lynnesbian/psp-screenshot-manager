package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"github.com/jessevdk/go-flags"
	"io"
	"log"
	"os"
	"path/filepath"
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
			PathName string
			Trailing []string
		} `positional-args:"yes" required:"yes"`
	}

	_, err := flags.ParseArgs(&opts, os.Args)
	if err != nil {
		panic(err)
	}

	gameDB := loadGames()

	gameDB = gameDB

	fmt.Printf("Scanning %v...\n", opts.Filepath.PathName)
	screenshots := make(map[string][]string)
	err = filepath.Walk(opts.Filepath.PathName, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", opts.Filepath.PathName, err)
			return err
		}
		if info.IsDir() {
			return nil
		}
		// fmt.Printf("visited file: %q\n", path)
		relpath := path[len(opts.Filepath.PathName):]
		_, filename := filepath.Split(relpath)
		screenshots[filepath.Dir(relpath)] = append(screenshots[filepath.Dir(relpath)], filename)
		return nil
	})

}
