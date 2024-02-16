package main

func main() {
	app := App{}
	app.Initialize(DbUser, DbPass, DbName)
	app.Run("localhost:9090")
}
