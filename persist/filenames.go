package persist

import "strings"

const optionsDir = "../../resources/store/options/"
const pricesDir = "../../resources/store/prices/"

func stockFilename(symbol string) string {
	return strings.Replace(symbol, "/", ".", 1)
}
