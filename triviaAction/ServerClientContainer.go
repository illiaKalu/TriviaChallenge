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
	Question_context string
	Answer  	 string
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

			if string(message) == question.Answer {
				// true means that question WAS guessed
				resetQuestion(true, hub, message)
			}else {
				convertAndSend(message, hub)
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

func broadcastNotGuessed(hub *Hub) {
	answerByteArr := []byte(" The correct answer was :  " + question.Answer)
	convertAndSend(answerByteArr, hub)
}
func modifyAndBroadcastHint(hub *Hub) {
	hintString = strings.Replace(hintString, string(hintString[hintCounter - 1]), string(question.Answer[hintCounter - 1]), 1)
	convertAndSend( []byte("TriviaHint: " + hintString), hub)
}

func sendCorrectAnswer(message []byte, hub *Hub) {
	answerByteArr := []byte(" <----- CORRECT ! ")
	message = append(message, answerByteArr...)
	convertAndSend(message, hub)
}

func convertAndSend(message []byte, hub *Hub) {
	jsonMsg, err := json.Marshal(MessageStruct{string(message)})

	if err != nil {
		panic(err)
	}
	broadcastMessage(hub, jsonMsg)
}

func broadcastMessage(hub *Hub, msg []byte) {
	for conn := range hub.Clients {
		conn.sendMessageChannel <- msg
	}
}

func makeHintString(size int) {
	hintString = strings.Repeat("*", size)
}
