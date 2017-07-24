package monitor

import (
	"github.com/hhy5861/logrus"
	"golang.org/x/net/context"
	"gopkg.in/olivere/elastic.v5"
)

var (
	err    error
	client *elastic.Client
)

func New(addres string) {
	client, err = elastic.NewClient(elastic.SetURL(addres))
	if err != nil {
		var ps logrus.Params
		logrus.Fatal(ps, err, "client elastic error")
	}
}

func selectCount(index string, query elastic.Query) (result *elastic.SearchResult, err error) {
	countService := client.Search().Index(index)
	ctx := context.Background()
	result, err = countService.Query(query).Pretty(true).Sort("time", false).Do(ctx)
	return
}
