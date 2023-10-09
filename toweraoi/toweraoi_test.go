package toweraoi

import (
	. "aoi/toweraoi/core"
	"fmt"
	"reflect"
	"testing"
)

type SimpleInterestTable struct {
	table map[string][]string
}

func (t *SimpleInterestTable) Init() {
	t.table = map[string][]string{}
	t.table["view"] = []string{"npc", "player", "item"}
	t.table["trigger"] = []string{"monster"}

}

func (t *SimpleInterestTable) GetWatchInterests(watchType string) []string {
	return t.table[watchType]
}

func TestTower(t *testing.T) {
	tower := NewTower()

	obj1 := &SimpleMarker{ID: 100, Type: "monster"}
	tower.AddMarker(obj1)
	obj2 := &SimpleMarker{ID: 101, Type: "npc"}
	tower.AddMarker(obj2)

	fmt.Printf("%v\n", tower.GetIds())
	fmt.Printf("%v\n", tower.GetIdsByTypes([]string{"monster"}))
	fmt.Printf("%v\n", tower.GetIdsByTypes([]string{"monster", "npc"}))

	fmt.Println("remove2......")
	tower.RemoveMarker(obj2)
	fmt.Printf("%v\n", tower.GetIdsByTypes([]string{"npc"}))
	fmt.Println("remove1......")
	tower.RemoveMarker(obj1)
	fmt.Printf("%v\n", tower.GetIdsByTypes([]string{"monster"}))
	fmt.Printf("%v\n", tower.GetIds())

	fmt.Println("waters.....")
	watcher1 := &SimpleWatcher{ID: 201, WatchType: "monster"}
	tower.AddWatcher(watcher1)
	watcher2 := &SimpleWatcher{ID: 202, WatchType: "npc"}
	tower.AddWatcher(watcher2)
	fmt.Printf("%v\n", tower.GetWatchers([]string{"monster", "npc"}))

	tower.RemoveWatcher(watcher2)
	fmt.Printf("%v\n", tower.GetWatchers([]string{"monster", "npc"}))

	fmt.Println(tower)
}

func TestTowerAOI(t *testing.T) {
	config := &Config{
		M: &MapConfig{
			&Rectangle{
				50,
				50,
			},
		},
		T: &TowerConfig{
			&Rectangle{
				10,
				10,
			},
		},
	}
	towerAOI := NewTowerAOI(config, func(msg interface{}) {
		fmt.Printf("event:%v=>%+v\n", reflect.TypeOf(msg), msg)
	})
	_ = towerAOI
	npcs := []*SimpleMarker{
		&SimpleMarker{ID: 1, Type: "npc", Pos: Position{11, 31}},
		&SimpleMarker{ID: 2, Type: "npc", Pos: Position{19, 31}},
		&SimpleMarker{ID: 3, Type: "npc", Pos: Position{21, 31}},
		&SimpleMarker{ID: 4, Type: "npc", Pos: Position{21, 31}},
	}

	watcher1 := &SimpleWatcher{ID: 101, WatchType: "npc", Pos: Position{11, 31}}

	for _, v := range npcs {
		towerAOI.AddMarker(v, v.Pos)
	}

	towerAOI.AddWatcher(watcher1, watcher1.Pos, 1)

	towerAOI.UpdateWatcher(watcher1, watcher1.Pos, Position{31, 21}, 1, 1)

	fmt.Println("maker0 move in")
	towerAOI.UpdateMarker(npcs[0], npcs[0].Pos, Position{31, 31})

	fmt.Println("maker0 move out")
	towerAOI.UpdateMarker(npcs[0], Position{31, 31}, Position{11, 11})

	fmt.Println("maker0 remove")
	towerAOI.RemoveMarker(npcs[0], Position{11, 11})

	towerAOI.RemoveWatcher(watcher1, Position{31, 21}, 1)

}

func TestTowerAOISceneForward(t *testing.T) {
	config := &Config{
		M: &MapConfig{
			&Rectangle{
				50,
				50,
			},
		},
		T: &TowerConfig{
			&Rectangle{
				10,
				10,
			},
		},
	}

	//添加兴趣过滤
	it := &SimpleInterestTable{}
	it.Init()
	SetWatchInterrestTable(it)

	scene := NewAOIScene(config)
	npcs := []*SimpleMarker{
		&SimpleMarker{ID: 1, Type: "npc", Pos: Position{11, 31}},
		&SimpleMarker{ID: 2, Type: "npc", Pos: Position{19, 31}},
		&SimpleMarker{ID: 3, Type: "npc", Pos: Position{21, 31}},
		&SimpleMarker{ID: 4, Type: "monster", Pos: Position{21, 31}},
	}

	watcher1 := &SimpleWatcher{ID: 101, WatchType: "player", Pos: Position{11, 31}, ViewRange: 1}

	for _, v := range npcs {
		scene.AddMarker(v.ID, v, v.GetPos(), 0)
	}

	fmt.Println("add watcher")

	scene.AddWatcher(watcher1.GetID(), "view", watcher1, watcher1.GetPos(), 10, 0)
	scene.AddWatcher(watcher1.GetID(), "trigger", watcher1, watcher1.GetPos(), 20, 0)
	//update
	for i := 0; i < 3; i++ {
		//scene.BeforeUpdate()
		fmt.Println("frame:", i)

		if i == 1 {
			fmt.Println("add m5")
			m5 := &SimpleMarker{ID: 5, Type: "npc", Pos: Position{21, 31}}
			scene.AddMarker(m5.GetID(), m5, m5.GetPos(), 0)
		}

		watcher1.Pos.X = watcher1.Pos.X + 10
		fmt.Println("w move to:", watcher1.Pos)
		scene.MoveWatcher(watcher1.GetID(), watcher1.GetPos())
		//scene.Update()
		//fmt.Println(i," view:",scene.GetViewMarkers(watcher1.GetID()))
	}
}
