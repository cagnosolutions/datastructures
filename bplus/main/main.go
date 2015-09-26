package main

import (
	"bytes"
	"fmt"

	"github.com/cagnosolutions/bplus"
)

func cmp(a, b []byte) int {
	return bytes.Compare(a, b)
}

func main() {

	bpt := bplus.NewTree(cmp)
	defer bpt.Close()

	bpt.Set([]byte("foo"), []byte("bar"))
	bpt.Set([]byte("baz"), []byte("nar"))
	bpt.Set([]byte("zap"), []byte("rad"))

	ret, ok := bpt.Get([]byte("baz"))
	fmt.Printf("%s (%v)\n", ret, ok)
}
