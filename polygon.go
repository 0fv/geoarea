package geoarea

import (
	"errors"
	"strings"
)

type Polygon struct {
	points []*Point
	lines  []*Line
}

func (p Polygon) Empty() bool {
	return len(p.points) == 0
}
func (p Polygon) String() string {
	str := make([]string, len(p.points))
	for i, v := range p.points {
		str[i] = v.String()
	}
	return "[" + strings.Join(str, ",") + "]"
}

func NewPolygon(p []*Point) (Polygon, error) {

	polygon := Polygon{
		points: p,
	}
	if len(p) < 3 {
		return polygon, errors.New("point < 3")
	}
	length := len(p)
	for i := 0; i < length; i++ {
		s := i
		e := i + 1
		if e == length {
			e = 0
		}
		startpoint := p[s]
		endpoint := p[e]
		polygon.lines = append(polygon.lines, NewLine(startpoint, endpoint))

	}
	return polygon, nil
}

func (p Polygon) Bound() (maxlat, minlat, maxlng, minlng float64) {
	if len(p.points) < 3 {
		return
	}
	minlat = p.points[0].y
	minlng = p.points[0].x
	maxlat = p.points[0].y
	maxlng = p.points[0].x
	for _, v := range p.points[1:] {

		if minlat > v.y {
			minlat = v.y
		}
		if minlng > v.x {
			minlng = v.x
		}
		if maxlat < v.y {
			maxlat = v.y
		}
		if maxlng < v.x {
			maxlng = v.x
		}
	}

	return
}

//查询线段交叉的，包含的 hashcode
func (p Polygon) Geohash() (cross []string, in []string) {
	//小框框获取
	squareData := SquareDivided(p.Bound())
	//遍历获取是否在框框中或者路过框
flag:
	for _, squarehashcode := range squareData {
		box := NewBox(squarehashcode)
		for _, polygonLine := range p.lines {
			if box.IsLinesIntersected(polygonLine) {
				cross = append(cross, squarehashcode)
				continue flag
			}
		}
		if p.SquareIn(box) {
			in = append(in, squarehashcode)
		}
	}
	return
}

type CrossStatus uint8

const (
	KeepInner CrossStatus = iota
	KeepOut
	InnnerToOut
	OutToInnner
)

func (p Polygon) CheckCross(startPoint, endPoint Point) CrossStatus {
	startStatus := p.PointIn(startPoint)
	endStatus := p.PointIn(endPoint)
	return CrossStatusCal(startStatus, endStatus)
}

func CrossStatusCal(startStatus, endStatus bool) CrossStatus {
	if (!startStatus) && endStatus {
		return OutToInnner
	}
	if startStatus && endStatus {
		return KeepInner
	}
	if startStatus && (!endStatus) {
		return InnnerToOut
	}
	return KeepOut
}

//点位是否在多边形内 多线程实现
func (p Polygon) PointIn(point Point) (intersected bool) {
	length := len(p.points)
	x := point.x
	y := point.y
	for i := 0; i < length; i++ {
		s := i
		e := i + 1
		if e == length {
			e = 0
		}
		start := p.points[s]
		end := p.points[e]
		xmin := start.x
		xmax := end.y
		if xmin > xmax {
			xmin, xmax = xmax, xmin
		}
		ymin := start.y
		ymax := end.y
		if ymin > ymax {
			ymax, ymin = ymin, ymax
		}
		if euqal(start.y, end.y) {
			if euqal(y, start.y) && little_equal(xmin, x) && little_equal(x, xmax) {
				return true
			}
			continue
		}
		xt := (end.x-start.x)*(y-start.y)/(end.y-start.y) + start.x
		if euqal(xt, x) && little_equal(ymin, y) && little_equal(y, ymax) {
			// on edge [vj,vi]
			return true
		}
		if little(x, xt) && little_equal(ymin, y) && little(y, ymax) {
			intersected = !intersected
		}
	}
	return
}

//矩形是否在多边形内
func (p Polygon) SquareIn(box Box) bool {
	return p.PointIn(Point{x: box.MaxLng, y: box.MaxLat}) &&
		p.PointIn(Point{x: box.MaxLng, y: box.MinLat}) &&
		p.PointIn(Point{x: box.MinLng, y: box.MaxLat}) &&
		p.PointIn(Point{x: box.MinLng, y: box.MinLat})
}
