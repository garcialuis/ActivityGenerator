package models

type FiveMinuteActivity struct {
	OriginID        uint64  `json:"originId"`
	StandHrs        int     `json:"standHours"`
	TotalSteps      int     `json:"steps"`
	ExerciseMinutes int     `json:"exerciseMinutes"`
	CaloriesBurned  float64 `json:"caloriesBurned"`
	Description     string  `json:"description"`
	Epochtime       int64   `json:"epochtime"`
}
