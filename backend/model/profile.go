package model

import "time"

type Profile struct {
	LastSmoke string `json:"last_smoke"` // TODO
	DailyAverage uint32 `json:"daily_average"`
	SticksPerPack uint32 `json:"sticks_per_pack"`
	PricePerPack uint32 `json:"price_per_pack"`
	StartYear uint32 `json:"start_year"`
}

func (p Profile) Validate() (errors []error) {
	_, err := time.Parse("2006-01-02", p.LastSmoke)
	if err != nil { errors = append(errors, err) }

	return
}
