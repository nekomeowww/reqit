package reqit

import (
	"encoding/json"
	"net/http"
	"testing"

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
	t.Run("GetData", func(t *testing.T) {
		assert := assert.New(t)
		require := require.New(t)

		handleGet := func(writer http.ResponseWriter, request *http.Request) {
			writer.Header().Add("Content-Type", "application/json")
			writer.WriteHeader(http.StatusOK)

			resp := map[string]interface{}{
				"message": "Hello, reqit!",
			}
			jsonData, _ := json.Marshal(resp)
			_, _ = writer.Write(jsonData)
		}

		server := http.NewServeMux()
		server.Handle("/get", http.HandlerFunc(handleGet))

		serverAddr := MockServer(server, 0)

		c := NewClient()
		var resp map[string]interface{}
		err := c.Get(serverAddr+"/get", nil).Result(&resp).Error()
		require.NoError(err)
		assert.Equal("Hello, reqit!", resp["message"])
	})

	t.Run("GetWithBaseURLData", func(t *testing.T) {
		assert := assert.New(t)
		require := require.New(t)

		handleGet := func(writer http.ResponseWriter, request *http.Request) {
			writer.Header().Add("Content-Type", "application/json")
			writer.WriteHeader(http.StatusOK)

			resp := map[string]interface{}{
				"message": "Hello, reqit!",
			}
			jsonData, _ := json.Marshal(resp)
			_, _ = writer.Write(jsonData)
		}

		server := http.NewServeMux()
		server.Handle("/get", http.HandlerFunc(handleGet))
		serverAddr := MockServer(server, 0)

		c := NewClient(ClientOptions{
			BaseURL: serverAddr,
		})

		var resp map[string]interface{}
		err := c.Get("/get", nil).Result(&resp).Error()
		require.NoError(err)
		assert.Equal("Hello, reqit!", resp["message"])
	})
}
