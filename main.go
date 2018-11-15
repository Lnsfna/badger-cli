package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/c-bata/go-prompt"
)

var dbPath = flag.String("db_path", "./badger", "badger db path")

var suggestions = []prompt.Suggest{
	{"keys", "list keys by key prefix"},
	{"get", "get key value"},
	{"set", "set value"},
	{"del", "delete key"},
	{"exit", "exit the program"},
	{"help", "help info"},
}

func completer(in prompt.Document) []prompt.Suggest {
	w := in.GetWordBeforeCursor()
	if w == "" {
		return []prompt.Suggest{}
	}
	return prompt.FilterHasPrefix(suggestions, w, true)
}

func pretty(data []byte) (string, error) {
	out := new(bytes.Buffer)
	if err := json.Indent(out, data, "", "  "); err != nil {
		return "", err
	}
	return out.String(), nil
}

func executor(in string) {
	in = strings.TrimSpace(in)

	blocks := strings.SplitN(in, " ", 3)
	switch blocks[0] {
	case "keys":
		if len(blocks) < 2 {
			fmt.Println("Please input key prefix.")
			return
		}
		if keys, err := ListKeys(blocks[1]); err != nil {
			fmt.Println(err)
		} else {
			for _, k := range keys {
				fmt.Println(k)
			}
		}
	case "get":
		if len(blocks) < 2 {
			fmt.Println("Please input a key.")
			return
		}
		if val, err := Get(blocks[1]); err != nil {
			fmt.Println(err)
		} else {
			if out, err := pretty(val); err != nil {
				fmt.Println(string(val))
			} else {
				fmt.Println(out)
			}
		}
	case "set":
		if len(blocks) < 3 {
			fmt.Println("Please input key and value.")
			return
		}
		if err := Set(blocks[1], []byte(blocks[2])); err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("Set success.")
		}
	case "del":
		if len(blocks) < 2 {
			fmt.Println("Please input a key.")
			return
		}
		if err := Delete(blocks[1]); err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("Delete success.")
		}
	case "help":
		fmt.Print("Available commands: \n\n")
		for _, s := range suggestions {
			fmt.Printf("%6s:\t\t%s\n\n", s.Text, s.Description)
		}
	case "exit":
		fmt.Println("Exit!")
		os.Exit(0)
	}

}

func main() {
	flag.Parse()

	var err error
	if err = InitDB(*dbPath); err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	p := prompt.New(
		executor,
		completer,
		prompt.OptionPrefix(">>> "),
		prompt.OptionTitle("badger-cli"),
	)
	p.Run()
}
