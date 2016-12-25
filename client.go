package main

import (
	"fmt"
	"bufio"
	"strings"
	"net/http"
	"os"
	"log"
	"encoding/json"
)

var FACE = [7]string{"list","join","create","logout", "leave", "help", "message"} //client api
var user string

type Result struct {
	Message string
	Status string
	Body []string
}

func main() {

	// print the welcome message
	fmt.Println("Welcome.")
	fmt.Print("Login with your username: ")

	// listen on the standard input and read username
	reader := bufio.NewReader(os.Stdin)
   	input, err := reader.ReadString('\n')
	if err != nil {
		log.Println(err)
	}
	login(strings.TrimSpace(input))
}

func login(username string) {
	// make http request
	url := "http://localhost:3000/users"
	entry := "username="+username
	
	payload := strings.NewReader(entry)
	req, _ := http.NewRequest("POST", url, payload)
	req.Header.Add("cache-control", "no-cache")
	req.Header.Add("content-type", "application/x-www-form-urlencoded")
	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()

	body := Result{}

	dec := json.NewDecoder(res.Body)
	err := dec.Decode(&body)
	if err != nil {
		log.Println(err)
	}

	if (body.Status == "okay") {
		
		// save username state
		user = username

		// print menu
		helpmenu()
		
		reader := bufio.NewReader(os.Stdin)
		for { 
	    	// read input from stdin
	    	input, err := reader.ReadString('\n')
			if err != nil {
				log.Println(err)
			}

			// split input into list then decode and call appropriate method
			msgClean := strings.TrimSpace(input)
			if len(msgClean) == 0 {
				continue
			}

			content := strings.Fields(msgClean) //TODO: work on double word room names or double anything
			//TODO: fix the error from failed array access
			switch content[0] {
				case "create": createRoom(content[1])
				case "list": listRooms()
				case "join": joinRoom(content[1])
				case "leave": leaveRoom(content[1])
				case "message": messageRoom(content[1], content[2])
				case "logout": logout()
				case "help": helpmenu()
				default: fmt.Println("Unknown request")
				fmt.Print(">>> ")
			}
		}
	} else {
		fmt.Println(body.Message)	
		fmt.Print(">>> ")
	}
}

func createRoom(roomname string) {
	url := "http://localhost:3000/chatrooms"
	entry := "roomname="+roomname

	payload := strings.NewReader(entry)

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("cache-control", "no-cache")
	req.Header.Add("content-type", "application/x-www-form-urlencoded")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body := Result{}

	dec := json.NewDecoder(res.Body)
	err := dec.Decode(&body)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(body.Message)
	fmt.Print(">>> ")
}


func listRooms() {
	
	url := "http://localhost:3000/chatrooms"

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("cache-control", "no-cache")
	req.Header.Add("postman-token", "202dc6c7-4eda-161b-b1a2-b6ecd7e100da")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body := Result{}

	dec := json.NewDecoder(res.Body)
	err := dec.Decode(&body)
	if err != nil {
		log.Println(err)
	}
	fmt.Println("Available rooms:", body.Body)
	fmt.Print(">>> ")
}

func joinRoom(roomname string) {

	url := "http://localhost:3000/users/"+user+"/chatrooms"
	entry := "roomname="+roomname+"&action=join"

	payload := strings.NewReader(entry)

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("cache-control", "no-cache")
	req.Header.Add("content-type", "application/x-www-form-urlencoded")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body := Result{}

	dec := json.NewDecoder(res.Body)
	err := dec.Decode(&body)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(body.Message)
	for i:=0; i<len(body.Body); i++ {
		fmt.Println(body.Body[i])	
	}	
	fmt.Print(">>> ")
}

func leaveRoom(roomname string) {

	url := "http://localhost:3000/users/"+user+"/chatrooms"
	entry := "roomname="+roomname+"&action=leave"

	payload := strings.NewReader(entry)

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("cache-control", "no-cache")
	req.Header.Add("content-type", "application/x-www-form-urlencoded")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body := Result{}

	dec := json.NewDecoder(res.Body)
	err := dec.Decode(&body)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(body.Message)
	fmt.Print(">>> ")

}

func messageRoom(roomname string, message string) {
	url := "http://localhost:3000/users/"+user+"/"+roomname+"/messages"

	entry := "message='"+message+"'"
	payload := strings.NewReader(entry)

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("content-type", "application/x-www-form-urlencoded")
	req.Header.Add("cache-control", "no-cache")
	req.Header.Add("postman-token", "d1c79124-b867-09ab-982e-aa601da0f0f7")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body := Result{}

	dec := json.NewDecoder(res.Body)
	err := dec.Decode(&body)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(body.Body[len(body.Body)-1])
	fmt.Print(">>> ")
}

func logout() {
	
	url := "http://localhost:3000/users"
	entry := "username="+user+"&action=logout"
	
	payload := strings.NewReader(entry)
	req, _ := http.NewRequest("POST", url, payload)
	req.Header.Add("cache-control", "no-cache")
	req.Header.Add("content-type", "application/x-www-form-urlencoded")
	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()

	body := Result{}

	dec := json.NewDecoder(res.Body)
	err := dec.Decode(&body)
	if err != nil {
		log.Println(err)
	}

	fmt.Println(body.Message)	
	os.Exit(0)
}

func helpmenu() {
	fmt.Println("Usage instructions:")
	fmt.Println("create AbC -> creates a chat room and set name to AbC")
	fmt.Println("list -> list the existing rooms")
	fmt.Println("join AbC -> join chatroom AbC")
	fmt.Println("leave AbC -> leave chatroom AbC")
	fmt.Println("message AbC [message] -> send message to chatroom AbC")
	fmt.Println("logout -> disconnect")
	fmt.Println("help -> usage instructions")
	fmt.Print(">>> ")
}
