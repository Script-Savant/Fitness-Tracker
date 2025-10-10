package utils

import "github.com/gin-contrib/multitemplate"

func SetupTemplates() multitemplate.Renderer {
	renderer := multitemplate.NewRenderer()

	// home
	renderer.AddFromFiles("home", "templates/layout.html", "templates/home.html")

	// auth
	renderer.AddFromFiles("register", "templates/layout.html", "templates/auth/register.html")
	renderer.AddFromFiles("login", "templates/layout.html", "templates/auth/login.html")

	return renderer
}