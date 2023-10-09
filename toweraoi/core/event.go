package core

type AddMarkerEvent struct {
	Id IDType
	Stype string
	Watchers map[string]map[IDType]IDType
}

type RemoveMarkerEvent struct {
	Id IDType
	Stype string
	Watchers map[string]map[IDType]IDType
}

type UpdateMarkerEvent struct {
	Id IDType
	Stype string
	OldWatchers map[string]map[IDType]IDType
	NewWatchers map[string]map[IDType]IDType
	SameMap map[IDType]string
}


type UpdateWatcherEvent struct {
	WatcherId IDType
	WatcherType string
	AddMarkers []IDType
	RemoveMarkers []IDType
}

type AddWatcherEvent struct {
	WatcherId IDType
	WatcherType string
	AddMarkers []IDType
}


type RemoveWatcherEvent struct {
	WatcherId IDType
	WatcherType string
	RemoveMarkers []IDType
}