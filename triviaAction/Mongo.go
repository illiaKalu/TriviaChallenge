package triviaAction

import (
	"gopkg.in/mgo.v2"
)

var skipFactor int
var questionsCount int = -1
var loadedQuestion QuestionStruct
var collection = &mgo.Collection{}

func init() {

	skipFactor = 0

	session, err := mgo.Dial("mongodb://final_masquerade:pass@ds139470.mlab.com:39470/trivia_challenge")

	if err != nil {
		panic(err)
	}

	session.SetSafe(&mgo.Safe{})

	collection = session.DB("trivia_challenge").C("questionare")
	questionsCount, err = collection.Count()

	if err != nil {
		panic(err)
	}

}

func LoadQuestion() QuestionStruct{

	err := collection.Find(nil).Skip(skipFactor).One(&loadedQuestion)

	if err != nil {
		panic(err)
	}

	skipFactor++

	// if it was last question -> start from beginning
	if skipFactor >= questionsCount {
		skipFactor = 0
	}

	return loadedQuestion

}
