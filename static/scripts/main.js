var socket = new WebSocket("ws://localhost:8080/ws");

var button = document.getElementById("button");

button.addEventListener("click", function(event){
    var text = document.getElementById("textbox").value;
    socket.send(text);
});

socket.onopen = function(){
    console.log("Socket opened successfully");
}

socket.onmessage = function(event){
    console.log("Message comes !");

    var box = document.createElement("div");
    var jsonData = JSON.parse(event.data);


    if (jsonData.Msg != undefined) {
        box.innerHTML = jsonData.Msg;
    }

    if (jsonData.Question_context != undefined) {
        box.innerHTML = "<!>Trivia<!>: Question is: " + jsonData.Question_context;
    }

    document.getElementById("box").appendChild(box);
}

window.onbeforeunload = function(event){
    console.log("Socket CLOSED successfully");
    socket.close();
}
