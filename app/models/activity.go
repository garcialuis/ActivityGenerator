package models

import (
	"math"
	"math/rand"
	"time"
)

type Activity struct {
	Config ActivityConfig
}

type ActivityConfig struct {
	intensity   int
	startHr     int
	EndHr       int
	eventHrs    int
	eventType   int
	description string
}

const (
	sleep = iota
	work
	exercise
	relax
)

const (
	lightWork = iota
	moderateWork
	heavyWork
)

func reserveTimes(startHr, totalHrs int) (int, int) {

	if startHr == -1 {
		startHr = RandomNumInRage(0, 23)
	}

	endHr := (startHr + totalHrs) % 24

	return startHr, endHr
}

func (a *ActivityConfig) Configure(startHr int, eventHrs int, eventType int, intensity int, description string) {

	a.startHr, a.EndHr = reserveTimes(startHr, eventHrs)
	a.eventType = eventType
	a.intensity = intensity
	a.eventHrs = eventHrs
	a.description = description
}

func (activity *Activity) GenerateEventIntervals(dayLayout *[288]FiveMinuteActivity) {

	switch activity.Config.eventType {
	case sleep:
		generateSleepEvent(dayLayout, activity.Config.startHr, activity.Config.EndHr, activity.Config.eventHrs)
	case work:
		generateWorkDay(dayLayout, activity.Config.startHr, activity.Config.EndHr, activity.Config.eventHrs, activity.Config.intensity)
	default:
		generateDayEvent(dayLayout, activity.Config.startHr, activity.Config.EndHr, activity.Config.eventHrs, activity.Config.eventType)
	}
}

func generateDayEvent(dayLayout *[288]FiveMinuteActivity, startHr int, endHr int, eventHrs int, eventType int) {

	exercising := false
	caloriesPerStep := 0.35
	description := "relaxing"

	if eventType == exercise {
		exercising = true
		caloriesPerStep = 0.45
	}

	fiveMinIntervals := eventHrs * 12
	startIndex := startHr * 12
	endIndex := startIndex + fiveMinIntervals

	for i := startIndex; i < endIndex; i++ {

		interval := FiveMinuteActivity{}
		intervalStand := exercising
		intervalStand = (RandomNumInRage(0, 1) == 1)
		intervalSteps := RandomNumInRage(5, 20)

		startOfHour := (i%12 == 0)

		if exercising {
			interval.ExerciseMinutes = 5
			intervalSteps = RandomNumInRage(2000, 2500) / 12
			description = "exercising"
		}

		if exercising || intervalStand {
			interval.totalSteps = intervalSteps

			if startOfHour {
				interval.standHrs = 1
			}
		}

		caloriesBurned := float64(intervalSteps) * caloriesPerStep
		interval.caloriesBurned = math.Round(caloriesBurned*100) / 100
		interval.description = description

		index := i % 288
		dayLayout[index] = interval
	}

}

func generateSleepEvent(dayLayout *[288]FiveMinuteActivity, startHr int, endHr int, eventHrs int) {

	fiveMinIntervals := eventHrs * 12
	startIndex := startHr * 12
	endIndex := startIndex + fiveMinIntervals

	interval := FiveMinuteActivity{description: "sleeping"}

	for i := startIndex; i < endIndex; i++ {
		index := i % 288
		dayLayout[index] = interval
	}
}

func generateWorkDay(dayLayout *[288]FiveMinuteActivity, startHr int, endHr int, eventHrs int, intensity int) {

	fiveMinIntervals := eventHrs * 12
	startIndex := startHr * 12
	endIndex := startIndex + fiveMinIntervals

	var (
		caloriesPerStep float64
		standing        bool
		stepRange       [2]int
	)

	switch intensity {
	case lightWork:
		stepRange = [2]int{50, 224}
		caloriesPerStep = 0.40
	case moderateWork:
		stepRange = [2]int{200, 850}
		caloriesPerStep = 0.45
	case heavyWork:
		standing = true
		stepRange = [2]int{900, 2000}
		caloriesPerStep = 0.50
	default:
		stepRange = [2]int{40, 250}
		caloriesPerStep = 0.40
	}

	for i := startIndex; i < endIndex; i++ {

		interval := FiveMinuteActivity{description: "working"}
		intervalStand := (RandomNumInRage(0, 1) == 1)

		startOfHour := (i%12 == 0)

		if standing && intervalStand {
			interval.ExerciseMinutes = 5
		}

		if standing || intervalStand {
			stepsTaken := RandomNumInRage(stepRange[0], stepRange[1]) / 12
			interval.totalSteps = stepsTaken
			caloriesBurned := float64(stepsTaken) * caloriesPerStep
			interval.caloriesBurned = math.Round(caloriesBurned*100) / 100

			if startOfHour {
				interval.standHrs = 1
			}
		}

		index := i % 288
		dayLayout[index] = interval
	}

}

func RandomNumInRage(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min+1) + min
}
