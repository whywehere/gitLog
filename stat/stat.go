package stat

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"sort"
	"time"
)

const outOfRange = 99999
const daysInLastSixMonths = 183
const weeksInLastSixMonths = 26

type column []int

// Stats calculates and prints the stats.
func Stats(email string) {
	commits := processRepositories(email)
	printCommitsStats(commits)
}

// getBeginningOfDay given a time.Time calculates the start time of that day
func getBeginningOfDay(t time.Time) time.Time {
	year, month, day := t.Date()
	startOfDay := time.Date(year, month, day, 0, 0, 0, 0, t.Location())
	return startOfDay
}

// countDaysSinceDate 计算自`data`以来经过了多少天
func countDaysSinceDate(date time.Time) int {
	days := 0
	now := getBeginningOfDay(time.Now())
	for date.Before(now) {
		date = date.Add(time.Hour * 24)
		days++
		if days > daysInLastSixMonths {
			return outOfRange
		}
	}
	return days
}

// fillCommits 给定在“path”中找到repository，获取提交并将它们放入“commits”映射中，完成后返回它
func fillCommits(email string, path string, commits *map[int]int) {
	// 实例化 git repo 对象repo
	repo, err := git.PlainOpen(path)
	if err != nil {
		panic(err)
	}

	// 获得repo的头部
	ref, err := repo.Head()
	if err != nil {
		panic(err)
	}
	// 获取从 HEAD 开始的提交历史记录
	iterator, err := repo.Log(&git.LogOptions{From: ref.Hash()})
	if err != nil {
		panic(err)
	}

	offset := calcOffset()
	// 遍历提交历史记录
	err = iterator.ForEach(func(c *object.Commit) error {
		daysAgo := countDaysSinceDate(c.Author.When) + offset

		if c.Author.Email != email {
			return nil
		}
		//如果提交时间在规定之内
		if daysAgo != outOfRange {
			(*commits)[daysAgo]++
		}

		return nil
	})
	if err != nil {
		panic(err)
	}
}

// processRepositories 给定用户电子邮件，返回过去 6 个月内所做的提交
func processRepositories(email string) map[int]int {
	// 获取git文件存储路径
	//filePath := scan.GetDotFilePath()

	// 根据filePath 获取repositories
	//repos := scan.ParseFileLinesToSlice(filePath)
	repos := []string{"C:\\Users\\19406\\Desktop\\go\\tta\\cinx"}
	daysInMap := daysInLastSixMonths

	//初始化 在过去 6 个月内每天的contributions
	commits := make(map[int]int, daysInMap)
	for i := daysInMap; i > 0; i-- {
		commits[i] = 0
	}
	//统计每个repository的contributions
	for _, path := range repos {
		fillCommits(email, path, &commits)
	}

	return commits
}

// calcOffset 确定并返回填充统计图最后一行所缺少的天数
func calcOffset() int {
	var offset int
	weekday := time.Now().Weekday()

	switch weekday {
	case time.Sunday:
		offset = 7
	case time.Monday:
		offset = 6
	case time.Tuesday:
		offset = 5
	case time.Wednesday:
		offset = 4
	case time.Thursday:
		offset = 3
	case time.Friday:
		offset = 2
	case time.Saturday:
		offset = 1
	}

	return offset
}

// printCell given a cell value prints it with a different format
// based on the value amount, and on the `today` flag.
func printCell(val int, today bool) {
	escape := "\033[0;37;30m"
	switch {
	case val > 0 && val < 5:
		escape = "\033[1;30;47m"
	case val >= 5 && val < 10:
		escape = "\033[1;30;43m"
	case val >= 10:
		escape = "\033[1;30;42m"
	}

	if today {
		escape = "\033[1;37;45m"
	}

	if val == 0 {
		fmt.Printf(escape + "  - " + "\033[0m")
		return
	}

	str := "  %d "
	switch {
	case val >= 10:
		str = " %d "
	case val >= 100:
		str = "%d "
	}

	fmt.Printf(escape+str+"\033[0m", val)
}

// printCommitsStats prints the commits stats
func printCommitsStats(commits map[int]int) {
	keys := sortMapIntoSlice(commits)
	columns := buildCols(keys, commits)
	printCells(columns)
}

// sortMapIntoSlice 返回已排序map的key的切片
func sortMapIntoSlice(m map[int]int) []int {
	// To store the keys in slice in sorted order
	var keys []int
	for k := range m {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	return keys
}

// buildCols 生成一个包含行和列的地图，可以打印到屏幕上
func buildCols(keys []int, commits map[int]int) map[int]column {
	columns := make(map[int]column) //图表的列 每列有七天

	col := column{}

	for _, k := range keys {
		week := k / 7      //26,25...1
		dayOfWeek := k % 7 // 0,1,2,3,4,5,6

		if dayOfWeek == 0 { //reset
			col = column{}
		}

		col = append(col, commits[k])

		if dayOfWeek == 6 {
			columns[week] = col
		}
	}

	return columns
}

// printCells prints the cells of the graph
func printCells(columns map[int]column) {
	printMonths()
	for j := 6; j >= 0; j-- {
		for i := weeksInLastSixMonths + 1; i >= 0; i-- {
			if i == weeksInLastSixMonths+1 {
				printDayCol(j)
			}
			if col, ok := columns[i]; ok {
				if len(col) > j {
					//special case today
					if i == weeksInLastSixMonths && j == calcOffset()-1 {
						printCell(col[j], false)

					} else {
						printCell(col[j], false)

					}
					continue
				}
			}
			printCell(0, false)
		}
		fmt.Printf("\n")
	}
}

// printMonths prints the month names in the first line, determining when the month
// changed between switching weeks
func printMonths() {
	week := getBeginningOfDay(time.Now()).Add(-(daysInLastSixMonths * time.Hour * 24))
	month := week.Month()
	fmt.Printf("         ")
	for {
		if week.Month() != month {
			fmt.Printf("%s ", week.Month().String()[:3])
			month = week.Month()
		} else {
			fmt.Printf("    ")
		}

		week = week.Add(7 * time.Hour * 24)
		if week.After(time.Now()) {
			break
		}
	}
	fmt.Printf("\n")
}

// printDayCol given the day number (0 is Sunday) prints the day name,
// alternating the rows (prints just 2,4,6)
func printDayCol(day int) {
	out := "     "
	switch day {
	case 1:
		out = " Mon "
	case 3:
		out = " Wed "
	case 5:
		out = " Fri "
	}

	fmt.Print(out)
}
