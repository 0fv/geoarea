package geoarea

import "testing"

type PolyData struct {
	Id int
}

func (p PolyData) GetKey() interface{} {
	return p.Id
}

func TestPointIn(t *testing.T) {
	geoTire := NewGeoTire()
	var p []*Point
	for _, v := range test01 {
		p = append(p, NewPoint(v[0], v[1]))
	}
	polygon, _ := NewPolygon(p)
	polygonData := PolyData{Id: 1}
	geoTire.Set(polygon, &polygonData)
	point := NewPoint(103.85714275409597, 30.68129910640877)
	ret := geoTire.Get(point)
	if len(ret) == 1 {
		t.Error("expect 1, but get ", len(ret))
	}
	point = NewPoint(105.85714275409597, 30.68129910640877)
	ret = geoTire.Get(point)
	if len(ret) != 0 {
		t.Error("expect 0, but get ", len(ret))
	}
}
