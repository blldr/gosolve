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
		points [][][]float64
		resolution int
	}

)

func NewGraph(resolution int) *Graph{
	return &Graph{
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

func (g *Graph) RemoveEquation(eqId int)  {
	if eqId >= len(g.points) {
		return
	}
	g.points = append(g.points[:eqId], g.points[eqId + 1:]...)
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
	xCount := 20*g.resolution
	pointsSlice := make([][]float64, 0, xCount)
	for i := 0.0; i < float64(xCount); i++ {
		x := i/float64(g.resolution) - 10.0
		x = math.Round(x * 100) / 100
		var eq string
		xFloat := strconv.FormatFloat(x, 'f', 2, 64)
		if x < 0 {
			xFloat = fmt.Sprintf("%s%s%s","(0-",xFloat[1:], ")")
		}
		eq = strings.ReplaceAll(eqString, "x", xFloat)
		yValues, err := mathparse.FindRoots(eq, -10, 10, 'y')
		if err != nil {
			return nil, err
		}
		pointsSlice = append(pointsSlice, yValues)

	}
	return pointsSlice, nil
}


