package main

import (
	"aoi"
	"aoi/toweraoi"
	"aoi/toweraoi/core"
	"fmt"
	"github.com/hajimehoshi/ebiten/examples/resources/images/blocks"
	g "github.com/magicsea/gosprite"
	"github.com/magicsea/gosprite/ui"
	"image/color"
	"math/rand"
	"strconv"
)

type ToolScane struct {
	g.Scene
	aoiScene aoi.IAOI
	watchers map[uint64]*Watcher
	makers   map[uint64]*Marker

	idgen  uint64
	notice *ui.TextBox
}

func (s *ToolScane) initBg() error {

	//text
	//fontbt,errlaodfd := ioutil.ReadFile("msyh.ttc")
	//if errlaodfd != nil {
	//	fmt.Println("loadfont:",errlaodfd)
	//	return errlaodfd
	//}
	//g.LoadFont("msyh",fontbt)
	//

	//bg
	cell := 100
	var fcellsize = float64(cell)
	bg2, bgErr2 := g.NewSprite(blocks.Background_png)
	if bgErr2 != nil {
		fmt.Println("bg load error:", bgErr2)
		return bgErr2
	}
	bg2.SetSize(g.NewVector(screenW, screenH))
	bg2.SetSpriteType(g.SpriteTypeSlice)
	bg2.SetScale(g.NewVector(fcellsize/32.0, fcellsize/32.0))
	s.AddNode(bg2)

	//load line
	for i := 0; i < screenW/cell; i++ {
		from := g.NewVector(float64(i*cell), 0)
		to := g.NewVector(float64(i*cell), screenH)
		line := g.NewLine(from, to, 1, color.Black)
		line.SetDepth(2)
		s.AddNode(line)
	}
	for j := 0; j < screenH/cell; j++ {
		from := g.NewVector(0, float64(j*cell))
		to := g.NewVector(screenW, float64(j*cell))
		line := g.NewLine(from, to, 1, color.Black)
		line.SetDepth(2)
		s.AddNode(line)
	}
	return nil
}

func (s *ToolScane) initUI() error {
	var fromx float64 = screenW - 130
	bg := ui.NewTextBox(g.NewVector(fromx, 0), g.NewVector(130, 260), "", ui.AliVertical_Mid, ui.AliHorizontal_Right, false)
	bg.SetColor(color.RGBA{128, 128, 0, 128})
	s.AddUINode(bg)

	fromx += 10
	tb := ui.NewTextBox(g.NewVector(fromx, 0), g.NewVector(70, 40), "addmaker:", ui.AliVertical_Mid, ui.AliHorizontal_Right, false)
	tb.SetColor(color.RGBA{128, 128, 128, 0})
	s.AddUINode(tb)
	inf := ui.NewInputField(g.NewVector(screenW-50, 0), g.NewVector(50, 40), "100", ui.AliVertical_Mid, ui.AliHorizontal_Right, false)
	inf.SetColor(color.RGBA{128, 128, 128, 128})
	s.AddUINode(inf)
	btn := ui.NewButton(g.NewVector(fromx, 50), g.NewVector(120, 40), "AddMarker")
	btn.SetColor(color.RGBA{100, 100, 100, 128})
	btn.SetOnPressed(func(b *ui.Button) {
		num, err := strconv.Atoi(inf.ValueText)
		if num <= 0 || num > 200000 {
			num = 100
			fmt.Println("input num invalid,", err)
		}
		for i := 0; i < num; i++ {
			s.newMarker()
		}
	})
	s.AddUINode(btn)

	tbw := ui.NewTextBox(g.NewVector(fromx, 100), g.NewVector(70, 40), "addwatch:", ui.AliVertical_Mid, ui.AliHorizontal_Right, false)
	tbw.SetColor(color.RGBA{128, 128, 128, 0})
	s.AddUINode(tbw)
	infw := ui.NewInputField(g.NewVector(screenW-50, 100), g.NewVector(50, 40), "1", ui.AliVertical_Mid, ui.AliHorizontal_Right, false)
	infw.SetColor(color.RGBA{128, 128, 128, 128})
	s.AddUINode(infw)
	btnW := ui.NewButton(g.NewVector(fromx, 150), g.NewVector(120, 40), "AddWatcher")
	btnW.SetColor(color.RGBA{100, 100, 100, 128})
	btnW.SetOnPressed(func(b *ui.Button) {
		num, err := strconv.Atoi(infw.ValueText)
		if num <= 0 || num > 200000 {
			num = 1
			fmt.Println("input num invalid,", err)
		}
		for i := 0; i < num; i++ {
			s.newWatcher()
		}
	})
	s.AddUINode(btnW)

	s.notice = ui.NewTextBox(g.NewVector(fromx, 200), g.NewVector(120, 50), "", ui.AliVertical_Mid, ui.AliHorizontal_Left, false)
	s.notice.SetColor(color.RGBA{128, 128, 128, 128})
	s.AddUINode(s.notice)

	return nil
}

//	func (s *ToolScane) OnEnterWatcherView(watcherId uint64, makerId uint64, stype string) {
//		//fmt.Println("OnEnterWatcherView:",watcherId,makerId)
//		m := s.makers[makerId]
//		m.SetInView(true)
//	}
//
//	func (s *ToolScane) OnLeaveWatcherView(watcherId uint64, makerId uint64, stype string) {
//		//fmt.Println("OnLeaveWatcherView:",watcherId,makerId)
//		m := s.makers[makerId]
//		m.SetInView(false)
//	}
func (s *ToolScane) Init() error {
	s.watchers = make(map[uint64]*Watcher)
	s.makers = make(map[uint64]*Marker)

	config := &core.Config{
		M: &core.MapConfig{
			&core.Rectangle{
				screenW,
				screenH,
			},
		},
		T: &core.TowerConfig{
			&core.Rectangle{
				100,
				100,
			},
		},
		Limit: 1000,
	}
	s.aoiScene = toweraoi.NewAOIScene(config)

	w := NewWatcher(s, s.GenWID(), 1, g.NewVector(400, 300))
	w.SetControl(true)
	s.addWatcher(w)

	for i := 0; i < 100; i++ {
		mk := s.newMarker()
		_ = mk
		//mk.SetInView(true)
	}

	s.initBg()
	s.initUI()
	return nil
}

func (s *ToolScane) GenWID() uint64 {
	s.idgen++
	return s.idgen + 10000
}
func (s *ToolScane) GenMID() uint64 {
	s.idgen++
	return s.idgen
}

func (s *ToolScane) addWatcher(w *Watcher) {
	s.watchers[w.GetID()] = w
	s.aoiScene.AddWatcher(w.GetID(), "view", w, w.GetPos(), 1, 0)
}

func (s *ToolScane) addMarker(m *Marker) {
	s.makers[m.GetID()] = m
	s.aoiScene.AddMarker(m.GetID(), m, m.GetPos(), 0)
}

func (s *ToolScane) newMarker() *Marker {
	fx := rand.Float64() * screenW
	fy := rand.Float64() * screenH
	m := NewMarker(s, s.GenMID(), g.Vector{fx, fy})
	s.addMarker(m)
	return m
}

func (s *ToolScane) newWatcher() {
	fx := rand.Float64() * screenW
	fy := rand.Float64() * screenH
	w := NewWatcher(s, s.GenWID(), 1, g.Vector{fx, fy})
	w.SetControl(false)
	s.addWatcher(w)
}

func (s *ToolScane) Update(detaTime float64) {
	//fmt.Println("dt:",detaTime)

	for _, w := range s.watchers {
		w.Update(detaTime)
	}

	for _, m := range s.makers {
		m.Update(detaTime)
	}

	s.notice.SetText(fmt.Sprintf("maker:%d\nwatcher:%d", len(s.makers), len(s.watchers)))
}
