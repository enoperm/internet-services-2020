package model

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"math"
	"time"
)

type Profile struct {
	gorm.Model

	UserID uint `gorm:"unique;not null"`

	LastSmoke     string `form:"last_smoke" binding:"required"`
	DailyAverage  uint32 `form:"daily_average" binding:"required"`
	SticksPerPack uint32 `form:"sticks_per_pack" binding:"required"`
	PricePerPack  uint32 `form:"price_per_pack" binding:"required"`
	StartYear     uint32 `form:"start_year" binding:"required"`
}

var (
	ErrNegativeSmokeSpan = errors.New("start year may not be lower than time of last stick")
	ErrLastStickInFuture = errors.New("if your last stick happens in the future, it is not your last as of now")
)

func (p Profile) Validate() (errors []error) {
	t, err := time.Parse("2006-01-02", p.LastSmoke)
	if err != nil {
		errors = append(errors, err)
	}

	if t.Year()-int(p.StartYear) < 0 {
		errors = append(errors, ErrNegativeSmokeSpan)
	}

	if time.Since(t) < 0 {
		errors = append(errors, ErrLastStickInFuture)
	}

	return
}

func (p Profile) LastSmokeTime() time.Time {
	t, err := time.Parse("2006-01-02", p.LastSmoke)
	if err != nil {
		panic(err)
	}
	return t
}

func (p Profile) Stats() Stats {
	daysSinceLast := math.Ceil(time.Since(p.LastSmokeTime()).Truncate(24*time.Hour).Hours() / 24)
	sticksNotConsumed := uint32(daysSinceLast) * p.DailyAverage
	packsNotConsumed := sticksNotConsumed / p.SticksPerPack
	startTime, _ := time.Parse("2006", fmt.Sprintf("%d", p.StartYear))
	smokedSpan := p.LastSmokeTime().Sub(startTime)
	smokedSpan = smokedSpan.Round(24 * time.Hour)
	smokedYears := smokedSpan.Hours() / (365 * 24) // TODO: Leap years and seconds are not accounted for

	return Stats{
		DaysSinceLastSmoke: uint32(daysSinceLast),
		SticksNotConsumed:  sticksNotConsumed,
		PacksNotConsumed:   packsNotConsumed,
		MoneySpared:        packsNotConsumed * p.PricePerPack,
		YearsSmoked:        uint32(smokedYears),
	}
}

type Stats struct {
	DaysSinceLastSmoke uint32 `json:"days_since_last_smoke"`
	SticksNotConsumed  uint32 `json:"sticks_not_consumed"`
	PacksNotConsumed   uint32 `json:"packs_not_consumed"`
	MoneySpared        uint32 `json:"money_spared"`
	YearsSmoked        uint32 `json:"years_smoked"`

	RankBelow uint32 `json:"rank_below"`
}

type ProfileWithStats struct {
	Profile
	Stats
}

func (p Profile) WithStats() ProfileWithStats {
	return ProfileWithStats{
		Profile: p,
		Stats:   p.Stats(),
	}
}
