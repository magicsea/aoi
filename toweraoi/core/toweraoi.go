package core

import (
	"fmt"
	"math"
)

type IDType = uint64

type Watcher interface {
	GetID() IDType
	GetType() string
}

type Marker interface {
	GetID() IDType
	GetType() string
}

type MapConfig struct {
	*Rectangle
}

type TowerConfig struct {
	*Rectangle
}

type Config struct {
	//地图配置
	M *MapConfig
	//格子配置
	T     *TowerConfig
	Limit int
}

type TowerAOI struct {
	config     *Config
	towers     [][]*Tower
	rangeLimit int
	max        *Position

	emit func(msg interface{})
}

func NewTowerAOI(conf *Config, emitFun func(msg interface{})) *TowerAOI {
	ta := &TowerAOI{
		config: conf,
		emit:   emitFun,
	}

	var m int = int(math.Ceil(float64(ta.config.M.Width/ta.config.T.Width))) + 1
	var n int = int(math.Ceil(float64(ta.config.M.Height/ta.config.T.Height))) + 1
	fmt.Println("make map:", m, n)
	ta.max = &Position{
		X: m - 1,
		Y: n - 1,
	}

	ta.towers = make([][]*Tower, m)
	for i := 0; i < m; i++ {
		ta.towers[i] = make([]*Tower, n)
		for j := 0; j < n; j++ {
			ta.towers[i][j] = NewTower()
		}
	}
	return ta
}

func (t *TowerAOI) GetConfig() *Config {
	return t.config
}

// Get given type object ids from tower aoi by range and types
// @param pos {Object} The pos to find objects
// @param range {Number} The range to find the object, in tower aoi, it means the tower number from the pos
// @param types {Array} The types of the object need to find
func (t *TowerAOI) GetIdsByRange(pos *Position, r int, types []string) map[string]map[IDType]IDType {
	if !t.checkPos(pos) || r < 0 || r > t.rangeLimit {
		return nil
	}

	p := t.transPos(pos)
	if p == nil {
		fmt.Println("p value is nil.")
	}
	start, end := t.getPosLimit(p, r, t.max)

	result := make(map[string]map[IDType]IDType)
	for i := start.X; i <= end.X; i++ {
		for j := start.Y; j <= end.Y; j++ {
			result = t.addMapByTypes(result, t.towers[i][j].GetIdsByTypes(types), types)
		}
	}
	return result
}

// Get all object ids from tower aoi by pos and range
//func (t *TowerAOI) GetIdsByPos(pos *Position, r int) []IDType {
//	if t.checkPos(pos) || r < 0 {
//		return nil
//	}
//
//	if r > 5 {
//		r = 5
//	}
//
//	p := t.transPos(pos)
//	start, end := t.getPosLimit(p, r, t.max)
//	var result []IDType
//	for i := start.X; i <= end.X; i++ {
//		for j := start.Y; j <= end.Y; j++ {
//			result = AddMap(result, t.towers[i][j].getIds())
//		}
//	}
//	return result
//}

// Add an object to tower aoi at given pos
func (t *TowerAOI) AddMarker(obj Marker, pos Position) bool {
	if t.checkPos(&pos) {
		p := t.transPos(&pos)
		t.towers[p.X][p.Y].AddMarker(obj)
		t.emit(AddMarkerEvent{Id: obj.GetID(), Stype: obj.GetType(), Watchers: t.towers[p.X][p.Y].watcherTypeMap})
		return true
	}
	return false
}

// Remove object from aoi module
func (t *TowerAOI) RemoveMarker(obj Marker, pos Position) bool {
	if t.checkPos(&pos) {
		p := t.transPos(&pos)
		t.towers[p.X][p.Y].RemoveMarker(obj)
		t.emit(RemoveMarkerEvent{Id: obj.GetID(), Stype: obj.GetType(), Watchers: t.towers[p.X][p.Y].watcherTypeMap})
		return true
	}
	return false
}

func (t *TowerAOI) UpdateMarker(obj Marker, oldPos Position, newPos Position) bool {
	if !t.checkPos(&oldPos) || !t.checkPos(&newPos) {
		return false
	}

	p1 := t.transPos(&oldPos)
	p2 := t.transPos(&newPos)

	if p1.X == p2.X && p1.Y == p2.Y {
		return true
	} else {
		if t.towers[p1.X] == nil || t.towers[p2.X] == nil {
			fmt.Printf("AOI pos error ! oldPos : %v, newPos : %v, p1 : %v, p2 : %v\n", oldPos, newPos, p1, p2)
			return false
		}

		oldTower := t.towers[p1.X][p1.Y]
		newTower := t.towers[p2.X][p2.Y]

		oldTower.RemoveMarker(obj)
		newTower.AddMarker(obj)

		oldmap := oldTower.watcherTypeMap
		newmap := newTower.watcherTypeMap
		//删除视野内格子间移动的节点
		samePart := map[IDType]string{}
		for tp, oldm := range oldTower.watcherTypeMap {
			newm, foundt := newTower.watcherTypeMap[tp]
			if foundt {
				for old, _ := range oldm {
					if _, b := newm[old]; b {
						samePart[old] = tp
					}
				}
			}
		}

		t.emit(UpdateMarkerEvent{Id: obj.GetID(), Stype: obj.GetType(), OldWatchers: oldmap, NewWatchers: newmap, SameMap: samePart})
		return true
	}
	return false
}

// Check if the pos is valid;
// @return {Boolean} Test result
func (t *TowerAOI) checkPos(pos *Position) bool {
	if pos == nil {
		return false
	}
	if pos.X < 0 || pos.Y < 0 || pos.X >= t.config.M.Width || pos.Y >= t.config.M.Height {
		return false
	}
	return true
}

// Trans the absolut pos to tower pos. For example : (210, 110} -> (1, 0), for tower width 200, height 200
func (t *TowerAOI) transPos(pos *Position) *Position {
	return &Position{
		X: int(math.Floor(float64(pos.X) / float64(t.config.T.Width))),
		Y: int(math.Floor(float64(pos.Y) / float64(t.config.T.Height))),
	}
}

func (t *TowerAOI) AddWatcher(watcher Watcher, pos Position, r int) {
	if r < 0 {
		return
	}

	if r > 5 {
		r = 5
	}

	p := t.transPos(&pos)
	start, end := t.getPosLimit(p, r, t.max)

	var addObjs []IDType
	for i := start.X; i <= end.X; i++ {
		for j := start.Y; j <= end.Y; j++ {
			t.towers[i][j].AddWatcher(watcher)
			ids := t.towers[i][j].GetIds()
			addObjs = AddMap(addObjs, ids)
		}
	}

	t.emit(AddWatcherEvent{WatcherId: watcher.GetID(), WatcherType: watcher.GetType(), AddMarkers: addObjs})
}

func (t *TowerAOI) RemoveWatcher(watcher Watcher, pos Position, r int) {
	if r < 0 {
		return
	}

	if r > 5 {
		r = 5
	}

	p := t.transPos(&pos)
	start, end := t.getPosLimit(p, r, t.max)

	var removeObjs []IDType
	for i := start.X; i <= end.X; i++ {
		for j := start.Y; j <= end.Y; j++ {
			t.towers[i][j].RemoveWatcher(watcher)
			ids := t.towers[i][j].GetIds()
			removeObjs = AddMap(removeObjs, ids)
		}
	}

	t.emit(RemoveWatcherEvent{WatcherId: watcher.GetID(), WatcherType: watcher.GetType(), RemoveMarkers: removeObjs})
}

func (t *TowerAOI) UpdateWatcher(watcher Watcher, oldPos Position, newPos Position, oldRange int, newRange int) bool {
	if !t.checkPos(&oldPos) || !t.checkPos(&newPos) {
		return false
	}

	p1 := t.transPos(&oldPos)
	p2 := t.transPos(&newPos)

	if p1.X == p2.X && p1.Y == p2.Y && oldRange == newRange {
		return true
	} else {
		if oldRange < 0 || newRange < 0 {
			return false
		}

		if oldRange > 5 {
			oldRange = 5
		}
		if newRange > 5 {
			newRange = 5
		}

		removeTowers, addTowers, unChangeTowers := t.getChangedTowers(p1, p2, oldRange, newRange, t.towers, t.max)
		var addObjs []IDType
		var removeObjs []IDType

		for i := 0; i < len(addTowers); i++ {
			addTowers[i].AddWatcher(watcher)
			ids := addTowers[i].GetIds()
			addObjs = AddMap(addObjs, ids)
		}

		for i := 0; i < len(removeTowers); i++ {
			removeTowers[i].RemoveWatcher(watcher)
			ids := removeTowers[i].GetIds()
			removeObjs = AddMap(removeObjs, ids)
		}
		_ = unChangeTowers
		//fmt.Printf("unChangeTowers: %v ?????", unChangeTowers)

		t.emit(UpdateWatcherEvent{WatcherId: watcher.GetID(), WatcherType: watcher.GetType(), AddMarkers: addObjs, RemoveMarkers: removeObjs})
		return true
	}
	return false
}

// Get changed towers for girven pos
// @param p1 {Object} The origin position
// @param p2 {Object} The now position
// @param oldRange {Number} The old range
// @param newRange {Number} The new range
// @param towers {Object} All towers of the aoi
// @param max {Object} The position limit of the towers
func (t *TowerAOI) getChangedTowers(p1 *Position, p2 *Position, oldRange int, newRange int, towers [][]*Tower, max *Position) (removeTowers []*Tower, addTowers []*Tower, unChangeTowers []*Tower) {
	start1, end1 := t.getPosLimit(p1, oldRange, max)
	start2, end2 := t.getPosLimit(p2, newRange, max)

	for x := start1.X; x <= end1.X; x++ {
		for y := start1.Y; y <= end1.Y; y++ {
			if isInRect(&Position{x, y}, start2, end2) {
				//if unChangeTowers == nil {
				//	unChangeTowers = make([]*Tower, 1)
				//}

				unChangeTowers = append(unChangeTowers, towers[x][y])
			} else {
				//if removeTowers == nil {
				//	removeTowers = make([]*Tower, 1)
				//}

				removeTowers = append(removeTowers, towers[x][y])
			}
		}
	}

	for x := start2.X; x <= end2.X; x++ {
		for y := start2.Y; y <= end2.Y; y++ {
			if !isInRect(&Position{x, y}, start1, end1) {
				//if addTowers == nil {
				//	addTowers = make([]*Tower, 1)
				//}

				addTowers = append(addTowers, towers[x][y])
			}
		}
	}
	return
}

// Get the postion limit of given range
// @param pos {Object} The center position
// @param range {Number} The range
// @param max {max} The limit, the result will not exceed the limit
// @return The pos limitition
func (t *TowerAOI) getPosLimit(pos *Position, r int, max *Position) (start *Position, end *Position) {
	if start == nil {
		start = &Position{}
	}
	if end == nil {
		end = &Position{}
	}

	if pos.X-r < 0 {
		start.X = 0
		end.X = 2 * r
	} else if pos.X+r > max.X {
		end.X = max.X
		start.X = max.X - 2*r
	} else {
		start.X = pos.X - r
		end.X = pos.X + r
	}

	if pos.Y-r < 0 {
		start.Y = 0
		end.Y = 2 * r
	} else if pos.Y+r > max.Y {
		end.Y = max.Y
		start.Y = max.Y - 2*r
	} else {
		start.Y = pos.Y - r
		end.Y = pos.Y + r
	}

	if start.X < 0 {
		start.X = 0
	}
	if end.X > max.X {
		end.X = max.X
	}

	if start.Y < 0 {
		start.Y = 0
	}
	if end.Y > max.Y {
		end.Y = max.Y
	}
	return
}

// Check if the pos is in the rect
func isInRect(pos *Position, start *Position, end *Position) bool {
	return (pos.X >= start.X && pos.X <= end.X && pos.Y >= start.Y && pos.Y <= end.Y)
}

func (t *TowerAOI) GetWatchers(pos Position, types []string) map[string]map[IDType]IDType {
	if t.checkPos(&pos) {
		p := t.transPos(&pos)
		return t.towers[p.X][p.Y].GetWatchers(types)
	}
	return nil
}

// Combine map to arr
// @param arr {Array} The array to add the map to
// @param map {Map} The map to add to array
func AddMap(result []IDType, m map[IDType]IDType) []IDType {
	r := make([]IDType, len(m))
	i := 0
	for _, v := range m {
		r[i] = v
		i++
	}
	return append(result, r...)
}

func (t *TowerAOI) addMapByTypes(result map[string]map[IDType]IDType, m map[string]map[IDType]IDType, types []string) map[string]map[IDType]IDType {
	for i := 0; i < len(types); i++ {
		objType := types[i]

		_, r1 := m[objType]
		if !r1 {
			continue
		}

		_, r2 := result[objType]
		if !r2 {
			result[objType] = make(map[IDType]IDType)
		}

		for k, v := range m[objType] {
			result[objType][k] = v
		}
	}
	return result
}
