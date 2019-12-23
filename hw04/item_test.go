package hw04

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestItem_Value(t *testing.T) {
	item := Item{nil, nil, nil}
	require.Nil(t, item.Value())

	item = Item{nil, nil, 100}
	require.IsType(t, int(0), item.Value())
	require.Equal(t, 100, item.Value())

	item = Item{nil, nil, "string"}
	require.IsType(t, string(""), item.Value())
	require.Equal(t, "string", item.Value())
}

func TestItem_Prev(t *testing.T) {
	item := Item{nil, nil, nil}
	require.Nil(t, item.Prev())

	prev := &item
	item = Item{prev, nil, 100}
	require.Equal(t, prev, item.Prev())
}

func TestItem_Next(t *testing.T) {
	item := Item{nil, nil, nil}
	require.Nil(t, item.Next())

	next := &item
	item = Item{nil, next, 100}
	require.Equal(t, next, item.Next())
}
