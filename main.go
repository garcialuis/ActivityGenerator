package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

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

type Activity struct {
	config ActivityConfig
}

type ActivityConfig struct {
	intensity   int
	startHr     int
	endHr       int
	eventHrs    int
	eventType   int
	description string
}

type FiveMinuteActivity struct {
	standHrs        int
	totalSteps      int
	ExerciseMinutes int
	caloriesBurned  float64
	description     string
}

func reserveTimes(startHr, totalHrs int) (int, int) {

	if startHr == -1 {
		startHr = getRandomNumInRage(0, 23)
	}

	endHr := (startHr + totalHrs) % 24

	return startHr, endHr
}

func (a *ActivityConfig) Configure(startHr int, eventHrs int, eventType int, intensity int, description string) {

	a.startHr, a.endHr = reserveTimes(startHr, eventHrs)
	a.eventType = eventType
	a.intensity = intensity
	a.eventHrs = eventHrs
	a.description = description
}

func (activity *Activity) GenerateEventIntervals(dayLayout *[288]FiveMinuteActivity) {

	switch activity.config.eventType {
	case sleep:
		generateSleepEvent(dayLayout, activity.config.startHr, activity.config.endHr, activity.config.eventHrs)
	case work:
		generateWorkDay(dayLayout, activity.config.startHr, activity.config.endHr, activity.config.eventHrs, activity.config.intensity)
	default:
		generateDayEvent(dayLayout, activity.config.startHr, activity.config.endHr, activity.config.eventHrs, activity.config.eventType)
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
		intervalStand = (getRandomNumInRage(0, 1) == 1)
		intervalSteps := getRandomNumInRage(5, 20)

		startOfHour := (i%12 == 0)

		if exercising {
			interval.ExerciseMinutes = 5
			intervalSteps = getRandomNumInRage(2000, 2500) / 12
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
		intervalStand := (getRandomNumInRage(0, 1) == 1)

		startOfHour := (i%12 == 0)

		if standing && intervalStand {
			interval.ExerciseMinutes = 5
		}

		if standing || intervalStand {
			stepsTaken := getRandomNumInRage(stepRange[0], stepRange[1]) / 12
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

func getRandomNumInRage(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min+1) + min
}

func main() {

	var fiveDaysOfActivity [5][288]FiveMinuteActivity

	for i := 0; i < 5; i++ {
		oneDay := ActivityGenerator()
		fiveDaysOfActivity[i] = oneDay
	}

	randomDay := getRandomNumInRage(0, 4)
	fmt.Println("Total days generated: ", len(fiveDaysOfActivity))
	fmt.Println("Picked day: ", randomDay)

	// Temporary: Displaying results to be inspected
	for i, v := range fiveDaysOfActivity[randomDay] {
		fmt.Printf("%d : %v,\n", i, v)
	}

}

func ActivityGenerator() [288]FiveMinuteActivity {

	var events [288]FiveMinuteActivity

	sleepActivity := Activity{}
	sleepActivity.config.Configure(-1, 8, sleep, -1, "Sleep")
	sleepActivity.GenerateEventIntervals(&events)

	nextStartHr := sleepActivity.config.endHr
	dressedActivity := Activity{}
	dressedActivity.config.Configure(nextStartHr, 2, relax, -1, "Relax")
	dressedActivity.GenerateEventIntervals(&events)

	nextStartHr = dressedActivity.config.endHr
	workActivity := Activity{}
	workActivity.config.Configure(nextStartHr, 8, work, 2, "Work")
	workActivity.GenerateEventIntervals(&events)

	toExercise := getRandomNumInRage(0, 1) == 1

	exerciseTime := 0
	relaxTime := 6

	if toExercise {
		exerciseTime = 2
		relaxTime = 4
	}

	nextStartHr = workActivity.config.endHr
	exerciseActivity := Activity{}
	exerciseActivity.config.Configure(nextStartHr, exerciseTime, exercise, -1, "Exercise")
	exerciseActivity.GenerateEventIntervals(&events)

	nextStartHr = exerciseActivity.config.endHr
	relaxActivity := Activity{}
	relaxActivity.config.Configure(nextStartHr, relaxTime, relax, -1, "Relax")
	relaxActivity.GenerateEventIntervals(&events)

	return events
}
