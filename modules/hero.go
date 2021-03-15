package modules

type Hero struct {
	Name  string `json:"name,omitempty"`
	Detail string
	AttackPower int
	DefensePower int
	Blood int
}