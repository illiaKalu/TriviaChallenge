package triviaAction

import (
	"net/http"
	"github.com/gorilla/websocket"
	"encoding/json"
	"strings"
)
var upgrader = websocket.Upgrader{}
var question QuestionStruct
var questionByteArray []byte

var hintCounter int = 0
var hintString string

type Hub struct {
	Clients map[*Client]bool
	Broadcast     	   chan []byte
	AddClient     	   chan *Client
	RemoveClient  	   chan *Client
	Hint     	   chan bool
}

type QuestionStruct struct {
	Question_context string "question"
	Answer  	 string "answer"
}

type MessageStruct struct {
	Msg string
}

func HandleWsPage (res http.ResponseWriter, req *http.Request, hub *Hub) {

	conn, err := upgrader.Upgrade(res, req, nil)

	if err != nil {
		http.NotFound(res, req)
		return
	}

	client := &Client{
		ws:   conn,
		sendMessageChannel: make(chan []byte),
	}

	hub.AddClient <- client

	go client.Write()
	go client.Read(hub)

}

func getQuestion() {

	if (question == QuestionStruct{}) {
		question = LoadQuestion()
		makeHintString(len(question.Answer))
		stringify, err := json.Marshal(question)
		if err != nil {
			panic(err)
		}
		questionByteArray = stringify
	}
}


func (hub *Hub) Start() {

	go StartTimer(hub)

	for {

		select {

		case conn := <- hub.AddClient: {
			hub.Clients[conn] = true
			getQuestion()
			conn.sendMessageChannel <- questionByteArray
		}

		case conn := <- hub.RemoveClient:
			if _, ok := hub.Clients[conn]; ok {
				delete(hub.Clients, conn)
				close(conn.sendMessageChannel)
			}

		case message := <- hub.Broadcast: {

			if strings.Contains(string(message), question.Answer) {
				// true means that question WAS guessed
				resetQuestion(true, hub, message)
			}else {
				convertAndSend(string(message), hub)
			}
		}

		case <- hub.Hint: {

			// if we didn't have question, that means that user still on login page
			// ignore ticks =/
			if (question != QuestionStruct{}) {
				if (hintCounter + 1 == len(question.Answer)) {
					// false means that question was NOT guessed
					resetQuestion(false, hub, []byte(""))
				} else {
					hintCounter += 1
					modifyAndBroadcastHint(hub)

				}
			}
		}
		}

	}

	}

// reset called in 2 cases:
// 1) question was guessed ( isGuessed = true )
// 2) question was not guessed and hint service opened all letters ( isGuessed = false )
/* anyway, clear hintcounter, hint, reset timer, load nextquestion and broadcast new question
   actions should be done
*/
func resetQuestion(isGuessed bool, hub *Hub, message []byte) {

	if (isGuessed){
		sendCorrectAnswer(message, hub)
	}else {
		broadcastNotGuessed(hub)
	}

	hintCounter = 0
	hintString = ""
	RestartTimer()
	question = QuestionStruct{}
	getQuestion()
	broadcastMessage(hub, questionByteArray)
}

// modify message by adding "correct answer was" mask and broadcast to clients
func broadcastNotGuessed(hub *Hub) {
	answerByteArr := ("The correct answer was :  " + question.Answer)
	convertAndSend(answerByteArr, hub)
}

// every time method called, replace next "*" character to actual letter
func modifyAndBroadcastHint(hub *Hub) {
	hintString = strings.Replace(hintString, string(hintString[hintCounter - 1]), string(question.Answer[hintCounter - 1]), 1)
	convertAndSend( ("TriviaHint: " + hintString), hub)
}

// modify message by adding "correct" mask and broadcast to clients
func sendCorrectAnswer(message []byte, hub *Hub) {
	answerByteArr := []byte(" CORRECT ------ > ")
	message = append(answerByteArr, message ... )
	convertAndSend(string(message), hub)
}

// converting message from message_context|nickname to "nickname: message_context"
// marshal to JSON format
func convertAndSend(message string, hub *Hub) {

	nickIndex := strings.Index(message, "|")

	if nickIndex != -1 {
		message = message[nickIndex + 1:] + ": " + message[:nickIndex]
	}


	if jsonMsg, err := json.Marshal(MessageStruct{string(message)}); err != nil {
		panic(err)
	}else {
		broadcastMessage(hub, jsonMsg)
	}


}

// transform hint to string of "*"
func makeHintString(size int) {
	hintString = strings.Repeat("*", size)
}

func broadcastMessage(hub *Hub, msg []byte) {
	for conn := range hub.Clients {
		conn.sendMessageChannel <- msg
	}
}


