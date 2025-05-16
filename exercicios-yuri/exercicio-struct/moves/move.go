package moves

import (
	"fmt"
	"strings"
)

type Move struct {
	Name  string
	Power int
	Type  string
}

func (m *Move) CalculateDamage(targetType string) int {
	fireXGrass := strings.EqualFold(m.Type, "Fire") && strings.EqualFold(targetType, "Grass")
	waterXFire := strings.EqualFold(m.Type, "Water") && strings.EqualFold(targetType, "Fire")

	if fireXGrass || waterXFire {
		return m.Power * 2
	}
	return m.Power
}

func (p *Move) Infos() {
	fmt.Printf("Nome: %s\nPoder: %d\nTipo: %s\n", p.Name, p.Power, p.Type)
}
