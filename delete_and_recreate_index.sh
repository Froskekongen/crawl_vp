curl -XDELETE "http://${1}:9200/wines"

curl -XPOST "http://${1}:9200/wines" -d '{
    "settings" : {
        "number_of_shards" : 3,
        "number_of_replicas" : 0
    }
}'

curl -XPUT "http://${1}:9200/wines/product/_mapping" -d '{
    "product" : {
        "properties" : {
            "Url" : {"type" : "string", "store" : true },
            "Name" : {"type" : "string", "store" : true },
            "WineType" : {"type" : "string", "store" : true },
            "Producer" : {"type" : "string", "store" : true },
            "Wholesaler" : {"type" : "string", "store" : true },
            "Material" : {"type" : "string", "store" : true },
            "Country" : {"type" : "string", "store" : true },
            "Subcountry" : {"type" : "string", "store" : true },
            "Store" : {"type" : "string", "store" : true },
            "Distributor" : {"type" : "string", "store" : true },

            "Prodnum" : {"type" : "integer", "store" : true },
            "Vintage" : {"type" : "integer", "store" : true },
            "Price" : {"type" : "integer", "store" : true },

            "Alcohol" : {"type" : "float", "store" : true },
            "Sugar" : {"type" : "float", "store" : true },
            "Acid" : {"type" : "float", "store" : true },

            "Soldout" : {"type" : "boolean", "store" : true },
            "Obsoleteproduct" : {"type" : "boolean", "store" : true },
            "Deeplookup" : {"type" : "boolean", "store" : true },

            "LookupTimes" : {"type" : "date", "store" : true }
        }
    }
}'
