package testutil

import "fmt"

// Setup 测试设定
func Setup() *App {
	// Initialize an in-memory database for full integration testing.
	app := &App{}
	app.Initialize()
	return app
}

// Teardown 测试后的一些清理工作
func Teardown(app *App) {
	// Closing the connection discards the in-memory database.
	err := app.DB.Close()
	if err != nil {
		fmt.Printf("teardown, %+v", err)
	}
}
