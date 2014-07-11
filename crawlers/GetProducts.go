package crawl_vp

import (
    "fmt"
    "net/http"
    "io/ioutil"
    "regexp"
    "strconv"
    //"sync"
    "html"
)

//type Fetcher interface {
//    // Fetch returns the body of URL and
//    // a slice of URLs found on that page.
//    Fetch(url string) (body []byte, err error)
//}


//type crawlShared struct {
//    mapAccess chan map[string]bool
//    printAccess chan bool
//}

//func (cs *crawlShared) Fetch(url string) ([]byte,error){
//    resp,err1:=http.Get(url)    
//    defer resp.Body.Close()
//    if err1!=nil{
//        var b []byte
//        return b,err1
//    }
//    body, err2 := ioutil.ReadAll(resp.Body)
//    if err2!=nil{
//        var b []byte
//        return b,err2
//    }
//    return body,nil
//}

type WineRep struct{
    //strings
    Url string `json:"url"`
    Name string `json:"name"`
    WineType string `json:"WineType"`
    Producer string `json:"producer"`
    Wholesaler string `json:"wholesaler"`
    Material string `json:"material"`
    Country string `json:"country"`
    Subcountry1 string `json:"subcountry1"`
    Subcountry2 string `json:"subcountry2"`
    Store string `json:"store"`

    //ints
    Prodnum uint64 `json:"prodnum"`
    Vintage uint64 `json:"vintage"`
    Distributor uint64 `json:"distributor"`

    //floats
    Alcohol float64 `json:"alcohol"`
    Sugar float64 `json:"sugar"`
    Acid float64 `json:"acid"`
    Price float64 `json:"price"`

    //Flags
    Soldout bool `json:"soldout"`
    Obsoleteproduct bool `json:"obsoleteproduct"`
    Deeplookup bool `json:"deeplokup"`
    
}










func GetProducts(url string,reg *regexp.Regexp,urlChan chan string,retryChan chan map[string]int) {
//    maxRetries:=5
    resp,err1:=http.Get(url)
    defer resp.Body.Close()
    if err1!=nil{
//        m:=<-retryChan
//        m[url]++
//        retries:=m[url]
//        retryChan <- m
//        if retries<maxRetries{
//            urlChan <- url
//        }
        return
    }
    

    body, err2 := ioutil.ReadAll(resp.Body)
    if err2!=nil{
//        m:=<-retryChan
//        m[url]++
//        retries:=m[url]
//        retryChan <- m
//        if retries<maxRetries{
//            urlChan <- url
//        }
        return
    }
    //fmt.Println(string(body))

    mm:=reg.FindAllSubmatch(body,-1)
    lenM:=len(mm)
    if lenM==0{
//        m:=<-retryChan
//        m[url]++
//        retries:=m[url]
//        retryChan <- m
//        if retries<maxRetries{
//            urlChan <- url
//        }
        return
    }
    m:=make([]WineRep,lenM)


    //fmt.Println("\n"+url+"\n")
    for iii,val:=range mm{
//        fmt.Println(string(val[0])+"\n")
//        fmt.Println(string(val[1])+"\n\n")
        for jjj,vv:=range val{
            
            if jjj!=0{
                //fmt.Println(string(vv))
                switch{
                    case jjj==1:
                        m[iii].Url=string(vv)//+"?ShowShopsWithProdInStock=true&sku=1492601&fylke_id=*"
                    case jjj==2:
                        m[iii].Name=html.UnescapeString(string(vv))
                    case jjj==3:
                        m[iii].WineType=html.UnescapeString(string(vv))
                    case jjj==4:
                        pf,_:=strconv.ParseUint(string(vv),10,64)
                        //pf,_:=strconv.ParseUint("10",10,64)
                        m[iii].Prodnum=pf
                    case jjj==5:
                        pf,_:=strconv.ParseFloat(string(vv),64)
                        m[iii].Price=pf
                }
            }
        }
        //fmt.Println(iii,"\n\n")
        //fmt.Println(val[1])
    }
    fmt.Println(url)
    fmt.Println(lenM)
    return 
}

