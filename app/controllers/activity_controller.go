package controllers

import (
	"fmt"

	"github.com/garcialuis/ActivityGenerator/app/models"
	"github.com/garcialuis/ActivityGenerator/app/services"
)

type Controller struct {
	FiveDaysOfActivity [5][288]models.FiveMinuteActivity
}

func (controller *Controller) GenerateDayActivities() {

	var fiveDaysOfActivity [5][288]models.FiveMinuteActivity

	for i := 0; i < 5; i++ {
		oneDay := services.ActivityGenerator()
		fiveDaysOfActivity[i] = oneDay
	}

	randomDay := models.RandomNumInRage(0, 4)
	fmt.Println("Total days generated: ", len(fiveDaysOfActivity))
	fmt.Println("Picked day: ", randomDay)

	// Temporary: Displaying results to be inspected
	for i, v := range fiveDaysOfActivity[randomDay] {
		fmt.Printf("%d : %v,\n", i, v)
	}

	controller.FiveDaysOfActivity = fiveDaysOfActivity
}
