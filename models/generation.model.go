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
	gorm.Model

	GenNumber uint
	Games     []Game `json:"-"`
	Status    uint
}

// //Generations hold a group of generations type
// type Generations *[]Generation

//BulkCreateRecords multiple Generation instance at once
func BulkCreateRecords(newGenrationsRecords *[]Generation) (*[]Generation, error) {

	valueStrings := []string{}
	valueArgs := []interface{}{}
	for _, gen := range *newGenrationsRecords {

		valueStrings = append(valueStrings, "(?, ?)")
		valueArgs = append(valueArgs, gen.GenNumber)
		valueArgs = append(valueArgs, gen.Status)
	}

	smt := `INSERT INTO generations(gen_number, status) VALUES %s`

	smt = fmt.Sprintf(smt, strings.Join(valueStrings, ","))
	if err := DB.GetDB().Exec(smt, valueArgs...).Error; err != nil {
		DB.GetDB().Rollback()
		return nil, err
	}

	return newGenrationsRecords, nil
}
