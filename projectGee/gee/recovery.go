package gee

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strings"
)

/*
以中间件形式实现 recovery: 错误恢复
*/

// 获取触发 panic 的堆栈信息
func trace(message string) string {
	/*
			uintptr: 无符号整数类型
			其大小足以存储一个指针值
			在 Go 的运行时系统中用于 存储程序计数器 的值
	*/
	var pcs [32]uintptr
	/*
	Callers 用来返回调用栈的程序计数器
	第 0 为: Callers 本身
	第 1 为: 上一层的 tarce
	第 2 为: 再上一层的 defer func
	为保持日志简洁, 跳过前三个
	*/
	n := runtime.Callers(3, pcs[:]) // skip first 3 callers

	/*
	strings.Builder 类型的变量是一个结构体
	包含了一个 内部缓冲区 和一些 操作缓冲区的方法
	可以有效的构建大量的字符串,省去连接字符串时创建新字符串类型对象的开销
	*/
	var str strings.Builder
	str.WriteString(message + "\nTraceback:")
	for _, pc := range pcs[:n] {
		// 通过 FuncForPC 获取对应的函数
		fn := runtime.FuncForPC(pc)
		// 跳过 FileLine 获取对应函数的文件名和行号
		file, line := fn.FileLine(pc)
		str.WriteString(fmt.Sprintf("\n\t%s:%d\n", file, line))
	}
	return str.String()
}


func Recovery() HandlerFunc {
	return func(c *Context) {
		defer func() {
			if err := recover(); err != nil {
				message := fmt.Sprintf("%s", err)
				log.Printf("%s\n\n", trace(message))
				c.Fail(http.StatusInternalServerError, "Internal Server Error")
			}
		}()

		c.Next()
	}
}