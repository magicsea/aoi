package aoi

import "aoi/linmath"

// AOILayer 层级
type AOILayer uint8

// IAOI AOI系统对外提供的接口
type IAOI interface {
	AddMarker(id uint64, marker IAOIMarker, pos linmath.Vector2, layer AOILayer) error
	RemoveMarker(id uint64) error
	MoveMarker(id uint64, pos linmath.Vector2) error
	MoveMarkerLayer(layer AOILayer) error

	AddWatcher(id uint64, watchType string, watcher IAOIWatcher, pos linmath.Vector2, visualField float32, layer AOILayer) error
	RemoveWatcher(id uint64, watchType string) error
	MoveWatcher(id uint64, pos linmath.Vector2) error
	MoveWatcherLayer(layer AOILayer) error
	AddExtraMarker(watcher uint64, marker uint64) error
	RemoveExtraMarker(watcher uint64, marker uint64) error
	UpdateWatcherVisualField(id uint64, typ string, visualField float32) error

	AddGlobalMarker(id uint64, marker IAOIMarker, pos linmath.Vector2, layer AOILayer) error
	SetGlobalMarker(id uint64) error
	UnsetGlobalMarker(id uint64) error

	TravsalAOI(IAOIMarker, []string, func(IAOIWatcher))
}

// IAOIMarker AOI中一个可被别人看见的物件
type IAOIMarker interface {
	GetID() uint64
	GetType() string
	GetPos() linmath.Vector2

	GetProps() []byte
}

// IAOIWatcher AOI中一个可以观察别人的物件
type IAOIWatcher interface {
	GetID() uint64
	GetType() string
	GetPos() linmath.Vector2

	OnMarkerEnterAOI(watchType string, marker IAOIMarker)
	OnMarkerLeaveAOI(watchType string, marker IAOIMarker)
}

// 观察兴趣列表 (view:npc,player,mapitem)
type IAOIWatchInterestTable interface {
	GetWatchInterests(watchType string) []string
}
