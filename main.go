package main

import (
	"flag"
	"fmt"
	"go_cli/scan"
	"go_cli/stat"
	"time"
)

func main() {
	var folder string
	var email string
	flag.StringVar(&folder, "add", "", "add a new folder to scan for Git repositories")
	flag.StringVar(&email, "email", "104337697+whywehere@users.noreply.github.com", "the email to scan")
	flag.Parse()
	startingTime := time.Now().UTC()

	if folder != "" {
		scan.Scan(folder)
		endingTime := time.Now().UTC()
		fmt.Println(endingTime.Sub(startingTime))
		return
	}

	stat.Stats(email)
	endingTime := time.Now().UTC()
	fmt.Println(endingTime.Sub(startingTime))
}
