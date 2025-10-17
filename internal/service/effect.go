package sugoroku

type Effect interface {
	Apply(player *Player)
}

type ProfitEffect struct {
	Amount int `json: "amount"`
}

type LossEffect struct {
	Amount int `json: "amount"`
}

type QuizEffect struct {
	QuizID int `json:"quiz_id"`
}

type OverallEffect struct {
	Amount int `json: "amount"`
}

type NeighborEffect struct {
	Amount int `json: "amount"`
}

type RequireEffect struct {
	RequireValue int `json: "require_value"`
}
