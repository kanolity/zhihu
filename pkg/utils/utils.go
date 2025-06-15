package utils

import (
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func EndOfDay(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 23, 59, 59, 0, t.Location())
}

func RandomNumber(size int) string {
	if size <= 0 {
		panic("{ size : " + strconv.Itoa(size) + " } must be more than 0 ")
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	value := make([]byte, size)

	for index := 0; index < size; index++ {
		value[index] = byte('0' + r.Intn(10))
	}
	var b strings.Builder
	b.Write(value)
	return b.String()
}
