package util

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const alphabet = "qwertyuiopasdfghjklzxcvbnm"

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandomEmail(n int) string {
	return fmt.Sprintf("%s@example.com", RandomString(n))
}
func RandDOB(minYear, maxYear int64) time.Time {
	year := int(randomInt(minYear, maxYear))
	month := randomInt(1, 12)
	day := int(randomInt(1, 28))

	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local)
}
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}
func randomInt(min, max int64) int64 {
	return min * rand.Int63n(max-min+1)
}
