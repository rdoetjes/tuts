package main

import (
	"math"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	screenWidth   = 1024
	screenHeight  = 768
	starCount     = 400
	heartbeatSize = 60
	fireWidth     = 320
	fireHeight    = 100
)

var (
	heartbeatTable [heartbeatSize]float32
	bgMusic        rl.Music
)

func initHeartbeat() {
	for i := 0; i < heartbeatSize; i++ {
		t := float64(i) / float64(heartbeatSize)
		var val float64

		if t >= 0.1 && t <= 0.2 {
			val = 0.15 * math.Sin((t-0.1)*10*math.Pi)
		} else if t > 0.22 && t <= 0.24 {
			val = -0.2 * math.Sin((t-0.22)*50*math.Pi)
		} else if t > 0.24 && t <= 0.28 {
			val = 1.0 * math.Sin((t-0.24)*25*math.Pi)
		} else if t > 0.28 && t <= 0.30 {
			val = -0.4 * math.Sin((t-0.28)*50*math.Pi)
		} else if t >= 0.4 && t <= 0.6 {
			val = 0.25 * math.Sin((t-0.4)*5*math.Pi)
		}

		heartbeatTable[i] = float32(val)
	}
}

func initAudio() {
	rl.InitAudioDevice()

	// Load Background Music
	bgMusic = rl.LoadMusicStream("assets/music/DumDum.mp3")
	if bgMusic.Stream.Buffer != nil {
		rl.PlayMusicStream(bgMusic)
	}
}

func main() {
	rl.InitWindow(screenWidth, screenHeight, "D'ELITE - 1989 CRACKTRO")
	defer rl.CloseWindow()

	initAudio()
	defer rl.CloseAudioDevice()
	defer rl.UnloadMusicStream(bgMusic)

	font := rl.LoadFont("assets/fonts/Impact.ttf")
	defer rl.UnloadFont(font)

	rl.SetTargetFPS(60)

	initHeartbeat()

	// Initialize Effects
	starfield := NewStarfield(starCount)
	fire := NewFireEffect()
	copperBars := &CopperBars{}
	balls := &OrbitingBalls{}
	scrollText := "-----------------CRU JONES & PHONAX PRESENTS... THE 1989 ULTIMATE CRACKTRO DEMO!    CODED IN GO USING RAYLIB-GO...    GREETINGS TO: FAIRLIGHT - RAZOR 1911 - SKID ROW - GENESIS - TRSI - THE SILENTS - PHENOMENA - ANTHROX - TITAN...    WE BRING YOU THE BEST RELEASES, CRACKED AND PACKED FOR YOUR PLEASURE! ...... AND  REMEMBER: LIVE FAST, DIE YOUNG, AND LEAVE A GOOD LOOKING BODY!!! ..... AND ALWAYS, STAY RAD!!! --------------------------------------------"
	scroller := NewScroller(font, scrollText)
	monitor := &HeartMonitor{}
	header := &Header{font: font}

	phrase := "LEAVE A GOOD LOOKING BODY"
	phraseIndex := strings.Index(scrollText, phrase)
	isFlatline := false

	scrollTimer := float64(0)
	var timer float64 = 0

	for !rl.WindowShouldClose() {
		dt := float64(rl.GetFrameTime())
		timer += dt

		rl.UpdateMusicStream(bgMusic)

		scrollTimer += dt
		pulseIndex := int(timer*60) % heartbeatSize

		if !isFlatline {
			scrollTimer += dt
		}

		scrollPos := float32(screenWidth) - float32(scrollTimer*300.0)

		// Update
		starfield.Update(float32(dt))
		fire.Update(isFlatline)

		// Check Flatline
		phraseX := scrollPos + float32(phraseIndex*50)
		if phraseX < float32(screenWidth/2+200) && phraseX > float32(screenWidth/2-400) {
			isFlatline = true
		}

		if scrollPos < -float32(len(scrollText)*50) {
			scrollTimer = 0
			isFlatline = false // Reset flatline only on scroll wrap
		}

		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)

		var pulse float32 = 0
		if !isFlatline {
			pulse = heartbeatTable[pulseIndex]
		}

		starfield.Draw()
		fire.Draw(false)
		fire.Draw(true)
		copperBars.Draw(timer)
		balls.Draw(timer)
		monitor.Draw(screenWidth/2, screenHeight/2, timer, isFlatline)
		scroller.Draw(timer, scrollPos)
		header.DrawLogo(timer, pulse)
		header.DrawSubHeader(timer, isFlatline)
		DrawBorder()

		rl.EndDrawing()
	}
}
