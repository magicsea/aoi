package main

import (
	"aoi/toweraoi"
	"fmt"
	g "github.com/magicsea/gosprite"
	"image/color"
	"math/rand"
)

type Marker struct {
	toweraoi.SimpleMarker
	root  *g.EmptyNode
	scene *ToolScane
	body  *g.Circle

	moveVel g.Vector
	aiTimer int

	inView bool
}

func NewMarker(scene *ToolScane, id uint64, pos g.Vector) *Marker {
	w := &Marker{
		scene: scene,
	}
	w.ID = id
	w.Type = "npc"
	w.Init(pos)
	return w
}

func (m *Marker) Init(pos g.Vector) {
	m.root = g.NewEmptyNode()
	m.root.SetPosition(pos)
	m.root.SetDepth(100)
	m.scene.AddNode(m.root)
	//load ani
	c := color.RGBA{0, 128, 0, 128}
	m.body = g.NewCircle(10, c)
	m.body.SetDepth(1)
	m.body.SetParent(m.root)
	m.body.SetLocalPosition(g.VectorZero())
	//txt
	txt := g.NewText(fmt.Sprintf("%v", m.GetID()), 8, color.RGBA{255, 0, 0, 255})
	txt.SetDepth(2)
	txt.SetParent(m.root)
	txt.SetLocalPosition(g.NewVector(-4, 4))

	m.TransPos()
}

func (m *Marker) TransPos() {
	m.Pos.X = int(m.root.GetPosition().X)
	m.Pos.Y = int(m.root.GetPosition().Y)

	//fmt.Println("maker:",m.ID," pos:",m.Pos)
}

func (m *Marker) Update(detaTime float64) {

	m.TransPos()

	c := color.RGBA{0, 128, 0, 128}
	if m.inView {
		c = color.RGBA{128, 0, 0, 128}
	}
	m.body.SetColor(c)

	m.aiTimer--
	if m.aiTimer < 0 {
		m.RecountAI()
		m.aiTimer = rand.Int()%100 + 200
	}

	pos := m.root.GetPosition().Add(m.moveVel.Mul(detaTime))
	if pos.X > screenW || pos.X < 0 {
		m.RecountAI()
		return
	}

	if pos.Y > screenH || pos.Y < 0 {
		m.RecountAI()
		return
	}

	m.root.SetPosition(pos)

	m.TransPos()

	m.scene.aoiScene.MoveMarker(m.GetID(), m.GetPos())
}

func (m *Marker) RecountAI() {
	var speed float64 = 60
	m.moveVel = g.NewVector(rand.Float64()*2-1, rand.Float64()*2-1).Normal().Mul(speed)
}

func (m *Marker) SetInView(b bool) {
	m.inView = b
}
