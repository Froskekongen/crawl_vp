package crawl_vp

import (
    //"io"
    "fmt"
    "regexp"
    "net/http"
    "strconv"
    "io/ioutil"
    "strings"
    "html"
    "log"
    "reflect"
    "time"
    "sync"
    elastigo "github.com/mattbaird/elastigo/lib"
    "encoding/json"
)


func fixDecimals(dd string)float64{
    dd=strings.Replace(dd,",",".",-1)
    aa,err:=strconv.ParseFloat(dd,64)
    if err!=nil{
        log.Println(err)
        return 0.
    }
    return aa
}

func countryRegion(cc string)(string,string){
    ccs:=strings.Split(cc,",")
    N:=len(ccs)
    if N==0{return "",""
    } else if N==1{return ccs[0],""
    } else {
        for iii:=0;iii<N;iii++{
            ccs[iii]=strings.TrimSpace(ccs[iii])        
        }
        return ccs[0],strings.Join(ccs[1:],", ")
    }
}

func countrySplit(dd []string){}

var deepRegexes map[string]*regexp.Regexp = map[string]*regexp.Regexp {
    "Alcohol":regexp.MustCompile(`Alkohol\s*(\d{1,2},?\d?)`),
    //"Producer":regexp.MustCompile(`Produsent:\s*</strong>\s*<span class="data">\s*([\s.]+?)\s*<`),
    "Producer":regexp.MustCompile(`Produsent:\s*</strong>\s*<span class="data">((?s).*?)</span>`),
    "Sugar":regexp.MustCompile(`Sukker:?\s*(\d{1,3},?\d{0,2})`),
    "Acid":regexp.MustCompile(`Syre:?\s*(\d{1,3},?\d{0,2})`),
    "Distributor":regexp.MustCompile(`Distribut&oslash;r:?\s*</strong>\s*<span class="data">((?s).*?)</span>`),
    "Wholesaler":regexp.MustCompile(`Grossist:?\s*</strong>\s*<span class="data">((?s).*?)</span>`),
    "Material":regexp.MustCompile(`Råstoff:?\s*</strong>\s*<span class="data">((?s).*?)</span>`),
    "Selection":regexp.MustCompile(`Produktutvalg:?\s*</strong>\s*<span class="data">((?s).*?)</span>`),
    "Vintage":regexp.MustCompile(`rgang:?\s*</strong>\s*<span class="data">((?s).*?)</span>`),
    "CC":regexp.MustCompile(`Land/distrikt:?\s*</strong>\s*<span class="data">((?s).*?)</span>`),
}

var decimalMap map[string]bool = map[string]bool{
    "Alcohol":true,
    "Sugar":true,
    "Acid":true,
}

var intMap map[string]bool = map[string]bool{
    "Vintage":true,
}







func DeepLookupRegex(url string,wr *WineRep){

    var ff float64
    var ctry,region string
    var ok1,ok2 bool
    var fint uint64
    resp, err := http.Get(url)
    defer func() {
        if err = resp.Body.Close(); err != nil {
            panic(err)
        }
    }()

    body, err := ioutil.ReadAll(resp.Body)
    if err!=nil{
        fmt.Println(err)
        return    
    }
    wr.Deeplookup=true
    for key,val:=range deepRegexes{
        mtch_:=val.FindSubmatch(body)
        var mtch string
        if len(mtch_)<2{
            fmt.Println(len(mtch_))
            continue        
        } else {
            mtch=strings.TrimSpace( html.UnescapeString(string(mtch_[1])) )
            fmt.Println(mtch)
        }
        _,ok1=decimalMap[key]
        _,ok2=intMap[key]

        if ok1{
            ff=fixDecimals(mtch)
            reflect.ValueOf(wr).Elem().FieldByName(key).SetFloat(ff)
        } else if ok2{
            fint,_=strconv.ParseUint(mtch,10,64)
            reflect.ValueOf(wr).Elem().FieldByName(key).SetUint(fint)
        } else if key=="CC" {
            ctry,region=countryRegion(mtch)
            wr.Country=ctry
            wr.Subcountry=region
        } else {
            reflect.ValueOf(wr).Elem().FieldByName(key).SetString(mtch)
        }

        //fmt.Println(key,strings.TrimSpace(mtch))

        

    }
    wr.LastWritten=time.Now()
    return
}



func EsScanForDeep(esConn chan *elastigo.Conn){
    searchJson := `{
        "query": {
            "filtered": {
                "filter": {
                    "term": {"Deeplookup":false}
                }
            }
        }
    }`
    sargs:=map[string]interface{} { "search_type" : `scan`, `scroll`:`20m`}
    c:= <-esConn
    sResp,err:=c.Search("wines","product",sargs,searchJson)
    if err!=nil{
        log.Println(err)
        return
    }
    esConn <- c

    tasks := make(chan WineRep,20000)

    wg:= new(sync.WaitGroup)

    for i := 0; i < 6; i++ {
        wg.Add(1)
        go UpdateWineRep(tasks,wg,esConn)
    }

    ScrollIterator(sResp.ScrollId,esConn,tasks)

    close(tasks)
    wg.Wait()
    close(esConn)
    
}

func TaskCreator(sResp elastigo.SearchResult,tasks chan WineRep){
    for _,hit:=range sResp.Hits.Hits{
        var wr WineRep
        err := json.Unmarshal(*hit.Source, &wr)
        if err!=nil{
            log.Println(err)
            continue
        }
        tasks <- wr
    }
}

func ScrollIterator(scroll_id string,esChan chan *elastigo.Conn,tasks chan WineRep){
    sargs:=map[string]interface{} { "search_type" : `scan`, `scroll`:`20m`}
    var err error
    var sResp elastigo.SearchResult
    var moreRes bool = true
    for moreRes{
        c:= <- esChan
        sResp,err = c.Scroll(sargs,scroll_id)
        esChan <- c
        if err!=nil{
            log.Println(err)
            moreRes=false
            break        
        }

        if sResp.Hits.Total==0{
            moreRes=false
            break
        }

        TaskCreator(sResp,tasks)
    }
}

func UpdateWineRep(tasks chan WineRep,wg *sync.WaitGroup,esChan chan *elastigo.Conn){
    defer wg.Done()
    for wr:=range tasks{
        DeepLookupRegex(wr.Url,&wr)
        c:= <- esChan
        c.Index("wines","product",strconv.FormatUint(wr.Prodnum,10),nil,wr)
        esChan <- c
    }
}


