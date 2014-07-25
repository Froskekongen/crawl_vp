package crawl_vp

import (
    "code.google.com/p/go.net/html"
    "io"
    "fmt"
)

func parseHtml(r io.Reader) {
    d := html.NewTokenizer(r)
    for { 
        // token type
        tokenType := d.Next() 
        if tokenType == html.ErrorToken {
            return     
        }       
        token := d.Token()
        
        switch tokenType {
            case html.StartTagToken: // <tag>
                // type Token struct {
                //     Type     TokenType
                //     DataAtom atom.Atom
                //     Data     string
                //     Attr     []Attribute
                // }
                //
                // type Attribute struct {
                //     Namespace, Key, Val string
                // }
            case html.TextToken: fmt.Println(token) // text between start and end tag
            case html.EndTagToken: // </tag>
            case html.SelfClosingTagToken: // <tag/>

        }
    }
}
