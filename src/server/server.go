package server

import (
	"github.com/petruki/gitops-poc/src/services"
)

func Init() {
	dummy := services.Add(1, 2)
	println(dummy)
}
