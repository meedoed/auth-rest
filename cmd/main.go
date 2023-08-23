package main

import "github.com/meedoed/auth-rest/internal/app"

const configDir = "configs"

func main() {
	app.Run(configDir)
}
