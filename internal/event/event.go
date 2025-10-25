package event

type Type string

const (
	MoneyChanged Type = "moneyChanged"
	Gambled      Type = "Gambled"
)

type Event struct {
	Type     Type
	PlayerID string
	Data     map[string]interface{}
}
