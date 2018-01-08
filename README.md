# K-d trees in Golang

This package provides a simple implementation of k-d trees in Go.

Usage:
```go
// Create a new tree
tree := NewKdTree([]Point{
		&point{[]float64{2, 3}},
		&point{[]float64{4, 7}},
		&point{[]float64{5, 4}},
		&point{[]float64{7, 2}},
		&point{[]float64{8, 1}},
		&point{[]float64{9, 6}},
	}, 0)

// Insert a new value
tree.Insert(&point{[]float64{2, 3}})

// Find the nearest neighbor node
nn := tree.NN(&point{[]float64{2, 4}}
```

Elements in the tree must implement the `Point` interface


```go
type Point interface {
	DimCount() int                              // Dimension count in the vectors
	Val(i int) float64                          // Retrieve value in dimension i
	String() string                             // String representation of the point
	Distance(Point) float64                     // Distance to another point
	PlaneDistance(val float64, dim int) float64 // Distance to an hyperplane in dimension dim
}
```