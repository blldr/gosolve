package grapher

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/blldr/gosolve/gosolve/mathparse"
)

type (
	Graph struct {
		width      int
		points     [][][]float64
		resolution int
	}
)

func NewGraph(width int, resolution int) *Graph {
	return &Graph{
		width,
		make([][][]float64, 0, 10),
		resolution,
	}
}

func (g *Graph) AddEquation(eqString string) (int, error) {
	pointsSlice, err := g.createPointsSlice(eqString)
	if err != nil {
		return 0, err
	}
	g.points = append(g.points, pointsSlice)
	return len(g.points) - 1, nil
}

func (g *Graph) RemoveEquation(eqId int) {
	if eqId >= len(g.points) {
		return
	}
	g.points = append(g.points[:eqId], g.points[eqId+1:]...)
}

func (g *Graph) ChangeEquation(eqId int, eqString string) error {
	if eqId >= len(g.points) {
		return nil
	}
	pointsSlice, err := g.createPointsSlice(eqString)
	if err != nil {
		return err
	}
	g.points[eqId] = pointsSlice
	return nil
}

func (g *Graph) GetPointsSlice() [][][]float64 {
	return g.points
}

func (g *Graph) createPointsSlice(eqString string) ([][]float64, error) {
	xCount := g.width * g.resolution
	pointsSlice := make([][]float64, 0, xCount)
	for i := 0.0; i < float64(xCount); i++ {
		x := i/float64(g.resolution) - float64(g.width/2)
		x = math.Round(x*100) / 100
		var eq string
		xFloat := strconv.FormatFloat(x, 'f', 2, 64)
		if x < 0 {
			xFloat = fmt.Sprintf("%s%s%s", "(0-", xFloat[1:], ")")
		}
		eq = strings.ReplaceAll(eqString, "x", xFloat)
		yValues, err := mathparse.FindRoots(eq, -float64(g.width/2), float64(g.width/2), 'y')
		if err != nil {
			return nil, err
		}
		pointsSlice = append(pointsSlice, yValues)
	}
	for i := 0.0; i < float64(xCount); i++ {
		y := float64(g.width/2) - i/float64(g.resolution)
		y = math.Round(y*100) / 100
		var eq string
		yFloat := strconv.FormatFloat(y, 'f', 2, 64)
		if y < 0 {
			yFloat = fmt.Sprintf("%s%s%s", "(0-", yFloat[1:], ")")
		}
		eq = strings.ReplaceAll(eqString, "y", yFloat)
		fmt.Println(eq)
		xValues, err := mathparse.FindRoots(eq, -float64(g.width/2), float64(g.width/2), 'x')
		if err != nil {
			fmt.Println("fuck")
			return nil, err
		}
		fmt.Println("buckaroo")
		pointsSlice = append(pointsSlice, xValues)

	}
	return pointsSlice, nil
}
