#curl -XGET 'http://172.30.31.203:9200/wines/_search?search_type=scan&scroll=5m' -d '{
#    "query": {
#        "match_all":{}
#    }
#}'

#curl -XGET 'http://172.30.31.203:9200/_search/scroll?scroll=5m&pretty' -d 'c2NhbjszOzEyNDM2OjMxcjNucmlBUmJ1Zkhvc2VmczgtRFE7MTI0Mzc6MzFyM25yaUFSYnVmSG9zZWZzOC1EUTsxMjQzODozMXIzbnJpQVJidWZIb3NlZnM4LURROzE7dG90YWxfaGl0czo0ODA7'


curl -XGET 'http://172.30.31.203:9200/wines/_search?scroll=5m' -d '{
        "query": {
            "term":{"Name":"tripel"}
        }
    }'


