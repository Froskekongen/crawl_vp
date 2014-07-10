package crawl_vp

import (
    "fmt"
//    "IHeap"
//    "container/heap"
    "net/http"
    "io/ioutil"
    "regexp"
//    "sync"
//    "os/exec"
    "strconv"
    //"encoding/binary"
)

func nPage(baseUrl string,pType string,reg *regexp.Regexp,mc chan map[string]int,pDone chan bool){
    url:=baseUrl+pType
    resp,err1:=http.Get(url)
    defer resp.Body.Close()
    if err1!=nil{
        pDone<-true
        return
    }
    body, err2 := ioutil.ReadAll(resp.Body)
    if err2!=nil{
        pDone<-true
        return
    }
    reM:=reg.FindAllSubmatch(body,-1)
    largest:=1
    for _,val := range reM{
        valInt,_:=strconv.Atoi(string(val[1]))
        if valInt>largest{
            largest=valInt
        } 
    }
    m:= <- mc
    m[pType]=largest
    fmt.Println(largest)
    mc <- m
    pDone<-true
    return
} 

func GetNPages()map[string]int {
    m:=make(chan map[string]int,1)
    baseUrl:="http://www.vinmonopolet.no/vareutvalg/sok?query=*&sort=2&sortMode=0&page=1&filterIds=25&filterValues="
    wineTypes:=[]string{"R%C3%B8dvin","Musserende+vin","Hvitvin","Ros%C3%A9vin","Fruktvin","Sterkvin","Brennevin","%C3%B8l"}

    rr:=regexp.MustCompile(`query=\*&amp;sort=2&amp;sortMode=0&amp;page=\d+&amp;filterIds=25&amp;filterValues=\S{1,15}">(\d+)`)
    m <- map[string]int{"tt":1}
    done:=make(chan bool,len(wineTypes))
    for _,tVal:=range wineTypes{
        fmt.Println(tVal)
        go nPage(baseUrl,tVal,rr,m,done)
    }
    //time.Sleep(time.Second)
    for iii:=0;iii<len(wineTypes);iii++{
        <-done
    }
    mm:=<-m
    delete(mm,"tt")

    return mm
}





//func main() {
////    baseUrl:="http://www.vinmonopolet.no/vareutvalg/sok?query=*&sort=2&sortMode=0&page="
////    filterUrl:="&filterIds=25&filterValues="
//    wineTypes:=[]string{"R%C3%B8dvin","Musserende+vin","Hvitvin","Ros%C3%A9vin","Fruktvin","Sterkvin","Brennevin","Øl"}

//    fmt.Println(wineTypes)

//    typeNPages:=getNPages()
//    fmt.Println(typeNPages)


////    tasks := make(chan *exec.Cmd, 64)

////    // spawn four worker goroutines
////    var wg sync.WaitGroup
////    for i := 0; i < 4; i++ {
////        wg.Add(1)
////        go func() {
////            for cmd := range tasks {
////                cmd.Run()
////            }
////            wg.Done()
////        }()
////    }

////    // generate some tasks
////    for i := 0; i < 10; i++ {
////        tasks <- exec.Command("zenity", "--info", "--text='Hello from iteration n."+strconv.Itoa(i)+"'")
////    }
////    close(tasks)

////    // wait for the workers to finish
////    wg.Wait()

//}
