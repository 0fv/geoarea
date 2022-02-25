package geoarea

import (
	"fmt"
	"image/color"
	"math"

	sm "github.com/flopp/go-staticmaps"
	"github.com/golang/geo/s2"
)

//直线
type Line struct {
	start, end *Point
	x, y       float64
	*MapAbbr
}

func NewLine(start, end *Point) *Line {
	return &Line{
		start: start,
		end:   end,
		x:     end.x - start.x,
		y:     end.y - start.y,
		MapAbbr: &MapAbbr{
			Color: color.RGBA{
				R: 0,
				G: 89,
				B: 121,
			},
			Size: 2,
		},
	}
}

func (l Line) String() string {
	return fmt.Sprintf("Data====\nstart:%v end:%v x:%v,y:%v\n====\n", l.start, l.end, l.x, l.y)
}

//向量叉乘
func CrossProduct(line1, line2 *Line) float64 {
	return line1.x*line2.y - line2.x*line1.y
}

//向量点乘
func DotProduct(line1, line2 *Line) float64 {
	return line1.x*line2.x + line1.y*line2.y
}

//是否排斥
func IsIncluded(line1, line2 *Line) bool {
	return math.Min(line1.start.x, line1.end.x) <= math.Max(line2.start.x, line2.end.x) &&
		math.Min(line2.start.x, line2.end.x) <= math.Max(line1.start.x, line1.end.x) &&
		math.Min(line1.start.y, line1.end.y) <= math.Max(line2.start.y, line2.end.y) &&
		math.Min(line2.start.y, line2.end.y) <= math.Max(line1.start.y, line1.end.y)
}

//是否互相跨立
func IsCrossed(line1, line2 *Line) bool {
	line_ac := NewLine(line1.start, line2.start)
	line_ad := NewLine(line1.start, line2.end)
	line_bc := NewLine(line1.end, line2.start)
	line_bd := NewLine(line1.end, line2.end)
	return CrossProduct(line_ac, line_ad)*CrossProduct(line_bc, line_bd) <= 0 &&
		CrossProduct(line_ac, line_bc)*CrossProduct(line_ad, line_bd) <= 0
}

// 线段是否与相交
func IsLinesIntersected(line1, line2 *Line) bool {
	return IsIncluded(line1, line2) && IsCrossed(line1, line2)
}

func (l Line) Export() [][]float64 {
	return [][]float64{
		{l.start.x, l.start.y},
		{l.end.x, l.end.y},
	}
}

func (l Line) RenderToMap(ctx *sm.Context) {
	ctx.AddObject(
		sm.NewPath([]s2.LatLng{
			s2.LatLngFromDegrees(l.start.y, l.start.x),
			s2.LatLngFromDegrees(l.end.y, l.end.x),
		}, l.Color, l.Size),
	)
}
