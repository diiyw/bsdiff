package bsdiff

import (
	"math/rand"
	"reflect"
	"testing"
)

func TestPatch_Diff(t *testing.T) {
	tests := []struct {
		a, b, c []byte
	}{
		{[]byte{}, []byte{}, []byte{}},
		{[]byte{1, 2, 3}, []byte{1, 2, 3}, []byte{}},
		{[]byte{1, 2, 3}, []byte{1, 2, 4}, []byte{4}},
		{[]byte{1, 2, 3}, []byte{1, 2, 3, 4}, []byte{4}},
		{[]byte{1, 2, 3, 4}, []byte{1, 2, 3}, []byte{}},
		{[]byte{1, 5, 3, 4}, []byte{1, 2, 3}, []byte{2}},
	}
	for i, tt := range tests {
		got := Diff(tt.a, tt.b)
		if !reflect.DeepEqual(got.diff, tt.c) {
			t.Errorf("#%d: Diff() = %v, want %v", i, got.diff, tt.c)
		}
	}
}

func TestPatch_DiffSize(t *testing.T) {
	tests := []struct {
		a, b []byte
		c    int64
	}{
		{[]byte{}, []byte{}, 0},
		{[]byte{1, 2, 3}, []byte{1, 2, 3}, 3},
		{[]byte{1, 2, 3}, []byte{1, 2, 4}, 3},
		{[]byte{1, 2, 3}, []byte{1, 2, 3, 4}, 4},
		{[]byte{1, 2, 3, 4}, []byte{1, 2, 3}, 3},
		{[]byte{1, 5, 3, 4}, []byte{1, 2, 3}, 3},
	}
	for i, tt := range tests {
		got := Diff(tt.a, tt.b)
		if got.size != tt.c {
			t.Errorf("#%d: Diff() = %v, want %v", i, got.size, tt.c)
		}
	}
}

func TestPatch_Bytes(t *testing.T) {
	tests := []struct {
		a, b, d []byte
		c       int64
	}{
		{[]byte{}, []byte{}, []byte{}, 0},
		{[]byte{1, 2, 3}, []byte{1, 2, 3}, []byte{}, 3},
		{[]byte{1, 2, 3}, []byte{1, 2, 4}, []byte{4}, 3},
		{[]byte{1, 2, 3}, []byte{1, 2, 3, 4}, []byte{4}, 4},
		{[]byte{1, 2, 3, 4}, []byte{1, 2, 3}, []byte{}, 3},
		{[]byte{1, 5, 3, 4}, []byte{1, 2, 3}, []byte{2}, 3},
	}
	for i, tt := range tests {
		got := Diff(tt.a, tt.b)
		b := got.ToBytes()
		c := FromBytes(b)
		if c == nil {
			t.Errorf("#%d: FromBytes() = nil", i)
			continue
		}
		if c.size != got.size {
			t.Errorf("#%d: Bytes() = %v, want %v", i, c.size, got.size)
		}
		if !reflect.DeepEqual(c.diff, tt.d) {
			t.Errorf("#%d: Bytes() = %v, want %v", i, c.diff, tt.d)
		}
	}
}

func TestPatch_Patch(t *testing.T) {
	tests := []struct {
		a, b, c []byte
	}{
		{[]byte{}, []byte{}, []byte{}},
		{[]byte{1, 2, 3}, []byte{1, 2, 3}, []byte{1, 2, 3}},
		{[]byte{1, 2, 3}, []byte{1, 2, 4}, []byte{1, 2, 4}},
		{[]byte{1, 2, 3}, []byte{1, 2, 3, 4}, []byte{1, 2, 3, 4}},
		{[]byte{1, 2, 3, 4}, []byte{1, 2, 3}, []byte{1, 2, 3}},
		{[]byte{1, 5, 3, 4}, []byte{1, 2, 3}, []byte{1, 2, 3}},
	}
	for i, tt := range tests {
		got := Diff(tt.a, tt.b)
		c := got.Apply(tt.a)
		if !reflect.DeepEqual(c, tt.c) {
			t.Errorf("#%d: Apply() = %v, want %v", i, c, tt.c)
		}
	}
}

func TestPatch_Patch2(t *testing.T) {
	tests := []struct {
		a, b []byte
	}{
		{[]byte("hello"), []byte{}},
		{[]byte("hello"), []byte("hella")},
		{[]byte("hello world"), []byte("hella")},
		{[]byte("hello world"), []byte("hello world")},
		{[]byte("hello world"), []byte("你好 hello world")},
		{[]byte("你好 hello world"), []byte("hello warld")},
	}
	for i, tt := range tests {
		got := Diff(tt.a, tt.b)
		c := got.Apply(tt.a)
		if !reflect.DeepEqual(c, tt.b) {
			t.Errorf("#%d: Apply() = %s, want %s", i, c, tt.b)
		}
	}
}

func BenchmarkPatch_Diff16B(b *testing.B) {
	a := []byte("hello world! 123")
	c := []byte("hello warld? 123")
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Diff(a, c)
	}
}

func BenchmarkPatch_Diff1KB(b *testing.B) {
	a := make([]byte, 1024)
	c := make([]byte, 1024)
	for i := 0; i < 200; i++ {
		a[rand.Intn(1000)] = 1
	}
	for i := 0; i < 200; i++ {
		c[rand.Intn(1000)] = 1
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Diff(a, c)
	}
}

func BenchmarkPatch_Diff64KB(b *testing.B) {
	a := make([]byte, 1024*64)
	c := make([]byte, 1024*64)
	for i := 0; i < 2000; i++ {
		a[rand.Intn(1000)] = 1
	}
	for i := 0; i < 2000; i++ {
		c[rand.Intn(1000)] = 1
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Diff(a, c)
	}
}

func BenchmarkPatch_Diff1MB(b *testing.B) {
	a := make([]byte, 1024*1024)
	c := make([]byte, 1024*1024)
	for i := 0; i < 2000; i++ {
		a[rand.Intn(2000)] = 1
	}
	for i := 0; i < 2000; i++ {
		c[rand.Intn(2000)] = 1
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Diff(a, c)
	}
}

func BenchmarkPatch_Diff6MB(b *testing.B) {
	a := make([]byte, 1024*1024*6)
	c := make([]byte, 1024*1024*6)
	for i := 0; i < 20000; i++ {
		a[rand.Intn(20000)] = 1
	}
	for i := 0; i < 20000; i++ {
		c[rand.Intn(20000)] = 1
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Diff(a, c)
	}
}

func BenchmarkPatch_Diff8MB(b *testing.B) {
	a := make([]byte, 1024*1024*8)
	c := make([]byte, 1024*1024*8)
	for i := 0; i < 20000; i++ {
		a[rand.Intn(20000)] = 1
	}
	for i := 0; i < 20000; i++ {
		c[rand.Intn(20000)] = 1
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Diff(a, c)
	}
}
