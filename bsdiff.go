package bsdiff

import (
	"encoding/binary"

	"github.com/RoaringBitmap/roaring/v2"
)

func Diff(a, b []byte) *Patch {
	var patch = &Patch{
		rb:   roaring.New(),
		diff: make([]byte, 0),
	}
	al := len(a)
	bl := len(b)
	if al == 0 {
		// `a` is empty, `b` is the new bytes
		patch.diff = b
		return patch
	}
	if bl == 0 {
		// `b` is empty, `a` is the old bytes
		return patch
	}
	// `a` and `b` are not empty len(a) <= len(b)
	if al <= bl {
		for i := 0; i < al; i++ {
			if a[i] != b[i] {
				patch.rb.Add(uint32(i))
				patch.diff = append(patch.diff, b[i])
			}
		}
		if al < bl {
			for i, bs := range b[al:] {
				patch.rb.Add(uint32(al + i))
				patch.diff = append(patch.diff, bs)
			}
		}
		patch.size = int64(bl)
		return patch
	}
	// `a` and `b` are not empty and have different length, len(a) > len(b)
	for i, bs := range b {
		if bs != a[i] {
			patch.rb.Add(uint32(i))
			patch.diff = append(patch.diff, bs)
		}
	}
	patch.size = int64(bl)
	return patch
}

type Patch struct {
	rb   *roaring.Bitmap
	size int64
	diff []byte
}

func (p *Patch) ToBytes() []byte {
	var buf = make([]byte, 8)
	n := binary.PutVarint(buf, int64(p.size))
	rbs, _ := p.rb.ToBytes()
	var rbuf = make([]byte, 8)
	rbn := binary.PutVarint(rbuf, int64(len(rbs)))
	return append(buf[:n], append(rbuf[:rbn], append(rbs, p.diff...)...)...)
}

func FromBytes(b []byte) *Patch {
	var size int64
	size, n1 := binary.Varint(b)
	rn, n2 := binary.Varint(b[n1:])
	rb := roaring.New()
	rn, err := rb.FromBuffer(b[n1+n2 : n1+n2+int(rn)])
	if err != nil {
		return nil
	}
	return &Patch{
		rb:   rb,
		size: size,
		diff: b[n1+n2+int(rn):],
	}
}

func (p *Patch) Apply(a []byte) []byte {
	if p.size == 0 {
		return []byte{}
	}
	var b = make([]byte, p.size)
	var j int
	al := len(a)
	if int64(al) >= p.size {
		a = a[:p.size]
	}
	for i, bs := range a {
		if p.rb.Contains(uint32(i)) {
			b[i] = p.diff[j]
			j++
			continue
		}
		b[i] = bs
	}
	if j < len(p.diff) {
		for i, ds := range p.diff[j:] {
			b[al+i] = ds
		}
	}
	return b
}
