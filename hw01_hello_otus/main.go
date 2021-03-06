package main

import (
	"fmt"
)

func main() {
	runes := []rune("Hello, OTUS!")
	iS, iE := 0, len(runes)-1

	for iS < iE {
		runes[iS] = (runes[iS] ^ runes[iE])
		runes[iE] = (runes[iS] ^ runes[iE])
		runes[iS] = (runes[iS] ^ runes[iE])
		iS++
		iE--
	}

	fmt.Println(string(runes))
}
