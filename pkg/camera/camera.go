package camera

import (
	"math"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/go-gl/mathgl/mgl64"
)

type Direction string

const (
	Forward  = "FORWARD"
	Backward = "BACKWARD"
	Left     = "LEFT"
	Right    = "RIGHT"
	Up       = "UP"
	Down     = "DOWN"
)

// TODO: Add optionals to receive custom values for position, front, up, fov etc.
func New() *Camera {
	camera := &Camera{
		Position:         mgl32.Vec3{0.0, 0.0, 3.0},
		Front:            mgl32.Vec3{0.0, 0.0, -1.0},
		WorldUp:          mgl32.Vec3{0.0, 1.0, 0.0},
		Yaw:              -90.0,
		Pitch:            0.0,
		MovementSpeed:    10.0,
		MouseSensitivity: 0.1,
		Fov:              45.0,
	}

	camera.updateVectors()

	return camera
}

type Camera struct {
	Position         mgl32.Vec3
	Front            mgl32.Vec3
	Up               mgl32.Vec3
	Right            mgl32.Vec3
	WorldUp          mgl32.Vec3
	Yaw              float64
	Pitch            float64
	MovementSpeed    float64
	MouseSensitivity float64
	Fov              float64
}

// updateVectors calculates the fron vector from the camera's (updated) euler angles.
func (c *Camera) updateVectors() {
	// Normalize the vectors, because their length gets closer to 0 the more you look up or
	// down which results in slower movement.

	// fmt.Printf("before update vectors:\ncamera vectors: front: %v right: %v up: %v\n", c.Front, c.Right, c.Up)

	front := mgl32.Vec3{
		float32(math.Cos(mgl64.DegToRad(c.Yaw)) * math.Cos(mgl64.DegToRad(c.Pitch))),
		float32(math.Sin(mgl64.DegToRad(c.Pitch))),
		float32(math.Sin(mgl64.DegToRad(c.Yaw)) * math.Cos(mgl64.DegToRad(c.Pitch))),
	}
	c.Front = front.Normalize()

	right := c.Front.Cross(c.WorldUp)
	c.Right = right.Normalize()

	up := c.Right.Cross(c.Front)
	c.Up = up.Normalize()

	// fmt.Printf("after update vectores:\ncamera vectors: front: %v right: %v up: %v\n", c.Front, c.Right, c.Up)
}

func (c *Camera) ViewMatrix() mgl32.Mat4 {
	return mgl32.LookAtV(c.Position, c.Position.Add(c.Front), c.Up)
}

func (c *Camera) ProcessKeyboard(direction Direction, deltaTime float64) {
	velocity := c.MovementSpeed * deltaTime

	switch direction {
	case Forward:
		c.Position = c.Position.Add(c.Front.Mul(float32(velocity)))
	case Backward:
		c.Position = c.Position.Sub(c.Front.Mul(float32(velocity)))
	case Left:
		c.Position = c.Position.Sub(c.Front.Cross(c.Up).Normalize().Mul(float32(velocity)))
	case Right:
		c.Position = c.Position.Add(c.Front.Cross(c.Up).Normalize().Mul(float32(velocity)))
	case Up:
		c.Position = c.Position.Add(c.Up.Mul(float32(velocity)))
	case Down:
		c.Position = c.Position.Sub(c.Up.Mul(float32(velocity)))
	}

	// fmt.Printf("camera position: %v\n", c.Position)
}

func (c *Camera) ProcessMouseMovement(xoffset, yoffset float64, constrainPitch bool) {
	xoffset *= c.MouseSensitivity
	yoffset *= c.MouseSensitivity

	c.Yaw += xoffset
	c.Pitch += yoffset

	if constrainPitch {
		if c.Pitch > 89.0 {
			c.Pitch = 89.0
		}

		if c.Pitch < -89.0 {
			c.Pitch = -89.0
		}
	}

	c.updateVectors()
}

func (c *Camera) ProcessMouseScroll(yoffset float64) {
	c.Fov -= yoffset

	if c.Fov < 1.0 {
		c.Fov = 1.0
	}

	if c.Fov > 45.0 {
		c.Fov = 45.0
	}
}
