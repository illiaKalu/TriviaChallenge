package main

import (
	"net/http"
	"log"
)
import (
	trivia "github.com/illiaKalu/GoTriviaChallenge/triviaAction"
//	"os"
)

func homePage(res http.ResponseWriter, req *http.Request) {
	http.ServeFile(res, req, "templates/index.html")
}

func faviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/favicon.ico")
}

func quizHandler(res http.ResponseWriter, req *http.Request) {

	req.ParseForm()

	if nick := req.Form.Get("nickName"); nick == "" {
		nick = "Anonymus"
	}else {

	}

	http.ServeFile(res, req, "templates/quiz.html")
}

func wsPage(res http.ResponseWriter, req *http.Request) {
	trivia.HandleWsPage(res, req, &hub)
}

var hub = trivia.Hub{
	Broadcast:     		make(chan []byte),
	AddClient:     		make(chan *trivia.Client),
	RemoveClient:  		make(chan *trivia.Client),
	Clients:       		make(map[*trivia.Client]bool),
	Hint: 			make(chan bool),
}


func main() {

	//port := os.Getenv("PORT")
	port := "8080"

	//if port == "" {
	//	log.Fatal("$PORT must be set")
	//
	//}

	// JS, css and other static files handling
	ScriptsDirectory := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", ScriptsDirectory))

	go hub.Start()

	// routes
	http.HandleFunc("/ws", wsPage)
	http.HandleFunc("/", homePage)
	http.HandleFunc("/quiz", quizHandler)
	http.HandleFunc("/favicon.ico", faviconHandler)

	log.Println("Listening at port 8080")
	if err := http.ListenAndServe(":" + port, nil); err != nil {
		log.Fatal("Listen And Serve faile ! ", err)
	}
}
