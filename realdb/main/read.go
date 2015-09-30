package main

import (
	"fmt"
	"io/ioutil"
	"log"
)

func main() {

	fi, err := ioutil.ReadDir("db")
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < len(fi); i++ {
		data := read("db/" + fi[i].Name())
		fmt.Println(data)
	}
	fmt.Printf("Found %d files to read...\n", len(fi))
}

func read(path string) string {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("reading file: %s -> %v\n", path, err)
	}
	return string(b)
}
