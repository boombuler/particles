package main

import (
	"image/color"
	"math/rand"
	"time"

	"github.com/boombuler/particles"

	"engo.io/ecs"
	"engo.io/engo"
	"engo.io/engo/common"
)

type ParticleScene struct{}

func (*ParticleScene) Type() string { return "particle" }

type partSrc struct {
	ecs.BasicEntity
	common.SpaceComponent
	particles.Component
}

var sampleParticleEffect *particles.Effect = &particles.Effect{
	MaxParticles:        3000,
	ParticleEmitRate:    200,
	MaxLife:             1.6,
	MinLife:             0.15,
	StartSize:           engo.Point{0.5, 0.5},
	EndSize:             engo.Point{1, 1},
	SizeEasing:          particles.OutQuad,
	StartColor:          color.NRGBA{255, 0, 0, 255},
	EndColor:            color.NRGBA{0, 0, 255, 128},
	ColorEasing:         particles.InQuad,
	StartTranslation:    particles.OffsetInRect(-5, -5, 5, 5),
	StartVelocity:       func() engo.Point { return engo.Point{0, 0} },
	Acceleration:        engo.Point{0, 300},
	Drawable:            common.Triangle{},
	MinRotation:         0,
	MaxRotation:         360,
	MinRotationVelocity: -20,
	MaxRotationVelocity: 20,
}

func newSource(x, y float32) *partSrc {
	return &partSrc{
		ecs.NewBasic(),
		common.SpaceComponent{
			Position: engo.Point{x, y},
			Width:    10,
			Height:   10,
		},
		particles.Component{
			Effect: sampleParticleEffect,
		},
	}
}

func (*ParticleScene) Preload() {}

func (*ParticleScene) Setup(u engo.Updater) {
	particleSystem := new(particles.System)

	world, _ := u.(*ecs.World)
	world.AddSystem(new(common.RenderSystem))
	world.AddSystem(particleSystem)

	particleSystem.AddEntity(newSource(200, 600))
	particleSystem.AddEntity(newSource(400, 600))
	particleSystem.AddEntity(newSource(600, 600))
}

func main() {
	rand.Seed(time.Now().Unix())
	engo.Run(engo.RunOptions{
		Title:  "Particle Sample",
		Width:  800,
		Height: 800,
	}, new(ParticleScene))
}
