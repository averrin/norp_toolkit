package dicespy

import "html/template"

type ConfigStruct struct {
	HistoryCount int `default:"1"`
}

type Roll struct {
	Type  string `json:"type"`
	Rolls []struct {
		Type string `json:"type"`
		Dice int    `json:"dice,omitempty"`
		// Fate bool   `json:"fate,omitempty"`
		Mods struct {
		} `json:"mods,omitempty"`
		Sides   int `json:"sides,omitempty"`
		Results []struct {
			V int `json:"v"`
		} `json:"results,omitempty"`
		Expr string `json:"expr,omitempty"`
		Text string `json:"text,omitempty"`
	} `json:"rolls"`
	ResultType string `json:"resultType"`
	Total      int    `json:"total"`
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
