package rebound

import "testing"

func TestNewTransformationMatrix(t *testing.T) {
	trans := [3]float64{
		-17.7082,
		-11.4156,
		2.0922,
	}

	scale := [3]float64{
		1,
		1,
		1,
	}

	rot := [4]float64{
		0,
		0,
		0,
		1,
	}

	expected := [16]float64{
		-0.99975,
		-0.00679829,
		0.0213218,
		0,
		0.00167596,
		0.927325,
		0.374254,
		0,
		-0.0223165,
		0.374196,
		-0.927081,
		0,
		-0.0115543,
		0.194711,
		-0.478297,
		1,
	}

	result := NewTransformationMatrix(trans, rot, scale)

	var check bool = false
	for index, f := range expected {
		if float32(f) != result[index] {
			check = true
		}
	}

	if check {
		t.Errorf("NewTransformationmatrix returns wrong result. Expected :\n%v\n Got:\n%v", expected, result)
	}
}
