package main

import (
	"fmt"

	"github.com/Depado/periodic-file-fetcher/external"
)

func main() {
	ft := external.Fetcher{
		ConfigurationDir: "conf/active/",
		BackupDir:        "content/backup/",
		ContentDir:       "content/",
	}
	ft.Start()
	fmt.Scanln()
}
