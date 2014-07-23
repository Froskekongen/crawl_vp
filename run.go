package main


import (
	"flag"
    elastigo "github.com/mattbaird/elastigo/lib"
    //"regexp"
    "strconv"
    "sync"
    "github.com/Froskekongen/crawl_vp/crawlers"
    "fmt"
    //"time"
)

var (
	eshost *string = flag.String("host", "172.30.31.203", "Elasticsearch Server Host Address")
)


func main(){

  

    mapPages:=crawl_vp.GetNPages()

    base1:="http://www.vinmonopolet.no/vareutvalg/sok?query=*&sort=2&sortMode=0&page="
    base2:="&filterIds=25&filterValues="
    //wineTypes:=[]string{"R%C3%B8dvin","Musserende+vin","Hvitvin","Ros%C3%A9vin","Fruktvin","Sterkvin","Brennevin","%C3%B8l"}
    maxPPerType:=10000
    
    

    tasks := make(chan string,1000)

    changedChan:=make(chan *[]crawl_vp.WineRep,1)
    newChan:=make(chan *[]crawl_vp.WineRep,1)
    cW:=make([]crawl_vp.WineRep,0,1000)
    nW:=make([]crawl_vp.WineRep,0,1000)
    changedChan <- &cW
    newChan <- &nW




    retry_url:=make(chan map[string]int,1)
    retry_url <- map[string]int{"base":1}

    elastChan:=make(chan *elastigo.Conn,1)
    c:= elastigo.NewConn()
    c.Domain = *eshost

    crawl_vp.TestIndex1(c)
    crawl_vp.TestIndex2(c)
    elastChan <- c
    wr,_:=crawl_vp.EsSearch(elastChan,89101)
    fmt.Println(wr)
    wr.UpdatePrice(200)
    fmt.Println(wr)

    
    
    // spawn four worker goroutines
    var wg sync.WaitGroup

//    for i := 0; i < 4; i++ {
//        wg.Add(1)
//        go func() {
//            for url := range tasks {
//                crawl_vp.GetProducts(url,retry_url)
//            }
//            wg.Done()
//        }()
//    }

    for i := 0; i < 4; i++ {
        wg.Add(1)
        go func() {
            for url := range tasks {
                crawl_vp.GetProductsWithES(url,elastChan,changedChan,newChan)
            }
            wg.Done()
        }()
    }

    // generate some tasks
    kkk:=0
    for key,val:=range mapPages{
        for iii:=1;iii<=val;iii++{
            ss:=base1+strconv.Itoa(iii)+base2+key
            kkk++
            fmt.Println(kkk,ss)
            tasks <- ss
            if iii>maxPPerType{break}
        }
    }
    
    
    close(tasks)
    wg.Wait()
    close(elastChan)
    
}
