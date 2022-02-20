package geometry

import (
	"fmt"
	"math"
)
type Point struct {
	X, Y  float64
	Label string
}

func (p *Point) Stringer() string {
	return fmt.Sprintf("%s: %.3f,%.3f", p.Label, p.X, p.Y)
}

type Path struct {
	Waypoints []Point
	Label string
}

func (p *Path) AziStart() float64 {
	p0 := p.Waypoints[0]
	p1 := p.Waypoints[1]
	// corr:=math.Cos(p0.y)
	return 90.0 - math.Atan2((p1.Y-p0.Y), p1.X-p0.X)*180/math.Pi - 15
}


type ScreenTransform struct {
	Minx, Maxx, Miny, Maxy float64
	Scale, Xc, Yc, W, H   float64
}

func (t *ScreenTransform) Stringer() string {
	return fmt.Sprintf("Transform: scale=%.1f centre(%.3f,%.3f) for screen %.1f x %.1f", t.Scale, t.Xc, t.Yc, t.W, t.H)
}

func (t *ScreenTransform) NewWindowSize(w, h float64) {
	t.W = w
	t.H = h
}

func (t *ScreenTransform) ToWorld(x, y float64) (a, b float64) {
	a = x - t.W/2
	a /= t.Scale
	a += t.Xc
	b = y - t.H/2
	b /= (-t.Scale)
	b += t.Yc
	// println(fmt.Sprintf("%+v",t))
	return a, b
}
func (t *ScreenTransform) ToScreen(x, y float64) (a, b float64) {
	a = x - t.Xc
	a *= t.Scale
	a += t.W / 2
	b = y - t.Yc
	b *= (-t.Scale)
	b += t.H / 2
	// println(fmt.Sprintf("%+v",t))
	return a, b
}
