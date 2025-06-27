package utils

import (
	"fmt"
	"math/rand"
	"time"
)

func GenerateOTP() string {
	seed := time.Now().UnixNano()
	source := rand.NewSource(seed)
	r := rand.New(source)
	return fmt.Sprintf("%06d", r.Intn(1000000))
}
