package main

import (
	"fmt"
	"github.com/jphsd/graphics2d"
)

func main() {
	path := graphics2d.NewPath([]float64{0, 50})
	fmt.Printf("%s\n", path.String())
	path.AddStep([]float64{12, 100, 25, 100, 75, 0, 88, 0, 100, 50})
	fmt.Printf("%s\n", path.String())
	path.Close()
	fmt.Printf("%s\n", path.String())
	fmt.Printf("%s\n", path.Flatten(.3).String())
}
