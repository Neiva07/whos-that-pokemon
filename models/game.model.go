package models

import (
	"log"
	u "whos-that-pokemon/utils"

	"github.com/jinzhu/gorm"
)

//GameStatus show possibilities from a game status
type GameStatus uint

const (
	//Open game started and not finished yet
	Open GameStatus = iota + 1
	//Finished game
	Finished
)

//Game model to serve while the game is being played.
type Game struct {
	gorm.Model
	//UserID start the game
	UserID uint
	//FriendID is who receive the invitation
	FriendID    uint
	UserScore   uint
	FriendScore uint
	Generations []Generation
	Timer       uint       // has to have a range(15,30,45,60) ideally
	User        User       `gorm:"foreignkey:id;association_foreignkey:UserID" json:"-"`
	Friend      User       `gorm:"foreignkey:id;association_foreignkey:FriendID" json:"-"`
	Status      GameStatus `sql:"default:1"`
}

func (game *Game) validate() (map[string]interface{}, bool) {

	supposedFriend := &User{}
	supposedUser := &User{}

	err := DB.GetDB().Table("users").Where("id = ?", game.UserID).First(supposedUser).Error

	if err == gorm.ErrRecordNotFound {
		return u.Message(false, "User not found in the database."), false
	}
	if err != nil {
		return u.Message(false, "Database fail to connect. Try again later."), false
	}
	err = DB.GetDB().Table("users").Where("id = ?", game.FriendID).First(supposedFriend).Error
	if err == gorm.ErrRecordNotFound {
		return u.Message(false, "Friend not found in the database."), false
	}
	if err != nil {
		return u.Message(false, "Database fail to connect. Try again later."), false
	}

	return u.Message(true, "Game valid."), true
}

//Create a game that are being played
func (game *Game) Create() map[string]interface{} {

	if response, ok := game.validate(); !ok {
		return response
	}
	DB.GetDB().Create(game)

	response := u.Message(true, "Game created successfully")

	response["game"] = game
	return response
}

//Find search for a specfic game and return a error if something went wrong
func (game *Game) Find(gameID uint) error {

	err := DB.GetDB().Preload("Generations").First(game, "id = ?", gameID).Error
	return err
}

//AddGenerations adds the relation between a game and a slice of generations
func (game *Game) AddGenerations(generations *[]Generation) error {
	err := DB.GetDB().Model(game).Association("generations").Append(&generations).Error
	return err
}

//ReplaceGenerations update the correct generations data
func (game *Game) ReplaceGenerations(generations *[]Generation) error {
	err := DB.GetDB().Model(game).Association("generations").Replace(generations).Error
	return err
}

//Update specific fields from a game
func (game *Game) Update(gameUpdated *Game) error {

	err := DB.GetDB().Model(game).
		Updates(&Game{UserScore: gameUpdated.UserScore, FriendScore: gameUpdated.FriendScore, Status: gameUpdated.Status}).Error
	if err != nil {
		return err
	}
	for i := range gameUpdated.Generations {
		gameUpdated.Generations[i].ID = game.Generations[i].ID
		log.Println(game.Generations[i].ID, gameUpdated.Generations[i].ID)
	}

	generations, err := BulkUpdateRecords(&gameUpdated.Generations)
	if err != nil {
		return err
	}
	err = game.ReplaceGenerations(generations)
	return err
}
