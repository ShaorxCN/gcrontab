package middleware

import (
	"bytes"
	"fmt"
	"gcrontab/utils"
	"gcrontab/web/response"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var (
	dunno     = []byte("???")
	centerDot = []byte("·")
	dot       = []byte(".")
	slash     = []byte("/")
)

// 返回当前程序指针指向信息 也就是函数的信息
func function(pc uintptr) []byte {
	fn := runtime.FuncForPC(pc)

	// 如果没有找到函数，返回 "???"
	if fn == nil {
		return dunno
	}

	name := []byte(fn.Name())

	// 根据`/`去除路径
	if lastslash := bytes.LastIndex(name, slash); lastslash >= 0 {
		name = name[lastslash+1:]
	}

	// 去除包名
	if period := bytes.Index(name, dot); period >= 0 {
		name = name[period+1:]
	}
	name = bytes.Replace(name, centerDot, dot, -1)
	return name
}

// 获取 nth 行信息
func source(lines [][]byte, n int) []byte {
	// 行信息1开始 数组是0开始 所以先减1
	n--
	if n < 0 || n >= len(lines) {
		return dunno
	}
	return bytes.TrimSpace(lines[n])
}

// 前三帧一般是main runtime 以及部分init skip 3
func stack(skip int) []byte {
	buf := new(bytes.Buffer)
	var lines [][]byte
	var lastFile string
	for i := skip; ; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		fmt.Fprintf(buf, "%s:%d (0x%x)\n", file, line, pc)
		if file != lastFile {
			data, err := os.ReadFile(file)
			if err != nil {
				continue
			}
			lines = bytes.Split(data, []byte{'\n'})
			lastFile = file
		}
		fmt.Fprintf(buf, "\t%s: %s\n", function(pc), source(lines, line))
	}
	return buf.Bytes()
}

func Recovery(log *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				stack := stack(3)
				httprequest, _ := httputil.DumpRequest(c.Request, false)
				log.Errorf("[Recovery] %s panic recovered:\n%s\n%s\n%s", utils.Now(), string(httprequest), err, stack)
				errorContainer := response.NewSystemFailedBaseResponse()
				c.AbortWithStatusJSON(http.StatusOK, errorContainer)
			}
		}()
		c.Next()
	}
}
