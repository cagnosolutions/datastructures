package main

import (
	"fmt"
	"time"

	"github.com/cagnosolutions/datastructures/realdb"
)

var n = 1000000 / 4
var o int

func main() {

	ts := time.Now().Unix()
	arch := realdb.NewArchiver(1)

	/*
		for i := 0; i < n; i++ {
			k, v := fmt.Sprintf("key-%d", i), fmt.Sprintf("val-%d", i)
			arch.Encode(&M{time.Now().UnixNano(), "SET", "widgets", k, []byte(v)})
		}
	*/

	arch.Decode()

	fmt.Printf("Took %dms (%d entries)\n", time.Now().Unix()-ts, n)

}

type M struct {
	Timestamp int64  `json:"ts,omitempty"`
	Command   string `json:"op,omitempty"`
	Key       string `json:"k,omitempty"`
	Field     string `json:"f,omitempty"`
	Value     []byte `json:"v,omitempty"`
}
