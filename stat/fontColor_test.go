package stat

import (
	"fmt"
	"testing"
)

func TestFontColor(t *testing.T) {
	escape := "\033[42m"
	reset := "\033[0m" // 用于重置文本样式和颜色

	fmt.Print(escape)
	fmt.Println("This text is bold with black text on a green background.")
	fmt.Print(reset)
	fmt.Println("This text has the default style.")
}
