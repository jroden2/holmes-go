package main

import (
	"github.com/jroden2/holmes-go/pkg/controllers"
	"github.com/jroden2/holmes-go/pkg/utils"
)

func main() {
	logger := utils.InitLogger()
	controllers.Initialise(logger)
}
