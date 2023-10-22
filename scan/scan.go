package scan

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"log/slog"
	"os"
	"strings"
)

/*
Scan
扫描给定路径，抓取该路径及其子文件夹,搜索 Git 存储库
*/
func Scan(path string) {
	fmt.Printf("Found folders:\n\n")
	repositories := recursiveScanFolder(path)
	filePath := GetDotFilePath()
	addNewSliceElementsToFile(filePath, repositories)
	fmt.Printf("\n\nSuccessfully added\n\n")
}

// recursiveScanFolder starts the recursive search of git repositories
// living in the `folder` subtree
func recursiveScanFolder(folder string) []string {
	return scanGitFolders(make([]string, 0), folder)
}

// scanGitFolders 返回以“.git”结尾的“folder”子文件夹列表.
// Returns the base folder of the repo, the .git folder parent.
// Recursively searches in the subfolders by passing an existing `folders` slice.
func scanGitFolders(folders []string, folder string) []string {
	// trim the last `/`
	folder = strings.TrimSuffix(folder, "/")

	f, err := os.Open(folder)
	if err != nil {
		log.Fatal(err)
	}
	files, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		log.Fatal(err)
	}

	var path string

	for _, file := range files {
		if file.IsDir() {
			path = folder + "\\" + file.Name()
			if file.Name() == ".git" {
				path = strings.TrimSuffix(path, "\\.git")
				fmt.Println(path)
				folders = append(folders, path)
				continue
			}
			if file.Name() == "vendor" || file.Name() == "node_modules" {
				continue
			}
			folders = scanGitFolders(folders, path)

		}
	}
	return folders
}

// GetDotFilePath returns the dot file for the repos list.
// Crt and the enclosing folder if it does nt exist.
func GetDotFilePath() string {
	//usr, err := user.Current()
	//if err != nil {
	//	log.Fatal(err)
	//}

	dotFile := "C:\\Users\\19406\\Desktop\\go\\store.txt"

	return dotFile
}

// addNewSliceElementsToFile given a slice of strings representing paths, stores them
// to the filesystem
func addNewSliceElementsToFile(filePath string, newRepos []string) {
	existingRepos := ParseFileLinesToSlice(filePath)
	repos := joinSlices(newRepos, existingRepos)
	dumpStringsSliceToFile(repos, filePath)
}

// ParseFileLinesToSlice given a file path string, gets the content
// of each line and parses it to a slice of strings.
func ParseFileLinesToSlice(filePath string) []string {
	f := openFile(filePath)
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

// openFile opens the file located at `filePath`. Creates it if not existing.
func openFile(filePath string) *os.File {
	f, err := os.Open(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			// file does not exist
			_, err = os.Create(filePath)
			if err != nil {
				panic(err)
			}
		} else {
			// other error
			panic(err)
		}
	}

	return f
}

// joinSlices adds the element of the `new` slice
// into the `existing` slice, only if not already there
func joinSlices(new []string, existing []string) []string {
	for _, i := range new {
		if !sliceContains(existing, i) {
			existing = append(existing, i)
		}
	}
	return existing
}

// sliceContains returns true if `slice` contains `value`
func sliceContains(slice []string, value string) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

// dumpStringsSliceToFile writes content to the file in path `filePath` (overwriting existing content)
func dumpStringsSliceToFile(repos []string, filePath string) {
	content := strings.Join(repos, "\n")
	os.WriteFile(filePath, []byte(content), 0755)
}
