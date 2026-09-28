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

	// Letter P
	letters["P"] = Letter{
		Lines: []Line{
			{P1: Vec3{X: -1, Y: -1, Z: 0}, P2: Vec3{X: -1, Y: 1, Z: 0}},       // Vertical left
			{P1: Vec3{X: -1, Y: 1, Z: 0}, P2: Vec3{X: 0.5, Y: 1, Z: 0}},       // Top
			{P1: Vec3{X: 0.5, Y: 1, Z: 0}, P2: Vec3{X: 1, Y: 0.5, Z: 0}},      // Curve top-right
			{P1: Vec3{X: 1, Y: 0.5, Z: 0}, P2: Vec3{X: 1, Y: 0, Z: 0}},        // Right side
			{P1: Vec3{X: 1, Y: 0, Z: 0}, P2: Vec3{X: 0.5, Y: -0.5, Z: 0}},     // Curve bottom-right
			{P1: Vec3{X: 0.5, Y: -0.5, Z: 0}, P2: Vec3{X: -1, Y: -0.5, Z: 0}}, // Middle cross
		},
	}

	// Letter H
	letters["H"] = Letter{
		Lines: []Line{
			{P1: Vec3{X: -1, Y: -1, Z: 0}, P2: Vec3{X: -1, Y: 1, Z: 0}}, // Left vertical
			{P1: Vec3{X: 1, Y: -1, Z: 0}, P2: Vec3{X: 1, Y: 1, Z: 0}},   // Right vertical
			{P1: Vec3{X: -1, Y: 0, Z: 0}, P2: Vec3{X: 1, Y: 0, Z: 0}},   // Crossbar
		},
	}

	// Letter O
	letters["O"] = Letter{
		Lines: []Line{
			{P1: Vec3{X: -1, Y: 0.5, Z: 0}, P2: Vec3{X: -0.5, Y: 1, Z: 0}},   // Top-left corner
			{P1: Vec3{X: -0.5, Y: 1, Z: 0}, P2: Vec3{X: 0.5, Y: 1, Z: 0}},    // Top
			{P1: Vec3{X: 0.5, Y: 1, Z: 0}, P2: Vec3{X: 1, Y: 0.5, Z: 0}},     // Top-right corner
			{P1: Vec3{X: 1, Y: 0.5, Z: 0}, P2: Vec3{X: 1, Y: -0.5, Z: 0}},    // Right side
			{P1: Vec3{X: 1, Y: -0.5, Z: 0}, P2: Vec3{X: 0.5, Y: -1, Z: 0}},   // Bottom-right corner
			{P1: Vec3{X: 0.5, Y: -1, Z: 0}, P2: Vec3{X: -0.5, Y: -1, Z: 0}},  // Bottom
			{P1: Vec3{X: -0.5, Y: -1, Z: 0}, P2: Vec3{X: -1, Y: -0.5, Z: 0}}, // Bottom-left corner
			{P1: Vec3{X: -1, Y: -0.5, Z: 0}, P2: Vec3{X: -1, Y: 0.5, Z: 0}},  // Left side
		},
	}

	// Letter N
	letters["N"] = Letter{
		Lines: []Line{
			{P1: Vec3{X: -1, Y: -1, Z: 0}, P2: Vec3{X: -1, Y: 1, Z: 0}}, // Left vertical
			{P1: Vec3{X: -1, Y: 1, Z: 0}, P2: Vec3{X: 1, Y: -1, Z: 0}},  // Diagonal
			{P1: Vec3{X: 1, Y: -1, Z: 0}, P2: Vec3{X: 1, Y: 1, Z: 0}},   // Right vertical
		},
	}

	// Letter A
	letters["A"] = Letter{
		Lines: []Line{
			{P1: Vec3{X: 0, Y: 1, Z: 0}, P2: Vec3{X: -1, Y: -1, Z: 0}},    // Left leg
			{P1: Vec3{X: 0, Y: 1, Z: 0}, P2: Vec3{X: 1, Y: -1, Z: 0}},     // Right leg
			{P1: Vec3{X: -0.5, Y: 0, Z: 0}, P2: Vec3{X: 0.5, Y: 0, Z: 0}}, // Crossbar
		},
	}

	// Letter X
	letters["X"] = Letter{
		Lines: []Line{
			{P1: Vec3{X: -1, Y: 1, Z: 0}, P2: Vec3{X: 1, Y: -1, Z: 0}}, // Backslash
			{P1: Vec3{X: 1, Y: 1, Z: 0}, P2: Vec3{X: -1, Y: -1, Z: 0}}, // Forwardslash
		},
	}

	return letters
}

func GetWord() Letter {
	letters := GetLetters()
	word := Letter{}

	// Spacing between letters
	spacing := 2.5
	chars := []string{"P", "H", "O", "N", "A", "X"}

	// Center the word (approx -2.5 * length / 2)
	startPos := -float64(len(chars)-1) * spacing / 2.0

	for i, char := range chars {
		offset := startPos + float64(i)*spacing
		for _, line := range letters[char].Lines {
			word.Lines = append(word.Lines, Line{
				P1: Vec3{X: line.P1.X + offset, Y: line.P1.Y, Z: line.P1.Z},
				P2: Vec3{X: line.P2.X + offset, Y: line.P2.Y, Z: line.P2.Z},
			})
		}

		// Add connecting lines between letters
		if i < len(chars)-1 {
			word.Lines = append(word.Lines, Line{
				P1:          Vec3{X: offset + 0.8, Y: 0, Z: 0},
				P2:          Vec3{X: offset + spacing - 0.8, Y: 0, Z: 0},
				IsConnector: true,
			})
		}
	}

	return word
}
