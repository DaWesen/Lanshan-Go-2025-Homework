package main

import (
	"md6/api"
	"md6/dao"
)

func main() {
	dao.Init()
	api.InitRouter()
}
