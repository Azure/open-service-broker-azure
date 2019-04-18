package amqp

import (
	"math"
	"math/bits"
	"reflect"
	"strconv"
	"testing"
)

func TestBitmap(t *testing.T) {
	type (
		add  uint32 // call add with this value
		rem  uint32 // call remove with this value
		next int64  // call next this many time
	)

	tests := []struct {
		max uint32
		ops []interface{}

		nextFail bool
		next     uint32
		count    uint32
	}{
		{
			max: 9,
			ops: []interface{}{
				add(0), add(1), add(2), add(3), add(4), add(5), add(6), add(7), add(8), add(9),
				rem(3), rem(7),
			},

			next:  3,
			count: 9,
		},
		{
			max: 9,

			next:  0,
			count: 1,
		},
		{
			max: math.MaxUint32,
			ops: []interface{}{
				add(13000),
			},

			next:  0,
			count: 2,
		},
		{
			max: math.MaxUint32,
			ops: []interface{}{
				next(64),
			},

			next:  64,
			count: 65,
		},
		{
			max: math.MaxUint32,
			ops: []interface{}{
				next(65535),
			},

			next:  65535,
			count: 65536,
		},
		{
			max: math.MaxUint32,
			ops: []interface{}{
				next(300),
				rem(32), rem(78), rem(13),
				next(1),
			},

			next:  32,
			count: 299,
		},
		{
			max: 63,
			ops: []interface{}{
				next(64),
			},

			nextFail: true,
			count:    64,
		},
		{
			max: 31,
			ops: []interface{}{
				next(32),
			},

			nextFail: true,
			count:    32,
		},
		{
			max: 31,
			ops: []interface{}{
				add(32),
			},

			next:  0,
			count: 1,
		},
		{
			max: 63,
			ops: []interface{}{
				next(64),
				rem(64),
			},

			nextFail: true,
			count:    64,
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			bm := &bitmap{max: tt.max}

			for _, op := range tt.ops {
				switch op := op.(type) {
				case add:
					bm.add(uint32(op))
				case rem:
					bm.remove(uint32(op))
				case next:
					for i := int64(0); i < int64(op); i++ {
						bm.next()
					}
				default:
					panic("unhandled op " + reflect.TypeOf(op).String())
				}
			}

			next, ok := bm.next()
			if ok == tt.nextFail {
				t.Errorf("next() failed with %d", next)
			}

			if tt.next != next && !tt.nextFail {
				t.Errorf("expected next() to be %d, but it was %d", tt.next, next)
			}

			count := countBitmap(bm)
			if tt.count != count {
				t.Errorf("expected count() to be %d, but it was %d", tt.count, count)
			}
		})
	}
}

func TestBitmap_Sequence(t *testing.T) {
	const max = 1024
	bm := &bitmap{max: max}

	for i := uint32(0); i <= max; i++ {
		next, ok := bm.next()
		if !ok {
			t.Errorf("next() failed with %d", next)
		}

		if i != next {
			t.Errorf("expected next() to be %d, but it was %d", i, next)
		}
	}

	count := countBitmap(bm)
	if want := uint32(max + 1); want != count {
		t.Errorf("expected count() to be %d, but it was %d", want, count)
	}
}

func countBitmap(bm *bitmap) uint32 {
	var count uint32
	for _, v := range bm.bits {
		count += uint32(bits.OnesCount64(v))
	}
	return count
}
