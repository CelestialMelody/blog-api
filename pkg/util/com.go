package util

import (
	"github.com/unknwon/com"
	"strconv"
)

// com 包没有StrTo.MustUint() 添加自己的方法
// 现在不考虑使用uint了

type StrTo com.StrTo

func (f StrTo) Uint() (uint, error) {
	v, err := strconv.ParseUint(f.String(), 10, 0)
	return uint(v), err
}

func (f StrTo) MustUint() uint {
	v, _ := f.Uint()
	return v
}

func (f StrTo) Exist() bool {
	return string(f) != string(0x1E) // unknwon/com: 0x1E is 30 and is the record separator in ASCII table
}

func (f StrTo) String() string {
	if f.Exist() {
		return string(f)
	}
	return ""
}
