package main

import (
	"engo.io/ecs"
	"engo.io/engo/common"
)

type RotateSystem struct {
	items map[uint64]*common.SpaceComponent
}

func (rs *RotateSystem) Update(dt float32) {
	angle := dt * 360
	for _, itm := range rs.items {
		itm.Rotation = (itm.Rotation + angle)
		for itm.Rotation > 360 {
			itm.Rotation -= 360
		}
		for itm.Rotation < 0 {
			itm.Rotation += 360
		}
	}
}

func (rs *RotateSystem) Add(basic *ecs.BasicEntity, space *common.SpaceComponent) {
	if rs.items == nil {
		rs.items = make(map[uint64]*common.SpaceComponent)
	}
	rs.items[basic.ID()] = space
}

func (rs *RotateSystem) Remove(e ecs.BasicEntity) {
	if rs.items != nil {
		delete(rs.items, e.ID())
	}
}
