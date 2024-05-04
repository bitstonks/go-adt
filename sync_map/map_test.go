package sync_map_test

import (
	"testing"

	"github.com/bitstonks/go-adt/sync_map"
	"github.com/stretchr/testify/assert"
)

func TestMap_Basic(t *testing.T) {
	t.Parallel()
	var m sync_map.Map[int, string]
	t.Run("load from empty", func(t *testing.T) {
		v, ok := m.Load(1)
		assert.Equal(t, "", v)
		assert.Equal(t, false, ok)
	})
	t.Run("store&load", func(t *testing.T) {
		m.Store(1, "Hello")
		v, ok := m.Load(1)
		assert.Equal(t, "Hello", v)
		assert.Equal(t, true, ok)
	})
	t.Run("delete", func(t *testing.T) {
		m.Delete(1)
		v, ok := m.Load(1)
		assert.Equal(t, "", v)
		assert.Equal(t, false, ok)
	})
}

func TestCompareAndDelete(t *testing.T) {
	t.Parallel()
	var m sync_map.Map[int, string]

	assert.False(t, m.CompareAndDelete(1, "Hello"))

	m.Store(1, "Hello")
	assert.True(t, m.CompareAndDelete(1, "Hello"))
	v, ok := m.Load(1)
	assert.False(t, ok)
	assert.Equal(t, "", v)

	m.Store(1, "World")
	assert.False(t, m.CompareAndDelete(1, "Hello"))
	v, ok = m.Load(1)
	assert.True(t, ok)
	assert.Equal(t, "World", v)
}

func TestCompareAndSwap(t *testing.T) {
	t.Parallel()
	var m sync_map.Map[int, string]

	assert.False(t, m.CompareAndSwap(1, "Hello", "World"))

	m.Store(1, "Hello")
	assert.True(t, m.CompareAndSwap(1, "Hello", "World"))
	v, ok := m.Load(1)
	assert.True(t, ok)
	assert.Equal(t, "World", v)

	assert.False(t, m.CompareAndSwap(1, "Hello", "Earth"))
	v, ok = m.Load(1)
	assert.True(t, ok)
	assert.Equal(t, "World", v)
}

func TestLoadAndDelete(t *testing.T) {
	t.Parallel()
	var m sync_map.Map[int, string]

	v, ok := m.LoadAndDelete(1)
	assert.False(t, ok)
	assert.Equal(t, "", v)

	m.Store(1, "Hello")
	v, ok = m.LoadAndDelete(1)
	assert.True(t, ok)
	assert.Equal(t, "Hello", v)

	v, ok = m.Load(1)
	assert.False(t, ok)
	assert.Equal(t, "", v)
}

func TestMap_LoadOrStore(t *testing.T) {
	t.Parallel()
	var m sync_map.Map[int, string]

	actual, loaded := m.LoadOrStore(1, "Hello")
	assert.Equal(t, "Hello", actual)
	assert.False(t, loaded)

	actual, loaded = m.LoadOrStore(1, "World")
	assert.Equal(t, "Hello", actual)
	assert.True(t, loaded)

	v, ok := m.Load(1)
	assert.True(t, ok)
	assert.Equal(t, "Hello", v)

	actual, loaded = m.LoadOrStore(2, "Goodbye")
	assert.Equal(t, "Goodbye", actual)
	assert.False(t, loaded)

	v, ok = m.Load(2)
	assert.True(t, ok)
	assert.Equal(t, "Goodbye", v)
}

func TestRange(t *testing.T) {
	t.Parallel()
	var m sync_map.Map[int, string]
	m.Store(1, "Hello")
	m.Store(2, "World")
	m.Store(30, "Goodbye")
	got := make(map[int]string)
	m.Range(func(k int, v string) bool {
		got[k] = v
		return true
	})
	assert.Equal(t, map[int]string{1: "Hello", 2: "World", 30: "Goodbye"}, got)
}

func TestSwap(t *testing.T) {
	t.Parallel()
	var m sync_map.Map[int, string]
	v, ok := m.Swap(1, "Hello")
	assert.Equal(t, "", v)
	assert.False(t, ok)

	v, ok = m.Load(1)
	assert.Equal(t, "Hello", v)
	assert.True(t, ok)

	v, ok = m.Swap(1, "World")
	assert.Equal(t, "Hello", v)
	assert.True(t, ok)
}
