package common

import (
	"crypto/rand"
	"encoding/base64"
	"log"
	"time"
)

// From https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go

// const letterBytes = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
// const letterLen = int64(len(letterBytes))
// for i := range buf {
// 	buf[i] = letterBytes[rand.Int63()%letterLen]
// }

func RandStringBytes(n int) string {
	return string(RandBytes(n))
}

func RandBytes(n int) []byte {
	now := time.Now()
	buf := make([]byte, n)
	rand.Read(buf)
	since := time.Since(now).String()
	res := base64.StdEncoding.EncodeToString(buf)
	b := "# " + since + "\n" + res
	log.Printf("Buffer of %d bytes generated in %s", n, since)

	return []byte(b[:n])
}
