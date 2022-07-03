/*
=== Утилита wget ===

Реализовать утилиту wget с возможностью скачивать сайты целиком

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/
package main

import (
	"os"

	"go-wget/wget"
)

func main() {
	os.Exit(wget.CLI(os.Args[1:]))
}
