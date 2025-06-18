package main

import "testing"

func TestCosts(t *testing.T) {
	tests := []struct {
		name         string
		shape        Shape
		expectedCost float64
	}{
		{
			name:         "Square_Cost",
			shape:        square{l: 5},
			expectedCost: 25 * 100,
		},
		{
			name:         "Rectangle_Cost",
			shape:        rectangle{l: 4.222222, b: 1},
			expectedCost: 4.222222 * 20,
		},
		{
			name:         "Square_ZeroArea",
			shape:        square{l: 0},
			expectedCost: 0,
		}, {
			name:         "react_ZeroArea",
			shape:        rectangle{l: 0, b: 0},
			expectedCost: 0,
		}}

	for _, tt := range tests {
		t.Logf("Running test case: %s \n", tt.name)

		gotCost := costs(tt.shape)
		tol := 1e-4
		currentToll := gotCost - tt.expectedCost

		if currentToll > tol {
			t.Errorf("FAIL: Test %s  For shape %T (Area: %f), expected cost %f, but got %f",
				tt.name, tt.shape, tt.shape.Area(), tt.expectedCost, gotCost)
		} else {
			t.Logf("PASS: Test %s", tt.name)
		}
	}
}
