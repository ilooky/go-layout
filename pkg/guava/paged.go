package guava

type Paged struct {
	Page  string `json:"page"`
	Start string `json:"start"`
	Limit string `json:"limit"`
}

func (p Paged) Pag() int {
	pag := ToInt(p.Page)
	if pag > 0 {
		return pag - 1
	}
	return 0
}

func (p Paged) Lim() int {
	return ToInt(p.Limit)
}
