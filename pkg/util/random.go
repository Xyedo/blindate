package util

import (
	"fmt"
	"math"
	"math/rand"
	"strings"
	"time"

	"github.com/google/uuid"
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

func RandomUUID() string {
	uid := uuid.New()
	return uid.String()
}

func RandomPoint(precision int) string {
	return fmt.Sprintf("POINT(%f %f)", randomFloat(-90, 90, precision), randomFloat(-180, 180, precision))
}

func RandomLat() float64 {
	return float64(-90) + rand.Float64()*float64(180)
}
func RandomLng() float64 {
	return float64(-180) + rand.Float64()*float64(360)
}
func randomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}
func randomFloat(min, max, precision int) float64 {
	output := math.Pow(10, float64(precision))
	num := float64(min) + rand.Float64()*float64(max-min)
	round := int(num*output + math.Copysign(0.5, num*output))
	return float64(round) / output
}
