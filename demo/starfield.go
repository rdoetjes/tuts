package main

import (
	"math/rand/v2"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Star struct {
	Pos        rl.Vector2
	Size       uint8
	Brightness uint8
	Speed      uint8
}

type StarField struct {
	Stars []Star
}

func NewStarField(nrStars int) *StarField {
	starfield := new(StarField)

	starfield.Stars = make([]Star, 0, nrStars)

	for i := 0; i < nrStars; i++ {
		p := rl.NewVector2(float32(rand.IntN(rl.GetScreenWidth()+1)), float32(rand.IntN(rl.GetScreenHeight()+1)))
		z := uint8(rand.IntN(2) + 1)
		b := uint8(rand.IntN(5) + 1)
		s := uint8(rand.IntN(5) + 1)
		star := Star{Pos: p, Size: z, Brightness: b, Speed: s}

		starfield.Stars = append(starfield.Stars, star)
	}
	return starfield
}

func (s *StarField) Process() {
	for i := range s.Stars {
		if s.Stars[i].Pos.X > 0 {
			s.Stars[i].Pos.X -= float32(s.Stars[i].Speed)
		} else {
			s.Stars[i].Pos.X = float32(rl.GetScreenWidth())
		}
	}
}

func (s *StarField) Draw() {
	for i := range s.Stars {
		star := s.Stars[i]
		color := rl.NewColor(star.Brightness*80, star.Brightness*80, star.Brightness*80, 255)
		rl.DrawCircleV(star.Pos, float32(star.Size), color)
	}
}
