package excel

import (
	"fmt"
	"gitee.com/lshsuper/go-boot/utils"
)

type MergeRule struct {
	StartCol int
	EndCol   int
	StartRow int
	EndRow   int
}

func (m MergeRule) Start() string {
	return fmt.Sprintf("%s%d", utils.GetCol(m.StartCol-1), m.StartRow)
}
func (m MergeRule) End() string {
	return fmt.Sprintf("%s%d", utils.GetCol(m.EndCol-1), m.EndRow)
}
