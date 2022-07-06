package reqit

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMapToQuery(t *testing.T) {
	assert := assert.New(t)

	queryMap := make(map[string]string)
	queryMap["param1"] = "val1"
	queryMap["param2"] = "1"

	query := mapToQuery(queryMap)
	assert.Equal(query.Get("param1"), "val1")
	assert.Equal(query.Get("param2"), "1")
}
