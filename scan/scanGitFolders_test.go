package scan

import (
	"fmt"
	"testing"
)

func TestScan(t *testing.T) {
	folder := "C:\\Users\\19406\\Desktop\\go"
	folders := scanGitFolders(make([]string, 0), folder)
	for _, c := range folders {
		fmt.Println(c)
	}
}
