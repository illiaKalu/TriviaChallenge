var HOST = location.origin.replace(/^http/, 'ws')
var socket = new WebSocket(HOST + "/ws");
var text;
var button = document.getElementById("sendButton");

button.addEventListener("click", sendAnswer);

socket.onopen = function(){
    console.log("Socket opened successfully");
}

socket.onmessage = function(event){
    console.log("Message comes !");

    box = document.createElement("div");
    var jsonData = JSON.parse(event.data);


    if (jsonData.Msg != undefined) {
        box.innerHTML = jsonData.Msg;
    }

    if (jsonData.Question_context != undefined) {
        box.innerHTML = "Question is: ".bold() + jsonData.Question_context;
    }

    document.getElementById("box").appendChild(box);

    // TODO bad code
    $('#box').scrollTop(Number.MAX_SAFE_INTEGER);
}

window.onbeforeunload = function(event){
    console.log("Socket CLOSED successfully");
    socket.close();
}

function sendAnswer(event) {

    text = document.getElementById("answerInputField").value;

    if (text != '') {
        console.log('message sent. - ' + text)
        socket.send(text);
    } else {
        //TODO
    }

    $('#answerInputField').val('');

}


$('#answerInputField').keypress(function(e){

    if(e.keyCode==13){
        sendAnswer(e);
    }

});