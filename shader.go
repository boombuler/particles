package particles

import (
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo/common"
)

// particleShader wraps an existing shader to prevent the sorting of thousands of particles within the rendersystem.
type particleShader struct {
	wrappedShader common.Shader
	culling       common.CullingShader
}

func newParticleShader(shader common.Shader) *particleShader {
	if culling, ok := shader.(common.CullingShader); ok {
		return &particleShader{shader, culling}
	}
	return &particleShader{shader, nil}
}

func (ps *particleShader) PrepareCulling() {
	if ps.culling != nil {
		ps.culling.PrepareCulling()
	}
}

func (ps *particleShader) ShouldDraw(*common.RenderComponent, *common.SpaceComponent) bool {
	// Todo: Implement View-Culling.
	return true
}

func (ps *particleShader) Setup(*ecs.World) error {
	return nil
}

func (ps *particleShader) Pre() {
	ps.wrappedShader.Pre()
}

func (ps *particleShader) Draw(rc *common.RenderComponent, sc *common.SpaceComponent) {
	data := rc.Drawable.(*particleData)
	data.allDead = true
	for _, p := range data.particles {
		if !p.Hidden {
			ps.wrappedShader.Draw(&p.RenderComponent, &p.SpaceComponent)
			data.allDead = false
		}
	}
}
func (ps *particleShader) Post() {
	ps.wrappedShader.Post()
}

func (ps *particleShader) SetCamera(cs *common.CameraSystem) {
}
