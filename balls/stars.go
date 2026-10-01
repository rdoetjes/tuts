package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Star struct {
	x, y, z float32
}

type Starfield struct {
	stars []Star
}

func NewStarfield(count int) *Starfield {
	stars := make([]Star, count)
	for i := range stars {
		stars[i] = Star{
			x: float32(rl.GetRandomValue(-screenWidth, screenWidth)),
			y: float32(rl.GetRandomValue(-screenHeight, screenHeight)),
			z: float32(rl.GetRandomValue(1, screenWidth)),
		}
	}
	return &Starfield{stars: stars}
}

func (s *Starfield) Update(dt float32) {
	speed := float32(480.0)
	for i := range s.stars {
		s.stars[i].z -= speed * dt

		isOffScreen := false
		if s.stars[i].z > 0 {
			sx := (s.stars[i].x/s.stars[i].z)*100 + float32(screenWidth/2)
			sy := (s.stars[i].y/s.stars[i].z)*100 + float32(screenHeight/2)
			if sx < 0 || sx > screenWidth || sy < 0 || sy > screenHeight {
				isOffScreen = true
			}
		}

		if s.stars[i].z <= 0 || isOffScreen {
			s.stars[i].x = float32(rl.GetRandomValue(-screenWidth, screenWidth))
			s.stars[i].y = float32(rl.GetRandomValue(-screenHeight, screenHeight))
			s.stars[i].z = float32(screenWidth)
		}
	}
}

func (s *Starfield) Draw() {
	for _, star := range s.stars {
		sx := (star.x/star.z)*100 + float32(screenWidth/2)
		sy := (star.y/star.z)*100 + float32(screenHeight/2)

		if sx < 0 || sx > screenWidth || sy < 0 || sy > screenHeight {
			continue
		}

		size := (1.0 - star.z/screenWidth) * 4
		alpha := uint8((1.0 - star.z/screenWidth) * 255)
		rl.DrawCircle(int32(sx), int32(sy), size, rl.NewColor(255, 255, 255, alpha))
	}
}
