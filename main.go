package main


func main(){

	go StartHTTPListener(8080)
	go RemoveInactiveAgents()
	StartTerminal()
}
