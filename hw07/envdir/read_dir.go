package envdir

import (
	"io/ioutil"
	"path/filepath"
)

// ReadDir read envdir
func ReadDir(dir string) (map[string]string, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var env map[string]string
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		data, err := ioutil.ReadFile(filepath.Join(dir, file.Name()))
		if err != nil {
			return nil, err
		}

		if len(data) == 0 {
			continue
		}

		if env == nil {
			env = make(map[string]string)
		}

		env[file.Name()] = string(data)
	}

	return env, nil
}
