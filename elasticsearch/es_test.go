package elasticsearch

import (
	"log"
	"testing"
)

func TestFetch(t *testing.T) {
	es := NewESearch("127.0.0.1:9200")

	index := "test_info"
	doc_type := "external"

	str := `{
			  "query": {
			    "bool": {
			      "must": [
			        {
			          "query_string": {
			            "default_field": "external.id_number",
			            "query": "320981198007010001"
			          }
			        }
			      ]
			    }
			  },
			  "from": 0,
			  "size": 10
			}`
	data, err := es.Fetch(index, doc_type, str)
	log.Println("TestFetch: data:", data, ", err:", err)
	if data == nil || err != nil {
		t.Error()
	}

	for _, src := range data {
		//debug
		str, err := src.MarshalJSON()
		log.Println("simplejson|", string(str), "|err|", err)
	}
}
