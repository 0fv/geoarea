package geoarea

import (
	"sync"
)

type geoTireNode struct {
	Element   rune //
	Child     geoTireList
	CrossData map[interface{}]*Polygon
	InData    map[interface{}]*Polygon
}

type IGetKey interface {
	GetKey() interface{}
}

type geoTireList []*geoTireNode

func (g *geoTireList) Add(ploygon *Polygon, value IGetKey) {
	if ploygon.Empty() {
		return
	}
	corss, in := ploygon.Geohash()
	for _, v := range corss {
		n := g.GetAndSetNode(v)
		n.CrossData[value.GetKey()] = ploygon
	}
	for _, v := range in {
		n := g.GetAndSetNode(v)
		n.InData[value.GetKey()] = ploygon
	}
}

func (g *geoTireList) Get(point *Point) []interface{} {
	if point == nil {
		return nil
	}
	lng, lat := point.x, point.y
	hash := PointGeoHash(lat, lng)
	if len(hash) == 0 {
		return nil
	}
	return g.GetNodeData(hash, point)
}

func (g *geoTireNode) Del(key interface{}) {
	delete(g.InData, key)
	delete(g.CrossData, key)
	for _, v := range g.Child {
		v.Del(key)
	}
}

func (g *geoTireList) GetAndSetNode(hash string) *geoTireNode {
	t := g
	var cn *geoTireNode
	for _, v := range hash {
		cn = t.GetAndSetNodebyRune(v)
		t = &cn.Child
	}
	return cn
}

func (g *geoTireList) GetNodeData(hash string, point *Point) []interface{} {
	var ret []interface{}
	t := g
	var cn *geoTireNode
	for _, v := range hash {
		cn = t.GetNodebyRune(v)
		if cn == nil {
			break
		}
		//在内部的
		for k := range cn.InData {
			ret = append(ret, k)
		}
		//在边界的
		for k, v := range cn.CrossData {
			if v.PointIn(point) {
				ret = append(ret, k)
			}
		}
		t = &cn.Child
	}
	return ret
}

func (g *geoTireList) GetNodebyRune(r rune) *geoTireNode {
	for _, v := range *g {
		if v.Element == r {
			return v
		}
	}
	return nil
}

func (g *geoTireList) GetAndSetNodebyRune(r rune) *geoTireNode {
	for _, v := range *g {
		if v.Element == r {
			return v
		}
	}
	gt := &geoTireNode{
		Element:   r,
		Child:     make(geoTireList, 0),
		CrossData: make(map[interface{}]*Polygon),
		InData:    make(map[interface{}]*Polygon),
	}
	*g = append(*g, gt)
	return gt
}

type GeoTire struct {
	root geoTireList
	data map[interface{}]IGetKey
	rmux sync.RWMutex
}

func NewGeoTire() *GeoTire {
	return &GeoTire{
		root: make(geoTireList, 0),
		data: make(map[interface{}]IGetKey),
	}
}

// 添加
func (g *GeoTire) Add(ploygon *Polygon, value IGetKey) {
	g.rmux.Lock()
	defer g.rmux.Unlock()
	g.root.Add(ploygon, value)
	g.data[value.GetKey()] = value
}

// 存在创建，不存在更新
func (g *GeoTire) Set(ploygon *Polygon, value IGetKey) {
	g.rmux.Lock()
	defer g.rmux.Unlock()
	key := value.GetKey()
	delete(g.data, key)
	for _, v := range g.root {
		v.Del(key)
	}
	g.root.Add(ploygon, value)
	g.data[value.GetKey()] = value
}

func (g *GeoTire) SetMulti(ploygons []*Polygon, value IGetKey) {
	g.rmux.Lock()
	defer g.rmux.Unlock()
	key := value.GetKey()
	delete(g.data, key)
	for _, v := range g.root {
		v.Del(key)
	}
	for _, v := range ploygons {
		g.root.Add(v, value)
	}
	g.data[value.GetKey()] = value
}

// 获取在某区域的数据
func (g *GeoTire) Get(point *Point) []IGetKey {
	g.rmux.RLock()
	defer g.rmux.RUnlock()
	keys := g.root.Get(point)
	ret := make([]IGetKey, 0)
	existMap := map[interface{}]struct{}{}
	for _, v := range keys {
		if _, ok := existMap[v]; ok {
			continue
		}
		existMap[v] = struct{}{}
		ret = append(ret, g.data[v])
	}
	return ret
}

// 删除
func (g *GeoTire) Del(value IGetKey) {
	g.rmux.Lock()
	defer g.rmux.Unlock()
	key := value.GetKey()
	delete(g.data, key)
	for _, v := range g.root {
		v.Del(key)
	}
}
