package excel

import (
	"fmt"
	"gitee.com/lshsuper/go-boot/utils"
	"github.com/xuri/excelize/v2"
	"io"
)

//writerProvider 写提供器
type writerProvider struct {
	fs *excelize.File
}

//NewWriterProvider 构建写器
func NewWriterProvider() *writerProvider {
	return &writerProvider{fs: excelize.NewFile()}
}

//simpleStyle 基础样式
func (e *writerProvider) simpleStyle() int {

	style, _ := e.fs.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			WrapText:   true,
			Vertical:   "center",
			Horizontal: "center",
		},
	})

	return style

}

//Save 保存
func (e *writerProvider) Save(fullPath string) {
	e.fs.Path = fullPath
	e.fs.Save()
}

//Buffer 转流
func (e *writerProvider) Buffer(r io.Reader) (err error) {

	r, err = e.fs.WriteToBuffer()
	return

}

//Writer 写入
func (e *writerProvider) Writer(header []Header, body [][]string, rules []MergeRule) {

	sName := e.fs.GetSheetName(0)
	startIndex := 1
	//1.填充头
	simpleStyle := e.simpleStyle()
	e.fs.SetColStyle(sName, fmt.Sprintf("A:%s", utils.GetCol(len(header)-1)), simpleStyle)
	for k, v := range header {
		e.fs.SetCellValue(sName, fmt.Sprintf("%s1", utils.GetCol(k)), v.Name)
		e.fs.SetColWidth(sName, utils.GetCol(k), utils.GetCol(k), v.Width)
	}
	e.fs.SetRowHeight(sName, 1, 20)

	//2.填充主题元素
	for _, v := range body {

		startIndex++
		if err := e.fs.SetSheetRow(sName, fmt.Sprintf("A%d", startIndex), &v); err != nil {
			panic(err)
		}
		e.fs.SetRowHeight(sName, startIndex, 20)

	}

	if rules != nil {
		for _, v := range rules {
			e.fs.MergeCell(sName, v.Start(), v.End())
		}
	}

}
