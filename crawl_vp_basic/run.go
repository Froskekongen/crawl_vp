package main


import (
	"flag"
    elastigo "github.com/mattbaird/elastigo/lib"
    //"regexp"
    //"strconv"
    "sync"
    //"github.com/Froskekongen/crawl_vp/crawlers"
    "github.com/Froskekongen/crawl_vp/crawlers"
    //"fmt"
    //"time"
    //"encoding/json"
)

var (
	eshost *string = flag.String("host", "172.30.31.203", "Elasticsearch Server Host Address")
)


func main(){

  

    mapPages:=crawl_vp.GetNPages()
    maxPPerType:=100000


    tasks := make(chan string,1000)

    changedChan:=make(chan crawl_vp.ListOfWines,1)
    newChan:=make(chan crawl_vp.ListOfWines,1)
    cW := crawl_vp.ListOfWines( make([]crawl_vp.WineRep,0,100) )
    //cW := crawl_vp.ListOfWines(cw)
    nW := crawl_vp.ListOfWines( make([]crawl_vp.WineRep,0,100) )
    //nW := crawl_vp.ListOfWines(nw)
    changedChan <- cW
    newChan <- nW




    retry_url:=make(chan map[string]int,1)
    retry_url <- map[string]int{"base":1}

    elastChan:=make(chan *elastigo.Conn,1)
    c:= elastigo.NewConn()
    c.Domain = *eshost

    crawl_vp.TestIndex1(c)
    crawl_vp.TestIndex2(c)
    elastChan <- c
//    wr,_:=crawl_vp.EsSearch(elastChan,89101)
//    fmt.Println(wr)
//    wr.UpdatePrice(200)
//    fmt.Println(wr)

    
    
    // spawn four worker goroutines
    var wg sync.WaitGroup

    for i := 0; i < 6; i++ {
        wg.Add(1)
        go func() {
            for url := range tasks {
                crawl_vp.GetProductsWithES(url,elastChan,changedChan,newChan)
            }
            wg.Done()
        }()
    }

    crawl_vp.GenerateTasks(mapPages,tasks,maxPPerType)
    
    
    close(tasks)
    wg.Wait()
    close(elastChan)

    newWines:= <- newChan
    changedWines:= <-changedChan

    if len(newWines)>0 || len(changedWines)>0{
        crawl_vp.SendChangedAndNew(changedWines,newWines)
    }

    close(newChan)
    close(changedChan)

    
}
