var HOST = location.origin.replace(/^http/, 'ws')
var socket = new WebSocket(HOST + "/ws");
var answerContext;
var nickName;
var button = document.getElementById("sendButton");

button.addEventListener("click", sendAnswer);

socket.onopen = function(){
    nickName =  $('#nickNameHolder').text();
    console.log("Socket opened successfully ");
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

    answerContext = $('#answerInputField').val() + "|" + nickName;

    if (answerContext.replace('|' + nickName, '') != '' && answerContext.length < 100) {
        console.log('message sent. - ' + answerContext);
        socket.send(answerContext);
    } else {
        //TODO
        console.log('smth wrong with message');
    }

    $('#answerInputField').val('');

}


$('#answerInputField').keypress(function(e){

    if(e.keyCode==13){
        sendAnswer(e);
    }

});