package kdtree

import (
	"bytes"
	"fmt"
	"reflect"
	"sort"
)

type Point interface {
	DimCount() int
	Val(int) float64
	String() string
	Distance(Point) float64
	PlaneDistance(float64, int) float64
}

type KdTree struct {
	leftChild  *KdTree
	rightChild *KdTree
	points     []Point
	axis       int
	depth      int
}

type ByIthDim struct {
	points []Point
	dim    int
}

func (a ByIthDim) Len() int           { return len(a.points) }
func (a ByIthDim) Swap(i, j int)      { a.points[i], a.points[j] = a.points[j], a.points[i] }
func (a ByIthDim) Less(i, j int) bool { return a.points[i].Val(a.dim) < a.points[j].Val(a.dim) }

func NewKdTree(points []Point, depth int) *KdTree {
	if len(points) == 0 {
		return nil
	}
	axis := depth % points[0].DimCount()
	if len(points) == 1 {
		return &KdTree{
			points: points,
			axis:   axis,
			depth:  depth,
		}
	}

	// Find the median point
	d := ByIthDim{points, axis}
	sort.Sort(d)

	medianPoint := points[len(points)/2]
	return &KdTree{
		NewKdTree(points[:len(points)/2], depth+1),
		NewKdTree(points[len(points)/2+1:], depth+1),
		[]Point{medianPoint},
		axis,
		depth,
	}
}

func (k *KdTree) Insert(pts ...Point) {
	for _, p := range pts {
		k.insert(p)
	}
}

func (k *KdTree) insert(p Point) {
	targetNode := &k.rightChild

	if reflect.DeepEqual(k.points[0], p) {
		k.points = append(k.points, p)
	}
	if p.Val(k.axis) < k.points[0].Val(k.axis) {
		targetNode = &k.leftChild
	}
	if *targetNode == nil {
		*targetNode = NewKdTree([]Point{p}, k.depth)
		return
	}
	(*targetNode).Insert(p)
}

func (k *KdTree) NN(p Point) []Point {
	smallestDistance := k.points[0].Distance(p)
	nn := k

	// Examine children
	if k.leftChild != nil {
		if k.leftChild.points[0].Distance(p) < smallestDistance {
			smallestDistance = k.leftChild.points[0].Distance(p)
			nn, smallestDistance = k.leftChild.nn(p, smallestDistance, k.leftChild)
		}

		// Check if we should look in the other leaf by seeing if the hypersphere centered in p of radius smalledDistance
		// intersects with the hyperplane. If so it means we must look in the leaf.
		if k.rightChild != nil {
			if k.rightChild.points[0].PlaneDistance(nn.points[0].Val(k.axis), k.axis) < smallestDistance {

				nnRightLeaf, d := k.rightChild.nn(p, smallestDistance, nn)
				if d < smallestDistance {
					nn = nnRightLeaf
				}
			}
		}
	}
	return nn.points
}

func (k *KdTree) nn(p Point, smallestDistance float64, nNode *KdTree) (*KdTree, float64) {

	nn := nNode

	// Examine children
	if k.leftChild != nil {
		if k.leftChild.points[0].Distance(p) < smallestDistance {
			smallestDistance = k.leftChild.points[0].Distance(p)
			nn, smallestDistance = k.leftChild.nn(p, smallestDistance, k.leftChild)
		}
		// Check if we should look in the other leaf by seeing if the hypersphere centered in p of radius smalledDistance
		// intersects with the hyperplane
		if k.rightChild != nil {
			if k.rightChild.points[0].PlaneDistance(nn.points[0].Val(k.axis), k.axis) < smallestDistance {
				nnRightLeaf, d := k.rightChild.nn(p, smallestDistance, nn)
				if d < smallestDistance {
					nn = nnRightLeaf
				}
			}
		}
	}

	return nn, smallestDistance
}

func (k *KdTree) String() string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("(%v", k.points))

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
	buf.WriteString("))")
	return buf.String()
}
