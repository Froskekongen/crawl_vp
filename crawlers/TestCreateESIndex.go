package crawl_vp

import (
    "fmt"
    elastigo "github.com/mattbaird/elastigo/lib"
)


func TestIndex1(c *elastigo.Conn)bool{
    searchJson := `{
    "query": {
        "filtered": {
            "filter": {
                "term": { "Prodnum": 89101 }
            }
        }
    }
}`
    searchresponse, err := c.Search("wines", "product", nil, searchJson)
    if err!=nil{
        fmt.Println(err)
        return false
    }
    fmt.Println("Length of response:",len(searchresponse.Hits.Hits))
	if len(searchresponse.Hits.Hits) >= 1 {
		fmt.Println(string( *searchresponse.Hits.Hits[0].Source))
    }
    if len(searchresponse.Hits.Hits)==1{
        return true
    }
    return false
}

func TestIndex2(c *elastigo.Conn)bool{
    countResponse, err := c.Count("wines", "product", nil)
    if err!=nil{
        fmt.Println(err)
        return false
    }
    fmt.Println("Number of products found:",countResponse.Count)
    if countResponse.Count>2{
        return true
    }
    return false
}



