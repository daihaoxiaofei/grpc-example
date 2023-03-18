// Package closehelp 用于优雅关闭已开启的服务
package closehelp

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type closeFunc func(context.Context)

var closeFuncArr []closeFunc

func init() {
	closeFuncArr = make([]closeFunc, 0)
}

// Stop 手动调用 用于正常执行完成就退出的程序
func Stop() {
	// 优雅退出
	// 创建一个超时的context
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	// 根据加入顺序倒序关闭
	for i := len(closeFuncArr) - 1; i > -1; i-- {
		select {
		case <-ctx.Done():
			return
		default:
			closeFuncArr[i](ctx)
		}
	}
}

// Register 注册优雅关闭
func Register(f func(context.Context)) {
	closeFuncArr = append(closeFuncArr, f)
}

// SignalClose 只能由系统信号关闭的程序 如网络服务等
// 这个方法是阻塞的
func SignalClose() {
	// 关闭信号
	ch := make(chan os.Signal, 1)
	// kill (no param) default send signal.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can"t be caught, so don't need to add it
	// syscall.SIGHUP 连接终端退出时
	// syscall.SIGQUIT 和SIGINT类似, 但由QUIT字符(通常是Ctrl+)来控制. 进程在因收到SIGQUIT退出时会产生core文件, 在这个意义上类似于一个程序错误信号。
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP, syscall.SIGQUIT)

	fmt.Println(`接收到系统停止信号`, <-ch, `逐步关闭程序中...`) // 接收信号之前阻塞在这里
	Stop()
}
