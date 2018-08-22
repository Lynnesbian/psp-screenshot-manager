package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"github.com/jessevdk/go-flags"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

type Game struct {
	serial   string
	name     string
	longname string
	category string
}

func UserHomeDir() string { //https://stackoverflow.com/a/7922977/4480824
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}
	return os.Getenv("HOME")
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
		Overwrite       bool   `short:"w" long:"overwrite" description:"Overwrite existing screenshots of the same name. If not provided, will rename files to avoid overwriting."`
		OutputDirectory string `short:"o" long:"output-directory" default:"/home/lynne/Pictures/PSP Screenshots" description:"Directory to output to. Defaults to ~/Pictures/PSP Screenshots"`

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
	fmt.Println(opts.OutputDirectory)
	screenshots := make(map[string][]string)
	err = filepath.Walk(opts.Filepath.PathName, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatal(err)
			return err
		}
		if info.IsDir() {
			return nil
		}
		// fmt.Printf("visited file: %q\n", path)
		relpath := path[len(opts.Filepath.PathName):]
		_, filename := filepath.Split(relpath)
		//[1:] to trim the trailing slash
		screenshots[filepath.Dir(relpath)[1:]] = append(screenshots[filepath.Dir(relpath)], filename)
		return nil
	})
	for folder, files := range screenshots {
		//determine the game's name from its serial
		gameName := "Unknown"
		for _, game := range gameDB {
			if game.serial == folder {
				gameName = game.name
				break
			}
		}
		//make the folder to put the screenshots in
		saveLocation := fmt.Sprintf("%v/%v", opts.OutputDirectory, gameName)
		err = os.MkdirAll(saveLocation, os.ModePerm)
		if err != nil {
			panic(err) //todo: not this
		}
		for _, file := range files {
			fullpath := fmt.Sprintf("%v/%v/%v", opts.Filepath.PathName, folder, file)
			fullpath = fullpath
			cmd := exec.Command("convert", fullpath, saveLocation+file+".png")
			err := cmd.Run()
			fmt.Println(err)
		}
	}
}
