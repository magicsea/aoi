package toweraoi

import (
	. "aoi"
	"aoi/linmath"
	"aoi/toweraoi/core"
	"errors"
)

type TIDType = core.IDType
type AOIScene struct {
	aoiArea *core.TowerAOI

	//上层对象缓存
	watchers map[TIDType]IAOIWatcher
	markers  map[TIDType]IAOIMarker
	//临时数据
	watchersTemp map[TIDType]map[string]*WatcherDataSnap
	markerTemp   map[TIDType]*MarkerDataSnap

	//全局可见marker
	globalMarkers map[TIDType]*MarkerDataSnap

	//视野传递:[key]传递给[value]
	watchLinkMap map[TIDType][]TIDType
}

// var ao IAOI = &AOIScene{}
type AOIDataSnap struct {
	id    TIDType
	pos   core.Position
	layer AOILayer
	stype string
}

type MarkerDataSnap struct {
	AOIDataSnap
}
type WatcherDataSnap struct {
	AOIDataSnap
	viewRange int
}

func ToAOIPos(v linmath.Vector2) core.Position {
	return core.Position{int(v.X), int(v.Y)}
}

//type WatcherEntity struct {
//	watchTypes map[string]*WatcherDataSnap
//}

func (snap *AOIDataSnap) GetID() TIDType {
	return snap.id
}
func (snap *AOIDataSnap) GetPos() core.Position {
	return snap.pos
}
func (snap *AOIDataSnap) GetType() string {
	return snap.stype
}

var (
	ERR_TODO          = errors.New("todo")
	ERR_EXIST_WATCHER = errors.New("exist watcher type")
	ERR_EXIST_MAKER   = errors.New("exist marker")
	ERR_NOT_FOUND     = errors.New("not found")
)
var watchInterrestTable IAOIWatchInterestTable

func SetWatchInterrestTable(table IAOIWatchInterestTable) {
	watchInterrestTable = table
}

func NewAOIScene(conf *core.Config) *AOIScene {

	s := &AOIScene{
		watchers: make(map[TIDType]IAOIWatcher),
		markers:  make(map[TIDType]IAOIMarker),

		watchersTemp: make(map[TIDType]map[string]*WatcherDataSnap),
		markerTemp:   make(map[TIDType]*MarkerDataSnap),
	}

	aoiArea := core.NewTowerAOI(conf, s.onRecvEvent)
	s.aoiArea = aoiArea
	return s
}

func (s *AOIScene) onRecvEvent(msg interface{}) {
	//fmt.Printf("event:%v=>%+v\n",reflect.TypeOf(msg),msg)
	switch tm := msg.(type) {
	case core.AddMarkerEvent:
		for stype, m := range tm.Watchers {
			for wid, _ := range m {
				s.onEnterWatcherView(wid, tm.Id, stype)
			}
		}
	case core.RemoveMarkerEvent:
		for stype, m := range tm.Watchers {
			for wid, _ := range m {
				s.onLeaveWatcherView(wid, tm.Id, stype)
			}
		}
	case core.UpdateMarkerEvent:
		for stype, m := range tm.OldWatchers {
			for wid, _ := range m {
				if _, b := tm.SameMap[wid]; b {
					//fmt.Println("ign cell leave:", wid)
					continue
				}
				s.onLeaveWatcherView(wid, tm.Id, stype)
			}
		}
		for stype, m := range tm.NewWatchers {
			for wid, _ := range m {
				if _, b := tm.SameMap[wid]; b {
					//fmt.Println("ign cell enter:", wid)
					continue
				}
				s.onEnterWatcherView(wid, tm.Id, stype)
			}
		}
	case core.AddWatcherEvent:
		for _, mkid := range tm.AddMarkers {
			s.onEnterWatcherView(tm.WatcherId, mkid, tm.WatcherType)
		}
	case core.RemoveWatcherEvent:
		for _, mkid := range tm.RemoveMarkers {
			s.onLeaveWatcherView(tm.WatcherId, mkid, tm.WatcherType)
		}
	case core.UpdateWatcherEvent:
		for _, mkid := range tm.RemoveMarkers {
			s.onLeaveWatcherView(tm.WatcherId, mkid, tm.WatcherType)
		}
		for _, mkid := range tm.AddMarkers {
			s.onEnterWatcherView(tm.WatcherId, mkid, tm.WatcherType)
		}
	}
}

func checkTypeIn(all []string, v string) bool {
	for _, a := range all {
		if a == v {
			return true
		}
	}
	return false
}

func (s *AOIScene) onEnterWatcherView(watcherId TIDType, markerId TIDType, stype string) {
	//fmt.Println("onEnterWatcherView:", watcherId, markerId, stype)
	w := s.watchers[watcherId]
	m := s.markers[markerId]

	if watchInterrestTable != nil {
		tps := watchInterrestTable.GetWatchInterests(stype)
		if !checkTypeIn(tps, m.GetType()) {
			return
		}
	}

	w.OnMarkerEnterAOI(stype, m)
}

func (s *AOIScene) onLeaveWatcherView(watcherId TIDType, markerId TIDType, stype string) {
	//fmt.Println("onLeaveWatcherView:", watcherId, markerId, stype)

	w := s.watchers[watcherId]
	m := s.markers[markerId]

	if watchInterrestTable != nil {
		tps := watchInterrestTable.GetWatchInterests(stype)
		if !checkTypeIn(tps, m.GetType()) {
			return
		}
	}

	w.OnMarkerLeaveAOI(stype, m)
}

func (s *AOIScene) ViewToCellRange(visualField float32) int {
	cellRange := int(visualField) / s.aoiArea.GetConfig().T.Width
	if cellRange < 1 {
		cellRange = 1
	}
	return cellRange
}

func (s *AOIScene) AddMarker(id TIDType, marker IAOIMarker, pos linmath.Vector2, layer AOILayer) error {
	if _, b := s.markerTemp[id]; b {
		return ERR_EXIST_MAKER
	}

	s.markers[id] = marker

	temp := &MarkerDataSnap{AOIDataSnap{id: id, pos: ToAOIPos(pos), layer: layer, stype: marker.GetType()}}
	s.markerTemp[id] = temp
	s.aoiArea.AddMarker(temp, ToAOIPos(pos))
	return nil
}
func (s *AOIScene) RemoveMarker(id TIDType) error {
	temp, b := s.markerTemp[id]
	if !b {
		return ERR_NOT_FOUND
	}

	s.aoiArea.RemoveMarker(temp, temp.pos)
	delete(s.markerTemp, id)
	delete(s.markers, id)
	return nil
}
func (s *AOIScene) MoveMarker(id TIDType, pos linmath.Vector2) error {
	temp, b := s.markerTemp[id]
	if !b {
		return ERR_NOT_FOUND
	}
	s.aoiArea.UpdateMarker(temp, temp.pos, ToAOIPos(pos))
	temp.pos = ToAOIPos(pos)
	return nil
}
func (s *AOIScene) MoveMarkerLayer(layer AOILayer) error {
	return ERR_TODO
}

func (s *AOIScene) AddWatcher(id uint64, watcherType string, watcher IAOIWatcher, pos linmath.Vector2, visualField float32, layer AOILayer) error {
	cellRange := s.ViewToCellRange(visualField)

	temp := &WatcherDataSnap{AOIDataSnap: AOIDataSnap{id: id, pos: ToAOIPos(pos), layer: layer, stype: watcherType}, viewRange: cellRange}
	m, b := s.watchersTemp[id]
	if !b {
		s.watchersTemp[id] = map[string]*WatcherDataSnap{}
	} else {
		if _, exist := m[watcherType]; exist {
			return ERR_EXIST_WATCHER
		}
	}
	if _, bw := s.watchers[id]; !bw {
		s.watchers[id] = watcher
	}
	s.watchersTemp[id][watcherType] = temp
	s.aoiArea.AddWatcher(temp, ToAOIPos(pos), cellRange)
	return nil
}
func (s *AOIScene) RemoveWatcher(id uint64, watcherType string) error {
	m, b := s.watchersTemp[id]
	if !b {
		return ERR_NOT_FOUND
	}
	temp, bt := m[watcherType]
	if !bt {
		return ERR_NOT_FOUND
	}
	s.aoiArea.RemoveWatcher(temp, temp.pos, temp.viewRange)

	delete(s.watchersTemp[id], watcherType)
	if len(s.watchersTemp[id]) < 1 {
		delete(s.watchersTemp, id)
		delete(s.watchers, id)
	}

	return nil
}
func (s *AOIScene) MoveWatcher(id TIDType, pos linmath.Vector2) error {
	m, b := s.watchersTemp[id]
	if !b {
		return ERR_NOT_FOUND
	}
	for _, temp := range m {
		s.aoiArea.UpdateWatcher(temp, temp.pos, ToAOIPos(pos), temp.viewRange, temp.viewRange)
		temp.pos = ToAOIPos(pos)
	}

	return nil
}
func (s *AOIScene) UpdateWatcherVisualField(id uint64, typ string, visualField float32) error {
	cellRange := s.ViewToCellRange(visualField)

	m, b := s.watchersTemp[id]
	if !b {
		return ERR_NOT_FOUND
	}

	for _, temp := range m {
		s.aoiArea.UpdateWatcher(temp, temp.pos, temp.pos, temp.viewRange, cellRange)
		temp.viewRange = cellRange
	}

	return nil
}

func (s *AOIScene) MoveWatcherLayer(layer AOILayer) error {
	return ERR_TODO
}

func (s *AOIScene) AddExtraMarker(watcher TIDType, marker TIDType) error {
	return ERR_TODO
}
func (s *AOIScene) RemoveExtraMarker(watcher TIDType, marker TIDType) error {
	return ERR_TODO
}

func (s *AOIScene) AddGlobalMarker(id TIDType, marker IAOIMarker, pos linmath.Vector2, layer AOILayer) error {
	return ERR_TODO
}
func (s *AOIScene) SetGlobalMarker(id TIDType) error {
	return ERR_TODO
}
func (s *AOIScene) UnsetGlobalMarker(id TIDType) error {
	return ERR_TODO
}

func (s *AOIScene) TravsalAOI(marker IAOIMarker, tps []string, f func(IAOIWatcher)) {
	mmap := s.aoiArea.GetWatchers(ToAOIPos(marker.GetPos()), tps)
	for _, m := range mmap {
		for k, _ := range m {
			w := s.watchers[k]
			f(w)
		}
	}
}
