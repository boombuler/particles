package particles

import (
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
	"github.com/EngoEngine/gl"
)

type particleData struct {
	ecs.BasicEntity
	*common.SpaceComponent
	common.RenderComponent
	*Component
	particles  []*particle
	life       float32
	shader     common.Shader
	allDead    bool
	lastCenter engo.Point
	srcEntity  *ecs.BasicEntity
}

type particle struct {
	common.SpaceComponent
	common.RenderComponent
	moveRotation        float32
	Life, LifeRemaining float32
	Velocity            engo.Point
	RotationVelocity    float32
	Frame               int
}

type ParticleSource interface {
	common.BasicFace
	common.SpaceFace
	ParticleFace
}

func (pd *particleData) getBaseShader() common.Shader {
	if pd.Effect.Shader == nil {
		// find the default shader for the drawable
		rc := new(common.RenderComponent)
		rc.Drawable = pd.Effect.Drawable
		pd.Effect.Shader = rc.Shader()
	}
	return pd.Effect.Shader
}

func (pd *particleData) prepare() {
	if baseShader := pd.getBaseShader(); pd.shader != baseShader {
		pd.RenderComponent.SetShader(newParticleShader(baseShader))
		pd.shader = baseShader
	}

	pd.RenderComponent.Drawable = pd
	if maxParticles, count := int(pd.Effect.MaxParticles), len(pd.particles); count != maxParticles {
		if maxParticles < count {
			pd.particles = pd.particles[:maxParticles]
		} else {
			for i := count; i < maxParticles; i++ {
				p := new(particle)
				p.Hidden = true
				p.Drawable = pd.Effect.Drawable
				p.Width = 1
				p.Height = 1
				pd.particles = append(pd.particles, p)
			}
		}
	}
}

func (pd *particleData) Texture() *gl.Texture {
	return pd.Effect.Drawable.Texture()
}

func (pd *particleData) Width() float32 {
	return pd.Effect.Drawable.Width()
}

func (pd *particleData) Height() float32 {
	return pd.Effect.Drawable.Height()
}

func (pd *particleData) View() (float32, float32, float32, float32) {
	return pd.Effect.Drawable.View()
}

func (pd *particleData) Close() {}
