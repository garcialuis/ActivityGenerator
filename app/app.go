package app

import (
	"sync"

	"github.com/garcialuis/ActivityGenerator/app/controllers"
	"github.com/garcialuis/ActivityGenerator/app/services"
)

func Run() {

	var wg sync.WaitGroup

	// Get map of activity days
	generatorController := controllers.Controller{}
	generatorController.GenerateDayActivities()

	// Get list of active users
	originIDs := []int{10080, 10209, 10038}
	wg.Add(1)
	// Upstream activity for each user based on their chosen activity day
	go func() {
		services.PublishActivities(generatorController.FiveDaysOfActivity, originIDs)
		wg.Done()
	}()

	wg.Wait()
}
