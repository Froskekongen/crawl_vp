curl -XPUT 'http://172.30.31.203:9200/wines/product/89101' -d '{
    "Prodnum" : 89101,
    "Producer" : "Lanson Pére et Fils",
    "Price" : [369]
}'

sleep 2s
#curl -XGET 'http://172.30.31.203:9200/wines/_mapping/product'

curl -XGET 'http://172.30.31.203:9200/wines/_search' -d '{
    "query": {
        "filtered": {
            "filter": {
                "term": { "Prodnum": 89101 }
            }
        }
    }
}'


curl -XDELETE 'http://172.30.31.203:9200/wines/product/89101'
