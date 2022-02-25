package geoarea

import (
	"github.com/mmcloughlin/geohash"
)

var precision uint = 7

func SetPrecision(p uint) {
	precision = p
}

func PointGeoHash(lat, lng float64) string {
	return geohash.EncodeWithPrecision(lat, lng, precision)
}

//方形区域划分
func SquareDivided(maxlat, minlat, maxlng, minlng float64) []string {
	//左上开始
	result := make([]string, 0)
	lat := maxlat
	lng := minlng
	hash := geohash.EncodeWithPrecision(lat, lng, precision)
	result = append(result, hash)
	movedirection := geohash.East
	for { //开始从左至右进行查找
		hash = geohash.Neighbor(hash, movedirection)
		result = append(result, hash)
		switch movedirection {
		case geohash.East: //向东移动，若超过最大纬度，则向南一步，然后向西移动
			box := geohash.BoundingBox(hash)
			if box.MaxLng > maxlng {
				if box.MinLat < minlat { //当前是否大于最小经度，小于说明就结束了
					return result
				}
				hash = geohash.Neighbor(hash, geohash.South)
				result = append(result, hash)
				movedirection = geohash.West
			}
		case geohash.West: // 向西移动，若小于最小纬度，则向南一步，然后向西移动
			box := geohash.BoundingBox(hash)
			if box.MinLng < minlng {
				if box.MinLat < minlat {
					return result
				}
				hash = geohash.Neighbor(hash, geohash.South)
				result = append(result, hash)
				movedirection = geohash.East
			}

		}
	}
}
