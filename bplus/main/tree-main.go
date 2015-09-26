package main

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/cagnosolutions/bplus"
)

var count int = 1000000 / 32

func cmp(a, b []byte) int {
	return bytes.Compare(a, b)
}

func makeBytes(n int) []byte {
	b := make([]byte, 10)
	if _, err := rand.Read(b); err != nil {
		log.Fatal(err)
	}
	return b
}

func main() {

	bpt := bplus.NewTree(cmp)
	defer bpt.Close()

	fmt.Printf("adding %d entries to btree...\n", count)
	ts := time.Now().Unix()
	for i := 0; i < count; i++ {
		bpt.Set(makeBytes(10), makeBytes(50))
	}
	fmt.Printf("Finished adding %d entries. Took %dms\n", count, (time.Now().Unix() - ts))
	fmt.Println("B+Tree count:", bpt.Len())

	/*
		fmt.Println("Clearing tree, adding a few entries, and seeking first...")
		bpt.Clear()
		bpt.Set([]byte("g"), []byte("gg"))
		bpt.Set([]byte("h"), []byte("hh"))
		bpt.Set([]byte("i"), []byte("ii"))
		bpt.Set([]byte("a"), []byte("aa"))
		bpt.Set([]byte("f"), []byte("ff"))
		bpt.Set([]byte("k"), []byte("kk"))
		bpt.Set([]byte("j"), []byte("jj"))
		bpt.Set([]byte("a"), []byte("second a"))
		bpt.Set([]byte("e"), []byte("ee"))
		bpt.Set([]byte("b"), []byte("bb"))
		bpt.Set([]byte("1"), []byte("one"))
		bpt.Set([]byte("d"), []byte("dd"))
		bpt.Set([]byte("c"), []byte("cc"))
		enum, err := bpt.SeekFirst()
		defer enum.Close()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Enumerating leaf nodes in lexigraphical order...")
		for {
			k, v, err := enum.Next()
			if err != nil {
				break
			}
			fmt.Printf("Key: %s, Val: %s\n", k, v)
		}
		fmt.Println("Done")
	*/

	fmt.Print("\n\nPress any key to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')

}
