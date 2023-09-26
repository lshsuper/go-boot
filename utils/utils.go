package utils

import (
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

//全局标量

var (
	excelChar = []string{"", "A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}
)

func UUID() string {

	id := uuid.NewV4()
	return strings.ReplaceAll(id.String(), "-", "")

}

//ErrorToString recover错误，转string
func ErrorToString(r interface{}) string {
	switch v := r.(type) {
	case error:
		return v.Error()
	default:
		return r.(string)
	}
}

//GetBrowser 获取浏览器类型
func GetBrowser(c *gin.Context) string {

	userAgent := strings.ToLower(c.GetHeader("User-Agent"))
	if strings.Index(userAgent, "opera") >= 0 {
		return "opera"
	}

	if strings.Index(userAgent, "edge") >= 0 {
		return "edge"
	}

	if strings.Index(userAgent, "firefox") >= 0 {
		return "firefox"
	}

	if strings.Index(userAgent, "chrome") >= 0 {
		return "chrome"
	}

	if strings.Index(userAgent, "safari") >= 0 {
		return "safari"
	}

	//IE
	if strings.Index(userAgent, "msie") >= 0 &&
		strings.Index(userAgent, "compatible") >= 0 &&
		strings.Index(userAgent, "opera") < 0 {

		//IE
		regex := regexp.MustCompile("msie (\\d+\\.\\d+);")
		str := regex.FindString(userAgent)
		switch str {

		case "7":
			return "IE7"
		case "8":
			return "IE8"
		case "9":
			return "IE9"
		case "10":
			return "IE10"
		case "11":
			return "IE11"
		default:
			return "IE"

		}

	}

	return ""

}

func GetCol(num int) string {

	var cols string
	v := num
	for v > 0 {
		k := v % 26
		if k == 0 {
			k = 26
		}
		v = (v - k) / 26
		cols = excelChar[k] + cols
	}
	return cols
}

//FUpper 字符串首字母大写
func FUpper(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

// FLower 字符串首字母小写
func FLower(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToLower(s[:1]) + s[1:]
}
func GetCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0])) //返回绝对路径  filepath.Dir(os.Args[0])去除最后一个元素的路径
	if err != nil {
		log.Fatal(err)
	}
	return strings.Replace(dir, "\\", "/", -1) //将\替换成/
}
