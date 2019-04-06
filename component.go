package particles

import (
	"image/color"
	"math/rand"

	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
)

// Effect descripes a particle effect and can be reused by multiple particle emitting entities.
type Effect struct {
	MaxParticles                             uint
	ParticleEmitRate                         float32
	MaxLife, MinLife                         float32
	StartSize, EndSize                       engo.Point
	SizeEasing                               EasingFn
	StartColor, EndColor                     color.Color
	ColorEasing                              EasingFn
	StartTranslation                         OffsetFn
	StartVelocity                            OffsetFn
	Acceleration                             engo.Point
	MaxRotation, MinRotation                 float32
	MaxRotationVelocity, MinRotationVelocity float32
	Drawable                                 common.Drawable
	Shader                                   common.Shader
}

// OffsetFn should return an offset for calculating start velocity or start translation
type OffsetFn func() engo.Point

// OffsetInRect returns an OffsetFn which returns a random Point in the given range.
func OffsetInRect(x1, y1, x2, y2 float32) OffsetFn {
	return func() engo.Point {
		r1, r2 := rand.Float32(), rand.Float32()
		return engo.Point{
			X: x1*(1.0-r1) + x2*r1,
			Y: y1*(1.0-r2) + y2*r2,
		}
	}
}

// Component holds the values for an particle emitting entity
type Component struct {
	// DisableSpawning stops the particle source from emitting new particles
	DisableSpawning bool
	// RemoveWhenDead removes the entity from the world, once all particles are dead.
	RemoveWhenDead bool
	// Effect descripes the particle effect.
	Effect *Effect
}

type ParticleFace interface {
	GetParticleComponent() *Component
}

func (pc *Component) GetParticleComponent() *Component {
	return pc
}
