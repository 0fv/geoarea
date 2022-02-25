package geoarea

import (
	"image"
	"image/color"

	sm "github.com/flopp/go-staticmaps"
	"github.com/fogleman/gg"
)

type MapAbbr struct {
	Color color.Color
	FillColor color.Color
	Size  float64
}

func (m *MapAbbr) GetMapAbbr() *MapAbbr {
	return m
}

func (m *MapAbbr) SetMapAbbr(abbr *MapAbbr) {
	m = abbr
}

type IRender interface {
	RenderToMap(ctx *sm.Context)
	GetMapAbbr() *MapAbbr
	SetMapAbbr(abbr *MapAbbr)
}

type MapPic struct {
	m *sm.Context
}

func NewMapPic(width, height int) *MapPic {
	m := sm.NewContext()
	m.SetSize(width, height)
	return &MapPic{
		m: m,
	}
}

func (m *MapPic) SetMapData(geodata ...IRender) *MapPic {
	for _, v := range geodata {
		v.RenderToMap(m.m)
	}
	return m
}

func (m *MapPic) ToImage() (image.Image, error) {
	return m.m.Render()
}

func (m *MapPic) ToFile(filename string) error {
	img, err := m.ToImage()
	if err != nil {
		return err
	}
	if err := gg.SavePNG(filename, img); err != nil {
		return nil
	}
	return nil
}
