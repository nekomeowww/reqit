package reqit

import (
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

func MockServer(port int) (*echo.Echo, func() string) {
	server := echo.New()

	start := func() string {
		// 如果传入的参数 port 为 0，则随机分配端口
		l, err := net.Listen("tcp", net.JoinHostPort("127.0.0.1", strconv.Itoa(port)))
		if err != nil {
			panic(err)
		}

		go func() {
			err := http.Serve(l, server)
			if err != nil {
				panic(err)
			}
		}()

		time.Sleep(100 * time.Millisecond) // 防止访问先于路由运行
		return l.Addr().String()
	}

	return server, start
}
