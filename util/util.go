package util

import (
	"bytes"
	"crypto/md5"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"strings"
	"unsafe"

	"github.com/google/uuid"
)

func ShaString(s string) string {
	hash := sha256.Sum256([]byte(s))
	return fmt.Sprintf("%x", hash[:])
}
func ShaReader(reader io.Reader) ([]byte, int64, error) {
	return SaveAndSha(reader, "")
}

func SaveAndSha(reader io.Reader, dist string) ([]byte, int64, error) {
	var file *os.File
	var err error
	if dist != "" {
		file, err = os.Create(dist)
		if err != nil {
			return nil, 0, err
		}
		defer file.Close()
	}
	hash := sha256.New()
	size := int64(0)
	block := make([]byte, hash.BlockSize())
	for {
		i, err := reader.Read(block)
		if err != nil {
			if err != io.EOF {
				return nil, size, err
			}
			break
		}
		if dist != "" {
			file.Write(block[:i])
		}
		hash.Write(block[:i])
		size += int64(i)
	}
	return hash.Sum(nil), size, nil
}

func UUID() string {
	var buffer bytes.Buffer
	for _, chr := range uuid.New().String() {
		if "-" == string(chr) {
			continue
		}
		buffer.WriteString(string(chr))
	}
	return buffer.String()
}
func ByteArrayToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func isBlankString(s string) bool {
	return strings.Trim(s, " ") == ""
}
func Md5String(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}
