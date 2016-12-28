package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"strings"
	//"os"
	"bufio"
	"flag"
	"os"
)

var homePath string
var buildArray []string
var copyArray []string
var buildName *string
var copyName *string

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	homePath = os.Getenv("AUTODHOME")
	if homePath == "" {
		panic("You must set AUTODHOME!")
	}

	buildName = flag.String("b", "build.sh", "build script name")
	copyName=flag.String("c","copy.sh","copy script name")
	port := flag.String("port", "6000", "server port")
	flag.Parse()

	bytes, err := ioutil.ReadFile(homePath + "/" + *buildName)
	if err != nil {
		log.Fatal(err)
	}

	text := string(bytes)
	buildArray = strings.Split(text, "\n")
	log.Printf("%V", buildArray)

	cBytes,err:=ioutil.ReadFile(homePath+"/"+*copyName)
	if err != nil {
		log.Fatal(err)
	}
	cText:=string(cBytes)
	copyArray=strings.Split(cText,"\n")
	log.Printf("%V", copyArray)

	http.HandleFunc("/build", buildHandler)
	http.HandleFunc("/copy",copyHandler)
	http.HandleFunc("/assoc_www",assocWWWHandler)
	log.Println("Server Listening:" + *port)

	log.Fatal(http.ListenAndServe(":"+*port, nil))
}
func assocWWWHandler(w http.ResponseWriter,r *http.Request){
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}

	if r.Method == "POST" {
		log.Println("Method:POST")
	} else if r.Method == "GET" {
		log.Println("Method:GET")
		log.Println(r.Form)
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		w.WriteHeader(403)
		return
	}
	if !strings.Contains(string(body), "\"secret\": \"sksheixks\"") {
		w.WriteHeader(403)
		return
	}

	execCommand("/bin/sh", []string{homePath + "/assoc_www.sh" })

}
func copyHandler(w http.ResponseWriter, r *http.Request){
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}

	if r.Method == "POST" {
		log.Println("Method:POST")
	} else if r.Method == "GET" {
		log.Println("Method:GET")
		log.Println(r.Form)
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		w.WriteHeader(403)
		return
	}
	if strings.Contains(string(body), "\"secret\": \"xxxxxx\"") {
		log.Println("true request!")
	}

	execCommand("/bin/sh", []string{homePath + "/" + *copyName})
}

func buildHandler(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}

	if r.Method == "POST" {
		log.Println("Method:POST")
	} else if r.Method == "GET" {
		log.Println("Method:GET")
		log.Println(r.Form)
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		w.WriteHeader(403)
		return
	}
	if strings.Contains(string(body), "\"secret\": \"xxxxxx\"") {
		log.Println("true request!")
	}

	execCommand("/bin/sh", []string{homePath + "/" + *buildName})

	//for _,v:=range cmdArray{
	//	cmd:=strings.Split(v," ")
	//	command:=cmd[0]
	//	params:=cmd[1:]
	//	execCommand(command,params)
	//}

}

func execCommand(command string, params []string) bool {
	pwd, err := os.Getwd()
	log.Println("pwd:", pwd)
	if strings.HasPrefix(command, "cd") {
		err := os.Chdir(params[0])
		if err != nil {
			log.Println(err)
			return false
		}
		return true
	}

	cmd := exec.Command(command, params...)
	log.Println("cmd:", cmd.Args)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Println(err)
		return false
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Println(err)
		return false
	}

	err = cmd.Start()
	if err != nil {
		log.Println("cmd start err:", err)
		return false
	}

	in := bufio.NewScanner(stdout)
	for in.Scan() {
		log.Println(in.Text())
	}
	if err := in.Err(); err != nil {
		log.Println(err)
	}

	in2 := bufio.NewScanner(stderr)
	for in2.Scan() {
		log.Println(in2.Text())
	}
	if err := in2.Err(); err != nil {
		log.Println(err)
	}

	//outReader, err := ioutil.ReadAll(stdout)
	//errReader, err := ioutil.ReadAll(stderr)
	//if err != nil {
	//	log.Println("err:", err)
	//	return false
	//}
	//
	//if outReader != nil {
	//	log.Println("Stdout:", string(outReader))
	//}
	//if errReader != nil {
	//	log.Println("Stderr:", string(errReader))
	//}

	err = cmd.Wait()
	if err != nil {
		log.Println("stderr:", err)
		return false
	}

	return true

}
