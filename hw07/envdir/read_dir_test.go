package envdir

import (
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

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

func createFile(fileName string, data []byte) string {
	testDir := getTestDir()
	filePath := filepath.Join(testDir, fileName)
	f, err := os.Create(filePath)
	if err != nil {
		log.Panicf("failed to create: %v", err)
	}

	defer f.Close()

	_, err = f.Write(data)
	if err != nil {
		log.Panicf("failed to write: %v", err)
	}

	return filePath
}

func TestReadDir_IncorrectDir(t *testing.T) {
	testDir := getTestDir()
	incorrectDir := filepath.Join(testDir, "incorrect_dir")

	env, err := ReadDir(incorrectDir)
	assert.Error(t, err)
	assert.True(t, os.IsNotExist(err))
	assert.Nil(t, env)
}

func TestReadDir_PathIsNotDir(t *testing.T) {
	file := createFile("FILE", []byte(""))
	defer os.Remove(file)

	env, err := ReadDir(file)
	assert.Error(t, err)
	assert.Nil(t, env)
}

func TestReadDir_EmptyDir(t *testing.T) {
	testDir := getTestDir()

	env, err := ReadDir(testDir)
	assert.NoError(t, err)
	assert.Nil(t, env)
}

func TestReadDir_SomeEnv(t *testing.T) {
	testDir := getTestDir()
	file1 := createFile("ENV1", []byte("one"))
	defer os.Remove(file1)

	file2 := createFile("ENV2", []byte("two"))
	defer os.Remove(file2)

	env, err := ReadDir(testDir)
	assert.NoError(t, err)
	assert.NotNil(t, env)
	assert.Len(t, env, 2)
	assert.Equal(t, map[string]string{
		"ENV1": "one",
		"ENV2": "two",
	}, env)
}

func TestReadDir_EmptyFile(t *testing.T) {
	testDir := getTestDir()
	file1 := createFile("ENV1", []byte("one"))
	defer os.Remove(file1)

	file2 := createFile("ENV2", []byte(""))
	defer os.Remove(file2)

	env, err := ReadDir(testDir)
	assert.NoError(t, err)
	assert.NotNil(t, env)
	assert.Len(t, env, 1)
	assert.Equal(t, map[string]string{
		"ENV1": "one",
	}, env)
}
