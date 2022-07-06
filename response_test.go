package reqit

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strings"
	"testing"
)

func newRandomHashString(length ...int) string {
	bytes := make([]byte, 1024)
	_, _ = rand.Read(bytes)
	if len(length) != 0 {
		sliceLength := length[0]
		if length[0] > 64 {
			sliceLength = 64
		}
		if length[0] <= 0 {
			sliceLength = 64
		}

		return fmt.Sprintf("%x", sha256.Sum256(bytes))[:sliceLength]
	}

	return fmt.Sprintf("%x", sha256.Sum256(bytes))
}

func newHugeBenchmarkingJSONObjetc() []byte {
	object := make(map[string]interface{})

	for i := 0; i < 100; i++ {
		object[fmt.Sprintf("key:%d", i)] = fmt.Sprintf("value:%d:%s", i, newRandomHashString(64))
	}

	var err error
	jsonObjectData, err := json.Marshal(object)
	if err != nil {
		log.Fatal(err)
	}

	return jsonObjectData
}

var hugeJSONObjectReader = strings.NewReader(string(newHugeBenchmarkingJSONObjetc()))

func BenchmarkReadReaderWithIoReadAll(b *testing.B) {
	b.ReportAllocs()

	for n := 0; n < b.N; n++ {
		var m map[string]interface{}
		b, _ := io.ReadAll(hugeJSONObjectReader)
		_ = json.Unmarshal(b, &m)
	}
}

func BenchmarkReadReaderWithIoCopy(b *testing.B) {
	b.ReportAllocs()

	for n := 0; n < b.N; n++ {
		var m map[string]interface{}
		buffer := new(bytes.Buffer)
		_, _ = io.Copy(buffer, hugeJSONObjectReader)
		_ = json.Unmarshal(buffer.Bytes(), &m)
	}
}

func BenchmarkReadReaderWithBytesBufferReadFrom(b *testing.B) {
	b.ReportAllocs()

	for n := 0; n < b.N; n++ {
		var m map[string]interface{}
		buffer := new(bytes.Buffer)
		_, _ = buffer.ReadFrom(hugeJSONObjectReader)
		_ = json.Unmarshal(buffer.Bytes(), &m)
	}
}
