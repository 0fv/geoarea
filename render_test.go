package geoarea

import (
	"fmt"
	"testing"
)

var mapPic *MapPic

func TestMain(m *testing.M) {
	mapPic = NewMapPic(1000, 1000)
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
	b := NewBox("wx4ej")
	fmt.Println(b.Export())
	mapPic.SetMapData(b)
}
