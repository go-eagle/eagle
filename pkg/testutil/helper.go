package testutil

import "fmt"

func Setup() *App {
	// Initialize an in-memory database for full integration testing.
	app := &App{}
	app.Initialize()
	return app
}

func Teardown(app *App) {
	// Closing the connection discards the in-memory database.
	err := app.DB.Close()
	if err != nil {
		fmt.Errorf("teardown, %+v", err)
	}
}
