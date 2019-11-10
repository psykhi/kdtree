package kdtree

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"math"
	"math/rand"
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
	//fmt.Printf("%v <==> %v dim %d = %f\n", p, val, i, ret)
	return ret
}

func (p *point) Distance(a Point) float64 {
	var ret float64
	for i := 0; i < len(p.vals); i++ {
		tmp := p.vals[i] - a.Val(i)
		ret += tmp * tmp
	}
	//fmt.Printf("%v <==> %v =%f\n", p, a, ret)
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
	})
	assert.Equal(t, "([[7 2] [7 2]], (([[5 4]], (([[2 3]], (none, none)), ([[4 7]], (none, none)))), ([[9 6]], (([[8 1]], (none, none)), none))))", tree.String())
}

func TestKdTree_Insert(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		tree := NewKdTree([]Point{&point{[]float64{7, 2}}})
		assert.Equal(t, "([[7 2]], (none, none))", tree.String())

		tree = tree.Insert(&point{[]float64{5, 4}})
		assert.Equal(t, "([[7 2]], (([[5 4]], (none, none)), none))", tree.String())
		tree = tree.Insert(
			&point{[]float64{2, 3}},
			&point{[]float64{9, 6}},
			&point{[]float64{8, 1}},
			&point{[]float64{8, 1}},
			&point{[]float64{4, 7}})
		assert.Equal(t, "([[7 2]], (([[5 4]], (([[2 3]], (none, ([[4 7]], (none, none)))), none)), ([[9 6]], (([[8 1] [8 1]], (none, none)), none))))", tree.String())
	})
	t.Run("random data", func(t *testing.T) {
		tree := NewKdTree([]Point{&point{[]float64{rand.Float64(), rand.Float64()}}})
		for i := 0; i < 10000; i++ {
			tree = tree.Insert(&point{[]float64{rand.Float64(), rand.Float64()}})
		}
		d := tree.NN(&point{[]float64{0.7, 0.28}})[0]
		assert.True(t, math.Abs(d.(*point).vals[0]-0.7) < 0.1)
		assert.True(t, math.Abs(d.(*point).vals[1]-0.28) < 0.1)
	})
	t.Run("insert line", func(t *testing.T) {
		tree := NewKdTree([]Point{&point{[]float64{7, 2}}})
		for i := 0; i < 100; i++ {
			tree = tree.Insert(&point{vals: []float64{float64(i), float64(i)}})
			fmt.Println("after", tree)
		}
		assert.EqualValues(t, &point{[]float64{7, 7}}, tree.NN(&point{[]float64{7, 8}})[0])
	})
	t.Run("insert same values", func(t *testing.T) {
		tree := NewKdTree([]Point{
			&point{[]float64{1, 0, 7}},
		})
		tree = tree.Insert(
			&point{[]float64{0, 2, 5}},
			&point{[]float64{0, 0, 2}},
			&point{[]float64{0, 0, 2}},
			&point{[]float64{0, 0, 2}},
			&point{[]float64{0, 0, 3}},
			&point{[]float64{0, 1, 3}}, // swap b
			&point{[]float64{0, 1, 4}}, // swap a
			&point{[]float64{2, 0, 0}},
		)
		assert.EqualValues(t, &point{[]float64{0, 1, 4}}, tree.NN(&point{[]float64{0, 1, 4}})[0])
	})

}

func TestKdTree_NN(t *testing.T) {
	t.Run("Simple case ", func(t *testing.T) {
		tree := NewKdTree([]Point{
			&point{[]float64{2, 3}},
			&point{[]float64{4, 7}},
			&point{[]float64{5, 4}},
			&point{[]float64{7, 2}},
			&point{[]float64{8, 1}},
			&point{[]float64{9, 6}},
		})
		assert.EqualValues(t, &point{[]float64{5, 4}}, tree.NN(&point{[]float64{5, 5}})[0])
		assert.EqualValues(t, &point{[]float64{7, 2}}, tree.NN(&point{[]float64{8, 4}})[0])

		// Insert another point already in there
		tree = tree.Insert(
			&point{[]float64{2, 3}})
		assert.EqualValues(t, &point{[]float64{2, 3}}, tree.NN(&point{[]float64{2, 4}})[0])
		assert.EqualValues(t, &point{[]float64{2, 3}}, tree.NN(&point{[]float64{2, 4}})[1])
	})
	t.Run("Multiple values on median", func(t *testing.T) {
		tree := NewKdTree([]Point{
			&point{[]float64{1, 0, 7, 22}},
			&point{[]float64{0, 2, 5, 27}},
			&point{[]float64{0, 0, 2, 9}},
			&point{[]float64{0, 1, 3, 11}}, // swap a
			&point{[]float64{0, 1, 3, 8}},  // swap b
			&point{[]float64{2, 0, 0, 12}},
		})
		assert.EqualValues(t, &point{[]float64{0, 0, 2, 9}}, tree.NN(&point{[]float64{0, 0, 2, 9}})[0])
	})
	t.Run("Multiple values on median 2", func(t *testing.T) {
		tree := NewKdTree([]Point{
			&point{[]float64{1, 0, 7}},
			&point{[]float64{0, 2, 5}},
			&point{[]float64{0, 0, 2}},
			&point{[]float64{0, 1, 3}}, // swap b
			&point{[]float64{0, 1, 4}}, // swap a
			&point{[]float64{2, 0, 0}},
		})
		assert.EqualValues(t, &point{[]float64{0, 0, 2}}, tree.NN(&point{[]float64{0, 0, 2}})[0])
	})
	t.Run("Multiple values on median case 3", func(t *testing.T) {
		tree := NewKdTree([]Point{
			&point{[]float64{1, 0, 7}},
			&point{[]float64{0, 2, 5}},
			&point{[]float64{0, 0, 2}},
			&point{[]float64{0, 0, 2}},
			&point{[]float64{0, 0, 2}},
			&point{[]float64{0, 0, 3}},
			&point{[]float64{0, 1, 3}}, // swap b
			&point{[]float64{0, 1, 4}}, // swap a
			&point{[]float64{2, 0, 0}},
		})
		assert.EqualValues(t, &point{[]float64{0, 1, 4}}, tree.NN(&point{[]float64{0, 1, 4}})[0])
	})
}

func BenchmarkKdTree_NN(b *testing.B) {
	b.Run("10 elements in tree", func(b *testing.B) {
		bench(b, 10)
	})
	b.Run("100 elements in tree", func(b *testing.B) {
		bench(b, 100)
	})
	b.Run("1000 elements in tree", func(b *testing.B) {
		bench(b, 1000)
	})
	b.Run("10000 elements in tree", func(b *testing.B) {
		bench(b, 10000)
	})
	b.Run("100000 elements in tree", func(b *testing.B) {
		bench(b, 100000)
	})
}
func bench(b *testing.B, count int) {

	pts := make([]Point, 0)
	for i := 0; i < count; i++ {
		randP := make([]float64, 10)
		for j := 0; j < 10; j++ {
			randP[j] = rand.Float64()
		}
		pts = append(pts, &point{randP})
	}
	tree := NewKdTree(pts)
	b.ReportAllocs()
	randP := make([]float64, 10)
	for j := 0; j < 10; j++ {
		randP[j] = rand.Float64()
	}
	p := &point{randP}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tree.NN(p)
	}
}
