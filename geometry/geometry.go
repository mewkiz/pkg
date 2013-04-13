// Package geometry implements basic geometric types and operations.
package geometry

// A Rectangle contains the points with Min.X <= X < Max.X, Min.Y <= Y < Max.Y.
type Rectangle struct {
	Min, Max Point
}

// Rect is shorthand for Rectangle{Pt(x0, y0), Pt(x1, y1)}.
func Rect(x0, y0, x1, y1 float64) Rectangle {
	return Rectangle{Pt(x0, y0), Pt(x1, y1)}
}

// Dx returns r's width.
func (r Rectangle) Dx() float64 {
	return r.Max.X - r.Min.X
}

// Dy returns r's height.
func (r Rectangle) Dy() float64 {
	return r.Max.Y - r.Min.Y
}

// A Point is an X, Y coordinate pair. The axes increase right and down.
type Point struct {
	X, Y float64
}

// Pt is shorthand for Point{X, Y}.
func Pt(x, y float64) Point {
	return Point{x, y}
}
