package crawl_vp

import(
    "time"
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
    Selection string

    //ints
    Prodnum uint64 //`json:"prodnum"`
    Vintage uint16 //`json:"vintage"`
    Price []uint64 //`json:"price"` // change to uint64

    //floats
    Alcohol float64 //`json:"alcohol"`
    Sugar float64 //`json:"sugar"`
    Acid float64 //`json:"acid"`
    LastRelativePriceDifference float64
    

    //Flags
    Soldout bool //`json:"soldout"`
    Obsoleteproduct bool //`json:"obsoleteproduct"`
    Deeplookup bool //`json:"deeplokup"`

    //Datetimes
    LookupTimes []time.Time
    LastWritten time.Time
}


func (wr *WineRep)UpdatePrice(px uint64){
//    timeN:=len(wr.LookupTimes)
//    priceN:=len(wr.Price)
    N1:=len(wr.LookupTimes)
    N2:=len(wr.Price)
    if N1==0{
        wr.LookupTimes=[]time.Time{time.Now()}    
    } else {
        wr.LookupTimes=append(wr.LookupTimes,time.Now())
    }

    if N2==0{
       wr.Price=[]uint64{px} 
    } else {
        wr.Price=append(wr.Price,px)
    }
    if len(wr.Price)>1{
        wr.LastRelativePriceDifference=float64(wr.Price[len(wr.Price)-1])/float64(wr.Price[len(wr.Price)-2])
    }
}


type ListOfWines []WineRep


func (lv ListOfWines)Len() int{
    return len(lv)
}

func (lv ListOfWines)Less(i,j int) bool{
    return lv[i].LastRelativePriceDifference < lv[j].LastRelativePriceDifference
}



func (lv ListOfWines)Swap(i,j int){
    lv[i],lv[j]=lv[j],lv[i]
}






