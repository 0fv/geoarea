package geoarea

import (
	"encoding/json"
	"fmt"
	"image/color"
	"math"

	sm "github.com/flopp/go-staticmaps"
	"github.com/golang/geo/s2"
)

type Point struct {
	x float64
	y float64
	*MapAbbr
}

func NewPoint(x, y float64) *Point {
	return &Point{
		x: x,
		y: y,
		MapAbbr: &MapAbbr{
			Color: color.Black,
			Size:  12,
		},
	}
}

func (p Point) Export() (x, y float64) {
	return p.x, p.y
}

func (p Point) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("[%v,%v]", p.x, p.y)), nil
}

func (p *Point) UnmarshalJSON(data []byte) error {
	var arr []float64
	if err := json.Unmarshal(data, &arr); err != nil {
		return err
	}
	if len(arr) >= 2 {
		p.x = arr[0]
		p.y = arr[1]
	}
	return nil
}

func (p Point) Less(other Point) bool {
	if p.x != other.y {
		return p.x < other.x
	}
	return p.y < other.y
}

func (p Point) Equal(other Point) bool {
	return other.x == p.x && other.y == p.y
}

func (p Point) String() string {
	return fmt.Sprintf("x:%v,y:%v\n", p.x, p.y)
}

const PI float64 = 3.141592653589793

func (p Point) Distance(b Point) float64 {

	radlat1 := float64(PI * p.y / 180)
	radlat2 := float64(PI * b.y / 180)

	theta := float64(p.x - b.x)
	radtheta := float64(PI * theta / 180)

	dist := math.Sin(radlat1)*math.Sin(radlat2) + math.Cos(radlat1)*math.Cos(radlat2)*math.Cos(radtheta)

	if dist > 1 {
		dist = 1
	}
	dist = math.Acos(dist)
	dist = dist * 180 / PI
	dist = dist * 60 * 1.1515
	dist = dist * 1.609344 * 1000

	return dist
}

func (p Point) RenderToMap(ctx *sm.Context) {
	ctx.AddObject(
		sm.NewMarker(
			s2.LatLngFromDegrees(p.y, p.x),
			p.Color, p.Size,
		))
}
