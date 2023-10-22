package scan

import (
	"go_cli/utils"
	"log/slog"
	"os"
	"strings"
)

/*
Scan
扫描给定路径，抓取该路径及其子文件夹,搜索 Git 仓库
*/
func Scan(path string) {
	slog.Info("[Start matching .git files]")

	repositories := recursiveScanFolder(path)

	utils.AddNewPaths(repositories)

	slog.Info("[.git files have added successfully]")
}

// recursiveScanFolder 开始递归搜索位于“folder”子树中的 git 存储库
func recursiveScanFolder(folder string) []string {
	return scanGitFolders(make([]string, 0), folder)
}

// scanGitFolders 返回git库的基本文件夹，即 .git 文件隶属的文件夹目录
func scanGitFolders(folders []string, folder string) []string {
	folder = strings.TrimSuffix(folder, "/")
	f, err := os.Open(folder)
	if err != nil {
		slog.Error("[os.Open() ERROR] ", err)
		return nil
	}
	// Readdir 函数通常用于从一个目录中读取其内容，并返回一个目录项（文件和子目录）的切片。
	files, err := f.Readdir(-1)
	defer f.Close()
	if err != nil {
		slog.Error("[f.Readdir() ERROR] ", err)
		return nil
	}

	for _, file := range files {
		/*
			if 文件夹
				then
			    if .git文件，
					then 将.git文件的父级目录存入folders
					continue
				继续递归
			endif
		*/
		if file.IsDir() {
			path := folder + "\\" + file.Name()
			if file.Name() == ".git" {
				path = strings.TrimSuffix(path, "\\.git")
				slog.Info("[.git file] ", "path", path)
				folders = append(folders, path)
				continue
			}
			// `vendor`文件夹和·node_modules`文件夹不扫描
			if file.Name() == "vendor" || file.Name() == "node_modules" {
				continue
			}
			folders = scanGitFolders(folders, path)
		}
	}
	return folders
}
