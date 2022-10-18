package lru

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	cache := New(10)
	l0 := cache.Len()
	require.Equal(t, l0, 0)

	for i := 0; i < 1000; i++ {
		cache.Set(fmt.Sprint(i), fmt.Sprintf("yang:%d", i))
	}
	l1 := cache.Len()
	require.Equal(t, l1, 10)

	t.Run("get1", func(t *testing.T) {
		res, ok := cache.Get("noexist")
		require.False(t, ok)
		require.Nil(t, res)
	})

	t.Run("get2", func(t *testing.T) {
		res, ok := cache.Get("999")
		require.True(t, ok)
		require.Equal(t, res, "yang:999")
	})

	t.Run("remove", func(t *testing.T) {
		for i := 0; i < 9; i++ {
			cache.RemoveOldest()
		}
		l2 := cache.Len()
		require.Equal(t, l2, 1)
	})
	t.Run("update", func(t *testing.T) {
		for i := 0; i < 9; i++ {
			cache.Set(fmt.Sprint(i), fmt.Sprintf("yang:%d", i))
		}
		cache.Set(fmt.Sprint("999"), "updated")
	})
	t.Run("remove", func(t *testing.T) {
		cache.Remove("999")
		v, ok := cache.Get("999")
		require.False(t, ok)
		require.Nil(t, v)
	})
	t.Run("clear", func(t *testing.T) {
		cache.Clear()
		l3 := cache.Len()
		require.Equal(t, l3, 0)
	})
}
