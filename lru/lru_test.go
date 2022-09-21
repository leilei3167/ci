package lru

import (
	"fmt"
	"testing"
)

func TestNew(t *testing.T) {
	cache := New(10)
	l0 := cache.Len()
	if l0 != 0 {
		t.Fatalf("want Len =0,got %d", l0)
	}

	for i := 0; i < 1000; i++ {
		cache.Set(fmt.Sprint(i), fmt.Sprintf("yang:%d", i))
	}
	l1 := cache.Len()
	if l1 != 10 {
		t.Fatalf("want Len =10,got %d", l0)
	}

}
