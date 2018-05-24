package main

import (

	"github.com/gorilla/mux"
	"net/http"
	"fmt"
	"encoding/json"
	"os/exec"
	"strings"
)




var Commands []string
var Data []string
var port = 8081
var AccessOrigin = "*"

func GetAgentName() string{

	return "Agent1"
}


func GetIndex(w http.ResponseWriter, r *http.Request){
	w.Write([]byte("OK"))
}

func GetData(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Access-Control-Allow-Origin", AccessOrigin)
	json_data, _ := json.Marshal(Data)
	w.Write(json_data)
	// Clean the data array
	Data = []string{}
}

func ExecuteCommand(command string){

	args := strings.Split(command," ")
	if len(args) >= 1{
		out, err := exec.Command(args[0],args[1:]...).Output()
		if err != nil {
			fmt.Println(err)
		}
		Data = append(Data,string(out))
	}else{
		out, err := exec.Command(args[0]).Output()
		if err != nil {
			fmt.Println(err)
		}
		Data = append(Data,string(out))
	}



}

func HandleCommands(commands []string){

	for _, command := range commands{
		fmt.Println(command)
		ExecuteCommand(command)
	}


}

func AddCommand(w http.ResponseWriter, r *http.Request){

	var commands []string
	w.Header().Set("Access-Control-Allow-Origin", AccessOrigin)
	json_commands := r.FormValue("cmd")
	json.Unmarshal([]byte(json_commands),&commands)
	HandleCommands(commands)
	w.Write([]byte("OK"))
}

func main()  {

	listener := mux.NewRouter()
	listener.HandleFunc("/",GetIndex).Methods("GET")
	listener.HandleFunc("/data/",GetData).Methods("GET")
	listener.HandleFunc("/command/",AddCommand).Methods("POST")
	server := &http.Server{
		Addr:fmt.Sprintf("127.0.0.1:%d",port),
		Handler:listener,

	}
	server.ListenAndServe()

}