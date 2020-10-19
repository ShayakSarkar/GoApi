//Semaphores can be modelled using wait groups and mutex9es 
//are only some special cases of semaphores

package main
import(
	"reflect"
	"fmt"
	"context"
	"log"
	"time"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Person struct{
	name string
	weight string
}

func main(){
	fmt.Println("Started")
	client,err:=mongo.NewClient(options.Client().ApplyURI("mongodb+srv://ShayakSarkar:qwerty1234@meetingcluster.kgaaa.mongodb.net/MeetingDatabase?retryWrites=true&w=majority"))

	if err!=nil{
		log.Fatal(err)
	}
	ctx,_:=context.WithTimeout(context.Background(),10*time.Second)
	err=client.Connect(ctx)
	if err!=nil{
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	err=client.Ping(ctx,readpref.Primary())
	if(err!=nil){
		fmt.Println("Ping not successful")
		log.Fatal(err)
	}
	databases,err:=client.ListDatabaseNames(ctx,bson.D{})
	if err!=nil{
		log.Fatal(err)
	}
	fmt.Println(databases)
	md:=client.Database("MeetingDatabase")
	mc:=md.Collection("Meeting")
	fmt.Println(reflect.TypeOf(md))
	fmt.Println(reflect.TypeOf(mc));
	fmt.Println(reflect.TypeOf(ctx))
}

