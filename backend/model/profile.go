package model

type Profile struct {
	LastSmoke string `json:"last_smoke"` // TODO
	DailyAverage uint32 `json:"daily_average"`
	SticksPerPack uint32 `json:"sticks_per_pack"`
	PricePerPack uint32 `json:"price_per_pack"`
	StartYear uint32 `json:"start_year"`
}

func (p Profile) Validate() (errors []error) {
	// TODO
	return
}
