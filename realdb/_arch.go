package realdb

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"sync"
)

type Archiver struct {
	maxmb   int64
	count   int64
	written int64
	fd      *os.File
	enc     *json.Encoder
	sync.RWMutex
}

func NewArchiver(maxmb int) *Archiver {
	a := &Archiver{
		maxmb: int64((1024 * 1024) * maxmb),
	}
	a.fd = open(a.count)
	a.enc = json.NewEncoder(a.fd)
	return a
}

func (a *Archiver) Encode(v interface{}) {
	a.Lock()
	defer a.Unlock()
	a.checkUpdateEncoder()
	if err := a.enc.Encode(&v); err != nil {
		log.Fatal("Archiver.Encode(): ", err)
	}
}

func (a *Archiver) Decode() {
	a.Lock()
	defer a.Unlock()
	info, err := ioutil.ReadDir("db")
	if err != nil {
		log.Fatal("Archiver.Decode(): ", err)
	}
	for n, _ := range info {
		fd := open(int64(n))
		defer fd.Close()
		r := bufio.NewReader(fd)
		for {
			line, err := r.ReadBytes('\n')
			if err == io.EOF {
				break
			} else if err != nil {
				log.Fatal("Archiver.Decode(): ", err)
			}
			doc := line[:len(line)-1]
			fmt.Printf("Doc: %s\n", doc)
		}
	}
}

func (a *Archiver) Close() {
	if err := a.fd.Close(); err != nil {
		log.Fatal("Archiver.checkUpdateEncoder() (a.fd.Close()): ", err)
	}
}

func (a *Archiver) checkUpdateEncoder() {
	info, err := a.fd.Stat()
	if err != nil {
		log.Fatal("Archiver.checkUpdateEncoder() (a.fd.Stat()...): ", err)
	}
	if info.Size() < a.maxmb {
		return
	}
	a.count++
	a.fd = open(a.count)
	a.enc = json.NewEncoder(a.fd)
}

func open(fileno int64) *os.File {
	filepath := fmt.Sprintf("db/%d.dat", fileno)
	dir, file := path.Split(filepath)
	if dir != "" {
		if _, err := os.Stat(dir); err != nil && os.IsNotExist(err) {
			if err := os.MkdirAll(dir, 0755); err != nil {
				log.Fatal("Archiver.open(): ", err)
			}
		}
	}
	if file != "" {
		if _, err := os.Stat(filepath); err != nil && os.IsNotExist(err) {
			if _, err := os.Create(filepath); err != nil {
				log.Fatal("Archiver.open(): ", err)
			}
		}
	}
	fd, err := os.OpenFile(filepath, os.O_APPEND|os.O_RDWR, 0644)
	if err != nil {
		log.Fatal("Archiver.open(): ", err)
	}
	if err := fd.Sync(); err != nil {
		log.Fatal("Archiver.open(): ", err)
	}
	return fd
}
