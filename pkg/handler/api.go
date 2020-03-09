package handler

type user struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Balance int    `json:"balance"`
}

type tournamentData struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	Deposit int      `json:"deposit"`
	Prize   int      `json:"prize"`
	Users   []string `json:"users"`
	Winner  string   `json:"winner"`
}
