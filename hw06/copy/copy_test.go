package copy

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

var FromFilePath string
var ToFilePath string

func init() {
	dir := getTestDir()
	FromFilePath = filepath.Join(dir, "from_file.txt")
	ToFilePath = filepath.Join(dir, "to_file.txt")

	err := createFile(FromFilePath)
	if err != nil {
		log.Fatal(err)
	}
}

func getTestDir() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	dir = filepath.Join(dir, "test_dir")
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.Mkdir(dir, os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
	}
	return dir
}

func createFile(filePath string) error {
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}

	defer f.Close()

	return nil
}

func writeFile(filePath string, text []byte) {
	file, err := os.Create(filePath)
	if err != nil {
		log.Fatal(err)
	}
	_, err = file.Write(text)
	if err != nil {
		log.Panicf("failed to write: %v", err)
	}
	file.Close()
}

func readFile(filePath string) string {
	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	return string(b)
}

func TestCopy_NegativeLimit(t *testing.T) {
	err := Copy("", "", -5, 100)
	assert.Error(t, err)
}

func TestCopy_NegativeOffset(t *testing.T) {
	err := Copy("", "", 0, -10)
	assert.Error(t, err)
}

func TestCopy_ErrorFromFilePath(t *testing.T) {
	from := FromFilePath + "incorrect_path"
	err := Copy(from, ToFilePath, 0, 0)
	assert.Error(t, err)
}

func TestCopy_OffsetOverFromFile(t *testing.T) {
	err := Copy(FromFilePath, ToFilePath, 0, 10000)
	assert.Error(t, err)
}

func TestCopy_LimitOverFromFile(t *testing.T) {
	err := Copy(FromFilePath, ToFilePath, 10000, 0)
	assert.Error(t, err)
}

func TestCopy_AllFile(t *testing.T) {
	text := "123456789"
	writeFile(FromFilePath, []byte(text))
	err := Copy(FromFilePath, ToFilePath, 0, 0)
	assert.NoError(t, err)

	writedText := readFile(ToFilePath)
	assert.Equal(t, text, writedText)
}

func TestCopy_AllFileWithSmallBufferSize(t *testing.T) {
	realBufferSize := maxBufferSize
	maxBufferSize = 2
	defer func() {
		maxBufferSize = realBufferSize
	}()
	text := "123456789"
	writeFile(FromFilePath, []byte(text))
	err := Copy(FromFilePath, ToFilePath, 0, 0)
	assert.NoError(t, err)

	writedText := readFile(ToFilePath)
	assert.Equal(t, text, writedText)
}

func TestCopy_WithOffset(t *testing.T) {
	text := "123456789"
	writeFile(FromFilePath, []byte(text))
	err := Copy(FromFilePath, ToFilePath, 0, 2)
	assert.NoError(t, err)

	writedText := readFile(ToFilePath)
	assert.Equal(t, "3456789", writedText)
}

func TestCopy_WithLimit(t *testing.T) {
	text := "123456789"
	writeFile(FromFilePath, []byte(text))
	err := Copy(FromFilePath, ToFilePath, 5, 0)
	assert.NoError(t, err)

	writedText := readFile(ToFilePath)
	assert.Equal(t, "12345", writedText)
}

func TestCopy_WithLimitAndOffset(t *testing.T) {
	text := "123456789"
	writeFile(FromFilePath, []byte(text))
	err := Copy(FromFilePath, ToFilePath, 5, 2)
	assert.NoError(t, err)

	writedText := readFile(ToFilePath)
	assert.Equal(t, "34567", writedText)
}
