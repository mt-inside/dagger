package value

import (
	"fmt"
	"strconv"

	au "github.com/logrusorgru/aurora/v3"
)

type Mode int

const (
	Pending   Mode = 0
	Available Mode = iota
)

type Value struct {
	Mode Mode
	Val  int64
}

func NewPending() Value {
	return Value{Pending, 0}
}
func NewAvailable(val int64) Value {
	return Value{Available, val}
}

func (v Value) String() string {
	if v.Mode == Pending {
		return "PENDING"
	}
	return strconv.FormatInt(v.Val, 10)
}

func (v Value) Print() {
	if v.Mode == Pending {
		fmt.Println(au.Bold(au.Magenta("PENDING")))
	}
	fmt.Println(au.Bold(au.Cyan(strconv.FormatInt(v.Val, 10))))
}
