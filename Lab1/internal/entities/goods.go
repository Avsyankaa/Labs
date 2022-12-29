package entities

type GoodBrief struct {
	Id   int64  `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
	Cost int64  `json:"cost" db:"cost"`
}

type Good struct {
	BriefInfo   GoodBrief `json:"main_info"`
	Count       int64     `json:"count" db:"count"`
	Description string    `json:"description" db:"description"`
}

func (g *GoodBrief) Valid() bool {
	return g.Cost != 0 && g.Name != ""
}

func (g *Good) Valid() bool {
	return g.Count > 0 && g.Description != "" && g.BriefInfo.Valid()
}
