package hw05

import (
	"errors"
	"runtime"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRun_IncorrectInput(t *testing.T) {
	// empty func list
	err := Run([]func() error{}, 1, 1)
	require.Error(t, err)

	// N <= 0
	simpleFuncArray := []func() error{func() error { return nil }}
	err = Run(simpleFuncArray, -1, 10)
	require.Error(t, err)

	err = Run(simpleFuncArray, 0, 10)
	require.Error(t, err)

	// M <= 0
	err = Run(simpleFuncArray, 2, -1)
	require.Error(t, err)

	err = Run(simpleFuncArray, 2, 0)
	require.Error(t, err)

	// M > len(funcs)
	err = Run(simpleFuncArray, 2, 2)
	require.Error(t, err)
}

func TestRun_MErrors(t *testing.T) {
	M := 2
	errorFunc := func() error {
		return errors.New("Error")
	}

	correctFunc := func() error {
		return nil
	}

	funcArray := []func() error{
		errorFunc,
		errorFunc,
		correctFunc,
		correctFunc,
		correctFunc,
		correctFunc,
		correctFunc,
	}

	err := Run(funcArray, 1, M)
	require.Error(t, err)
	require.Equal(t, 2, runtime.NumGoroutine())

	err = Run(funcArray, 2, M)
	require.Error(t, err)
	require.Equal(t, 2, runtime.NumGoroutine())

	err = Run(funcArray, 100, M)
	require.Error(t, err)
	require.Equal(t, 2, runtime.NumGoroutine())
}

func TestRun_OK(t *testing.T) {
	counter := 0
	correctFunc := func() error {
		var m sync.Mutex
		m.Lock()
		defer m.Unlock()
		counter++
		return nil
	}

	funcArray := []func() error{
		correctFunc,
		correctFunc,
		correctFunc,
		correctFunc,
		correctFunc,
		correctFunc,
		correctFunc,
	}

	err := Run(funcArray, 1, 1)
	require.NoError(t, err)
	require.Equal(t, len(funcArray), counter)
	require.Equal(t, 2, runtime.NumGoroutine())

	counter = 0
	err = Run(funcArray, 2, 1)
	require.NoError(t, err)
	require.Equal(t, len(funcArray), counter)
	require.Equal(t, 2, runtime.NumGoroutine())

	counter = 0
	err = Run(funcArray, 100, 1)
	require.NoError(t, err)
	require.Equal(t, len(funcArray), counter)
	require.Equal(t, 2, runtime.NumGoroutine())
}
