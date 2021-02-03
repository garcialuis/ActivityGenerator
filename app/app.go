package app

import "github.com/garcialuis/ActivityGenerator/app/controllers"

func Run() {
	// Get map of activity days
	generatorController := controllers.Controller{}
	generatorController.GenerateDayActivities()
	// Get list of active users

	// Upstream activity for each user based on their chosen activity day
}
