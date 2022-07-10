package reqit

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMergeHeadersFrom(t *testing.T) {
	t.Run("MergeClientHeaders", func(t *testing.T) {
		assert := assert.New(t)
		require := require.New(t)

		// pre-configured gloabl headers in Client
		headerHName := "X-Header-H"
		headerHVal := []string{"Val1"}
		c := NewClient(ClientOptions{
			Headers: http.Header{
				headerHName: headerHVal,
			},
		})

		// new request
		request, err := http.NewRequest(http.MethodGet, "", nil)
		require.NoError(err)
		require.NotNil(request)

		// new headers specified per request
		headerXName := "X-Header-X"
		headerXVal := []string{"Val2"}
		headerHVal2 := []string{"Val1", "Val2"}
		request.Header = http.Header{
			headerXName: headerXVal,
			headerHName: headerHVal2,
		}

		mergeHeadersFrom(request.Header, c.Options.Headers)

		// assert for headers in request
		assert.Equal("Val1", request.Header.Get(headerHName))
		assert.Equal([]string{"Val1", "Val2"}, request.Header[headerHName])
		assert.Equal(headerXVal[0], request.Header.Get(headerXName))
		assert.Equal(headerXVal, request.Header[headerXName])
	})
}

func TestParseDesiredType(t *testing.T) {
	c := NewClient()
	request := newRequest(c, http.MethodGet, "", nil)

	t.Run("Struct", func(t *testing.T) {
		assert := assert.New(t)

		type tStruct struct {
			Key string `json:"key"`
		}

		var ts tStruct
		assert.Equal(reflect.Struct, request.parseDestKind(&ts))
	})

	t.Run("Map", func(t *testing.T) {
		assert := assert.New(t)

		var ts map[string]string
		assert.Equal(reflect.Map, request.parseDestKind(&ts))
	})

	t.Run("Slice", func(t *testing.T) {
		assert := assert.New(t)

		var ts []string
		assert.Equal(reflect.Slice, request.parseDestKind(&ts))
	})

	t.Run("String", func(t *testing.T) {
		assert := assert.New(t)

		var ts string
		assert.Equal(reflect.String, request.parseDestKind(&ts))
	})
}
