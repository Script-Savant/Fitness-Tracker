package utils

import "github.com/gin-contrib/multitemplate"

func SetupTemplates() multitemplate.Renderer {
	renderer := multitemplate.NewRenderer()
	base := "templates/layout.html"

	// home
	renderer.AddFromFiles("home", base, "templates/home.html")

	// auth
	renderer.AddFromFiles("register", base, "templates/auth/register.html")
	renderer.AddFromFiles("login", base, "templates/auth/login.html")

	workoutTemp := "templates/workout/"

	// workout
	renderer.AddFromFiles("create-workout", base, workoutTemp+"create.html")
	renderer.AddFromFiles("update-workout", base, workoutTemp+"update.html")

	return renderer
}
