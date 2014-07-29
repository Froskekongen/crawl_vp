package main

import(
    "net/http"
    "io"
    "os"
    "fmt"
    "encoding/csv"
)


func errf(err error){
    if err!=nil{
        fmt.Println(err)
    }
}

func main(){

    out, err := os.Create("vmp.csv")
    defer out.Close()
    errf(err)

    resp, err := http.Get("http://www.vinmonopolet.no/api/produkter")
    defer resp.Body.Close()
    errf(err)
    _, err = io.Copy(out, resp.Body)
    errf(err)
    
    ff,err := os.Open("vmp.csv")
    defer ff.Close()
    
    csv_reader:=csv.NewReader(ff)
    csv_reader.Comma=';'
    ll,_:=csv_reader.Read()
    fmt.Println(ll)
    var line []string
    for {
        line,err = csv_reader.Read()
        if err==io.EOF{break}
        fmt.Println(line)
    }
}
