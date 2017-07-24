package monitor

import (
	"encoding/json"
	"fmt"
	"github.com/bitly/go-simplejson"
	"go/parser"
	"reflect"
)

var (
	markTeble = map[int]string{
		34: "&&",
		40: "<",
		13: "-",
		41: ">",
		35: "||",
		12: "+",
		14: "*",
		16: "%",
		45: "<=",
		46: ">=",
	}
)

func FormulaAnalyze(formula string) {
	var formulaNumber int
	tr, _ := parser.ParseExpr(formula)
	jsonCode, _ := json.Marshal(tr)
	fmt.Println("jsonCode", string(jsonCode))
	simpleJsonData, err := simplejson.NewJson(jsonCode)
	if err == nil {
		resultMap, _ := simpleJsonData.Map()
		for k, v := range resultMap {
			f := reflect.ValueOf(v)
			switch f.Kind().String() {
			case "map":
				mapData, err := simpleJsonData.Get(k).Map()
				if err == nil {
					for _, value := range mapData {
						ff := reflect.ValueOf(value)
						switch ff.Kind().String() {
						case "map":
							formulaNumber++
						}
					}
				}
			}

		}

		fmt.Println("formulaNumber", formulaNumber)
		if formulaNumber == 1 {
			var value string
			var op, opIn int
			var field1, field2 string
			n := 0
			for k, v := range resultMap {
				op, _ = simpleJsonData.Get("Op").Int()
				inData, _ := simpleJsonData.Get(k).Map()
				f := reflect.ValueOf(v)
				if f.Kind().String() == "map" {
					fmt.Println("v", v)
					for kk, vv := range inData {
						ff := reflect.ValueOf(vv)
						if ff.Kind().String() == "map" {
							field1, _ = simpleJsonData.Get(k).Get(kk).Get("X").Get("Name").String()
							field2, _ = simpleJsonData.Get(k).Get(kk).Get("Y").Get("Name").String()
							opIn, _ = simpleJsonData.Get(k).Get(kk).Get("Op").Int()
						}

						if kk == "Value" {
							value, _ = simpleJsonData.Get(k).Get("Value").String()
						}
					}
				}
			}

			fmt.Println("==========================================", n)
			fmt.Println("Op", op)
			fmt.Println("value", value)
			fmt.Println("opIn", opIn)
			fmt.Println("field1", field1)
			fmt.Println("field2", field2)

		}

		if formulaNumber == 2 {
			var value, inValue string
			var op, opIn, inNumber int
			var field1 string
			for k, v := range resultMap {
				op, _ = simpleJsonData.Get("Op").Int()
				inData, _ := simpleJsonData.Get(k).Map()
				f := reflect.ValueOf(v)
				if f.Kind().String() == "map" {
					for _, vv := range inData {
						ff := reflect.ValueOf(vv)
						if ff.Kind().String() == "map" {
							inNumber++
						}
					}
					fmt.Println("inData", inData)
					if inNumber == 1 {
						for kk, vv := range inData {
							ff := reflect.ValueOf(vv)
							if ff.Kind().String() == "map" {
								field1, _ = simpleJsonData.Get(k).Get(kk).Get("X").Get("Name").String()
								if opIn == 0 {
									opIn, _ = simpleJsonData.Get(k).Get(kk).Get("Op").Int()
								}

								if inValue == "" {
									inValue, _ = simpleJsonData.Get(k).Get(kk).Get("Y").Get("Value").String()
								}
							}
						}
					}

					if inNumber == 2 {

						for _, vvv := range inData {
							//fmt.Println("vvv", vvv)
							ff := reflect.ValueOf(vvv)
							if ff.Kind().String() == "map" {
							}
						}
					}
				}

				if k == "Y" {
					value, _ = simpleJsonData.Get(k).Get("Name").String()
				}
			}

			fmt.Println("======================================================================")
			fmt.Println("Op", op)
			fmt.Println("value", value)
			fmt.Println("opIn", opIn)
			fmt.Println("field1", field1)
			fmt.Println("field2", inValue)
		}
	}
}
