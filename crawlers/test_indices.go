package crawl_vp


import (
	"flag"
    elastigo "github.com/mattbaird/elastigo/lib"
    //"regexp"
    //"strconv"
    "sync"
    "fmt"
    //"time"
    "testing"
)




func TestCrawlers(t *testing.T){
    var (
	t_eshost *string = flag.String("host", "172.30.31.203", "Elasticsearch Server Host Address")
    )
    mapPages:=GetNPages()
    maxPPerType:=0

    tasks := make(chan string,1000)

    changedChan:=make(chan ListOfWines,1)
    newChan:=make(chan ListOfWines,1)
    var cW ListOfWines
    var nW ListOfWines
    changedChan <- cW
    newChan <- nW




    retry_url:=make(chan map[string]int,1)
    retry_url <- map[string]int{"base":1}

    elastChan:=make(chan *elastigo.Conn,1)
    c:= elastigo.NewConn()
    c.Domain = *t_eshost

    TestIndex1(c)
    TestIndex2(c)
    elastChan <- c
    wr,_:=EsSearch(elastChan,89101)
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
                GetProductsWithES(url,elastChan,changedChan,newChan)
            }
            wg.Done()
        }()
    }

    // generate some tasks
    GenerateTasks(mapPages,tasks,maxPPerType)
    
    
    close(tasks)
    wg.Wait()
    close(elastChan)
    
}
