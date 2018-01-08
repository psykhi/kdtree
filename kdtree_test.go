package kdtree

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

type point struct {
	vals []float64
}

func (p *point) Equal(point Point) bool {
	return reflect.DeepEqual(p, point)
}

func (p *point) PlaneDistance(val float64, i int) float64 {
	tmp := p.vals[i] - val
	ret := tmp * tmp
	fmt.Printf("%v <==> %v dim %d = %f\n", p, val, i, ret)
	return ret
}

func (p *point) Distance(a Point) float64 {
	var ret float64
	for i := 0; i < len(p.vals); i++ {
		tmp := p.vals[i] - a.Val(i)
		ret += tmp * tmp
	}
	fmt.Printf("%v <==> %v =%f\n", p, a, ret)
	return ret
}

func (p *point) String() string {
	return fmt.Sprint(p.vals)
}

func (p *point) DimCount() int {
	return len(p.vals)
}

func (p *point) Val(i int) float64 {
	return p.vals[i]
}

func TestKdTree_New(t *testing.T) {
	tree := NewKdTree([]Point{
		&point{[]float64{2, 3}},
		&point{[]float64{4, 7}},
		&point{[]float64{5, 4}},
		&point{[]float64{7, 2}},
		&point{[]float64{7, 2}},
		&point{[]float64{8, 1}},
		&point{[]float64{9, 6}},
	}, 0)
	assert.Equal(t, "([[7 2] [7 2]], (([[5 4]], (([[2 3]], (none, none)), ([[4 7]], (none, none)))), ([[9 6]], (([[8 1]], (none, none)), none))))", tree.String())
}

func TestKdTree_Insert(t *testing.T) {
	tree := NewKdTree([]Point{
		&point{[]float64{7, 2}},
	}, 0)
	assert.Equal(t, "([[7 2]], (none, none))", tree.String())

	tree.Insert(&point{[]float64{5, 4}})
	assert.Equal(t, "([[7 2]], (([[5 4]], (none, none)), none))", tree.String())
	tree.Insert(
		&point{[]float64{2, 3}},
		&point{[]float64{9, 6}},
		&point{[]float64{8, 1}},
		&point{[]float64{8, 1}},
		&point{[]float64{4, 7}})
	assert.Equal(t, "([[7 2]], (([[5 4]], (([[2 3]], (none, ([[4 7]], (none, none)))), none)), ([[9 6]], (([[8 1] [8 1]], (none, ([[8 1]], (none, none)))), none))))", tree.String())
}

func TestKdTree_NN(t *testing.T) {
	tree := NewKdTree([]Point{
		&point{[]float64{2, 3}},
		&point{[]float64{4, 7}},
		&point{[]float64{5, 4}},
		&point{[]float64{7, 2}},
		&point{[]float64{8, 1}},
		&point{[]float64{9, 6}},
	}, 0)
	assert.EqualValues(t, &point{[]float64{5, 4}}, tree.NN(&point{[]float64{5, 5}})[0])
	assert.EqualValues(t, &point{[]float64{7, 2}}, tree.NN(&point{[]float64{8, 4}})[0])

	// Insert another point already in there
	tree.Insert(
		&point{[]float64{2, 3}})
	assert.EqualValues(t, &point{[]float64{2, 3}}, tree.NN(&point{[]float64{2, 4}})[0])
	assert.EqualValues(t, &point{[]float64{2, 3}}, tree.NN(&point{[]float64{2, 4}})[1])
}
