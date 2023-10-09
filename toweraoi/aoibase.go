package toweraoi

import (
	"aoi"
	"aoi/linmath"
	"aoi/toweraoi/core"
	"fmt"
)

type SimpleMarker struct {
	ID   uint64
	Type string
	Pos  core.Position
}

func (mk *SimpleMarker) GetID() uint64 {
	return mk.ID
}
func (mk *SimpleMarker) GetType() string {
	return mk.Type
}
func (mk *SimpleMarker) GetPos() linmath.Vector2 {
	return linmath.Vector2{float32(mk.Pos.X), float32(mk.Pos.Y)}
}
func (mk *SimpleMarker) GetProps() []byte {
	return nil
}

type SimpleWatcher struct {
	ID        uint64
	WatchType string
	Pos       core.Position
	ViewRange int
}

func (w *SimpleWatcher) GetID() uint64 {
	return w.ID
}
func (w *SimpleWatcher) GetType() string {
	return w.WatchType
}
func (w *SimpleWatcher) GetPos() linmath.Vector2 {
	return linmath.Vector2{float32(w.Pos.X), float32(w.Pos.Y)}
}
func (w *SimpleWatcher) GetViewRange() int {
	return w.ViewRange
}
func (w *SimpleWatcher) OnMarkerEnterAOI(watchType string, maker aoi.IAOIMarker) {
	fmt.Println("SimpleWatcher.OnMarkerEnterAOI=>", watchType, maker.GetID(), maker.GetType())
}
func (w *SimpleWatcher) OnMarkerLeaveAOI(watchType string, maker aoi.IAOIMarker) {
	fmt.Println("SimpleWatcher.OnMarkerLeaveAOI<=", watchType, maker.GetID(), maker.GetType())
}
