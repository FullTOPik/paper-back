package utils

import (
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func GenerateCode() string {
	number := rand.Intn(10)
	var chars strings.Builder

	if number < 4 {
		number = 4
	}

	for i := 0; i < number; i++ {
		register := rand.Intn(2)
		var startRune rune

		if register == 0 {
			startRune = 65
		} else {
			startRune = 97
		}
		chars.WriteRune(startRune + rune(rand.Intn(26)))
	}

	for _, value := range strings.Split(strconv.Itoa(int(time.Now().Unix())), "") {
		i, _ := strconv.ParseInt(value, 10, 32)
		chars.WriteRune(rune(65 + i))
	}

	return chars.String()
}
