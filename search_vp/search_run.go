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


var eshost *string = flag.String("host", "172.30.31.203", "Elasticsearch Server Host Address")
func searchResultHandler(w http.ResponseWriter, r *http.Request) {  

    //w.Header().Set("Content-Type", "application/json")
    Name := r.FormValue("Name")
    c:= elastigo.NewConn()
    
    c.Domain = *eshost
    sargs:=map[string]interface{} {`scroll`:`20m`}
    searchJson:=`{
        "query": {
            "match":{"Name":"%v"}
        },
        "size":100
    }`
    searchJson = fmt.Sprintf(searchJson,Name)

//    sResp,err:=c.Search("wines","product",nil,searchJson)
//    if err!=nil{
//        log.Println(err)
//    }

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
        temp=fmt.Sprintf(`<font size="5">Name:<a href="%v+?HideDropdownIfNotInStock=true&ShowShopsWithProdInStock=true" STYLE="text-decoration: none"> %v</a></font>`,wr.Url,wr.Name)
        str=str+temp+"<br></br>"
        str = str + fmt.Sprintf(`<font size="4">Producer: %v</font>`,wr.Producer) +"<br></br>"
        str = str + fmt.Sprintf(`<font size="4">Country: %v</font>`,wr.Country) +"<br></br>"
        str = str + fmt.Sprintf(`<font size="4">Region: %v</font>`,wr.Subcountry) +"<br></br>"
        str = str + fmt.Sprintf(`<font size="4">Vintage: %v</font>`,wr.Vintage) +"<br></br>"
        str = str + fmt.Sprintf(`<font size="4">Price: %v</font>`,wr.Price) +"<br></br>"
        str = str + fmt.Sprintf(`<font size="4">Alcohol: %v, Acid: %v g/l, Sugar: %v g/l</font>`,wr.Alcohol,wr.Acid,wr.Sugar) +"<br></br>"

        str = str + fmt.Sprintf(`<font size="4"><a href="https://www.google.com/search?q=%v+site:cellartracker.com">Search on cellartracker</a></font>`,wr.Name) +fmt.Sprintf(` <font size="4"><a href="https://www.google.com/search?q=%v+%v+site:cellartracker.com">(with producer)</a></font>`,wr.Name,wr.Producer) +"\t---\t"
        str = str + fmt.Sprintf(`<font size="4"><a href="http://www.vivino.com/search?q=%v+site:cellartracker.com">Search on vivino</a></font>`,wr.Name) +"<br></br>"
        str = str + fmt.Sprintf(`<font size="4"><a href="https://www.google.com/search?q=%v+site:ratebeer.com">Search on RateBeer</a></font>`,wr.Name) +"\t---\t"
        str = str + fmt.Sprintf(`<font size="4"><a href="https://www.google.com/search?q=%v+site:beeradvocate.com">Search on BeerAdvocate</a></font>`,wr.Name) +"<br></br>"


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

	http.ListenAndServe(":8080", nil)
}
