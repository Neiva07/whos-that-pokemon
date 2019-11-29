package models

import (
	"fmt"
	"strings"
	u "whos-that-pokemon/utils"

	"github.com/jinzhu/gorm"
)

const (
	NotPlayed = iota
	Playing
	Played
)

//Generation of pokemon that are being played
type Generation struct {
	gorm.Model

	GenNumber uint
	Games     []Game
	Status    uint
}

// //Generations hold a group of generations type
// type Generations *[]Generation

//BulkCreateRecords multiple Generation instance at once
func BulkCreateRecords(generations *[]Generation) map[string]interface{} {

	valueStrings := []string{}
	valueArgs := []interface{}{}
	print(*generations)
	for _, gen := range *generations {
		valueStrings = append(valueStrings, "(?, ?)")

		valueArgs = append(valueArgs, gen.GenNumber)
		valueArgs = append(valueArgs, gen.Status)
	}

	smt := `INSERT INTO generations VALUES %s`

	smt = fmt.Sprintf(smt, strings.Join(valueStrings, ","))

	if err := DB.GetDB().Exec(smt, valueArgs...).Error; err != nil {
		DB.GetDB().Rollback()
		return u.Message(false, "Something went wrong in the database.")
	}

	response := u.Message(true, "Generations created successfully")

	response["generations"] = generations

	return response
}
