package main

import (
	"bufio"
	"crypto/rand"
	"fmt"
	"log"
	"os"
	"time"
)

var count = 1000000 / 32

func makeBytes(n int) []byte {
	b := make([]byte, 10)
	if _, err := rand.Read(b); err != nil {
		log.Fatal(err)
	}
	return b
}

func main() {

	m := make(map[string][]byte, 0)

	fmt.Printf("Adding %d entries to map...\n", count)
	ts := time.Now().Unix()
	for i := 0; i < count; i++ {
		m[string(makeBytes(10))] = makeBytes(50)
	}
	fmt.Printf("Finished adding %d entries. Took %dms\n", count, (time.Now().Unix() - ts))
	fmt.Println("Map count:", len(m))

	fmt.Print("\n\nPress any key to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}
