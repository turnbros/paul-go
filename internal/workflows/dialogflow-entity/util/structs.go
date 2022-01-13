package util

type EntityOP int

const (
	SET EntityOP = iota
	ADD
	REMOVE
)

type EntityRequest struct {
	Operation    EntityOP `json:"request_operation"`
	EntityType   string   `json:"entity_type"`
	EntityValues []string `json:"entity_values"`
}
