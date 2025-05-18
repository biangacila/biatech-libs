package utils

import (
	"github.com/biangacila/luvungula-go/global"
	"testing"
)

func TestSearchInArray(t *testing.T) {
	type args struct {
		Name     string
		Age      float64
		Location string
	}
	data := []args{
		{
			Name:     "Biangacila",
			Age:      49,
			Location: "Cape town",
		},
		{
			Name:     "Phineas Nkuna",
			Age:      60,
			Location: "Pretoria",
		},
	}
	out, err := SearchInArray(data, "Nk", args{})
	if err != nil {
		t.Error(err)
		return
	}

	global.DisplayObject("result", out)
}
