package domain

type Impact struct {
	ID          int
	Description string
	Weight      int
	TypeID      int
	Default     bool
}
