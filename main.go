package main

func main() {
	app := App{}
	config := DbParams{
		DbName:     DbName,
		DbUser:     DbUser,
		DbPassword: DbPassword,
	}
	useDb := false
	app.Initialize(useDb, config)
	app.Run("localhost:8090")
}
