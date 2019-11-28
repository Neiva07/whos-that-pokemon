package models

const (
	NotPlayed = iota
	Playing
	Played
)

//Generation of pokemon that are being played
type Generation struct {
	GenNumber uint
	Games     []Game
	Status    uint
}
