package dicespy

import "html/template"
import "fmt"

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

type Template struct {
	templates *template.Template
}

func getTestRoll() *Roll {
	return &Roll{
		Type:       "V",
		ResultType: "sum",
		Total:      5,
		Player:     "NoRP Toolkit",
		Avatar:     fmt.Sprintf("%v/users/avatar/267336/30", avatarRoot),
		Skill:      "Roll for Dice Rolling",
		Mod:        "+3",
		OrigRoll:   "4df+3 Roll for Dice Rolling",
		Results: []struct {
			V int `json:"v"`
		}{{V: 1}, {V: 1}, {V: 0}, {V: -1}},
		Message: "TODO",
		Rolls: []RollResult{
			RollResult{
				Type:  "R",
				Dice:  4,
				Sides: 3,
				Fate:  true,
				Results: []struct {
					V int `json:"v"`
				}{{V: 1}, {V: 1}, {V: 0}, {V: -1}},
			},
			RollResult{
				Type: "M",
				Expr: "+3",
			},
			RollResult{
				Type: "C",
				Expr: "Roll for Dice Rolling",
			},
		},
	}

}
