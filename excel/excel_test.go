package excel

import (
	"log"
	"testing"
)

type TestModel struct {
	Name   string  `x-col:"姓名"`
	Phone  string  `x-col:"手机号"`
	Hobby  string  `x-col:"兴趣"`
	ID     int     `x-col:"ID"`
	Credit float64 `x-col:"学分"`
}

func TestRead(t *testing.T) {

	arr := make([]TestModel, 0)
	provider := NewReaderProvider("D:\\projects\\20230906\\test.xlsx")
	defer provider.Close()
	provider.ReadModel(&arr)

	log.Printf("%v \r\n", arr)

}
