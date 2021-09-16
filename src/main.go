package main

import (
	"app/config"
	"app/models"
	"app/task"
)

func main() {
	config.NewConfig()
	models.NewMysqlBase()
	task.StartBfService()
}
