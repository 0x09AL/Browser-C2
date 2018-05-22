package main

import (
	_"github.com/chzyer/readline"
	"github.com/gorilla/mux"
	"net/http"
	"fmt"
	"io/ioutil"
	"log"
	"time"
)

type Agent struct{
	Name string
	LastCallBack int64
	FirstCallBack int64
}

var Agents []Agent

func AddAgent(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	name := vars["agent"]
	currentTime := time.Now().Unix()
	newAgent := true
	for _,agent := range Agents{
		if(agent.Name == name){
			agent.LastCallBack = currentTime
			newAgent = false
		}
	}
	if(newAgent){
		Agents = append(Agents,Agent{Name:name,FirstCallBack:currentTime,LastCallBack:currentTime})
	}

}


func PrintAgents(){

	for i,agent := range Agents{
		// Removes the Agent from array if no callback was received last 30 seconds.
		if(agent.LastCallBack < (time.Now().Unix()-30)){
			Agents = append(Agents[:i], Agents[i+1:]...)
		}

		fmt.Println(agent.Name)
	}
}

func GetIndex(w http.ResponseWriter, r *http.Request){

	b, err := ioutil.ReadFile("static/index.html")
	if err != nil {
		log.Fatal(err)
	}
	w.Write(b)
}
func GetJS(w http.ResponseWriter, r *http.Request){

	b, err := ioutil.ReadFile("static/jquery.js")
	if err != nil {
		log.Fatal(err)
	}
	w.Write(b)
}

func StartHTTPListener(port int)  {

		listener := mux.NewRouter()
		listener.HandleFunc("/",GetIndex).Methods("GET")
		listener.HandleFunc("/jquery.js",GetJS).Methods("GET")
		listener.HandleFunc("/callback/{agent}",AddAgent).Methods("GET")
		server := &http.Server{
			Addr:fmt.Sprintf(":%d",port),
			Handler:listener,
		}
		server.ListenAndServe()
}

func main(){

	go StartHTTPListener(8080)
	for{
		time.Sleep(5 * time.Second)
		fmt.Println("Printing Agents")
		PrintAgents()
	}

}
