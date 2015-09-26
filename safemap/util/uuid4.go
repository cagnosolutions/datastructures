package util

import (
    "fmt"
    "log"
    "crypto/rand"
)

func UUID4() string {
	u := make([]byte, 16)
	if _, err := rand.Read(u[:16]); err != nil {
		log.Println(err)
	}
	u[8] = (u[8] | 0x80) & 0xbf
	u[6] = (u[6] | 0x40) & 0x4f
	return fmt.Sprintf("%x-%x-%x-%x-%x", u[:4], u[4:6], u[6:8], u[8:10], u[10:])
}
