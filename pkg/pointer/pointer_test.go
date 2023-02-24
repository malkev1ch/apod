package pointer

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_Int(t *testing.T) {
	var x = 10
	if *New(x) != x {
		t.Errorf("%v != %v", *New(x), x)
	}

	if From(&x) != x {
		t.Errorf("%v != %v", From(&x), x)
	}
}

func Test_Int32(t *testing.T) {
	var x int32 = 10
	if *New(x) != x {
		t.Errorf("%v != %v", *New(x), x)
	}

	if From(&x) != x {
		t.Errorf("%v != %v", From(&x), x)
	}
}

func Test_Uint(t *testing.T) {
	var x uint = 10
	if *New(x) != x {
		t.Errorf("%v != %v", *New(x), x)
	}

	if From(&x) != x {
		t.Errorf("%v != %v", From(&x), x)
	}
}

func Test_Uint32(t *testing.T) {
	var x uint32 = 10
	if *New(x) != x {
		t.Errorf("%v != %v", *New(x), x)
	}

	if From(&x) != x {
		t.Errorf("%v != %v", From(&x), x)
	}
}

func Test_String(t *testing.T) {
	var x = "string"
	if *New(x) != x {
		t.Errorf("%v != %v", *New(x), x)
	}

	if From(&x) != x {
		t.Errorf("%v != %v", From(&x), x)
	}
}

func Test_Bool(t *testing.T) {
	var x = true
	if *New(x) != x {
		t.Errorf("%v != %v", *New(x), x)
	}

	if From(&x) != x {
		t.Errorf("%v != %v", From(&x), x)
	}
}

func Test_Time(t *testing.T) {
	var x = time.Time{}
	if *New(x) != x {
		t.Errorf("%v != %v", *New(x), x)
	}

	if From(&x) != x {
		t.Errorf("%v != %v", From(&x), x)
	}
}

func Test_NewSlice(t *testing.T) {
	in := []int{1, 2, 3}
	out := NewSlice(in)
	for i := range in {
		if From(out[i]) != in[i] {
			t.Errorf("%v != %v", From(out[i]), in[i])
		}
	}
}

func Test_FromSlice(t *testing.T) {
	expected := []int{1, 2}
	in := []*int{New(1), nil, New(2)}
	out := FromSlice(in)
	assert.EqualValues(t, expected, out)
}
