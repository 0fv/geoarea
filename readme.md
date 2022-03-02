geohash+前缀树快速判断点是否在区域内

### 初始化

```go
geoTire := NewGeoTire()
```

### 添加多边形区域

```go
geoTire.Add(ploygon,igetKey)
```
值需要实现IGetKey接口

### 查询点位是否在某区域内

```go
geoTire.Get(point)
```

### 删除某区域

```go
geoTire.Del(igetKey)
```

#### 思路

将需要判断的多边形通过geohash划分，并将geohash分为两类，一种在多边形内，一种与多边形相交，把与此多边形相关的geohash放入前缀树中，若包含类型的，直接命中，若相交类型的，通过射线法对其判断，若射线法判断为相交，则命中，否则不命中

![成都区域例子](./cd.png)

上图为成都市范围内geohash方格划分示例