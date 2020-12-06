package model

import (
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
		Days:          uint32(daysSinceLast),
		TotalSticks:   sticksNotConsumed,
		PackCount:     packsNotConsumed,
		PackRemainder: sticksNotConsumed % p.SticksPerPack,
		TotalSavings:  packsNotConsumed * p.PricePerPack,
		SmokeSpan:     uint32(smokedYears),
	}
}

type Stats struct {
	Days          uint32
	TotalSticks   uint32
	PackCount     uint32
	PackRemainder uint32
	TotalSavings  uint32
	SmokeSpan     uint32
}

