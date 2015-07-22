package elasticsearch

import (
	"github.com/bitly/go-simplejson"
	elastigo "github.com/mattbaird/elastigo/lib"
	"strings"
)

type ESearch struct {
	ElasClient *elastigo.Conn
	Hosts      string
}

func NewESearch(hosts string) *ESearch {
	es := &ESearch{
		Hosts: hosts,
	}

	es.ElasClient = elastigo.NewConn()
	strs := strings.Split(hosts, ":")
	if len(strs) == 2 {
		es.ElasClient.Domain = strs[0]
		es.ElasClient.Port = strs[1]
	} else {
		es.ElasClient.Domain = strs[0]
	}

	return es
}

func (self *ESearch) Fetch(index, doc_type, search_json string) ([]simplejson.Json, error) {
	var jsons []simplejson.Json

	esout, err := self.ElasClient.Search(index, doc_type, nil, search_json)
	if err != nil {
		return jsons, err
	}

	for i := 0; i < len(esout.Hits.Hits); i++ {
		sjson, err := simplejson.NewJson([]byte(*(esout.Hits.Hits[i].Source)))
		if err != nil {
			return jsons, err
		}
		jsons = append(jsons, *sjson)
	}

	return jsons, nil
}

func (self *ESearch) Write(index, doc_type, doc_id, write_json string) error {
	_, err := self.ElasClient.Index(index, doc_type, doc_id, nil, []byte(write_json))
	if err != nil {
		return err
	}

	return nil
}
