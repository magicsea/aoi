package main

import (
	"aoi"
	"aoi/toweraoi"
	"fmt"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/examples/resources/images"
	g "github.com/magicsea/gosprite"
	"image"
	"image/color"
	"math/rand"
)

type Watcher struct {
	toweraoi.SimpleWatcher
	root  *g.EmptyNode
	scene *ToolScane

	isControl bool
	moveVel   g.Vector
	aiTimer   int
}

func NewWatcher(scene *ToolScane, id uint64, view int, pos g.Vector) *Watcher {
	w := &Watcher{
		scene: scene,
	}
	w.ID = id
	w.ViewRange = view
	w.WatchType = "npc"
	w.Init(pos)
	return w
}

func (w *Watcher) OnMarkerEnterAOI(watchType string, maker aoi.IAOIMarker) {
	//fmt.Println("w.OnMarkerEnterAOI=>", watchType, maker.GetID(), maker.GetType())
	maker.(*Marker).SetInView(true)
}
func (w *Watcher) OnMarkerLeaveAOI(watchType string, maker aoi.IAOIMarker) {
	//fmt.Println("w.OnMarkerLeaveAOI<=", watchType, maker.GetID(), maker.GetType())
	maker.(*Marker).SetInView(false)
}

func (w *Watcher) Init(pos g.Vector) {
	w.root = g.NewEmptyNode()
	w.root.SetPosition(pos)
	w.root.SetDepth(101)
	w.scene.AddNode(w.root)
	//load ani
	frameAniInfo := map[string]*g.FrameAniData{
		"idle": &g.FrameAniData{
			FrameInterval: 10,
			AniName:       "idle",
			FrameRects: []image.Rectangle{
				image.Rect(0, 0, 32, 32),
				image.Rect(1*32, 0, 2*32, 32),
				image.Rect(2*32, 0, 3*32, 32),
				image.Rect(3*32, 0, 4*32, 32),
			},
		},
		"run": &g.FrameAniData{
			FrameInterval: 10,
			AniName:       "run",
			FrameRects: []image.Rectangle{
				image.Rect(0, 32, 32, 64),
				image.Rect(1*32, 32, 2*32, 64),
				image.Rect(2*32, 32, 3*32, 64),
				image.Rect(3*32, 32, 4*32, 64),
				image.Rect(4*32, 32, 5*32, 64),
				image.Rect(5*32, 32, 6*32, 64),
				image.Rect(6*32, 32, 7*32, 64),
				image.Rect(7*32, 32, 8*32, 64),
			},
		},
		"jump": &g.FrameAniData{
			FrameInterval: 10,
			AniName:       "jump",
			FrameRects: []image.Rectangle{
				image.Rect(0, 64, 32, 96),
				image.Rect(1*32, 64, 2*32, 96),
				image.Rect(2*32, 64, 3*32, 96),
				image.Rect(3*32, 64, 4*32, 96),
			},
		},
	}

	as, _ := g.NewAniSprite(images.Runner_png, frameAniInfo)
	//as.SetScale(g.NewVector(3,3))
	as.SetDepth(1)
	as.SetParent(w.root)
	as.SetLocalPosition(g.VectorZero())
	as.Play("run")

	//txt
	txt := g.NewText(fmt.Sprintf("%v", w.GetID()), 8, color.RGBA{255, 0, 0, 255})
	txt.SetDepth(2)
	txt.SetParent(w.root)
	txt.SetLocalPosition(g.VectorZero())

	w.TransPos()
}

func (w *Watcher) SetControl(b bool) {
	w.isControl = b
}

func (w *Watcher) Update(detaTime float64) {
	w.TransPos()
	var speed float64 = 60 * detaTime
	if !w.isControl {
		//auto move
		w.aiTimer--
		if w.aiTimer < 0 {
			w.RecountAI()
			w.aiTimer = rand.Int()%100 + 200
		}

		pos := w.root.GetPosition().Add(w.moveVel.Mul(speed))
		if pos.X > screenW || pos.X < 0 {
			w.RecountAI()
			return
		}

		if pos.Y > screenH || pos.Y < 0 {
			w.RecountAI()
			return
		}
		w.root.SetPosition(pos)

		w.TransPos()
		w.scene.aoiScene.MoveWatcher(w.GetID(), w.GetPos())
		return
	}
	var charX float64 = w.root.GetPosition().X
	var charY float64 = w.root.GetPosition().Y
	if ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyLeft) {
		charX -= speed
	} else if ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyRight) {
		charX += speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyUp) {
		charY -= speed
	} else if ebiten.IsKeyPressed(ebiten.KeyS) || ebiten.IsKeyPressed(ebiten.KeyDown) {
		charY += speed
	}
	w.root.SetPosition(g.NewVector(charX, charY))

	w.TransPos()
	w.scene.aoiScene.MoveWatcher(w.GetID(), w.GetPos())
}

func (w *Watcher) TransPos() {
	w.Pos.X = int(w.root.GetPosition().X)
	w.Pos.Y = int(w.root.GetPosition().Y)
	//fmt.Println("*wathcer:",w.ID," pos:",w.Pos)
}

func (w *Watcher) RecountAI() {
	w.moveVel = g.NewVector(rand.Float64()*2-1, rand.Float64()*2-1).Normal()
}
