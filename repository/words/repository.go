package words

import (
  "time"
//   "log"
//   "errors"
//   "fmt"
//   "bytes"
)


type Word struct {
	WordId int `json:"wordId"`
	Name string `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
  }

  
var Words []Word


func checkWordExistance(words Word) (error) {
	for _, v := range Words {
	  if v.Name == words.Name {
		
	  } 
	}
  
	return nil
  }


func AddWord(wordsData Word) {

	
  }



