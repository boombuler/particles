package particles

import (
	"image/color"
	"math/rand"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
)

type System struct {
	w     *ecs.World
	rs    *common.RenderSystem
	items map[uint64]*particleData
	mat   *engo.Matrix
}

// New initializes the System
func (ps *System) New(w *ecs.World) {
	ps.w = w
	ps.mat = engo.IdentityMatrix()
	for _, s := range w.Systems() {
		if rendersystem, ok := s.(*common.RenderSystem); ok {
			ps.rs = rendersystem
		}
	}
}

// Update emits and updates the particles.
func (ps *System) Update(dt float32) {
	if len(ps.items) == 0 {
		return
	}
	for _, data := range ps.items {
		data.prepare()

		//number of new particles to spawn
		previousLife := data.life
		data.life = data.life + dt
		previouseSpawnCount := int(previousLife * data.Effect.ParticleEmitRate)
		newSpawnCount := int(data.life * data.Effect.ParticleEmitRate)
		spawnCount := newSpawnCount - previouseSpawnCount
		if !data.DisableSpawning {
			ps.mat.Identity()
			ps.mat.Rotate(data.SpaceComponent.Rotation + 180)
			c := data.Center()
			for i := 0; i < spawnCount; i++ {
				off := float32(i+1) / float32(spawnCount)
				pt := lerpPoint(c, data.lastCenter, off, nil)
				ps.spawnParticle(data, pt)
			}
		} else if data.allDead && data.Component.RemoveWhenDead {
			ps.w.RemoveEntity(*data.srcEntity)
		}
		//update all particles:
		for _, p := range data.particles {
			if !p.RenderComponent.Hidden {
				ps.updateParticle(p, data, dt)
			}
		}
		data.lastCenter = data.Center()
	}
}

func rndBetween(a, b float32) float32 {
	rnd := rand.Float32()
	return a*(1.0-rnd) + b*rnd
}

func lerpPoint(start, end engo.Point, amount float32, easing EasingFn) engo.Point {
	if easing != nil {
		amount = 1 - easing(1-amount)
	}
	end.MultiplyScalar(amount)
	start.MultiplyScalar(1.0 - amount)
	start.Add(end)
	return start
}

func toNRGBA(c color.Color) color.NRGBA {
	if n, ok := c.(color.NRGBA); ok {
		return n
	}
	return color.NRGBAModel.Convert(c).(color.NRGBA)
}

func lerpColor(c1, c2 color.Color, amount float32, easing EasingFn) color.NRGBA {
	if easing != nil {
		amount = 1 - easing(1-amount)
	}

	color1 := toNRGBA(c1)
	color2 := toNRGBA(c2)
	r := int(float32(color1.R)*(1.0-amount)) + int(float32(color2.R)*amount)
	g := int(float32(color1.G)*(1.0-amount)) + int(float32(color2.G)*amount)
	b := int(float32(color1.B)*(1.0-amount)) + int(float32(color2.B)*amount)
	a := int(float32(color1.A)*(1.0-amount)) + int(float32(color2.A)*amount)
	return color.NRGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
}

func (ps *System) spawnParticle(d *particleData, pos engo.Point) {
	//get first available inactive particle
	for _, p := range d.particles {
		if p.Hidden {
			//spawn partice
			p.Hidden = false
			p.Scale = d.Effect.StartSize
			p.Life = rndBetween(d.Effect.MaxLife, d.Effect.MinLife)
			p.Rotation = rndBetween(d.Effect.MaxRotation, d.Effect.MinRotation)
			p.LifeRemaining = p.Life
			trans := d.Effect.StartTranslation()
			trans.MultiplyMatrixVector(ps.mat)
			pos.Add(trans)

			p.SpaceComponent.Width = d.Effect.StartSize.X * p.RenderComponent.Drawable.Width()
			p.SpaceComponent.Height = d.Effect.StartSize.Y * p.RenderComponent.Drawable.Height()
			p.SetCenter(pos)

			p.Velocity = d.Effect.StartVelocity()
			p.RotationVelocity = rndBetween(d.Effect.MaxRotationVelocity, d.Effect.MinRotationVelocity)
			p.moveRotation = d.SpaceComponent.Rotation
			break
		}
	}
}

func (ps *System) updateParticle(p *particle, d *particleData, dt float32) {
	p.LifeRemaining = p.LifeRemaining - dt
	// set valuew based on life remaining
	lifeRatio := p.LifeRemaining / p.Life
	//is particle dead
	if p.LifeRemaining <= 0 {
		p.RenderComponent.Hidden = true
		return
	}

	//set translation
	ps.mat.Identity()
	ps.mat.Rotate(p.moveRotation + 180)

	mov := p.Velocity
	mov.MultiplyScalar(dt)
	mov.MultiplyMatrixVector(ps.mat)
	center := p.SpaceComponent.Center()
	center.Add(mov)

	acc := d.Effect.Acceleration
	acc.MultiplyScalar(dt)
	p.Velocity.Add(acc)

	//set orientation / rotation
	p.Rotation = (p.Rotation + (p.RotationVelocity * float32(dt)))
	for p.Rotation > 360 {
		p.Rotation -= 360
	}
	for p.Rotation < 0 {
		p.Rotation += 360
	}

	p.Scale = lerpPoint(d.Effect.EndSize, d.Effect.StartSize, lifeRatio, d.Effect.SizeEasing)

	textWid := p.RenderComponent.Drawable.Width()
	textHei := p.RenderComponent.Drawable.Height()
	if textWid == 0 && textHei == 0 {
		textWid, textHei = d.SpaceComponent.Width, d.SpaceComponent.Height
	}

	p.SpaceComponent.Width = textWid * p.Scale.X
	p.SpaceComponent.Height = textHei * p.Scale.Y
	p.SpaceComponent.SetCenter(center)

	p.Color = lerpColor(d.Effect.EndColor, d.Effect.StartColor, lifeRatio, d.Effect.ColorEasing)
}

// Remove from System interface
func (ps *System) Remove(e ecs.BasicEntity) {
	if ps.items != nil {
		if data, ok := ps.items[e.ID()]; ok {
			ps.rs.Remove(data.BasicEntity)
			delete(ps.items, e.ID())
		}
	}
}

// AddEntity adds an entity to the System. This is much like the AddByInterface but it checks if the type matches.
func (ps *System) AddEntity(i ecs.Identifier) {
	if m, ok := i.(ParticleSource); ok {
		ps.Add(m.GetBasicEntity(), m.GetSpaceComponent(), m.GetParticleComponent())
	}
}

// AddByInterface adds any ParticleSource to the particle system.
func (ps *System) AddByInterface(i ecs.Identifier) {
	o, _ := i.(ParticleSource)
	ps.Add(o.GetBasicEntity(), o.GetSpaceComponent(), o.GetParticleComponent())
}

// Add adds an entity to the particle system.
func (ps *System) Add(basic *ecs.BasicEntity, space *common.SpaceComponent, particle *Component) {
	if ps.items == nil {
		ps.items = make(map[uint64]*particleData)
	}

	data := &particleData{
		BasicEntity:     ecs.NewBasic(),
		SpaceComponent:  space,
		RenderComponent: common.RenderComponent{},
		Component:       particle,
		srcEntity:       basic,
	}
	data.lastCenter = data.Center()
	data.prepare()
	if ps.rs == nil {
		panic("the particle system needs a rendersystem!")
	}
	ps.rs.Add(&data.BasicEntity, &data.RenderComponent, data.SpaceComponent)
	ps.items[basic.ID()] = data
}
