// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"html/template"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"regexp"
    elastigo "github.com/mattbaird/elastigo/lib"
    "fmt"
    "encoding/json"
    "github.com/Froskekongen/crawl_vp/crawlers"
    //"strconv"
    "strings"
    "os"
)

var (
	addr = flag.Bool("addr", false, "find open address and print to final-port.txt")
)

type Page struct {
	Title string
	Body  []byte
}

func (p *Page) save() error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}





func searchHandler(w http.ResponseWriter, r *http.Request) {
    
//	if err != nil {
//		http.Redirect(w, r, "/search/", http.StatusFound)
//		return
//	}
    err := templates.ExecuteTemplate(w, "search.html", "Search")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
    http.Redirect(w, r, "/results/", http.StatusFound)
}


var eshost *string = flag.String("eshost", "localhost", "Elasticsearch Server Host Address")
func searchResultHandler(w http.ResponseWriter, r *http.Request) {  

    queries:=make([]string,0,5)
    stringMustMatch:=make(map[string]string,10)
    stringShouldMatch:=make(map[string]string,10)
    //w.Header().Set("Content-Type", "application/json")
    if r.FormValue("Name")!=""{
        stringMustMatch["Name"]=r.FormValue("Name")
    }
    if r.FormValue("Producer")!=""{
        stringMustMatch["Producer"]=r.FormValue("Producer")
    }
    if r.FormValue("Country")!=""{
        stringMustMatch["Country"]=r.FormValue("Country")
    }
    if r.FormValue("Subcountry")!=""{
        stringMustMatch["Subcountry"]=r.FormValue("Subcountry")
    }
    if r.FormValue("WineType")!=""{
        stringMustMatch["WineType"]=r.FormValue("WineType")
    }
    
    if r.FormValue("Material")!=""{
        stringShouldMatch["Material"]=r.FormValue("Material")
    }
    
    var tString string
    if len(stringMustMatch)>0{
        tString=""
        for key,val:=range stringMustMatch{
            tString=tString+fmt.Sprintf(`{ "match": { "%v": "%v" }},`,key,val)
        }
        tString=`"must":[`+tString[:len(tString)-1]+"]"
        queries=append(queries,tString)
    }

    if len(stringShouldMatch)>0{
        tString=""
        for key,val:=range stringShouldMatch{
            tString=tString+fmt.Sprintf(`{ "match": { "%v": "%v" }},`,key,val)
        }
        tString=`"should":[`+tString[:len(tString)-1]+"]"
        queries=append(queries,tString)
    }
    
    c:= elastigo.NewConn()
    
    c.Domain = *eshost
    sargs:=map[string]interface{} {`scroll`:`20m`}

//    searchJson:=`{
//    "query":
//        {
//            "bool": {
//                "must":     [{ "match": { "Name": "%v" }},{ "match": { "Producer": "%v" }}]
//            }
//        }
//    }`
    
    query:=strings.Join(queries,",\n")
    searchJson:=`{
    "query":
        {
            "bool": {
                    %v
            }
        }
    }`
    searchJson = fmt.Sprintf(searchJson,query)
    log.Println("ESQuery:"+"("+searchJson+")")

    sResp,err:=c.Search("wines","product",sargs,searchJson)
    if err!=nil{
        log.Println(err)
    }

    Nhits:=0
    wrs:=make([]crawl_vp.WineRep,0,300)
    var wrsTemp []crawl_vp.WineRep    
    for{
        Nhits=Nhits+len(sResp.Hits.Hits)
        //fmt.Println("Hits:",sResp.Hits.Total)
        if err!=nil{
            log.Println(err)
            break
        }
        if len(sResp.Hits.Hits)==0{
            break
        }
        wrsTemp=MakeWrs(&sResp.Hits)
        wrs=append(wrs,wrsTemp...)
        sResp,err = c.Scroll(sargs,sResp.ScrollId)
    }

    outStr:=ParseResponse(wrs)
    outStr = fmt.Sprintf(`<font size="4">N hits: %v</font>`,Nhits) +"<br></br>"+outStr

//    var wr crawl_vp.WineRep
//    json.Unmarshal(*sResp.Hits.Hits[0].Source, &wr)
//    OneResp,_:=json.MarshalIndent(wr,"","  ")
    //hitJson := json.MarshalIndent(sResp.Hits.Hits[:].Source,"","  ")
    fmt.Fprint(w, outStr)
//	if err != nil {
//		http.Redirect(w, r, "/search/", http.StatusFound)
//		return
//	}
//    err := templates.ExecuteTemplate(w, "search_results.html", data)
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//	}
}


func MakeWrs(rp *elastigo.Hits) []crawl_vp.WineRep{
    wrs:=make([]crawl_vp.WineRep,len(rp.Hits))
    for iii,hh :=range rp.Hits{
        json.Unmarshal(*hh.Source,&wrs[iii])
    }
    return wrs
}




func ParseResponse(wrs []crawl_vp.WineRep) string{
//    wrs:=make([]crawl_vp.WineRep,len(rp.Hits))
//    for iii,hh :=range rp.Hits{
//        json.Unmarshal(*hh.Source,&wrs[iii])
//    }
//    bb,_:=json.MarshalIndent(wrs,"","  ")
    str:=""
    var temp string    
    for _,wr := range wrs{
        temp=fmt.Sprintf(`<font size="5"><a href="%v+?HideDropdownIfNotInStock=true&ShowShopsWithProdInStock=true" STYLE="text-decoration: none" target="_blank"> %v</a></font>`,wr.Url,wr.Name)
        str=str+temp+"<br></br>"
        str = str + fmt.Sprintf(`<font size="4">Produsent: %v</font>`,wr.Producer) +"<br></br>"
        str = str + fmt.Sprintf(`<font size="4">Land: %v</font>`,wr.Country) +"<br></br>"
        str = str + fmt.Sprintf(`<font size="4">Region: %v</font>`,wr.Subcountry) +"<br></br>"
        str = str + fmt.Sprintf(`<font size="4">Årgang: %v</font>`,wr.Vintage) +"<br></br>"
        str = str + fmt.Sprintf(`<font size="4">Pris: %v</font>`,wr.Price) +"<br></br>"
        str = str + fmt.Sprintf(`<font size="4">Alkohol: %v, Syre: %v g/l, Sukker: %v g/l</font>`,wr.Alcohol,wr.Acid,wr.Sugar) +"<br></br>"
        str = str + fmt.Sprintf(`<font size="4">Produkttype: %v</font>`,wr.WineType) +"<br></br>"
        if wr.Deeplookup{
            str = str + fmt.Sprintf(`<font size="4">Råstoff: %v</font>`,wr.Material) +"<br></br>"
        }

        if wr.Obsoleteproduct{
            str=str + `<font size="4">Utgått fra sortimentet.</font><br></br>`
        }
        if wr.Soldout{
            str=str + `<font size="4">Utsolgt fra leverandør.</font><br></br>`
        }

        str = str + fmt.Sprintf(`<font size="4"><a href="https://www.google.com/search?q=%v+site:cellartracker.com" target="_blank">Search on cellartracker</a></font>`,wr.Name) +fmt.Sprintf(` <font size="4"><a href="https://www.google.com/search?q=%v+%v+site:cellartracker.com" target="_blank">(with producer)</a></font>`,wr.Name,wr.Producer) +"\t---\t"
        str = str + fmt.Sprintf(`<font size="4"><a href="http://www.vivino.com/search?q=%v" target="_blank">Search on vivino</a></font>`,wr.Name) +"<br></br>"
        str = str + fmt.Sprintf(`<font size="4"><a href="https://www.google.com/search?q=%v+site:ratebeer.com" target="_blank">Search on RateBeer</a></font>`,wr.Name) +"\t---\t"
        str = str + fmt.Sprintf(`<font size="4"><a href="https://www.google.com/search?q=%v+site:beeradvocate.com" target="_blank">Search on BeerAdvocate</a></font>`,wr.Name) +"<br></br>"


        str=str+"<br></br><br></br>"
    }
    return str
}


func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

var templates = template.Must(template.ParseFiles("edit.html", "view.html","search.html","search_results.html"))

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}

func main() {
    ff,err:= os.OpenFile("log.txt", os.O_RDWR|os.O_APPEND, 0660)
    if err!=nil{
        log.Println(err)
        ff,err=os.Create("log.txt")
        if err!=nil{
            log.Println(err)
            return
        }
    }
    log.SetOutput(ff)
	flag.Parse()
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))
    http.HandleFunc("/search/", searchHandler)
    http.HandleFunc("/results/", searchResultHandler)

	if *addr {
		l, err := net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			log.Fatal(err)
		}
		err = ioutil.WriteFile("final-port.txt", []byte(l.Addr().String()), 0644)
		if err != nil {
			log.Fatal(err)
		}
		s := &http.Server{}
		s.Serve(l)
		return
	}

	http.ListenAndServe("localhost:8080", nil)
}
