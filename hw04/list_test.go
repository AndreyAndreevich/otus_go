package hw04

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestList_Empty(t *testing.T) {
	list := List{}
	assert.Nil(t, list.First())
	assert.Nil(t, list.Last())
}

func TestList_Len(t *testing.T) {
	list := List{}
	assert.Equal(t, 0, list.Len())

	list.PushBack(1)
	assert.Equal(t, 1, list.Len())

	list.PushFront(2)
	assert.Equal(t, 2, list.Len())

	item := list.Last()
	list.Remove(item)
	assert.Equal(t, 1, list.Len())
}

func TestList_First_And_Last(t *testing.T) {
	first := Item{nil, nil, 1}
	last := Item{nil, nil, 2}
	list := List{&first, &last, 2}

	assert.Equal(t, first, *list.First())
	assert.Equal(t, &first, list.First())

	assert.Equal(t, last, *list.Last())
	assert.Equal(t, &last, list.Last())
}

func TestList_PushFront(t *testing.T) {
	list := List{}
	list.PushFront(1)

	assert.NotNil(t, list.First())
	assert.Equal(t, 1, list.First().Value())
	assert.Equal(t, list.First(), list.Last())
	assert.Nil(t, list.First().Prev())
	assert.Nil(t, list.First().Next())

	list.PushFront(2)
	assert.Equal(t, 2, list.First().Value())
	assert.NotEqual(t, list.First(), list.Last())
	assert.NotNil(t, list.First().Next())
	assert.NotNil(t, list.Last().Prev())
	assert.Nil(t, list.First().Prev())
	assert.Equal(t, list.Last(), list.First().Next())
}

func TestList_PushBack(t *testing.T) {
	list := List{}
	list.PushBack(1)

	assert.NotNil(t, list.Last())
	assert.Equal(t, 1, list.Last().Value())
	assert.Equal(t, list.First(), list.Last())
	assert.Nil(t, list.Last().Prev())
	assert.Nil(t, list.Last().Next())

	list.PushBack(2)
	assert.Equal(t, 2, list.Last().Value())
	assert.NotEqual(t, list.First(), list.Last())
	assert.NotNil(t, list.First().Next())
	assert.NotNil(t, list.Last().Prev())
	assert.Nil(t, list.Last().Next())
	assert.Equal(t, list.First(), list.Last().Prev())
}

func TestList_Remove(t *testing.T) {
	list := List{}

	list.PushFront(1)
	list.PushBack(2)
	list.PushBack(3)
	item := list.Last()
	list.PushBack(4)
	list.Remove(item)

	assert.Equal(t, 1, list.First().Value())
	assert.Equal(t, 4, list.Last().Value())
	assert.Equal(t, 3, list.Len())
	assert.Equal(t, 2, list.Last().Prev().Value())

	item = list.Last()
	list.Remove(item)

	assert.Equal(t, 2, list.Last().Value())
	assert.Equal(t, list.First(), list.Last().Prev())
	assert.Equal(t, list.Last(), list.First().Next())

	item = list.First()
	list.Remove(item)

	assert.Equal(t, 2, list.First().Value())
	assert.Equal(t, list.First(), list.Last())
	assert.Nil(t, list.First().Prev())
	assert.Nil(t, list.First().Next())

	item = list.First()
	list.Remove(item)

	assert.Nil(t, list.First())
	assert.Nil(t, list.Last())
}
