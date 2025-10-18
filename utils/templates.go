package utils

import (
	"html/template"

	"github.com/gin-contrib/multitemplate"
)

func SetupTemplates() multitemplate.Renderer {
	renderer := multitemplate.NewRenderer()
	base := "templates/layout.html"

	// home
	renderer.AddFromFiles("home", base, "templates/home.html")

	// auth
	renderer.AddFromFiles("register", base, "templates/auth/register.html")
	renderer.AddFromFiles("login", base, "templates/auth/login.html")

	// workout
	workoutTemp := "templates/workout/"
	renderer.AddFromFiles("create-workout", base, workoutTemp+"create.html")
	renderer.AddFromFiles("update-workout", base, workoutTemp+"update.html")

	// metrics

	funcMap := template.FuncMap {
		"add": func(a, b int) int {
			return a + b
		},
	}

	metricsPath := "templates/metrics/"
	renderer.AddFromFilesFuncs("display-metrics", funcMap, base, metricsPath+"display-metrics.html")
	renderer.AddFromFiles("create-metrics", base, metricsPath+"create-metrics.html")

	return renderer
}
