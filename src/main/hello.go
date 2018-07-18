package main

import (
	"fmt"
	"os"
	"net/http"
	"strings"
	"io/ioutil"
)

var host string
var uriUser string
var uriPass string

func main(){
  // register out handle request method
  fmt.Printf("> Starting up HeisnerServer\n")
  http.HandleFunc("/", handleRequest)
  if err := http.ListenAndServe(":8080", nil); err != nil {
    panic(err)
  }
}

func handleRequest(w http.ResponseWriter, r *http.Request){

	urlSuffix := strings.TrimPrefix(r.URL.Path, "/")
	q := r.URL.Query()
	if q == nil {
		writeMessageToFrontend(w, "Mandatory URI Parameters not found. Provide 'homepage', 'userid' and 'passwd'. Aborting.")
		return
	}
	
	if len(q["homepage"]) > 0 && host == "" {
		host = q["homepage"][0]
		// when switching to a new homepage root we ask for the user id
		if len(q["userid"]) > 0{
			uriUser = q["userid"][0]
		}else{
			writeMessageToFrontend(w, "Mandatory URI Parameter 'userid' not found. Aborting.")
			return
		}
		if len(q["passwd"]) > 0{
			uriPass = q["passwd"][0]
		}else{
			writeMessageToFrontend(w, "Mandatory URI Parameter 'passwd' not found. Aborting.")
			return
		}

		if !checkUserPass(uriUser, uriPass){
			writeMessageToFrontend(w, "Credential verification failed. Ensure you provide the secret values... &gt; Aborting.")
			return
		}
	}else{
		writeMessageToFrontend(w, "Mandatory URI Parameter 'homepage' not found. Aborting.")
		return
	}

	if(strings.Compare("exit", urlSuffix) == 0){
		fmt.Printf("Received Exit Command. Terminating ...\n")
		os.Exit(0)
	}else{
		html, _ := getResponse(host + "/" + urlSuffix)
		fmt.Printf("Delivering requested content for " + urlSuffix + "\n")
		w.Write(html)
	}	
}

func getResponse(urlSuffix string)([]byte, error){
	fmt.Printf("Asked to load: " + urlSuffix + "\n" )
	var urlPrefix = "http://"
	response, error := http.Get(urlPrefix + urlSuffix)
	if error != nil {
		var emptyByte []byte
		return emptyByte, error
	}else{
		body, error := ioutil.ReadAll(response.Body)
		return body, error
	}
}

func checkUserPass(user string, pass string) bool{
	userdata, err := ioutil.ReadFile("/etc/heisner-secret/username")
	if err != nil {
		panic(err)
	}
	passdata, err := ioutil.ReadFile("/etc/heisner-secret/password")
	if err != nil {
		panic(err)
	}
	username := string(userdata)
	password := string(passdata)
	
	if(strings.Compare(username, user) == 0 && strings.Compare(password, pass) == 0){
		fmt.Printf("Valid Username and Password provided\n")	
		return true
	}else{
		fmt.Printf("Provided Username and Password does not match expected value\nUser from Store: " + username + " vs. provided: " + user + "\nPass from Store: " + password + " vs. provided: " + pass)	
		return false
	}
}

func writeMessageToFrontend(w http.ResponseWriter, message string)(){
	w.Write([]byte("<html><head></head><body>" + message + "</body></html>"))
}

