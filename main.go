package main

import (
	"fmt"

	"github.com/Depado/periodic-file-fetcher/external"
)

func main() {
	external.LoadAndStart()
	fmt.Scanln()
}
