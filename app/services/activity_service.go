package services

import (
	"github.com/garcialuis/ActivityGenerator/app/models"
)

const (
	sleep = iota
	work
	exercise
	relax
)

func ActivityGenerator() [288]models.FiveMinuteActivity {

	var events [288]models.FiveMinuteActivity

	sleepActivity := models.Activity{}
	sleepActivity.Config.Configure(-1, 8, sleep, -1, "Sleep")
	sleepActivity.GenerateEventIntervals(&events)

	nextStartHr := sleepActivity.Config.EndHr
	dressedActivity := models.Activity{}
	dressedActivity.Config.Configure(nextStartHr, 2, relax, -1, "Relax")
	dressedActivity.GenerateEventIntervals(&events)

	nextStartHr = dressedActivity.Config.EndHr
	workActivity := models.Activity{}
	workIntensity := models.RandomNumInRage(0, 2)
	workActivity.Config.Configure(nextStartHr, 8, work, workIntensity, "Work")
	workActivity.GenerateEventIntervals(&events)

	toExercise := models.RandomNumInRage(0, 1) == 1

	exerciseTime := 0
	relaxTime := 6

	if toExercise {
		exerciseTime = 2
		relaxTime = 4
	}

	nextStartHr = workActivity.Config.EndHr
	exerciseActivity := models.Activity{}
	exerciseActivity.Config.Configure(nextStartHr, exerciseTime, exercise, -1, "Exercise")
	exerciseActivity.GenerateEventIntervals(&events)

	nextStartHr = exerciseActivity.Config.EndHr
	relaxActivity := models.Activity{}
	relaxActivity.Config.Configure(nextStartHr, relaxTime, relax, -1, "Relax")
	relaxActivity.GenerateEventIntervals(&events)

	return events
}
