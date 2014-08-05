package main

import(
	//"encoding/json"
	"flag"
	//"time"
    "fmt"
    //"os"
    "github.com/Froskekongen/crawl_vp/crawlers"
    elastigo "github.com/mattbaird/elastigo/lib"
)

//var (
//	eshost *string = flag.String("host", "172.30.31.203", "Elasticsearch Server Host Address")
//)


func errf(err error){
    if err!=nil{
        fmt.Println(err)
    }
}

func main(){
	//log.SetFlags(log.Ltime | log.Lshortfile)

    var eshost *string = flag.String("eshost", "localhost", "Elasticsearch Server Host Address")
    flag.Parse()

    elastChan:=make(chan *elastigo.Conn,1)
    c:= elastigo.NewConn()
    c.Domain = *eshost
    elastChan <- c
    crawl_vp.EsScanForDeep(elastChan)
}
