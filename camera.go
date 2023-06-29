package glHelper

import (
	"github.com/go-gl/mathgl/mgl32"
	"math"
)

type Camera struct {
	pos     mgl32.Vec3
	front   mgl32.Vec3
	up      mgl32.Vec3
	right   mgl32.Vec3
	worldUp mgl32.Vec3

	yaw           float32
	pitch         float32
	movementSpeed float32
	mouseSens     float32
	zoom          float32
}

type Direction int

const (
	Forward Direction = iota
	Backward
	Left
	Right
	Upward
	Downward
	Nowhere
)

func NewCamera(position, worldUp mgl32.Vec3, yaw, pitch, speed, sens float32) Camera {
	camera := Camera{}
	camera.pos = position
	camera.worldUp = worldUp
	camera.yaw = yaw
	camera.pitch = pitch
	camera.movementSpeed = speed
	camera.mouseSens = sens
	camera.updateVectors()
	return camera
}

func (Camera *Camera) updateVectors() {
	front := mgl32.Vec3{
		float32(math.Cos(float64(mgl32.DegToRad(Camera.yaw))) * math.Cos(float64(mgl32.DegToRad(Camera.pitch)))),
		float32(math.Sin(float64(mgl32.DegToRad(Camera.pitch)))),
		float32(math.Sin(float64(mgl32.DegToRad(Camera.yaw))) * math.Cos(float64(mgl32.DegToRad(Camera.pitch)))),
	}
	Camera.front = front.Normalize()
	Camera.right = Camera.front.Cross(Camera.worldUp).Normalize()
	Camera.up = Camera.right.Cross(Camera.front).Normalize()
}

func (Camera *Camera) GetViewMatrix() mgl32.Mat4 {
	center := Camera.pos.Add(Camera.front)
	return mgl32.LookAt(Camera.pos.X(), Camera.pos.Y(), Camera.pos.Z(), center.X(), center.Y(), center.Z(), Camera.up.X(), Camera.up.Y(), Camera.up.Z())
}

// GetCameraPosition function returns the current camera position as a mgl32.vec3
func (Camera *Camera) GetCameraPosition() mgl32.Vec3 {
	return Camera.pos
}

func (Camera *Camera) UpdateCamera(direction Direction, deltaT, xOffset, yOffest float32) {
	magnitude := Camera.movementSpeed * deltaT
	switch direction {
	case Forward:
		Camera.pos = Camera.pos.Add(Camera.front.Mul(magnitude))
	case Backward:
		Camera.pos = Camera.pos.Sub(Camera.front.Mul(magnitude))
	case Left:
		Camera.pos = Camera.pos.Sub(Camera.right.Mul(magnitude))
	case Right:
		Camera.pos = Camera.pos.Add(Camera.right.Mul(magnitude))
	case Upward:
		Camera.pos = Camera.pos.Add(Camera.up.Mul(magnitude))
	case Downward:
		Camera.pos = Camera.pos.Sub(Camera.up.Mul(magnitude))
	case Nowhere:
	}
	xOffset *= Camera.mouseSens
	yOffest *= Camera.mouseSens

	Camera.yaw += xOffset
	Camera.pitch += yOffest

	//if Camera.pitch > 90.0 {
	//	Camera.pitch = 90.0
	//}
	//if Camera.pitch < -90.0 {
	//	Camera.pitch = -90.0
	//}

	Camera.updateVectors()
}
