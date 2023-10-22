package stat

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"go_cli/utils"
	"sort"
	"time"
)

const outOfRange = 99999
const daysInLastSixMonths = 183
const weeksInLastSixMonths = 26

type column []int

// Stats 计算并打印统计数据。
func Stats(email string) {
	commits := processRepositories(email)
	printCommitsStats(commits)
}

// getBeginningOfDay 给定一个`time.Time` 计算当天的开始时间
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

	offset := GetOffset()
	// 遍历提交历史记录
	err = iterator.ForEach(func(c *object.Commit) error {
		daysAgo := countDaysSinceDate(c.Author.When) + offset
		println(daysAgo)
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

	// 获取repositories
	repos := utils.ParseFileToSlice()

	//初始化 在过去 6 个月内每天的contributions
	commits := make(map[int]int, daysInLastSixMonths)
	for i := daysInLastSixMonths; i > 0; i-- {
		commits[i] = 0
	}
	//统计每个repository的contributions
	for _, path := range repos {
		fillCommits(email, path, &commits)
	}

	return commits
}

// GetOffset 确定并返回填充统计图最后一行所缺少的天数
func GetOffset() (offset int) {
	switch time.Now().Weekday() {
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
	return
}

// printCell 打印一个单元格值，根据值contributions和`today`标志以不同的格式打印它。
func printCell(contributions int, today bool) {
	escape := "\033[0;37;30m"
	switch {
	case contributions > 0 && contributions < 5:
		escape = "\033[1;42m"
	case contributions >= 5 && contributions < 10:
		escape = "\033[1;30;42m"
	case contributions >= 10:
		escape = "\033[100m"
	}

	if today {
		escape = "\033[1;37;45m"
	}
	// 当天无contribution
	if contributions == 0 {
		fmt.Printf(escape + "  - " + "\033[0m")
		return
	}

	//控制对齐
	str := "  %d "
	switch {
	case contributions >= 10:
		str = " %d "
	case contributions >= 100:
		str = "%d "
	}

	fmt.Printf(escape+str+"\033[0m", contributions)
}

// printCommitsStats prints the commits stats
func printCommitsStats(commits map[int]int) {
	keys := sortMapIntoSlice(commits)
	columns := buildCols(keys, commits)
	printCells(columns)
}

// sortMapIntoSlice 返回已排序map的key的切片
func sortMapIntoSlice(commits map[int]int) []int {
	var keys []int
	for k := range commits {
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
					if i == weeksInLastSixMonths && j == GetOffset()-1 {
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

// printMonths 在第一行打印月份名称，确定月份在切换周之间发生变化的时间
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

// printDayCol 给定行号(0为周日）打印日期名称，交替行（仅打印 2,4,6）
func printDayCol(row int) {
	out := "     "
	switch row {
	case 1:
		out = " Mon "
	case 3:
		out = " Wed "
	case 5:
		out = " Fri "
	}
	fmt.Print(out)
}
