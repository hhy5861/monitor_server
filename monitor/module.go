package monitor

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/bitly/go-simplejson"
	"gitlab.pnlyy.com/monitor_server/email"
	"gitlab.pnlyy.com/monitor_server/model"
	"gitlab.pnlyy.com/monitor_server/utils"
	"gopkg.in/olivere/elastic.v5"
	"strings"
)

var (
	intervalTime   = int64(300)
	funcList       map[string]strategy
	remarkSendTime = make(map[string]int64)
)

type (
	data      []map[string]interface{}
	strategy  func(strategy *model.Strategy)
	whereType map[string]interface{}
)

func Select(strategy []*model.Strategy) {
	for _, v := range strategy {
		key := fmt.Sprintf("%s_%d", v.Op, v.Type)
		funcObj, ok := funcList[key]
		if ok {
			go funcObj(v)
		}
	}
}

func init() {
	funcList = map[string]strategy{
		">_1": func(strategy *model.Strategy) {
			var where = make(whereType)
			where["token"] = strategy.Token
			where["module"] = strategy.Module
			where["point"] = strategy.Point

			if strategy.Level >= 0 {
				where["level"] = strategy.Level
			}

			delete(where, "time")
			key := utils.GetMD5Hash(where)
			t := time.Now().Unix() - intervalTime
			if t >= remarkSendTime[key] {
				bq := elastic.NewBoolQuery()
				bq = bq.Must(elastic.NewTermQuery("point", strategy.Point))
				bq = bq.Must(elastic.NewTermQuery("level", strategy.Level))
				bq = bq.Must(elastic.NewTermQuery("module", strategy.Module))

				ranges := int64(strategy.Range)
				whereTime := time.Now().Unix() - ranges
				if strings.EqualFold(strategy.Field2, "now") {
					bq = bq.Must(elastic.NewRangeQuery(strategy.Field1).Gt(whereTime))
				} else if strategy.Field1 != "" && strategy.Field2 != "" {
					bq = bq.Must(elastic.NewRangeQuery(strategy.Field1).Gt(strategy.Field2))
					if ranges > 1 {
						bq = bq.Must(elastic.NewRangeQuery("time").Gt(whereTime))
					}
				} else if strategy.Field1 == "" && strategy.Field2 == "" {
					bq = bq.Must(elastic.NewRangeQuery("time").Gt(whereTime))
				}

				index := fmt.Sprintf("%s-*", strategy.Token)
				result, err := selectCount(index, bq)
				//fmt.Printf("op: > limit: %d count: %d \n", int64(strategy.Limit), result.Hits.TotalHits)
				if err == nil {
					if int64(strategy.Limit) <= result.Hits.TotalHits {
						var msg *email.Message
						for _, hit := range result.Hits.Hits {
							err := json.Unmarshal(*hit.Source, &msg)
							if err == nil {
								break
							}
						}

						msg.Limit = result.Hits.TotalHits
						go email.NewMail(msg).SendEmail(1, strategy.Group_id)
						remarkSendTime[key] = time.Now().Unix()
					}
				}
			}
		},

		"<_1": func(strategy *model.Strategy) {
			var where = make(whereType)
			where["token"] = strategy.Token
			where["module"] = strategy.Module
			where["point"] = strategy.Point

			if strategy.Level >= 0 {
				where["level"] = strategy.Level
			}

			delete(where, "time")
			key := utils.GetMD5Hash(where)
			t := time.Now().Unix() - intervalTime
			if t >= remarkSendTime[key] {
				bq := elastic.NewBoolQuery()
				bq = bq.Must(elastic.NewTermQuery("point", strategy.Point))
				bq = bq.Must(elastic.NewTermQuery("level", strategy.Level))
				bq = bq.Must(elastic.NewTermQuery("module", strategy.Module))

				ranges := int64(strategy.Range)
				whereTime := time.Now().Unix() - ranges
				if strings.EqualFold(strategy.Field2, "now") {
					bq = bq.Must(elastic.NewRangeQuery(strategy.Field1).Lt(whereTime))
				} else if strategy.Field1 != "" && strategy.Field2 != "" {
					bq = bq.Must(elastic.NewRangeQuery(strategy.Field1).Lt(strategy.Field2))
					if ranges > 1 {
						bq = bq.Must(elastic.NewRangeQuery("time").Lt(whereTime))
					}
				} else if strategy.Field1 == "" && strategy.Field2 == "" {
					bq = bq.Must(elastic.NewRangeQuery("time").Lt(whereTime))
				}

				index := fmt.Sprintf("%s-*", strategy.Token)
				result, err := selectCount(index, bq)
				if err == nil {
					if int64(strategy.Limit) <= result.Hits.TotalHits {
						var msg *email.Message
						for _, hit := range result.Hits.Hits {
							err := json.Unmarshal(*hit.Source, &msg)
							if err == nil {
								break
							}
						}

						msg.Limit = result.Hits.TotalHits
						go email.NewMail(msg).SendEmail(1, strategy.Group_id)
						remarkSendTime[key] = time.Now().Unix()
					}
				}
			}
		},

		">=_1": func(strategy *model.Strategy) {
			var where = make(whereType)
			where["token"] = strategy.Token
			where["module"] = strategy.Module
			where["point"] = strategy.Point

			if strategy.Level >= 0 {
				where["level"] = strategy.Level
			}

			delete(where, "time")
			key := utils.GetMD5Hash(where)
			t := time.Now().Unix() - intervalTime
			if t >= remarkSendTime[key] {
				bq := elastic.NewBoolQuery()
				bq = bq.Must(elastic.NewTermQuery("point", strategy.Point))
				bq = bq.Must(elastic.NewTermQuery("level", strategy.Level))
				bq = bq.Must(elastic.NewTermQuery("module", strategy.Module))

				ranges := int64(strategy.Range)
				whereTime := time.Now().Unix() - ranges
				if strings.EqualFold(strategy.Field2, "now") {
					bq = bq.Must(elastic.NewRangeQuery(strategy.Field1).Gte(whereTime))
				} else if strategy.Field1 != "" && strategy.Field2 != "" {
					bq = bq.Must(elastic.NewRangeQuery(strategy.Field1).Gte(strategy.Field2))
					if ranges > 1 {
						bq = bq.Must(elastic.NewRangeQuery("time").Gte(whereTime))
					}
				} else if strategy.Field1 == "" && strategy.Field2 == "" {
					bq = bq.Must(elastic.NewRangeQuery("time").Gte(whereTime))
				}

				index := fmt.Sprintf("%s-*", strategy.Token)
				result, err := selectCount(index, bq)
				if err == nil {
					if int64(strategy.Limit) <= result.Hits.TotalHits {
						var msg *email.Message
						for _, hit := range result.Hits.Hits {
							err := json.Unmarshal(*hit.Source, &msg)
							if err == nil {
								break
							}
						}

						msg.Limit = result.Hits.TotalHits
						go email.NewMail(msg).SendEmail(1, strategy.Group_id)
						remarkSendTime[key] = time.Now().Unix()
					}
				}
			}
		},

		"<=_1": func(strategy *model.Strategy) {
			var where = make(whereType)
			where["token"] = strategy.Token
			where["module"] = strategy.Module
			where["point"] = strategy.Point

			if strategy.Level >= 0 {
				where["level"] = strategy.Level
			}

			delete(where, "time")
			key := utils.GetMD5Hash(where)
			t := time.Now().Unix() - intervalTime
			if t >= remarkSendTime[key] {
				bq := elastic.NewBoolQuery()
				bq = bq.Must(elastic.NewTermQuery("point", strategy.Point))
				bq = bq.Must(elastic.NewTermQuery("level", strategy.Level))
				bq = bq.Must(elastic.NewTermQuery("module", strategy.Module))

				ranges := int64(strategy.Range)
				whereTime := time.Now().Unix() - ranges
				if strings.EqualFold(strategy.Field2, "now") {
					bq = bq.Must(elastic.NewRangeQuery(strategy.Field1).Lte(whereTime))
				} else if strategy.Field1 != "" && strategy.Field2 != "" {
					bq = bq.Must(elastic.NewRangeQuery(strategy.Field1).Lte(strategy.Field2))
					if ranges > 1 {
						bq = bq.Must(elastic.NewRangeQuery("time").Lte(whereTime))
					}
				} else if strategy.Field1 == "" && strategy.Field2 == "" {
					bq = bq.Must(elastic.NewRangeQuery("time").Lte(whereTime))
				}

				index := fmt.Sprintf("%s-*", strategy.Token)
				result, err := selectCount(index, bq)
				if err == nil {
					if int64(strategy.Limit) <= result.Hits.TotalHits {
						var msg *email.Message
						for _, hit := range result.Hits.Hits {
							err := json.Unmarshal(*hit.Source, &msg)
							if err == nil {
								break
							}
						}

						msg.Limit = result.Hits.TotalHits
						go email.NewMail(msg).SendEmail(1, strategy.Group_id)
						remarkSendTime[key] = time.Now().Unix()
					}
				}
			}
		},

		"+_1": func(strategy *model.Strategy) {
			var where = make(whereType)
			where["token"] = strategy.Token
			where["module"] = strategy.Module
			where["point"] = strategy.Point

			if strategy.Level >= 0 {
				where["level"] = strategy.Level
			}

			delete(where, "time")
			key := utils.GetMD5Hash(where)
			t := time.Now().Unix() - intervalTime
			if t >= remarkSendTime[key] {
				bq := elastic.NewBoolQuery()
				bq = bq.Must(elastic.NewTermQuery("point", strategy.Point))
				bq = bq.Must(elastic.NewTermQuery("level", strategy.Level))
				bq = bq.Must(elastic.NewTermQuery("module", strategy.Module))

				ranges := int64(strategy.Range)
				whereTime := time.Now().Unix() - ranges
				bq = bq.Must(elastic.NewRangeQuery("time").Gte(whereTime))

				index := fmt.Sprintf("%s-*", strategy.Token)
				result, err := selectCount(index, bq)
				if err == nil {
					if int64(strategy.Limit) <= result.Hits.TotalHits {
						var fieldOne, fieldTwo string
						fieldOneArray := strings.Split(strategy.Field1, ".")
						fieldTwoArray := strings.Split(strategy.Field2, ".")

						if len(fieldOneArray) > 1 {
							fieldOne = fieldOneArray[1]
						} else {
							fieldOne = fieldOneArray[0]
						}

						if len(fieldTwoArray) > 1 {
							fieldTwo = fieldTwoArray[1]
						} else {
							fieldTwo = fieldTwoArray[0]
						}

						for _, hit := range result.Hits.Hits {
							var msg *email.Message
							err := json.Unmarshal(*hit.Source, &msg)
							if err == nil {
								jsonMessage, err := json.Marshal(msg.Message)
								if err == nil {
									proJson, err := simplejson.NewJson([]byte(jsonMessage))
									if err == nil {
										valueOne, errOne := proJson.Get(fieldOne).Int()
										valueTwo, errTwo := proJson.Get(fieldTwo).Int()
										if errOne == nil && errTwo == nil {
											if (valueOne + valueTwo) >= strategy.Limit {
												msg.Limit = result.Hits.TotalHits
												go email.NewMail(msg).SendEmail(1, strategy.Group_id)
												remarkSendTime[key] = time.Now().Unix()
											}
										}
									}
								}
							}
						}
					}
				}
			}
		},

		"-_1": func(strategy *model.Strategy) {
			var where = make(whereType)
			where["token"] = strategy.Token
			where["module"] = strategy.Module
			where["point"] = strategy.Point

			if strategy.Level >= 0 {
				where["level"] = strategy.Level
			}

			delete(where, "time")
			key := utils.GetMD5Hash(where)
			t := time.Now().Unix() - intervalTime
			if t >= remarkSendTime[key] {
				bq := elastic.NewBoolQuery()
				bq = bq.Must(elastic.NewTermQuery("point", strategy.Point))
				bq = bq.Must(elastic.NewTermQuery("level", strategy.Level))
				bq = bq.Must(elastic.NewTermQuery("module", strategy.Module))

				ranges := int64(strategy.Range)
				whereTime := time.Now().Unix() - ranges
				bq = bq.Must(elastic.NewRangeQuery("time").Gte(whereTime))

				index := fmt.Sprintf("%s-*", strategy.Token)
				result, err := selectCount(index, bq)
				if err == nil {
					if int64(strategy.Limit) <= result.Hits.TotalHits {
						var fieldOne, fieldTwo string
						fieldOneArray := strings.Split(strategy.Field1, ".")
						fieldTwoArray := strings.Split(strategy.Field2, ".")

						if len(fieldOneArray) > 1 {
							fieldOne = fieldOneArray[1]
						} else {
							fieldOne = fieldOneArray[0]
						}

						if len(fieldTwoArray) > 1 {
							fieldTwo = fieldTwoArray[1]
						} else {
							fieldTwo = fieldTwoArray[0]
						}

						for _, hit := range result.Hits.Hits {
							var msg *email.Message
							err := json.Unmarshal(*hit.Source, &msg)
							if err == nil {
								jsonMessage, err := json.Marshal(msg.Message)
								if err == nil {
									proJson, err := simplejson.NewJson([]byte(jsonMessage))
									if err == nil {
										valueOne, errOne := proJson.Get(fieldOne).Int()
										valueTwo, errTwo := proJson.Get(fieldTwo).Int()
										if errOne == nil && errTwo == nil {
											if (valueOne - valueTwo) >= strategy.Limit {
												msg.Limit = result.Hits.TotalHits
												go email.NewMail(msg).SendEmail(1, strategy.Group_id)
												remarkSendTime[key] = time.Now().Unix()
												break
											}
										}
									}
								}
							}
						}
					}
				}
			}
		},

		"×_1": func(strategy *model.Strategy) {
			var where = make(whereType)
			where["token"] = strategy.Token
			where["module"] = strategy.Module
			where["point"] = strategy.Point

			if strategy.Level >= 0 {
				where["level"] = strategy.Level
			}

			delete(where, "time")
			key := utils.GetMD5Hash(where)
			t := time.Now().Unix() - intervalTime
			if t >= remarkSendTime[key] {
				bq := elastic.NewBoolQuery()
				bq = bq.Must(elastic.NewTermQuery("point", strategy.Point))
				bq = bq.Must(elastic.NewTermQuery("level", strategy.Level))
				bq = bq.Must(elastic.NewTermQuery("module", strategy.Module))

				ranges := int64(strategy.Range)
				whereTime := time.Now().Unix() - ranges
				bq = bq.Must(elastic.NewRangeQuery("time").Gte(whereTime))

				index := fmt.Sprintf("%s-*", strategy.Token)
				result, err := selectCount(index, bq)
				if err == nil {
					if int64(strategy.Limit) <= result.Hits.TotalHits {
						var fieldOne, fieldTwo string
						fieldOneArray := strings.Split(strategy.Field1, ".")
						fieldTwoArray := strings.Split(strategy.Field2, ".")

						if len(fieldOneArray) > 1 {
							fieldOne = fieldOneArray[1]
						} else {
							fieldOne = fieldOneArray[0]
						}

						if len(fieldTwoArray) > 1 {
							fieldTwo = fieldTwoArray[1]
						} else {
							fieldTwo = fieldTwoArray[0]
						}

						for _, hit := range result.Hits.Hits {
							var msg *email.Message
							err := json.Unmarshal(*hit.Source, &msg)
							if err == nil {
								jsonMessage, err := json.Marshal(msg.Message)
								if err == nil {
									proJson, err := simplejson.NewJson([]byte(jsonMessage))
									if err == nil {
										valueOne, errOne := proJson.Get(fieldOne).Int()
										valueTwo, errTwo := proJson.Get(fieldTwo).Int()
										if errOne == nil && errTwo == nil {
											if (valueOne * valueTwo) >= strategy.Limit {
												msg.Limit = result.Hits.TotalHits
												go email.NewMail(msg).SendEmail(1, strategy.Group_id)
												remarkSendTime[key] = time.Now().Unix()
											}
										}
									}
								}
							}
						}
					}
				}
			}
		},

		"÷_1": func(strategy *model.Strategy) {
			var where = make(whereType)
			where["token"] = strategy.Token
			where["module"] = strategy.Module
			where["point"] = strategy.Point

			if strategy.Level >= 0 {
				where["level"] = strategy.Level
			}

			delete(where, "time")
			key := utils.GetMD5Hash(where)
			t := time.Now().Unix() - intervalTime
			if t >= remarkSendTime[key] {
				bq := elastic.NewBoolQuery()
				bq = bq.Must(elastic.NewTermQuery("point", strategy.Point))
				bq = bq.Must(elastic.NewTermQuery("level", strategy.Level))
				bq = bq.Must(elastic.NewTermQuery("module", strategy.Module))

				ranges := int64(strategy.Range)
				whereTime := time.Now().Unix() - ranges
				bq = bq.Must(elastic.NewRangeQuery("time").Gte(whereTime))

				index := fmt.Sprintf("%s-*", strategy.Token)
				result, err := selectCount(index, bq)
				if err == nil {
					if int64(strategy.Limit) <= result.Hits.TotalHits {
						var fieldOne, fieldTwo string
						fieldOneArray := strings.Split(strategy.Field1, ".")
						fieldTwoArray := strings.Split(strategy.Field2, ".")

						if len(fieldOneArray) > 1 {
							fieldOne = fieldOneArray[1]
						} else {
							fieldOne = fieldOneArray[0]
						}

						if len(fieldTwoArray) > 1 {
							fieldTwo = fieldTwoArray[1]
						} else {
							fieldTwo = fieldTwoArray[0]
						}

						for _, hit := range result.Hits.Hits {
							var msg *email.Message
							err := json.Unmarshal(*hit.Source, &msg)
							if err == nil {
								jsonMessage, err := json.Marshal(msg.Message)
								if err == nil {
									proJson, err := simplejson.NewJson([]byte(jsonMessage))
									if err == nil {
										valueOne, errOne := proJson.Get(fieldOne).Int()
										valueTwo, errTwo := proJson.Get(fieldTwo).Int()
										if errOne == nil && errTwo == nil {
											if (valueOne / valueTwo) >= strategy.Limit {
												msg.Limit = result.Hits.TotalHits
												go email.NewMail(msg).SendEmail(1, strategy.Group_id)
												remarkSendTime[key] = time.Now().Unix()
											}
										}
									}
								}
							}
						}
					}
				}
			}
		},
	}
}
