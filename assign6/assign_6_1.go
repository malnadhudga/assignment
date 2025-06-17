package main

import "fmt"

// Shape is an interface that defines the Area method.
type Shape interface {
	Area() float64
}

// rectangle represents a rectangle with breadth (b) and length (l).
type rectangle struct {
	b, l float64
}

// square represents a square with a side length (l).
type square struct {
	l float64
}

// Area calculates the area of a square.
func (s square) Area() float64 {
	return s.l * s.l
}

// Area calculates the area of a rectangle.
func (r rectangle) Area() float64 {
	return r.l * r.b
}

// calArea takes any Shape and returns its area.
func calArea(s Shape) float64 {
	return s.Area()
}

func main() {

	rect := rectangle{b: 10, l: 20}
	// Calculate the area of the rectangle.
	rectArea := calArea(rect)
	fmt.Printf("Area of rectangle: %.2f\n", rectArea)

	// Create a new square with side length 4.
	sq := square{l: 10}
	squareArea := calArea(sq)
	fmt.Printf("Area of square: %.2f\n", squareArea)
}
