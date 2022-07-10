package reqit

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewClient(t *testing.T) {
	t.Run("EmptyParameters", func(t *testing.T) {
		require := require.New(t)
		require.NotPanics(func() {
			_ = NewClient()
		})
	})
}

func TestConfigureClientByClientOptions(t *testing.T) {
	t.Run("Headers", func(t *testing.T) {
		assert := assert.New(t)

		headers := http.Header{
			"X-A": []string{"A"},
		}

		c := NewClient()
		configureClientByClientOptions(c, ClientOptions{
			Headers: headers,
		})

		assert.Equal(headers.Get("X-A"), c.Options.Headers.Get("X-A"))
		assert.Equal(headers["X-A"], c.Options.Headers["X-A"])
	})

	t.Run("BaseURL", func(t *testing.T) {
		assert := assert.New(t)

		c := NewClient()
		configureClientByClientOptions(c, ClientOptions{
			BaseURL: "http://example.com",
		})

		assert.Equal("http://example.com", c.Options.BaseURL)
	})
}

func TestGet(t *testing.T) {
	t.Run("WithBaseURL", func(t *testing.T) {
		t.Run("application/json", func(t *testing.T) {
			assert := assert.New(t)
			require := require.New(t)

			handleGet := func(c echo.Context) error {
				return c.JSON(http.StatusOK, map[string]interface{}{
					"message": "Hello, reqit!",
				})
			}

			server, start := MockServer(0)
			server.GET("/get", handleGet)
			serverAddr := start()

			c := NewClient(ClientOptions{
				BaseURL: serverAddr,
			})

			var resp map[string]interface{}
			res := c.Get("/get", nil).Send()
			status, err := res.Result(&resp)
			require.NoError(err)
			require.Equal(http.StatusOK, status)
			assert.Equal("Hello, reqit!", resp["message"])
		})

		t.Run("text/plain", func(t *testing.T) {
			assert := assert.New(t)
			require := require.New(t)

			handleGet := func(c echo.Context) error {
				return c.String(http.StatusOK, "Hello, reqit!")
			}

			server, start := MockServer(0)
			server.GET("/get", handleGet)
			serverAddr := start()

			c := NewClient(ClientOptions{
				BaseURL: serverAddr,
			})

			var resp string
			res := c.Get("/get", nil).Send()
			status, err := res.Result(&resp)
			require.NoError(err)
			require.Equal(http.StatusOK, status)
			assert.Equal("Hello, reqit!", resp)
		})
	})
	t.Run("WithoutBaseURL", func(t *testing.T) {
		t.Run("application/json", func(t *testing.T) {
			assert := assert.New(t)
			require := require.New(t)

			handleGet := func(c echo.Context) error {
				return c.JSON(http.StatusOK, map[string]interface{}{
					"message": "Hello, reqit!",
				})
			}

			server, start := MockServer(0)
			server.GET("/get", handleGet)
			serverAddr := start()

			c := NewClient()
			var resp map[string]interface{}
			res := c.Get(serverAddr+"/get", nil).Send()
			status, err := res.Result(&resp)
			require.NoError(err)
			require.Equal(http.StatusOK, status)
			assert.Equal("Hello, reqit!", resp["message"])
		})

		t.Run("text/plain", func(t *testing.T) {
			assert := assert.New(t)
			require := require.New(t)

			handleGet := func(c echo.Context) error {
				return c.String(http.StatusOK, "Hello, reqit!")
			}

			server, start := MockServer(0)
			server.GET("/get", handleGet)
			serverAddr := start()

			c := NewClient()
			var resp string
			res := c.Get(serverAddr+"/get", nil).Send()
			status, err := res.Result(&resp)
			require.NoError(err)
			require.Equal(http.StatusOK, status)
			assert.Equal("Hello, reqit!", resp)
		})
	})
}

func TestPost(t *testing.T) {
	t.Run("WithBaseURL", func(t *testing.T) {
		t.Run("application/json", func(t *testing.T) {
			assert := assert.New(t)
			require := require.New(t)

			// 参数结构体定义
			type param struct {
				Message string `json:"message"`
			}

			type resp struct {
				Result string `json:"result"`
			}

			// mock POST 处理器
			handlePost := func(c echo.Context) error {
				var p param
				// 解析参数绑定
				err := c.Bind(&p)
				if err != nil {
					// 错误处理
					return c.JSON(http.StatusBadRequest, map[string]interface{}{
						"message": fmt.Sprintf("Invalid request, %v", err),
					})
				}

				// 响应
				return c.JSON(http.StatusOK, resp{
					Result: p.Message,
				})
			}

			// 创建 mock 服务器
			server, start := MockServer(0)
			// 路由挂载
			server.GET("/post", handlePost)
			// 开始服务器，获取 mock 服务器的地址
			serverAddr := start()

			// 新建 HTTP 客户端
			c := NewClient(ClientOptions{
				BaseURL: serverAddr, // 指定请求通用的基准 URL
			})

			// POST 参数
			p := param{Message: "Hello, reqit!"}
			// 响应结构体
			var r resp

			// 创建 POST 请求并尝试解析
			res := c.Post("/post", nil).WithBody(p).Send()
			status, err := res.Result(&r)
			// 断言错误
			require.NoError(err)
			// 断言状态码
			require.Equal(http.StatusOK, status)
			// 断言响应结构体
			assert.Equal("Hello, reqit!", r.Result)
		})
		t.Run("text/plain", func(t *testing.T) {
			assert := assert.New(t)
			require := require.New(t)

			// 参数结构体定义
			type param struct {
				Message string `json:"message"`
			}

			// mock POST 处理器
			handlePost := func(c echo.Context) error {
				var p param
				// 解析参数绑定
				err := c.Bind(&p)
				if err != nil {
					// 错误处理
					return c.JSON(http.StatusBadRequest, map[string]interface{}{
						"message": fmt.Sprintf("Invalid request, %v", err),
					})
				}

				// 响应
				return c.String(http.StatusOK, p.Message)
			}

			// 创建 mock 服务器
			server, start := MockServer(0)
			// 路由挂载
			server.GET("/post", handlePost)
			// 开始服务器，获取 mock 服务器的地址
			serverAddr := start()

			// 新建 HTTP 客户端
			c := NewClient(ClientOptions{
				BaseURL: serverAddr, // 指定请求通用的基准 URL
			})

			// POST 参数
			p := param{Message: "Hello, reqit!"}
			// 响应结构体
			var r string

			// 创建 POST 请求并尝试解析
			res := c.Post("/post", nil).WithBody(p).Send()
			status, err := res.Result(&r)
			// 断言错误
			require.NoError(err)
			// 断言状态码
			require.Equal(http.StatusOK, status)
			// 断言响应结构体
			assert.Equal("Hello, reqit!", r)
		})
	})
	t.Run("WithoutBaseURL", func(t *testing.T) {
		t.Run("application/json", func(t *testing.T) {
			assert := assert.New(t)
			require := require.New(t)

			// 参数结构体定义
			type param struct {
				Message string `json:"message"`
			}

			type resp struct {
				Result string `json:"result"`
			}

			// mock POST 处理器
			handlePost := func(c echo.Context) error {
				var p param
				// 解析参数绑定
				err := c.Bind(&p)
				if err != nil {
					// 错误处理
					return c.JSON(http.StatusBadRequest, map[string]interface{}{
						"message": fmt.Sprintf("Invalid request, %v", err),
					})
				}

				// 响应
				return c.JSON(http.StatusOK, resp{
					Result: p.Message,
				})
			}

			// 创建 mock 服务器
			server, start := MockServer(0)
			// 路由挂载
			server.GET("/post", handlePost)
			// 开始服务器，获取 mock 服务器的地址
			serverAddr := start()

			// 新建 HTTP 客户端
			c := NewClient()

			// POST 参数
			p := param{Message: "Hello, reqit!"}
			// 响应结构体
			var r resp

			// 创建 POST 请求并尝试解析
			res := c.Post(serverAddr+"/post", nil).WithBody(p).Send()
			status, err := res.Result(&r)
			// 断言错误
			require.NoError(err)
			// 断言状态码
			require.Equal(http.StatusOK, status)
			// 断言响应结构体
			assert.Equal("Hello, reqit!", r.Result)
		})
		t.Run("text/plain", func(t *testing.T) {
			assert := assert.New(t)
			require := require.New(t)

			// 参数结构体定义
			type param struct {
				Message string `json:"message"`
			}

			// mock POST 处理器
			handlePost := func(c echo.Context) error {
				var p param
				// 解析参数绑定
				err := c.Bind(&p)
				if err != nil {
					// 错误处理
					return c.JSON(http.StatusBadRequest, map[string]interface{}{
						"message": fmt.Sprintf("Invalid request, %v", err),
					})
				}

				// 响应
				return c.String(http.StatusOK, p.Message)
			}

			// 创建 mock 服务器
			server, start := MockServer(0)
			// 路由挂载
			server.GET("/post", handlePost)
			// 开始服务器，获取 mock 服务器的地址
			serverAddr := start()

			// 新建 HTTP 客户端
			c := NewClient()

			// POST 参数
			p := param{Message: "Hello, reqit!"}
			// 响应结构体
			var r string

			// 创建 POST 请求并尝试解析
			res := c.Post(serverAddr+"/post", nil).WithBody(p).Send()
			status, err := res.Result(&r)
			// 断言错误
			require.NoError(err)
			// 断言状态码
			require.Equal(http.StatusOK, status)
			// 断言响应结构体
			assert.Equal("Hello, reqit!", r)
		})
	})
}
