package crawl_vp

import (
    "fmt"
    "net/http"
    "io/ioutil"
    "regexp"
    "strconv"
    //"sync"
    "html"
    "time"
    "encoding/json"
    elastigo "github.com/mattbaird/elastigo/lib"
)






var productRegex *regexp.Regexp = regexp.MustCompile(`<a href="([\w\:\-/\.]+).*?" class="product">(.+)</a>\s*</h3>\s*<p>\s*(\S*\s?\S*\s?\S*\s?\S*)\s*\((\d+)\)(?s:.*?)<strong>Kr\.\s+(\d*\.?\d*).*?\s+</strong>`)
//var productRegex *regexp.Regexp = regexp.MustCompile(`<a href="([\w\:\-/\.]+).*?" class="product">(.+)</a>\s*</h3>\s*<p>\s*(\S*\s?\S*\s?\S*\s?\S*)\s*\((\d+)\)(?s:.*?)<strong>Kr\.\s+(\d*\.?\d*).*?\s+</strong>(?s:.*?)<h3 class="stock">([.\s]*?)<`)


func EsSearch(esConn chan *elastigo.Conn,prodNum uint64)(WineRep,bool){
    searchJson := `{
        "query": {
            "filtered": {
                "filter": {
                    "term": {"Prodnum":%d}
                }
            }
        }
    }`
    searchJson=fmt.Sprintf(searchJson,prodNum)
    c:= <-esConn
    sResp,err:=c.Search("wines","product",nil,searchJson)
    if err!=nil{
        fmt.Println(err)
        return WineRep{},false
    }
    esConn <- c
    if len(sResp.Hits.Hits) == 1{
        var wr WineRep
        err := json.Unmarshal(*sResp.Hits.Hits[0].Source, &wr)
        if err!=nil{
            fmt.Println(err)
            return WineRep{},false  
        }
        //fmt.Println(wr)
        return wr,true
    }

    //fmt.Println(searchJson)
    return WineRep{},false
}

func GetProductsWithES(url string,esConn chan *elastigo.Conn,changedChan chan ListOfWines,newChan chan ListOfWines){
    resp,err1:=http.Get(url)
    defer resp.Body.Close()
    if err1!=nil{
        return
    }
    body, err2 := ioutil.ReadAll(resp.Body)
    if err2!=nil{
        return
    }
    reMatch:=productRegex.FindAllSubmatch(body,-1)
    for _,m := range reMatch{
        pn,err3:=strconv.ParseUint(string(m[4]),10,64)
        if err3!=nil{
            continue
        }
        wr,exists:=EsSearch(esConn,pn)
        if exists{
            fmt.Println("wine exists in es",wr.Name)
            pf,err4:=strconv.ParseUint(string(m[5]),10,64)
            if err4!=nil{
                continue
            }
            if pf!=wr.Price[len(wr.Price)-1]{
                wr.UpdatePrice(pf)
                c:= <-esConn
                wr.LastWritten=time.Now()
                c.Index("wines","product",string(m[4]),nil,wr)
                esConn <- c

                change:= <- changedChan
                change=append(change,wr) //troublesome line
                changedChan <- change
            }
        } else {
            wr=ParseOneMatch(&m)
            c:= <-esConn
            wr.LastWritten=time.Now()
            c.Index("wines","product",string(m[4]),nil,wr)
            esConn <- c

            newWR:= <- newChan
            newWR=append(newWR,wr) //troublesome line
            changedChan <- newWR
        }
    }
    return   
}


func ParseOneMatch(val *[][]byte)WineRep{
    var wr WineRep
    for jjj,vv:=range *val{
        if jjj!=0{
            //fmt.Println(string(vv))
            switch{
                case jjj==1:
                    wr.Url=string(vv)//+"?ShowShopsWithProdInStock=true&sku=1492601&fylke_id=*"
                case jjj==2:
                    wr.Name=html.UnescapeString(string(vv))
                case jjj==3:
                    wr.WineType=html.UnescapeString(string(vv))
                case jjj==4:
                    pf,_:=strconv.ParseUint(string(vv),10,64)
                    //pf,_:=strconv.ParseUint("10",10,64)
                    wr.Prodnum=pf
                case jjj==5:
                    //pf,_:=strconv.ParseFloat(string(vv),64)
                    pf,_:=strconv.ParseUint(string(vv),10,64)
                    wr.Price=[]uint64{pf}
                    wr.LookupTimes=[]time.Time{time.Now()}
            }
        }
    }
    return wr
}



func ParseRegexMatch(mts *[][][]byte,lenM int)[] WineRep{
    m:=make([]WineRep,lenM)

    for iii,val:=range *mts{
        wr:=ParseOneMatch(&val)
        m[iii]=wr
    }
    return m
}


func GetProducts(url string,retryChan chan map[string]int) {
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

    mm:=productRegex.FindAllSubmatch(body,-1)
    lenM:=len(mm)
    if lenM==0{
//        m:= <- retryChan
//        m[url]++
//        retries:=m[url]
//        retryChan <- m
//        if retries<maxRetries{
//            urlChan <- url
//        }
        return
    }
    var mPtr *[][][]byte = &mm
    m:=ParseRegexMatch(mPtr,lenM)

    b,_:=json.Marshal(m[0])
    fmt.Println(string(b))
    fmt.Println(url)
    fmt.Println(lenM)
    return 
}

