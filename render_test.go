package geoarea

import (
	"fmt"
	"image/color"
	"testing"
)

var mapPic *MapPic

func TestMain(m *testing.M) {
	mapPic = NewMapPic(4000,
		4000)
	m.Run()
	mapPic.ToFile("point.png")
}

func TestRenderPoint(t *testing.T) {
	p := NewPoint(104.143552, 30.657429)
	mapPic.SetMapData(p)
}

func TestRenderLine(t *testing.T) {
	start := NewPoint(104.143552, 30.657429)
	end := NewPoint(104.76428, 31.091185)
	l := NewLine(start, end)
	mapPic.SetMapData(l)
}

func TestRenderBox(t *testing.T) {
	b := NewBox("wm6jc")
	mapPic.SetMapData(b)
}

func TestRenderPolygon(t *testing.T) {
	p2 := [][]float64{
		{103.85714275409597, 30.68129910640877},
		{103.85556814923899, 30.680076078427184},
		{103.85593660884881, 30.679673699758563},
		{103.8581310921787, 30.679142816468783},
		{103.8573851798921, 30.68053889681081},
		{103.85714275409597, 30.68129910640877},
	}
	var p []*Point
	for _, v := range p2 {
		p = append(p, NewPoint(v[0], v[1]))
	}
	ploy, _ := NewPolygon(p)
	mapPic.SetMapData(ploy)
}

var test01 = [][]float64{
	{103.93669026090713, 30.738044461790988},
	{104.0002291999821, 30.780356718629523},
	{104.12132545060979, 30.779675351496618},
	{104.19791284208355, 30.711624358611253},
	{104.2004728246817, 30.619255091476898},
	{104.14390629271686, 30.573242230616533},
	{104.02873294308358, 30.57797542280453},
	{103.9655281808921, 30.61621458064528},
	{103.93669026090713, 30.738044461790988},
}
var test02 = [][]float64{
	{103.85714275409597, 30.68129910640877},
	{103.85556814923899, 30.680076078427184},
	{103.85593660884881, 30.679673699758563},
	{103.8581310921787, 30.679142816468783},
	{103.8573851798921, 30.68053889681081},
	{103.85714275409597, 30.68129910640877},
}

func TestGeoIn(t *testing.T) {
	var p []*Point
	for _, v := range CDPoly {
		p = append(p, NewPoint(v[0], v[1]))
	}
	ploy, _ := NewPolygon(p)
	f := ploy.SquarePointIn(NewBox("wm6p"))
	fmt.Println(f)
}

func TestGeohash(t *testing.T) {
	var p []*Point
	for _, v := range CDPoly {
		p = append(p, NewPoint(v[0], v[1]))
	}
	ploy, _ := NewPolygon(p)
	mapPic.SetMapData(ploy)
	cross, in := ploy.Geohash()
	for _, v := range cross {
		b := NewBox(v)
		b.FillColor = color.RGBA{
			R: 255,
			G: 115,
			B: 0,
			A: 90,
		}
		mapPic.SetMapData(b)
	}
	for _, v := range in {
		b := NewBox(v)
		b.FillColor = color.RGBA{
			R: 111,
			G: 255,
			B: 0,
			A: 90,
		}
		mapPic.SetMapData(b)
	}
}
