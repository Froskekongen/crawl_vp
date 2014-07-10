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
    Url string `json:"url"`
    Name string `json:"name"`
    WineType string `json:"WineType"`
    Prodnum uint64 `json:"prodnum"`
    Price float64 `json:"price"`
}










func GetProducts(url string,reg *regexp.Regexp,urlChan chan string) {
    resp,err1:=http.Get(url)
    if err1!=nil{
        urlChan <- url
        return 
    }
    defer resp.Body.Close()

    body, err2 := ioutil.ReadAll(resp.Body)
    if err2!=nil{
        urlChan <- url
        return
    }
    //fmt.Println(string(body))

    mm:=reg.FindAllSubmatch(body,-1)
    lenM:=len(mm)
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

