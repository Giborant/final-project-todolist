package main

import (
	api "github.com/Giborant/final-project-todolist/pkg/api"
	db "github.com/Giborant/final-project-todolist/pkg/db"
	srv "github.com/Giborant/final-project-todolist/pkg/server"
)

func main() {
	api.Init()
	db.Init("scheduler.db")
	srv.Run()
}
