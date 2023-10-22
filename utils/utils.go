package utils

import (
	"bufio"
	"io"
	"log"
	"log/slog"
	"os"
	"os/user"
	"slices"
	"strings"
)

var DOTFILE string

// GetDotFilePath 返回存储仓库路径列表的txt文件
func init() {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	DOTFILE = usr.HomeDir + "\\" + "repositories.txt"
}

// ParseFileToSlice 根据给定存储仓库路径的文件路径, 获取每一行的内容并将其解析为字符串切片。
func ParseFileToSlice() []string {
	f := openFile(DOTFILE)
	defer f.Close()

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

	return lines
}

// openFile 打开位于`filePath`的文件, 如果不存在则创建它。
func openFile(filePath string) *os.File {
	f, err := os.Open(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			_, err = os.Create(filePath)
			if err != nil {
				panic(err)
			}
		} else {
			panic(err)
		}
	}
	return f
}

// AddNewPaths 将新扫描到的含.git文件的目录追加到存储文件中
func AddNewPaths(newRepos []string) {

	existingRepos := ParseFileToSlice()
	repos := joinSlices(newRepos, existingRepos)

	paths := strings.Join(repos, "\n")
	os.WriteFile(DOTFILE, []byte(paths), 0755)
}

// joinSlices 将`new`切片的元素添加到`existing`切片中，前提是尚不存在
func joinSlices(newPaths []string, existedPaths []string) []string {
	for _, path := range newPaths {
		if !slices.Contains(existedPaths, path) {
			existedPaths = append(existedPaths, path)
		}
	}
	return existedPaths
}
