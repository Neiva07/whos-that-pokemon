package models

import (
	"fmt"
	"strings"

	"github.com/jinzhu/gorm"
)

//GenerationStatus manage the possibilities of status from generations
type GenerationStatus int

const (
	//NotPlayed yet
	NotPlayed GenerationStatus = iota
	//Playing represent the gen current playing
	Playing
	//Played represent the gen already played
	Played
)

//Generation of pokemon that are being played
type Generation struct {
	gorm.Model  `json:"-"`
	Game        Game `gorm:"foreignkey:id" json:"-"`
	GenNumber   uint
	Status      uint
	UserScore   uint
	FriendScore uint
	GameID      uint `json:"-"`
}

// //Generations hold a group of generations type
// type Generations *[]Generation

//BulkCreateRecords multiple Generation instance at once
func BulkCreateRecords(newGenrationsRecords *[]Generation) (*[]Generation, error) {

	valueStrings := []string{}
	valueArgs := []interface{}{}
	for _, gen := range *newGenrationsRecords {

		valueStrings = append(valueStrings, "(?, ?, ?, ?)")
		valueArgs = append(valueArgs, gen.GenNumber)
		valueArgs = append(valueArgs, gen.Status)
		valueArgs = append(valueArgs, gen.FriendScore)
		valueArgs = append(valueArgs, gen.UserScore)
	}

	smt := `INSERT INTO generations(gen_number, status, friend_score, user_score) VALUES %s`

	smt = fmt.Sprintf(smt, strings.Join(valueStrings, ","))
	if err := DB.GetDB().Exec(smt, valueArgs...).Error; err != nil {
		DB.GetDB().Rollback()
		return nil, err
	}

	return newGenrationsRecords, nil
}

//BulkUpdateRecords update the generations while the game is being played
func BulkUpdateRecords(generations *[]Generation) (*[]Generation, error) {

	valueStrings := []string{}
	valueArgs := []interface{}{}
	updatedGens := &[]Generation{}

	for _, gen := range *generations {

		valueStrings = append(valueStrings, "(?, ?, ?, ?)")
		valueArgs = append(valueArgs, gen.ID)
		valueArgs = append(valueArgs, gen.Status)
		valueArgs = append(valueArgs, gen.FriendScore)
		valueArgs = append(valueArgs, gen.UserScore)
	}

	smt := `INSERT INTO generations(id, status, friend_score, user_score) VALUES %s ON DUPLICATE KEY UPDATE status=VALUES(status), friend_score=VALUES(friend_score), user_score=VALUES(user_score)`

	smt = fmt.Sprintf(smt, strings.Join(valueStrings, ","))

	err := DB.GetDB().Exec(smt, valueArgs...).Find(updatedGens).Error

	return updatedGens, err
}
