package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/cagnosolutions/datastructures/index"
)

func cmp(a, b int64) int {
	if a < b {
		return -1
	}
	if a > b {
		return 1
	}
	// a == b ??
	return 0
}

func main() {

	bpt := index.NewTree(cmp)
	defer bpt.Close()

	bpt.Set(int64(312842341), []byte("bar"))
	fmt.Println("Tree size:", bpt.Len())

	fmt.Print("\n\nPress any key to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')

}
