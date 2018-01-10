package kdtree

import (
	"bytes"
	"fmt"
	"math"
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

func NewKdTree(points []Point) *KdTree {
	return newKdTree(points, 0)
}

// Returns a new k-d tree.
func newKdTree(points []Point, depth int) *KdTree {
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

	beforePoints := make([]Point, 0)
	afterPoints := make([]Point, 0)
	splittingValue := medianPoint.Val(axis)

	for i := 0; i < len(points); i++ {
		if points[i].Equal(medianPoint) {
			medianPoints = append(medianPoints, points[i])
			continue
		}

		if points[i].Val(axis) <= splittingValue {
			beforePoints = append(beforePoints, points[i])
			continue
		}

		afterPoints = append(afterPoints, points[i])
	}

	return &KdTree{
		newKdTree(beforePoints, depth+1),
		newKdTree(afterPoints, depth+1),
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
		return
	}
	if p.Val(k.axis) <= k.points[0].Val(k.axis) {
		targetNode = &k.leftChild
	}
	if *targetNode == nil {
		*targetNode = newKdTree([]Point{p}, k.depth)
		return
	}
	(*targetNode).insert(p)
}

// Find the nearest neighboring node
func (k *KdTree) NN(p Point) []Point {
	pts, _ := k.nn(p, math.MaxFloat64, k)
	return pts.points
}

func (k *KdTree) nn(p Point, smallestDistance float64, nNode *KdTree) (*KdTree, float64) {

	nn := nNode
	d := k.points[0].Distance(p)
	if d < smallestDistance {
		nn = k
		smallestDistance = d
	}

	// Find where to look first
	targetNode := k.leftChild
	otherSideOfPlandeNode := k.rightChild
	if p.Val(k.axis) > k.points[0].Val(k.axis) {
		targetNode = k.rightChild
		otherSideOfPlandeNode = k.leftChild
	}

	// Examine children
	if targetNode != nil {
		nnLeftLeaf, d := targetNode.nn(p, smallestDistance, nn)
		if d < smallestDistance {
			nn = nnLeftLeaf
			smallestDistance = d
		}
	}
	// Check if we should look in the other leaf by seeing if the hypersphere centered in p of radius smalledDistance
	// intersects with the hyperplane
	if otherSideOfPlandeNode != nil {
		if k.points[0].PlaneDistance(p.Val(k.axis), k.axis) <= smallestDistance {
			nnRightLeaf, d := otherSideOfPlandeNode.nn(p, smallestDistance, nn)
			if d < smallestDistance {
				nn = nnRightLeaf
				smallestDistance = d
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
