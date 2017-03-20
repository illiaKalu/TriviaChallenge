package triviaAction

import (
//	"gopkg.in/mgo.v2"
//	"log"
)

var max int = 1

func LoadQuestion() QuestionStruct{

	//session, err := mgo.Dial("localhost")
	//if err != nil {
	//	panic(err)
	//}

	//questions := []QuestionStruct{}
	//
	//c := session.DB("trivia").C("questionare")
	//log.Print(c.Find(nil).All(&questions))

	if (max == 1) {
		max = max + 1
		return QuestionStruct{"Vita samaya krasivaya ?", "ans"}
	}
	if (max == 2) {
		max += 1
		return QuestionStruct{"quest 2 ?", "ans"}
	}

	if (max == 3) {
		max = 1
		return QuestionStruct{"quest 3 ?", "ans"}
	}

	return QuestionStruct{"quest LOL ?", "ans"}

}
