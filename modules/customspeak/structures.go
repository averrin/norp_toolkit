package customspeak

type ConfigStruct struct {
	Email string `default:""`
	Password string `default:""`
	Token string  `default:""`
}

type Event struct {
	Username string
	Speak    bool
}
