package dicespy

import (
	"fmt"
	"html/template"
	"math/rand"
	"time"
)

type ConfigStruct struct {
	HistoryCount int `default:"1"`
}

type RollResult struct {
	Type string `json:"type"`
	Dice int    `json:"dice,omitempty"`
	Fate bool   `json:"fate,omitempty"`
	Mods struct {
	} `json:"mods,omitempty"`
	Sides   int `json:"sides,omitempty"`
	Results []struct {
		V int `json:"v"`
	} `json:"results,omitempty"`
	Expr string `json:"expr,omitempty"`
	Text string `json:"text,omitempty"`
}

type Roll struct {
	Type       string       `json:"type"`
	Rolls      []RollResult `json:"rolls"`
	ResultType string       `json:"resultType"`
	Total      int          `json:"total"`
	Player     string
	Avatar     string
	OrigRoll   string
	Message    string
	Skill      string
	Mod        string
	Results    []struct {
		V int `json:"v"`
	}
}

type RollWrapper struct {
	P string `json:"p"`
	D struct {
		Content  string `json:"content"`
		Avatar   string `json:"avatar"`
		OrigRoll string `json:"origRoll"`
		Playerid string `json:"playerid"`
		Type     string `json:"type"`
		Who      string `json:"who"`
	} `json:"d"`
}

type MyTemplate struct {
	templates *template.Template
}

func getTestRoll() *Roll {
	rand.Seed(time.Now().Unix())
	v := []int{-1, 0, 1}
	r := []int{
		v[rand.Intn(3)],
		v[rand.Intn(3)],
		v[rand.Intn(3)],
		v[rand.Intn(3)],
	}
	s := 3
	for _, rr := range r {
		s += rr
	}
	return &Roll{
		Type:       "V",
		ResultType: "sum",
		Total:      s,
		Player:     "NoRP Toolkit",
		Avatar:     fmt.Sprintf("%v/users/avatar/267336/200", avatarRoot),
		OrigRoll:   "4df+3 Roll for Dice Rolling",
		Rolls: []RollResult{
			RollResult{
				Type:  "R",
				Dice:  4,
				Sides: 3,
				Fate:  true,
				Results: []struct {
					V int `json:"v"`
				}{{V: r[0]}, {V: r[1]}, {V: r[2]}, {V: r[3]}},
			},
			RollResult{
				Type: "M",
				Expr: "+3",
			},
			RollResult{
				Type: "C",
				Text: "Roll for Dice Rolling",
			},
		},
	}

}
