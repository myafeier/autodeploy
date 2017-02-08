package main

import (
	"bufio"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"sync"
)

var homePath string
var urlRegex = regexp.MustCompile(`^[a-zA-Z\-_]{2,20}$`)
var password string
var syncx sync.Mutex

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	var err error
	homePath = os.Getenv("AUTODHOME")
	if homePath == "" {
		homePath, err = os.Getwd()
		if err != nil {
			panic(err)
		}

	}

	passwd := flag.String("pwd", "", "server passwd")
	port := flag.String("port", "8100", "server port")
	flag.Parse()

	if *passwd == "" {
		panic("passwd of server must be setup!")
	}
	password = *passwd

	http.HandleFunc("/", deployHandler)
	log.Println("Server Listening:" + *port)

	log.Fatal(http.ListenAndServe(":"+*port, nil))
}

func deployHandler(w http.ResponseWriter, r *http.Request) {
	syncx.Lock()
	defer syncx.Unlock()
	queryForm, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		panic(err)
	}
	log.Printf("Client :%v\n", r.RequestURI)
	if len(queryForm["site"]) == 0 || len(queryForm["pwd"]) == 0 {
		w.WriteHeader(404)
		return
	}
	store := queryForm["site"][0]
	pwd := queryForm["pwd"][0]
	if pwd != password {
		w.WriteHeader(403)
		return
	}
	log.Printf("store = %+v\n", store)

	if !urlRegex.MatchString(store) {

		w.WriteHeader(400)
		w.Write([]byte("regex didn`t match"))
		return
	}

	_, err = ioutil.ReadFile(homePath + "/" + store + ".sh")
	if err != nil {
		panic(err)
	}

	execCommand("/bin/sh", []string{homePath + "/" + store + ".sh"})
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
	_, err = cmd.StderrPipe()
	if err != nil {
		log.Println(err)
		return false
	}
	err = cmd.Start()
	if err != nil {
		log.Println("cmd start err:", err)
		panic(err)
	}
	/*
		go func (kill chan bool,cmd *exec.Cmd){
			select {
			case <-kill:
				return
			default:
				cmd.Wait()
			}

		}(kill,cmd)
	*/
	err = cmd.Wait()
	if err != nil {
		log.Println(err)
		return false
	}

	in := bufio.NewScanner(stdout)
	for in.Scan() {
		log.Println(in.Text())
	}
	if err = in.Err(); err != nil {
		log.Println(err)
	}
	return true
}
