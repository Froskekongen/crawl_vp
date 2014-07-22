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



type WineRep struct{
    //strings
    Url string //`json:"url"`
    Name string //`json:"name"`
    WineType string //`json:"WineType"`
    Producer string //`json:"producer"`
    Wholesaler string //`json:"wholesaler"`
    Material string //`json:"material"`
    Country string //`json:"country"`
    Subcountry1 string //`json:"subcountry1"`
    Subcountry2 string //`json:"subcountry2"`
    Store string //`json:"store"`
    Distributor string //`json:"distributor"`

    //ints
    Prodnum uint64 //`json:"prodnum"`
    Vintage uint16 //`json:"vintage"`
    Price []uint64 //`json:"price"` // change to uint64

    //floats
    Alcohol float64 //`json:"alcohol"`
    Sugar float64 //`json:"sugar"`
    Acid float64 //`json:"acid"`
    

    //Flags
    Soldout bool //`json:"soldout"`
    Obsoleteproduct bool //`json:"obsoleteproduct"`
    Deeplookup bool //`json:"deeplokup"`

    //Datetimes
    LookupTimes []time.Time
}


func (wr *WineRep)UpdatePrice(px uint64){
//    timeN:=len(wr.LookupTimes)
//    priceN:=len(wr.Price)
    wr.LookupTimes=append(wr.LookupTimes,time.Now())
    wr.Price=append(wr.Price,px)
}


var productRegex *regexp.Regexp = regexp.MustCompile(`<a href="([\w\:\-/\.]+).*?" class="product">(.+)</a>\s*</h3>\s*<p>\s*(\S*\s?\S*\s?\S*\s?\S*)\s*\((\d+)\)(?s:.*?)<strong>Kr\.\s+(\d*\.?\d*).*?\s+</strong>`)


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






func ParseRegexMatch(mts *[][][]byte,lenM int)[] WineRep{
    m:=make([]WineRep,lenM)

    for iii,val:=range *mts{
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
                        //pf,_:=strconv.ParseFloat(string(vv),64)
                        pf,_:=strconv.ParseUint(string(vv),10,64)
                        m[iii].Price=[]uint64{pf}
                        m[iii].LookupTimes=[]time.Time{time.Now()}
                }
            }
        }
    }
    return m
}


func GetProducts(url string,urlChan chan string,retryChan chan map[string]int) {
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

