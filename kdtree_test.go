package kdtree

import (
	"testing"
	"fmt"
	"github.com/stretchr/testify/assert"
)

type point struct {
	vals []float64
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

func TestKdTree_Insert(t *testing.T) {
	tree := NewKdTree([]Point{
		&point{[]float64{2,3}},
		&point{[]float64{5,4}},
		&point{[]float64{9,6}},
		&point{[]float64{4,7}},
		&point{[]float64{8,1}},
		&point{[]float64{7,2}},
		},0)
	assert.Equal(t,"[[7 2]], ([[5 4]], ([[2 3]], (none, none), [[4 7]], (none, none)), [[9 6]], ([[8 1]], (none, none), none))",tree.String())
}
