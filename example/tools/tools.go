package tools

import (
	"strconv"
	"time"
)

//查看map[string]interface{}是否存在这个KEY
func In_mapSI(str string, maps map[string]interface{}) bool {
	if len(maps) < 1 {
		return false
	}
	for k, _ := range maps {
		if k == str {
			return true
		}
	}
	return false
}

//转换成字符串
func ToString(args ...interface{}) string {
	result := ""
	for _, arg := range args {
		switch val := arg.(type) {
		case int:
			result += strconv.Itoa(val)
		case int8:
			result += strconv.Itoa(int(val))
		case int64:
			result += strconv.Itoa(int(val))
		case string:
			result += val
		case float64:
			result += strconv.FormatFloat(val, 'f', -1, 64)
		case float32:
			result += strconv.FormatFloat(float64(val), 'f', -1, 64)
		case time.Time:
			result += val.Format("2006-01-02 15:04:05")
		}
	}
	return result
}

/*
* php的in_array  to go的版本
 */
func In_array(str string, slice []string) bool {
	for i := 0; i < len(slice); i++ {
		if slice[i] == str {
			return true
		}
	}
	return false
}

/*
* PHP的substr函数
 */
func Substr(str string, start, length int) string {
	rs := []rune(str)
	rl := len(rs)
	end := 0

	if start < 0 {
		start = rl - 1 + start
	}
	end = start + length

	if start > end {
		start, end = end, start
	}

	if start < 0 {
		start = 0
	}
	if start > rl {
		start = rl
	}
	if end < 0 {
		end = 0
	}
	if end > rl {
		end = rl
	}

	return string(rs[start:end])
}

func Test(str string) string {
	return "我是一个测试函数：" + str
}
