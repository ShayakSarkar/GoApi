package main

import(
	//"strings"
	"fmt"
	"net/http"
	"log"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"reflect"
)

func handler(w http.ResponseWriter, r *http.Request){
	var bodyMap bson.M=bson.M{}
	json.NewDecoder(r.Body).Decode(&bodyMap)

	list:=bodyMap["participants"].([]interface{})
	list2:=[]map[string]string{}
	for i:=0;i<len(list);i++{
		finalMap:=make(map[string]string)
		mp1:=list[i]
		mp2:=mp1.(map[string]interface{})
		email:=mp2["email"].(string)
		rsvp:=mp2["rsvp"].(string)
		reflect.TypeOf(email)
		finalMap["email"]=email
		finalMap["rsvp"]=rsvp
		list2=append(list2,finalMap)
	}
	fmt.Println(list2)
	fmt.Println((list2[0])["email"])
	fmt.Fprintf(w,"{\"name\":\"Shayak\",\"lastname\":\"Sarkar\"}")
}

func main(){
	http.HandleFunc("/",handler)
	//This is the handler that spawns the 
	//entire routing method. I.E, the HandeFunc
	//has the logic for routing all requests and
	//this function is presented with the Response
	//Writer and http.Request Object presumably
	//from the environment(need to see exactly how)
	//and then uses them to route requests effectively
	fmt.Println("Started Server")
	log.Fatal(http.ListenAndServe(":8080",nil))
}

