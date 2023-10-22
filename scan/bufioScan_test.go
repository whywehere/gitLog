package scan

import (
	"bufio"
	"fmt"
	"io"
	"log/slog"
	"os"
	"testing"
)

func TestScan2(t *testing.T) {
	filePath := "C:\\Users\\19406\\Desktop\\go\\go_cli\\scan\\data.txt"
	f, err := os.Open(filePath)
	defer f.Close()
	if err != nil {
		if os.IsNotExist(err) {
			_, err = os.Create(filePath)
			if err != nil {
				panic(err)
			}
		} else {
			// other error
			panic(err)
		}
	}

	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		if err != io.EOF {
			slog.Info("", err)
			panic(err)
		}
	}
	fmt.Println(lines)
}
