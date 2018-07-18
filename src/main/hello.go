package main

import (
	"fmt"
	"os"
	"net/http"
	"strings"
	"io/ioutil"
)

var host string

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
	if q := r.URL.Query(); q != nil && len(q["homepage"]) > 0{
		host = q["homepage"][0]
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


