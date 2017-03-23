package triviaAction

import (

	"gopkg.in/mgo.v2"
	"log"

)

var skipFactor int
var questions []QuestionStruct

func init() {
	skipFactor = 1

	session, err := mgo.Dial("mongodb://final_masquerade:pass@ds139470.mlab.com:39470/trivia_challenge")

	if err != nil {
		panic(err)
	}

	session.SetSafe(&mgo.Safe{})

	c := session.DB("trivia_challenge").C("questionare")

	c.Find(nil).All(&questions)
	log.Print(questions[0])


}

func LoadQuestion() QuestionStruct{

	//		question = QuestionStruct{"error question, see logs", "error occured !"}

	if (skipFactor == 1) {
		skipFactor = skipFactor + 1
		return QuestionStruct{"Vita samaya krasivaya ?", "ans"}
	}
	if (skipFactor == 2) {
		skipFactor += 1
		return QuestionStruct{"quest 2 ?", "ans"}
	}

	if (skipFactor == 3) {
		skipFactor = 1
		return QuestionStruct{"quest 3 ?", "ans"}
	}

	return QuestionStruct{"quest LOL ?", "ans"}

}
