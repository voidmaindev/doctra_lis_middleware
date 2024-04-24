package main

import (
	"github.com/voidmaindev/doctra_lis_middleware/inits"
	"github.com/voidmaindev/doctra_lis_middleware/model"
	"github.com/voidmaindev/doctra_lis_middleware/store"
)

func init() {
	inits.LoadEnvVars()
	store.ConnectToDB()
}

func main() {
	store.DB.AutoMigrate(&models.Hardware{})
}