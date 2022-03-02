package geoarea

import (
	"errors"
	"fmt"
	"image/color"
	"strings"
	"sync"

	sm "github.com/flopp/go-staticmaps"
	"github.com/golang/geo/s2"
)

type Polygon struct {
	points []*Point
	lines  []*Line
	*MapAbbr
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

func NewPolygon(p []*Point) (*Polygon, error) {

	polygon := Polygon{
		points: p,
		MapAbbr: &MapAbbr{
			Color: color.RGBA{
				R: 121,
				G: 0,
				B: 121,
			},
			FillColor: color.RGBA{
				R: 121,
				G: 0,
				B: 121,
				A: 124,
			},
			Size: 1,
		},
	}
	if len(p) < 3 {
		return &polygon, errors.New("point < 3")
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
	return &polygon, nil
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
	inmux := sync.RWMutex{}
	crossmux := sync.Mutex{}
	wg := sync.WaitGroup{}
	//小框框获取
	squareData := SquareDivided(p.Bound())
	//逐步扩大框大小
	containSet := make(map[string]struct{})
	//遍历获取是否在框框中或者路过框
	for i, code := range squareData {
		fmt.Println(i, len(squareData))
		wg.Add(1)
		go func(squarehashcode string) {
			defer wg.Done()
			//先判断是否有
			for i := 1; i < len(squarehashcode); i++ {
				key := squarehashcode[:i]
				inmux.RLock()
				_, ok := containSet[key]
				inmux.RUnlock()
				if ok {
					return
				}
			}
			box := NewBox(squarehashcode)
			boxStatus := p.SquareStatus(box)
			switch boxStatus {
			case boxStatusCross:
				crossmux.Lock()
				cross = append(cross, squarehashcode)
				crossmux.Unlock()
			case boxStatusInner:
				privLevel := squarehashcode
				for {
					//扩大方框等级
					squarehashcode = squarehashcode[:len(squarehashcode)-1]
					inmux.RLock()
					_, ok := containSet[squarehashcode]
					inmux.RUnlock()
					if !ok {
						if p.SquareStatus(NewBox(squarehashcode)) == boxStatusInner {
							privLevel = squarehashcode
						} else {
							inmux.Lock()
							containSet[privLevel] = struct{}{}
							inmux.Unlock()
							return
						}
					} else {
						break //已判断在内的，直接跳过
					}
				}
			}
		}(code)

	}
	wg.Wait()
	for k := range containSet {
		in = append(in, k)
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

func (p Polygon) CheckCross(startPoint, endPoint *Point) CrossStatus {
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

//点位是否在多边形内
func (p Polygon) PointIn(point *Point) (intersected bool) {
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

type boxStatus uint8

const (
	boxStatusInner boxStatus = iota + 1
	boxStatusOut
	boxStatusCross
)

//矩形是否在多边形内 初步判断
func (p Polygon) SquarePointIn(box Box) boxStatus {
	point1 := NewPoint(box.MaxLng, box.MaxLat)
	point2 := NewPoint(box.MaxLng, box.MinLat)
	point3 := NewPoint(box.MinLng, box.MaxLat)
	point4 := NewPoint(box.MinLng, box.MinLat)
	status1 := p.PointIn(point1)
	status2 := p.PointIn(point2)
	status3 := p.PointIn(point3)
	status4 := p.PointIn(point4)
	if status1 && status2 && status3 && status4 {
		return boxStatusInner
	}
	if status1 || status2 || status3 || status4 {
		return boxStatusCross
	}
	return boxStatusOut

}

func (p Polygon) SquareStatus(box Box) boxStatus {
	status := p.SquarePointIn(box)
	switch status {
	case boxStatusInner, boxStatusCross:
		for _, line := range p.lines {
			if box.IsLinesIntersected(line) {
				return boxStatusCross
			}
		}
		return boxStatusInner
	default:
		return boxStatusOut
	}
}

func (p Polygon) Export() [][]float64 {
	var res [][]float64
	for _, v := range p.points {
		x, y := v.Export()
		res = append(res, []float64{x, y})
	}
	return res
}

func (p Polygon) RenderToMap(ctx *sm.Context) {
	var dataList []s2.LatLng
	for _, v := range p.points {
		dataList = append(dataList, s2.LatLngFromDegrees(v.y, v.x))
	}
	ctx.AddObject(sm.NewArea(dataList,
		p.Color, p.FillColor, p.Size))
}
