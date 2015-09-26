package main

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"fmt"
	"log"
	"os"

	"github.com/cagnosolutions/bplus"
)

func makeBytes(n int) []byte {
	b := make([]byte, 10)
	if _, err := rand.Read(b); err != nil {
		log.Fatal(err)
	}
	return b
}

func cmp(a, b []byte) int {
	return bytes.Compare(a, b)
}

func main() {

	bpt := bplus.NewTree(cmp)
	defer bpt.Close()

	bpt.Set([]byte("foo"), []byte("bar"))
	fmt.Println("Tree size:", bpt.Len())

	fmt.Print("\n\nPress any key to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')

}
