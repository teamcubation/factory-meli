package models

type Move struct {
	Name  string
	Power int
	Type  string
}

func (m *Move) CalculateDamage(targetType string) int {
	dobroDano := (m.Type == "Fogo" && targetType == "Planta") || (m.Type == "Água" && targetType == "Fogo")

	if dobroDano {
		return m.Power * 2
	}

	return m.Power
}
