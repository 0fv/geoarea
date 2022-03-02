package geoarea

import (
	"image/color"

	sm "github.com/flopp/go-staticmaps"
	"github.com/golang/geo/s2"
	"github.com/mmcloughlin/geohash"
)

type Box struct {
	geohash.Box
	WestLine, EastLine, SouthLine, NorthLine *Line
	*MapAbbr
}

func NewBox(geohashcode string) Box {
	box := geohash.BoundingBox(geohashcode)
	//左上角 西北角 大lng 小lat
	northwestPoint := &Point{
		x: box.MaxLng,
		y: box.MinLat,
	}
	//右上角 东北角 大lng 大lat
	northeastPoint := &Point{
		x: box.MaxLng,
		y: box.MaxLat,
	}
	//右下角 东南角 小lng 大lat
	southeastPoint := &Point{
		x: box.MinLng,
		y: box.MaxLat,
	}
	//左下角 西南角  小lng 小lat
	southwestPoint := &Point{
		x: box.MinLng,
		y: box.MinLat,
	}

	return Box{
		Box: box,
		//西部线 左上，左下
		WestLine: NewLine(northwestPoint, southwestPoint),
		//东部线 右上，右下
		EastLine: NewLine(northeastPoint, southeastPoint),
		//北部线 左上，右上
		NorthLine: NewLine(northwestPoint, northeastPoint),
		//南部线 左下，右下
		SouthLine: NewLine(southwestPoint, southeastPoint),
		MapAbbr: &MapAbbr{
			Color: color.RGBA{
				R: 121,
				G: 89,
				B: 121,
			},
			FillColor: color.RGBA{
				R: 121,
				G: 89,
				B: 121,
				A: 124,
			},
			Size: 1,
		},
	}
}

func (b Box) IsLinesIntersected(line *Line) bool {
	switch {
	case IsLinesIntersected(b.EastLine, line):
		return true
	case IsLinesIntersected(b.WestLine, line):
		return true
	case IsLinesIntersected(b.NorthLine, line):
		return true
	case IsLinesIntersected(b.SouthLine, line):
		return true
	}
	return false
}

func (b Box) Inbox(p *Point) bool {
	return b.Box.Contains(p.y, p.x)
}

//从左上顺时针旋转
func (b Box) Export() [][]float64 {
	return [][]float64{
		b.NorthLine.Export()[0],
		b.NorthLine.Export()[1],
		b.SouthLine.Export()[1],
		b.SouthLine.Export()[0],
	}
}
func (b Box) RenderToMap(ctx *sm.Context) {
	ctx.AddObject(
		sm.NewArea([]s2.LatLng{
			s2.LatLngFromDegrees(b.NorthLine.Export()[0][1], b.NorthLine.Export()[0][0]),
			s2.LatLngFromDegrees(b.NorthLine.Export()[1][1], b.NorthLine.Export()[1][0]),
			s2.LatLngFromDegrees(b.SouthLine.Export()[1][1], b.SouthLine.Export()[1][0]),
			s2.LatLngFromDegrees(b.SouthLine.Export()[0][1], b.SouthLine.Export()[0][0]),
		}, b.Color, b.FillColor, b.Size),
	)
}
