curl -XGET 'http://172.30.31.203:9200/wines/_search?pretty' -d '{
    "query": {
        "match_all":{}
    }
}'
