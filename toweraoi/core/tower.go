package core

type Tower struct {
	markerIds      map[IDType]IDType
	markerTypeMap   map[string]map[IDType]IDType
	watcherTypeMap map[string]map[IDType]IDType
	//size     int
}

func NewTower() *Tower {
	t := &Tower{
		markerIds:make(map[IDType]IDType),
		watcherTypeMap: make(map[string]map[IDType]IDType),
		markerTypeMap:   make(map[string]map[IDType]IDType),
	}
	return t
}

func (t *Tower) AddMarker(obj Marker) bool {
	var id = obj.GetID()
	var tp = obj.GetType()

	t.markerIds[id] = id

	if tp != "" {
		if _, b := t.markerTypeMap[tp]; !b {
			t.markerTypeMap[tp] = map[IDType]IDType{}
		}
		if t.markerTypeMap[tp][id] == id {
			return false
		}

		t.markerTypeMap[tp][id] = id
		//t.size++
		return true
	} else {
		return false
	}
}

func  (t *Tower) RemoveMarker(obj Marker)  {
	var id = obj.GetID()
	var tp = obj.GetType()

	if _,b:=t.markerIds[id];b {
		delete(t.markerIds,id)
		if tp!="" {
			delete(t.markerTypeMap[tp],id)
		}
		//t.size--
	}
}

func  (t *Tower) AddWatcher(watcher Watcher)  {
	var tp = watcher.GetType()
	var id = watcher.GetID()

	if tp!="" {
		if _,b:=t.watcherTypeMap[tp];!b {
			t.watcherTypeMap[tp] = map[IDType]IDType{}
		}
		t.watcherTypeMap[tp][id]=id
	}
}

func  (t *Tower) RemoveWatcher(watcher Watcher)  {
	var tp = watcher.GetType()
	var id = watcher.GetID()


	if _,b:=t.watcherTypeMap[tp];tp!=""&&b {
		delete(t.watcherTypeMap[tp],id)
	}
}

//watchers
func (t *Tower) GetWatchers(types []string) map[string]map[IDType]IDType{
	var result = map[string]map[IDType]IDType{}
	if types!=nil&&len(types)>0 {
		for _, tp := range types {
			if m,b:=t.watcherTypeMap[tp];b {
				result[tp]=m
			}
		}
	}
	return result
}


func (t *Tower) GetIds() map[IDType]IDType {
	return t.markerIds
}


func (t *Tower) getIdsByType_(objType string) map[IDType]IDType {

	result := make(map[IDType]IDType)

	_, ok := t.markerTypeMap[objType]
	if ok {
		result = t.markerTypeMap[objType]
		}
	return result
}

// Get object ids of given types in this tower
func (t *Tower) GetIdsByTypes(types []string) map[string]map[IDType]IDType {
	result := make(map[string]map[IDType]IDType)

	for i := 0; i < len(types); i++ {
		objType := types[i]
		_, ok := t.markerTypeMap[objType]

		if ok {
			result[objType] = t.markerTypeMap[objType]
		}
	}

	return result
}

