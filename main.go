package main

import (
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	_ "io/ioutil"
	"log"
	"net/http"
	"time"

	worker "github.com/youngwook/security-gateway-apis/libraries"
)

type Message struct {
	Result string `json:"Result"`
}

var message Message
var insecureSkipVerify = flag.Bool("insureskipverify", true, "skip server side SSL verification, mainly for self-signed cert")

var configFileLocation = flag.String("configfile", "res/configuration.toml", "configuration file")

var config, err = worker.LoadTomlConfig(*configFileLocation)

var client = getNewClient(*insecureSkipVerify)

var er = worker.EdgeXRequestor{ProxyBaseURL: config.GetProxyBaseURL(), SecretSvcBaseURL: config.GetSecretSvcBaseURL(), Client: client}
var service = &worker.Service{Connect: &er, CertCfg: config, ServiceCfg: config}

type UserInfo struct {
	User  string `json:"User"`
	Group string `json:"Group"`
}

func main() {

	if err != nil {
		fmt.Println("failed to retrieve config data from local file. Please make sure res/configuration.toml file exists with correct formats")
		return
	}

	handleRequests()
}
func homePage(w http.ResponseWriter, r *http.Request) {
	err = service.CheckProxyServiceStatus()
	if err != nil {

		message.Result = "the service is closed!"
		json.NewEncoder(w).Encode(message)
	} else {
		message.Result = "the service is running"

		json.NewEncoder(w).Encode(message)
	}
	fmt.Println("Endpoint Hit: entry point")
}
func initKong(w http.ResponseWriter, r *http.Request) {
	err = service.Init()
	if err != nil {
		fmt.Println(err.Error())
		message.Result = "something wrong!"
		json.NewEncoder(w).Encode(message)
	} else {
		message.Result = "Init Success!"

		json.NewEncoder(w).Encode(message)
	}
	fmt.Println("Endpoint Hit: InitKong")
}
func resetKong(w http.ResponseWriter, r *http.Request) {
	err = service.ResetProxy()

	if err != nil {
		fmt.Println("error")
		message.Result = "something wrong!"
		json.NewEncoder(w).Encode(message)
	} else {
		message.Result = "Reset Success!"

		json.NewEncoder(w).Encode(message)
	}
	fmt.Println("Endpoint Hit: resetKong")
}

func getTocken(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	user := vars["user"]
	if user != "" {
		c := &worker.Consumer{Name: user, Connect: &er, Cfg: config}
		t, err := c.GetToken(user)
		if err != nil {
			fmt.Println(fmt.Sprintf("failed to retrive access token for edgex service due to error %s", err.Error()))
			return
		}
		message.Result = fmt.Sprintf(t)
		json.NewEncoder(w).Encode(message)
	} else {
		message.Result = "failed get the tocken"

		json.NewEncoder(w).Encode(message)
	}
	fmt.Println("Endpoint Hit: get tocken")
}
func createUser(w http.ResponseWriter, r *http.Request) {
	// get the body of our POST request
	// return the string response containing the request body
	cc := UserInfo{}
	json.NewDecoder(r.Body).Decode(&cc)

	if cc.User != "" && cc.Group != "" {
		c := &worker.Consumer{Name: cc.User, Connect: &er, Cfg: config}

		err := c.Create(worker.EdgeXService)
		if err != nil {
			message.Result = err.Error()

			json.NewEncoder(w).Encode(message)
			fmt.Println(err.Error())
			return
		}

		err = c.AssociateWithGroup(cc.Group)
		if err != nil {
			message.Result = err.Error()

			json.NewEncoder(w).Encode(message)
			fmt.Println(err.Error())
			return
		}

		t, err := c.CreateToken()
		if err != nil {
			fmt.Println(fmt.Sprintf("failed to create access token for edgex service due to error %s", err.Error()))
			return
		}

		fmt.Println(fmt.Sprintf("the access token for user %s is: %s. Please keep the token for accessing edgex services", cc.User, t))

		tf := &worker.TokenFileWriter{Filename: "accessToken.json"}
		err = tf.Save(cc.User, t)
		if err != nil {
			fmt.Println(err.Error())
		}
		message.Result = fmt.Sprintf("user %s created! the access token is %s", cc.User, t)
		json.NewEncoder(w).Encode(message)
	} else {
		message.Result = "failed to create user"

		json.NewEncoder(w).Encode(message)
	}
	fmt.Println("Endpoint Hit: createUser")
}
func deleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["user"]
	t := &worker.Consumer{Name: key, Connect: &er, Cfg: config}
	t.Delete()
	fmt.Fprintf(w, "Delete user %s", key)
}
func handleRequests() {
	// creates a new instance of a mux router
	myRouter := mux.NewRouter().StrictSlash(true)

	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/init", initKong)
	myRouter.HandleFunc("/reset", resetKong)
	myRouter.HandleFunc("/deleteUser/{user}", deleteUser).Methods("DELETE")

	myRouter.HandleFunc("/createUser", createUser).Methods("POST")
	myRouter.HandleFunc("/getTocken/{user}", getTocken)

	log.Fatal(http.ListenAndServe(":8088", myRouter))
}
func getNewClient(skipVerify bool) *http.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: skipVerify},
	}
	return &http.Client{Timeout: 10 * time.Second, Transport: tr}
}
