package main


import (
    //"regexp"
    "strconv"
    //"sync"
    "github.com/froskekongen/crawl_vp/crawlers"
    "time"
    "fmt"
)


func main(){

    mapPages:=crawl_vp.GetNPages()

    base1:="http://www.vinmonopolet.no/vareutvalg/sok?query=*&sort=2&sortMode=0&page="
    base2:="&filterIds=25&filterValues="
    //wineTypes:=[]string{"R%C3%B8dvin","Musserende+vin","Hvitvin","Ros%C3%A9vin","Fruktvin","Sterkvin","Brennevin","%C3%B8l"}
    
    

    tasks := make(chan string,1000)


//    rr2:=regexp.MustCompile(`<a href="([\w\:\-/\.]+).*?" class="product">(.+)</a>\s*</h3>\s*<p>\s*(\S*)\s*\((\d+)\)(?s:.*?)<strong>Kr\.\s+(\d+).*?\s+</strong>`)

//    retry_url:=make(chan map[string]int,1)
//    retry_url <- map[string]int{"base":1}
//    
//    // spawn four worker goroutines
//    var wg sync.WaitGroup
//    for i := 0; i < 4; i++ {
//        wg.Add(1)
//        go func() {
//            for url := range tasks {
//                crawl_vp.GetProducts(url,rr2,tasks,retry_url)
//            }
//            wg.Done()
//        }()
//    }

    // generate some tasks
    kkk:=0
    for key,val:=range mapPages{
        for iii:=1;iii<=val;iii++{
            ss:=base1+strconv.Itoa(iii)+base2+key
            kkk++
            fmt.Println(kkk,ss)
            tasks <- ss
        }
    }
    close(tasks)

    // wait for the workers to finish

    time.Sleep(1*time.Second)
//    fmt.Println(wg)
//    wg.Wait()
}
