// Package rand is a utility package to generate random
package rand

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var defaultInput = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func String(n int, in []rune) string {
	b := make([]rune, n)

	if in == nil || len(in) < 1 {
		in = defaultInput
	}

	for i := range b {
		b[i] = in[rand.Intn(len(in))]
	}
	return string(b)
}
