package main

import (
	"paper_back/config"
	routers "paper_back/routes"
)

func main() {
  config.Connect()
  defer config.Disconnect()

  r := routers.InitRouter()

  r.Run(":8000")
}