package	main

import (
	"net/http"
	"log"
	"io/ioutil"
	"strings"
	"os/exec"
	//"os"
	"os"
)

const ConfFile string ="conf"
var homePath string
var cmdArray []string

func main() {
	log.SetFlags(log.Ldate|log.Ltime|log.Lshortfile)

	homePath=os.Getenv("AUTODHOME")
	if homePath==""{
		panic("You must set AUTODHOME!")
	}
	bytes,err:=ioutil.ReadFile(homePath+"/"+ConfFile)
	if err!=nil{
		log.Fatal(err)
	}
	text:=string(bytes)
	cmdArray=strings.Split(text,"\n")
	log.Printf("%V",cmdArray)

	http.HandleFunc("/",defaultHandler)
	log.Println("Server Listening 9090")

	log.Fatal(http.ListenAndServe(":9090",nil))
}

func defaultHandler(w http.ResponseWriter,r *http.Request)  {

	err:=r.ParseForm()
	if err!=nil{
		log.Println(err)
	}


	if r.Method=="POST"{
		log.Println("Method:POST")
	}else if r.Method=="GET"{
		log.Println("Method:GET")
		log.Println(r.Form)
	}

	body,err:=ioutil.ReadAll(r.Body)
	if err!=nil{
		log.Println(err)
		w.WriteHeader(403)
		return
	}
	if strings.Contains(string(body),"\"secret\": \"ohmygod\""){
		log.Println("true request!")
	}

	execCommand("/bin/sh",[]string{homePath+"/conf"})

	//for _,v:=range cmdArray{
	//	cmd:=strings.Split(v," ")
	//	command:=cmd[0]
	//	params:=cmd[1:]
	//	execCommand(command,params)
	//}


}

func execCommand(command string,params []string)bool{
	pwd,err:=os.Getwd()
	log.Println("pwd:",pwd)
	if strings.HasPrefix(command,"cd"){
		err:=os.Chdir(params[0])
		if err!=nil{
			log.Println(err)
			return false
		}
		return true
	}


	cmd:=exec.Command(command,params...)
	log.Println("cmd:",cmd.Args)

	stdout,err:=cmd.StdoutPipe()
	if err!=nil{
		log.Println(err)
		return false
	}
	stderr,err:=cmd.StderrPipe()
	if err!=nil{
		log.Println(err)
		return false
	}

	err=cmd.Start()
	if err!=nil{
		log.Println("cmd start err:",err)
		return false
	}

	outReader,err:=ioutil.ReadAll(stdout)
	errReader,err:=ioutil.ReadAll(stderr)
	if err!=nil{
		log.Println("err:",err)
		return false
	}


	if outReader!=nil{
		log.Println("Stdout:",string(outReader))
	}
	if errReader!=nil{
		log.Println("Stderr:",string(errReader))
	}


	err=cmd.Wait()
	if err!=nil{
		log.Println("stderr:",err)
		return false
	}





	return true

}
