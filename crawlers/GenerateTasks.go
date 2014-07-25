package crawl_vp

import(
    "fmt"
    "strconv"
)

func GenerateTasks(mapPages map[string]int,tasks chan string,maxPPerType int){
    // generate some tasks
    base1:="http://www.vinmonopolet.no/vareutvalg/sok?query=*&sort=2&sortMode=0&page="
    base2:="&filterIds=25&filterValues="
    var ss string
    kkk:=0
    for key,val:=range mapPages{
        for iii:=1;iii<=val;iii++{
            ss=base1+strconv.Itoa(iii)+base2+key
            kkk++
            fmt.Println(kkk,ss)
            tasks <- ss
            if iii>maxPPerType{break}
        }
    }
}
