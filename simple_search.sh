#curl -XGET "http://${1}:9200/wines/_search?pretty" -d '{
#    "query": {
#        "match_all":{}
#    }
#}'

#curl -XGET "http://${1}:9200/wines/_search?pretty" -d '{
#        "query": {
#            "filtered": {
#                "filter": {
#                    "term": {"Deeplookup":true}
#                }
#            }
#        }
#    }'



curl -XGET "http://${1}:9200/wines/product/_search?pretty" -d "{
        \"query\": {
            \"term\":{\"Name\":\"${2}\"}
        }
    }"
