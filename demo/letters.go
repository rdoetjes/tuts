package main

type Line struct {
	P1, P2      Vec3
	IsConnector bool
}

type Letter struct {
	Lines []Line
}

func GetLetters() map[string]Letter {
	letters := make(map[string]Letter)

	// Letter R
	letters["R"] = Letter{
		Lines: []Line{
			{P1: Vec3{-1, -1, 0}, P2: Vec3{-1, 1, 0}},       // Vertical left
			{P1: Vec3{-1, 1, 0}, P2: Vec3{0.5, 1, 0}},       // Top
			{P1: Vec3{0.5, 1, 0}, P2: Vec3{1, 0.5, 0}},      // Top-right curve
			{P1: Vec3{1, 0.5, 0}, P2: Vec3{1, 0, 0}},        // Right side of loop
			{P1: Vec3{1, 0, 0}, P2: Vec3{0.5, -0.5, 0}},     // Bottom-right curve
			{P1: Vec3{0.5, -0.5, 0}, P2: Vec3{-1, -0.5, 0}}, // Middle cross
			{P1: Vec3{0, -0.5, 0}, P2: Vec3{1, -1, 0}},      // Leg
		},
	}

	// Letter A
	letters["A"] = Letter{
		Lines: []Line{
			{P1: Vec3{0, 1, 0}, P2: Vec3{-1, -1, 0}},    // Left leg
			{P1: Vec3{0, 1, 0}, P2: Vec3{1, -1, 0}},     // Right leg
			{P1: Vec3{-0.5, 0, 0}, P2: Vec3{0.5, 0, 0}}, // Crossbar
		},
	}

	// Letter Y
	letters["Y"] = Letter{
		Lines: []Line{
			{P1: Vec3{-1, 1, 0}, P2: Vec3{0, 0, 0}}, // Top left branch
			{P1: Vec3{1, 1, 0}, P2: Vec3{0, 0, 0}},  // Top right branch
			{P1: Vec3{0, 0, 0}, P2: Vec3{0, -1, 0}}, // Bottom stem
		},
	}

	return letters
}

func GetWord() Letter {
	letters := GetLetters()
	word := Letter{}

	// Spacing between letters
	spacing := 2.5

	// Add R
	for _, line := range letters["R"].Lines {
		word.Lines = append(word.Lines, Line{
			P1: Vec3{line.P1.X - spacing, line.P1.Y, line.P1.Z},
			P2: Vec3{line.P2.X - spacing, line.P2.Y, line.P2.Z},
		})
	}

	// Add A
	for _, line := range letters["A"].Lines {
		word.Lines = append(word.Lines, Line{
			P1: Vec3{line.P1.X, line.P1.Y, line.P1.Z},
			P2: Vec3{line.P2.X, line.P2.Y, line.P2.Z},
		})
	}

	// Add Y
	for _, line := range letters["Y"].Lines {
		word.Lines = append(word.Lines, Line{
			P1: Vec3{line.P1.X + spacing, line.P1.Y, line.P1.Z},
			P2: Vec3{line.P2.X + spacing, line.P2.Y, line.P2.Z},
		})
	}

	// Add connecting "structural" lines (Black/Invisible)
	// Connecting R to A
	word.Lines = append(word.Lines, Line{
		P1:          Vec3{-spacing + 1, 0, 0}, // Right side of R
		P2:          Vec3{-0.5, 0, 0},         // Left side of A
		IsConnector: true,
	})

	// Connecting A to Y
	word.Lines = append(word.Lines, Line{
		P1:          Vec3{0.5, 0, 0},         // Right side of A
		P2:          Vec3{spacing - 1, 0, 0}, // Left side of Y
		IsConnector: true,
	})

	return word
}
