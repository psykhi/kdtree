package kdtree

import (
	"sort"
	"fmt"
	"bytes"
)

type Point interface{
	DimCount() int
	Val(int)float64
	String()string
}

type KdTree struct {
	leftChild *KdTree
	rightChild *KdTree
	points []Point
}

type ByIthDim struct {
	points []Point
	dim int
}

func (a ByIthDim) Len() int           { return len(a.points) }
func (a ByIthDim) Swap(i, j int)      { a.points[i], a.points[j] = a.points[j], a.points[i]}
func (a ByIthDim) Less(i, j int) bool { return a.points[i].Val(a.dim) < a.points[j].Val(a.dim) }

func NewKdTree(points []Point, depth int) *KdTree{
	if len(points) == 0{
		return nil
	}
	if len(points)==1{
		return &KdTree{
			points:points,
		}
	}

	axis := depth%points[0].DimCount()
	// Find the median point
	d :=ByIthDim{points,axis}
	sort.Sort(d)

medianPoint :=points[len(points)/2]
	return &KdTree{
		NewKdTree(points[:len(points)/2],depth+1),
		NewKdTree(points[len(points)/2+1:],depth+1),
		[]Point{medianPoint},
	}
}

func (k*KdTree) Insert(p Point){

}

func (k*KdTree) String() string{
	var buf bytes.Buffer
	if len(k.points)==0{
		buf.WriteString("none")
		return buf.String()
	} else {
		buf.WriteString(fmt.Sprintf("%v",k.points))
	}

	buf.WriteString(", (")
	if k.leftChild != nil {
		buf.WriteString(k.leftChild.String())
	} else {
		buf.WriteString("none")
	}
	buf.WriteString(", ")
	if k.rightChild != nil {
		buf.WriteString(k.rightChild.String())
	} else {
		buf.WriteString("none")
	}
	buf.WriteString(")")
	return buf.String()
}