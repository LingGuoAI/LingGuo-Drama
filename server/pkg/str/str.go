// Package str 字符串辅助方法
package str

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"github.com/gertd/go-pluralize"
	"github.com/iancoleman/strcase"
	"github.com/w3liu/go-common/constant/timeformat"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

// Plural 转为复数 user -> users
func Plural(word string) string {
	return pluralize.NewClient().Plural(word)
}

// Singular 转为单数 users -> user
func Singular(word string) string {
	return pluralize.NewClient().Singular(word)
}

// Snake 转为 snake_case，如 TopicComment -> topic_comment
func Snake(s string) string {
	return strcase.ToSnake(s)
}

// Camel 转为 CamelCase，如 topic_comment -> TopicComment
func Camel(s string) string {
	return strcase.ToCamel(s)
}

// LowerCamel 转为 lowerCamelCase，如 TopicComment -> topicComment
func LowerCamel(s string) string {
	return strcase.ToLowerCamel(s)
}

var num int64

// 生成24位订单号
// 前面17位代表时间精确到毫秒，中间3位代表进程id，最后4位代表序号
// 修改版：17位时间  4位随机数
func Generate() string {
	t := time.Now()
	s := t.Format(timeformat.Continuity)
	m := t.UnixNano()/1e6 - t.UnixNano()/1e9*1e3
	ms := sup(m, 3)
	//p := os.Getpid() % 1000
	//ps := sup(int64(p), 3)
	//i := atomic.AddInt64(&num, 1)
	//r := i % 10000
	//rs := sup(r, 4)
	rs := Createcode()
	n := fmt.Sprintf("%s%s%s", s, ms, rs)
	return n
}

// 对长度不足n的数字前面补0
func sup(i int64, n int) string {
	m := fmt.Sprintf("%d", i)
	for len(m) < n {
		m = fmt.Sprintf("0%s", m)
	}
	return m
}

// PriceFloatToUint float 小数 乘 100
func PriceFloatToUint(nums float64) uint64 {
	// 处理负数情况
	if nums < 0 {
		return 0 // 价格不应该为负数，返回0
	}

	str := fmt.Sprintf("%.2f", nums)
	num, err := strconv.ParseUint(strings.Replace(str, ".", "", 1), 10, 64)
	if err != nil {
		// 解析失败时返回0
		return 0
	}
	return num
}

func Createcode() string {
	return fmt.Sprintf("%04v", rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(10000)) //这里面前面的04v是和后面的1000相对应的
}

// GenerateUniqueFileName 生成唯一文件名
func GenerateUniqueFileName(originalFileName string) string {
	// 获取当前时间的Unix时间戳，以确保文件名的唯一性
	timestamp := time.Now().Unix()

	// 生成一个随机数，以防止相同时间戳的文件名冲突
	randomNum := rand.Intn(10000) // 这里可以根据需要设置一个随机数的上限

	// 获取文件的扩展名
	fileExt := getFileExtension(originalFileName)

	// 拼接文件名
	uniqueFileName := fmt.Sprintf("%d_%d.%s", timestamp, randomNum, fileExt)

	return uniqueFileName
}

// getFileExtension 获取文件的扩展名
func getFileExtension(fileName string) string {
	// 根据最后一个点（.）来分割文件名，获取扩展名
	parts := strings.Split(fileName, ".")
	if len(parts) > 1 {
		return parts[len(parts)-1]
	}
	return ""
}

func GenerateRandomFourDigitString() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%04d", rand.Intn(10000))
}

func MaskTitle(title string, n int) string {
	runes := []rune(title)
	length := len(runes)
	if length <= n {
		return title
	}
	// 保留前n个字符，后面用*代替
	masked := string(runes[:n]) + strings.Repeat("*", length-n)
	return masked
}

// TrimQuotes 去除单双引号
func TrimQuotes(s string) string {
	// 定义要修剪的引号字符集
	quotes := `'"‘’“”`

	// 使用循环确保移除所有前后的引号字符
	for {
		if len(s) == 0 {
			break
		}
		first := s[0]
		if strings.ContainsRune(quotes, rune(first)) {
			s = s[1:]
		} else {
			break
		}
	}

	for {
		if len(s) == 0 {
			break
		}
		last := s[len(s)-1]
		if strings.ContainsRune(quotes, rune(last)) {
			s = s[:len(s)-1]
		} else {
			break
		}
	}

	return s
}

func SHA1(s string) string {
	o := sha1.New()
	o.Write([]byte(s))
	return hex.EncodeToString(o.Sum(nil))
}
