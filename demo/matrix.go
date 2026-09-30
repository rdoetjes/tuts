package main

import (
	"math"
)

type Vec3 struct {
	X, Y, Z float64
}

type Vec2 struct {
	X, Y float64
}

// RotateX rotates a vector around the X axis
func (v Vec3) RotateX(angle float64) Vec3 {
	cosA := math.Cos(angle)
	sinA := math.Sin(angle)
	return Vec3{
		v.X,
		v.Y*cosA - v.Z*sinA,
		v.Y*sinA + v.Z*cosA,
	}
}

// RotateY rotates a vector around the Y axis
func (v Vec3) RotateY(angle float64) Vec3 {
	cosA := math.Cos(angle)
	sinA := math.Sin(angle)
	return Vec3{
		v.X*cosA + v.Z*sinA,
		v.Y,
		-v.X*sinA + v.Z*cosA,
	}
}

// RotateZ rotates a vector around the Z axis
func (v Vec3) RotateZ(angle float64) Vec3 {
	cosA := math.Cos(angle)
	sinA := math.Sin(angle)
	return Vec3{
		v.X*cosA - v.Y*sinA,
		v.X*sinA + v.Y*cosA,
		v.Z,
	}
}

// Scale scales a vector
func (v Vec3) Scale(s float64) Vec3 {
	return Vec3{v.X * s, v.Y * s, v.Z * s}
}

// Matrix4x4 represents a 4x4 transformation matrix
type Matrix4x4 [4][4]float64

// Identity returns an identity matrix
func Identity() Matrix4x4 {
	return Matrix4x4{
		{1, 0, 0, 0},
		{0, 1, 0, 0},
		{0, 0, 1, 0},
		{0, 0, 0, 1},
	}
}

// Multiply multiplies two matrices
func (m Matrix4x4) Multiply(other Matrix4x4) Matrix4x4 {
	var result Matrix4x4
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			for k := 0; k < 4; k++ {
				result[i][j] += m[i][k] * other[k][j]
			}
		}
	}
	return result
}

// Transform transforms a Vec3 using the matrix
func (m Matrix4x4) Transform(v Vec3) Vec3 {
	x := v.X*m[0][0] + v.Y*m[0][1] + v.Z*m[0][2] + m[0][3]
	y := v.X*m[1][0] + v.Y*m[1][1] + v.Z*m[1][2] + m[1][3]
	z := v.X*m[2][0] + v.Y*m[2][1] + v.Z*m[2][2] + m[2][3]
	w := v.X*m[3][0] + v.Y*m[3][1] + v.Z*m[3][2] + m[3][3]

	if w != 0 && w != 1 {
		return Vec3{x / w, y / w, z / w}
	}
	return Vec3{x, y, z}
}

// Project projects a 3D point onto a 2D plane with perspective
func (v Vec3) Project(width, height, fov, viewerDistance float64) Vec2 {
	factor := fov / (viewerDistance + v.Z)
	return Vec2{
		v.X*factor + width/2,
		-v.Y*factor + height/2,
	}
}

func transform(v Vec3, rx, ry, rz float64) Vec3 {
	v = v.RotateX(rx)
	v = v.RotateY(ry)
	v = v.RotateZ(rz)
	return v
}
