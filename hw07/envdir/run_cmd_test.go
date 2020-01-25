package envdir

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunCmd_EmptyCmd(t *testing.T) {
	assert.Panics(t, func() { RunCmd(nil, nil) })
	assert.Panics(t, func() { RunCmd([]string{}, nil) })
}

func TestRunCmd_RunIncorrectCmd(t *testing.T) {
	assert.Panics(t, func() { RunCmd([]string{"incorrect_cmd"}, nil) })
}

func TestRunCmd_LSWithoutArgs(t *testing.T) {
	res := RunCmd([]string{"ls"}, nil)
	assert.Zero(t, res)
}

func TestRunCmd_CheckCorrectOut(t *testing.T) {
	fileName := filepath.Join(getTestDir(), "test")
	file, _ := os.Create(fileName)
	defer os.Remove(fileName)

	old := os.Stdout
	os.Stdout = file
	res := RunCmd([]string{"ls", getTestDir()}, nil)
	os.Stdout = old

	assert.Zero(t, res)

	data, _ := ioutil.ReadFile(fileName)
	assert.Equal(t, "test\n", string(data))
}

func TestRunCmd_CheckWithFlags(t *testing.T) {
	fileName := filepath.Join(getTestDir(), "test")
	file, _ := os.Create(fileName)
	defer os.Remove(fileName)

	old := os.Stdout
	os.Stdout = file
	res := RunCmd([]string{"ls", "-a", getTestDir()}, nil)
	os.Stdout = old

	assert.Zero(t, res)

	data, _ := ioutil.ReadFile(fileName)
	assert.Equal(t, ".\n..\ntest\n", string(data))
}

func TestRunCmd_CheckEnv(t *testing.T) {
	fileName := filepath.Join(getTestDir(), "test")
	file, _ := os.Create(fileName)
	defer os.Remove(fileName)

	old := os.Stdout
	os.Stdout = file
	res := RunCmd([]string{"env"}, map[string]string{
		"ENV1": "one",
	})
	os.Stdout = old

	assert.Zero(t, res)

	data, _ := ioutil.ReadFile(fileName)
	assert.Equal(t, "ENV1=one\n", string(data))
}
