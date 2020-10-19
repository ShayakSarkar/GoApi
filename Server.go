package main

import(
	"fmt"
	"net/http"
	//"net/url"
	"log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"reflect"
	"context"
	"time"
	"regexp"
	"encoding/json"
	"strings"
)

func getMeetingById(meetingDatabase *mongo.Database,w http.ResponseWriter,r *http.Request)int{
	splittedPath:=strings.Split(r.URL.Path,"/")
	id:=splittedPath[len(splittedPath)-1]
	fmt.Println("id",id)
	meetingCollection:=meetingDatabase.Collection("Meeting")
	var meetingDetails bson.M=bson.M{}
	meetingCollection.FindOne(context.TODO(),bson.M{"id":id}).Decode(&meetingDetails)
	fmt.Println(meetingDetails)
	rawJson,err:=json.Marshal(meetingDetails)
	if err!=nil{
		fmt.Println("Error marshalling json")
		return -1
	}
	fmt.Fprintf(w,string(rawJson))
	return 0
}

func getAllMeetingsOfParticipant(meetingDatabase *mongo.Database,w http.ResponseWriter, r *http.Request)int{
	/*
	email:=r.URL.Query()["participant"][0]
	fmt.Println(email)
	meetingCollection:=meetingDatabase.Collection("Meeting")

	finalJson:=make(map[string]interface{})
	finalJson["root"]=make([]interface{},0,1)
	fmt.Println(finalJson)

	cur,_:=meetingCollection.Find(context.TODO(),bson.D{})
	cur.Next(context.TODO())
	object:=bson.M{}
	cur.Decode(&object)
	par:=object["participants"].([]interface{})
	fmt.Println(par)

	cur.Next(context.TODO())
	cur.Decode(&object)
	fmt.Println(object)
	*/
	return 0
}

func createMeeting(meetingDatabase *mongo.Database,r *http.Request) int{

	meetingCollection:=meetingDatabase.Collection("Meeting")

	var meetingDetails bson.M=bson.M{}
	json.NewDecoder(r.Body).Decode(&meetingDetails)
	fmt.Println("Got meeting details")
	fmt.Println(meetingDetails)

	_,has:=meetingDetails["id"]
	if has!=true{
		fmt.Println("API error: id not found")
		return -1
	}

	_,has=meetingDetails["title"]
	if has!=true{
		fmt.Println("API error: title not fund")
		return -2
	}

	_,has=meetingDetails["participants"]
	if has!=true{
		fmt.Println("API error: participants not found")
		return -3
	}

	_,has=meetingDetails["start"]
	if has!=true{
		fmt.Println("API error: start time not found")
		return -4
	}
	_,has=meetingDetails["end"]
	if has!=true{
		fmt.Println("API error: end time not found")
		return -5
	}
	//Creation timestamp is needed

	if len(meetingDetails)!=5{
		fmt.Println("API error: extra params in json body")
		return -6
	}
	ctx,_:=context.WithTimeout(context.Background(),time.Second*10)

	var existingMeet bson.M=bson.M{}
	meetingCollection.FindOne(ctx,bson.M{"id":meetingDetails["id"]}).Decode(&existingMeet)

	fmt.Println("ex meet",existingMeet)
	if len(existingMeet)!=0{
		fmt.Println("Meeting already present")
		return -7
	}

	fmt.Println("attempting to insert")
	result,err:=meetingCollection.InsertOne(ctx,meetingDetails)
	if err!=nil{
		fmt.Println("error during insertion")
		return -6
	}
	fmt.Println("Inserted successfully",result)
	return 0
}

func createRouterWithDBAccess(meetingDatabase *mongo.Database) func(http.ResponseWriter,*http.Request){

	router:=func (w http.ResponseWriter, r *http.Request){
		fmt.Println(r.URL.Path)
		fmt.Println(r.URL.Query())
		createMeetingRegex:="/meeting"
		getMeetingRegex:="/meeting/[0-9 a-z]"
		getMeetingsRegex:="/meetings"
		switch r.Method{
			case "GET":
					match1,_:=regexp.MatchString (getMeetingRegex, r.URL.Path)
					match2,_:=regexp.MatchString (getMeetingsRegex, r.URL.Path)

					if match1{
						fmt.Println("get meeting by id")
						getMeetingById(meetingDatabase,w,r)
					}else if match2{
						params:=r.URL.Query()
						_,hasStart:=params["start"]
						_,hasEnd:=params["end"]
						_,hasParticipant:=params["participant"]

						if hasStart && hasEnd && len(params)==2{
							fmt.Println("get meetings by timing range")
						}else if hasParticipant && len(params)==1{
							fmt.Println("get meetings by participant")
							getAllMeetingsOfParticipant(meetingDatabase,w,r)
						}
					}
			case "POST":
					match4,_:=regexp.MatchString (createMeetingRegex, r.URL.Path)
					if match4 && len(r.URL.Query())==0{
						fmt.Println("create a meeting")
						createMeeting(meetingDatabase,r)
					}
				}
	}
	return router
}
func main(){
	fmt.Println("Setting up connection")
	client,err:=mongo.NewClient(options.Client().ApplyURI("mongodb+srv://ShayakSarkar:qwerty1234@meetingcluster.kgaaa.mongodb.net/MeetingDatabase?retryWrites=true&w=majority"))
	if err!=nil{
		fmt.Println("error in creating client")
	}
	ctx,_:=context.WithTimeout(context.Background(),10*time.Second)
	err=client.Connect(ctx)
	if err!=nil{
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	err=client.Ping(ctx,readpref.Primary())
	if err!=nil{
		fmt.Println("Ping not successful, connection is faulty")
		log.Fatal(err)
	}
	databases,err:=client.ListDatabaseNames(ctx,bson.D{})
	if err!=nil{
		log.Fatal(err)
	}
	fmt.Println(databases)
	meetingDatabase:=client.Database("MeetingDatabase")
	fmt.Println(reflect.TypeOf(meetingDatabase))

	routerFunction:=createRouterWithDBAccess(meetingDatabase)
	http.HandleFunc("/",routerFunction)
	http.ListenAndServe(":8080",nil)
}
