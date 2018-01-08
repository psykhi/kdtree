package kdtree

import (
	"bytes"
	"fmt"
	"sort"
)

type Point interface {
	DimCount() int                              // Dimension count in the vectors
	Val(i int) float64                          // Retrieve value in dimension i
	String() string                             // String representation of the point
	Distance(Point) float64                     // Distance to another point
	PlaneDistance(val float64, dim int) float64 // Distance to an hyperplane in dimension dim
	Equal(point Point) bool                     // Point comparison
}

// KdTree implements the k-d tree structure
type KdTree struct {
	leftChild  *KdTree
	rightChild *KdTree
	points     []Point
	axis       int
	depth      int
}

type byIthDim struct {
	points []Point
	dim    int
}

func (a byIthDim) Len() int           { return len(a.points) }
func (a byIthDim) Swap(i, j int)      { a.points[i], a.points[j] = a.points[j], a.points[i] }
func (a byIthDim) Less(i, j int) bool { return a.points[i].Val(a.dim) < a.points[j].Val(a.dim) }

// Returns a new k-d tree.
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

	// Find the median points
	d := byIthDim{points, axis}
	sort.Sort(d)
	medianPoint := points[len(points)/2]
	medianPoints := make([]Point, 0)
	medianPoints = append(medianPoints, medianPoint)
	i := 1
	j := 1
	beforeMedian := len(points) / 2
	afterMedian := len(points) / 2
	for i > 0 && j > 0 && len(points)/2+i < len(points) && len(points)/2-j >= 0 {
		if points[len(points)/2+i].Equal(medianPoint) {
			medianPoints = append(medianPoints, points[len(points)/2+i])
			afterMedian = len(points)/2 + i
			i++
		} else {
			i = -1
		}
		if points[len(points)/2-j].Equal(medianPoint) {
			medianPoints = append(medianPoints, points[len(points)/2-j])
			beforeMedian = len(points)/2 - j
			j++
		} else {
			j = -1
		}
	}

	return &KdTree{
		NewKdTree(points[:beforeMedian], depth+1),
		NewKdTree(points[afterMedian+1:], depth+1),
		medianPoints,
		axis,
		depth,
	}
}

// Insert points in the k-d tree. The tree might become unbalanced
func (k *KdTree) Insert(pts ...Point) {
	for _, p := range pts {
		k.insert(p)
	}
}

func (k *KdTree) insert(p Point) {
	targetNode := &k.rightChild

	if k.points[0].Equal(p) {
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

// Find the nearest neighboring node
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

// Pretty print the tree
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
