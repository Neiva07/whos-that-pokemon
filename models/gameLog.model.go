package models

import (
	u "whos-that-pokemon/utils"

	"github.com/jinzhu/gorm"
)

//GameLog register a game from two users
type GameLog struct {
	gorm.Model

	WinnerScore uint
	LoserScore  uint
	WinnerID    uint
	LoserID     uint
	Winner      User `gorm:"foreignkey:id;association_foreignkey:WinnerID"`
	Loser       User `gorm:"foreignkey:id;association_foreignkey:LoserID"`
}

//validate checks if the data passed is valid
func (gameLog *GameLog) validate() (map[string]interface{}, bool) {

	supposedWinner := &User{}
	supposedLoser := &User{}
	err := DB.GetDB().Table("users").Where("id = ?", gameLog.WinnerID).First(supposedWinner).Error
	if err == gorm.ErrRecordNotFound {
		return u.Message(false, "Wrong winner id. Do not exist in the Database"), false
	}
	if err != nil {
		return u.Message(false, "Database connection error"), false
	}

	err = DB.GetDB().Table("users").Where("id = ?", gameLog.LoserID).First(supposedLoser).Error
	if err == gorm.ErrRecordNotFound {
		return u.Message(false, "Wrong loser id. Do not exist in the Database"), false
	}
	if err != nil {
		return u.Message(false, "Database connection error"), false
	}

	if gameLog.WinnerScore <= gameLog.LoserScore {
		return u.Message(false, "Wrong login. Loser has greater score than winner"), false
	}

	return u.Message(true, "Requirement passed"), true

}

//Create a gamelog right after a game ended
func (gameLog *GameLog) Create() map[string]interface{} {

	if response, ok := gameLog.validate(); !ok {
		return response
	}
	DB.GetDB().Create(gameLog)

	if gameLog.ID <= 0 {

		return u.Message(false, "Something went wrong. Try again later")
	}

	return u.Message(true, "GameLog created successfully")
}
