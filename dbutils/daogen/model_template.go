package daogen

var model_template = `
package model

import (
	"%v/dao"
)

type Model struct {
	%v
}

var model *Model = nil

func ModelInit() {

	model = &Model{}

	%v
}


`
