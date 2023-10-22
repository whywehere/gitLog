package stat

import (
	"fmt"
	"testing"
	"time"
)

func TestTimePackage(t *testing.T) {
	fmt.Println(time.Now())
	println(time.Now().Weekday())
}
