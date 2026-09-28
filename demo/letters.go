package main

type Line struct {
	P1, P2 Vec3
}

type Letter struct {
	Lines []Line
}

func GetLetters() map[string]Letter {
	letters := make(map[string]Letter)

	// Letter R
	letters["R"] = Letter{
		Lines: []Line{
			{Vec3{-1, -1, 0}, Vec3{-1, 1, 0}},       // Vertical left
			{Vec3{-1, 1, 0}, Vec3{0.5, 1, 0}},       // Top
			{Vec3{0.5, 1, 0}, Vec3{1, 0.5, 0}},      // Top-right curve
			{Vec3{1, 0.5, 0}, Vec3{1, 0, 0}},        // Right side of loop
			{Vec3{1, 0, 0}, Vec3{0.5, -0.5, 0}},     // Bottom-right curve
			{Vec3{0.5, -0.5, 0}, Vec3{-1, -0.5, 0}}, // Middle cross
			{Vec3{0, -0.5, 0}, Vec3{1, -1, 0}},      // Leg
		},
	}

	// Letter A
	letters["A"] = Letter{
		Lines: []Line{
			{Vec3{0, 1, 0}, Vec3{-1, -1, 0}},    // Left leg
			{Vec3{0, 1, 0}, Vec3{1, -1, 0}},     // Right leg
			{Vec3{-0.5, 0, 0}, Vec3{0.5, 0, 0}}, // Crossbar
		},
	}

	// Letter Y
	letters["Y"] = Letter{
		Lines: []Line{
			{Vec3{-1, 1, 0}, Vec3{0, 0, 0}}, // Top left branch
			{Vec3{1, 1, 0}, Vec3{0, 0, 0}},  // Top right branch
			{Vec3{0, 0, 0}, Vec3{0, -1, 0}}, // Bottom stem
		},
	}

	return letters
}
