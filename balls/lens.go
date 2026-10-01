package main

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type LensEffect struct {
	Pos           rl.Vector2
	Vel           rl.Vector2
	Radius        float32
	Magnification float32
	Shader        rl.Shader
	Target        rl.RenderTexture2D

	// Uniform locations
	posLoc    int32
	radiusLoc int32
	magLoc    int32
	sizeLoc   int32
}

func NewLensEffect(radius float32) *LensEffect {
	shader := rl.LoadShader("", "assets/shaders/lens.fs")
	target := rl.LoadRenderTexture(int32(screenWidth), int32(screenHeight))

	return &LensEffect{
		Pos:           rl.NewVector2(float32(screenWidth)/2, float32(screenHeight)/2),
		Vel:           rl.NewVector2(200.0, 150.0), // Pixels per second
		Radius:        radius,
		Magnification: 1.5,
		Shader:        shader,
		Target:        target,

		posLoc:    rl.GetShaderLocation(shader, "lensPos"),
		radiusLoc: rl.GetShaderLocation(shader, "lensRadius"),
		magLoc:    rl.GetShaderLocation(shader, "magnification"),
		sizeLoc:   rl.GetShaderLocation(shader, "renderSize"),
	}
}

func (l *LensEffect) Update(dt float32, timer float64) {
	// Bounce movement
	l.Pos.X += l.Vel.X * dt
	l.Pos.Y += l.Vel.Y * dt

	if l.Pos.X-l.Radius < 0 || l.Pos.X+l.Radius > float32(screenWidth) {
		l.Vel.X *= -1
	}
	if l.Pos.Y-l.Radius < 0 || l.Pos.Y+l.Radius > float32(screenHeight) {
		l.Vel.Y *= -1
	}

	// Pulse magnification: Base 2.0 + range 2.0
	l.Magnification = 2.0 + float32(math.Abs(math.Sin(timer*2.0)))*2.0
}

func (l *LensEffect) Begin() {
	rl.BeginTextureMode(l.Target)
	rl.ClearBackground(rl.Black)
}

func (l *LensEffect) End() {
	rl.EndTextureMode()

	// Update shader uniforms
	// Note: Raylib's RenderTexture is flipped vertically in Y,
	// so the shader needs to know the correct position in screen space.
	// However, we draw the texture flipped later, so we pass standard coordinates.
	rl.SetShaderValue(l.Shader, l.posLoc, []float32{l.Pos.X, float32(screenHeight) - l.Pos.Y}, rl.ShaderUniformVec2)
	rl.SetShaderValue(l.Shader, l.radiusLoc, []float32{l.Radius}, rl.ShaderUniformFloat)
	rl.SetShaderValue(l.Shader, l.magLoc, []float32{l.Magnification}, rl.ShaderUniformFloat)
	rl.SetShaderValue(l.Shader, l.sizeLoc, []float32{float32(screenWidth), float32(screenHeight)}, rl.ShaderUniformVec2)

	rl.BeginShaderMode(l.Shader)
	// Draw the texture flipped back to normal
	rl.DrawTextureRec(l.Target.Texture, rl.NewRectangle(0, 0, float32(l.Target.Texture.Width), float32(-l.Target.Texture.Height)), rl.NewVector2(0, 0), rl.White)
	rl.EndShaderMode()
}

func (l *LensEffect) Unload() {
	rl.UnloadShader(l.Shader)
	rl.UnloadRenderTexture(l.Target)
}
