package stock

import (
	"io/ioutil"
	"log"
	"strings"
)

const symbolsFile = "../../resources/data/symbols.dat"

func GetSymbols() []string {
	allSymbols, err := ioutil.ReadFile(symbolsFile)
	if err != nil {
		log.Fatalln("Error reading list of stock symbols:", err)
	}

	return strings.Split(string(allSymbols), "\n")
}
