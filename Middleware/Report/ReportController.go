package report

import (
	conn "Amazon_Server/Config"
	Generic "Amazon_Server/Generic"
	"Amazon_Server/Model"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type StatsBasedOnMonth struct {
	ProductName   string   `json:"productName" bson:"productName"`
	ProductImage  []string `json:"productImage" bson:"productImage"`
	TotalUnitSold int      `json:"totalUnitSold" bson:"totalUnitSold"`
	AvergaePrice  float64  `json:"avergaePrice" bson:"avergaePrice"`

	January   int `json:"january" bson:"january"`
	February  int `json:"february" bson:"february"`
	March     int `json:"march" bson:"march"`
	April     int `json:"april" bson:"april"`
	May       int `json:"may" bson:"may"`
	June      int `json:"june" bson:"june"`
	July      int `json:"july" bson:"july"`
	August    int `json:"august" bson:"august"`
	September int `json:"september" bson:"september"`
	October   int `json:"october" bson:"october"`
	November  int `json:"november" bson:"november"`
	December  int `json:"december" bson:"december"`
}

func ViewReports(w http.ResponseWriter, r *http.Request) {
	Generic.SetupResponse(&w, r)
	w.Header().Set("Content-Type", "application/json")

	//taking the sellerid from the url
	var params = mux.Vars(r)
	ids := params["id"]
	sellerId, err := primitive.ObjectIDFromHex(ids)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//taking the statsBasedOn from body of postman
	data, err := ioutil.ReadAll(r.Body)
	asString := string(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	statsBasedOn := make(map[string]interface{})
	err = json.Unmarshal([]byte(asString), &statsBasedOn)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//findOne report using sellerId from report API
	var report Model.Report
	reportCollection := conn.ConnectDB("reports")
	reportFilter := bson.M{"sellerId": sellerId}
	errr := reportCollection.FindOne(context.TODO(), reportFilter).Decode(&report)
	if errr != nil {
		conn.GetError(errr, w)
		return
	}

	//looping through all the product to take report statsBasedOnMonth
	if statsBasedOn["statsBasedOn"].(string) == "month" || statsBasedOn["statsBasedOn"].(string) == "Month" {
		fmt.Fprintln(w, "Stats Based On Month For the Last One Year!!")
		for i := 0; i < len(report.SoldItems); i++ {
			var product StatsBasedOnMonth
			product.ProductName = report.SoldItems[i].ProductName
			product.ProductImage = report.SoldItems[i].ProductImage
			for j := 0; j < len(report.SoldItems[i].Quantity); j++ {
				product.TotalUnitSold += report.SoldItems[i].Quantity[j]
				product.AvergaePrice += float64(report.SoldItems[i].ProductPrice[j] * report.SoldItems[i].Quantity[j])

				monthCnt := report.SoldItems[i].DeliveryDate[j].Month()
				yearCnt := report.SoldItems[i].DeliveryDate[j].Year()
				currMonth := time.Now().Month()
				currYear := time.Now().Year()
				if monthCnt == 1 {
					if monthCnt <= currMonth && yearCnt == currYear {
						product.January += report.SoldItems[i].Quantity[j]
					} else if monthCnt > currMonth && yearCnt+1 == currYear {
						product.January += report.SoldItems[i].Quantity[j]
					}
				} else if monthCnt == 2 {
					if monthCnt <= currMonth && yearCnt == currYear {
						product.February += report.SoldItems[i].Quantity[j]
					} else if monthCnt > currMonth && yearCnt+1 == currYear {
						product.February += report.SoldItems[i].Quantity[j]
					}
				} else if monthCnt == 3 {
					if monthCnt <= currMonth && yearCnt == currYear {
						product.March += report.SoldItems[i].Quantity[j]
					} else if monthCnt > currMonth && yearCnt+1 == currYear {
						product.March += report.SoldItems[i].Quantity[j]
					}
				} else if monthCnt == 4 {
					if monthCnt <= currMonth && yearCnt == currYear {
						product.April += report.SoldItems[i].Quantity[j]
					} else if monthCnt > currMonth && yearCnt+1 == currYear {
						product.April += report.SoldItems[i].Quantity[j]
					}
				} else if monthCnt == 5 {
					if monthCnt <= currMonth && yearCnt == currYear {
						product.May += report.SoldItems[i].Quantity[j]
					} else if monthCnt > currMonth && yearCnt+1 == currYear {
						product.May += report.SoldItems[i].Quantity[j]
					}
				} else if monthCnt == 6 {
					if monthCnt <= currMonth && yearCnt == currYear {
						product.June += report.SoldItems[i].Quantity[j]
					} else if monthCnt > currMonth && yearCnt+1 == currYear {
						product.June += report.SoldItems[i].Quantity[j]
					}
				} else if monthCnt == 7 {
					if monthCnt <= currMonth && yearCnt == currYear {
						product.July += report.SoldItems[i].Quantity[j]
					} else if monthCnt > currMonth && yearCnt+1 == currYear {
						product.July += report.SoldItems[i].Quantity[j]
					}
				} else if monthCnt == 8 {
					if monthCnt <= currMonth && yearCnt == currYear {
						product.August += report.SoldItems[i].Quantity[j]
					} else if monthCnt > currMonth && yearCnt+1 == currYear {
						product.August += report.SoldItems[i].Quantity[j]
					}
				} else if monthCnt == 9 {
					if monthCnt <= currMonth && yearCnt == currYear {
						product.September += report.SoldItems[i].Quantity[j]
					} else if monthCnt > currMonth && yearCnt+1 == currYear {
						product.September += report.SoldItems[i].Quantity[j]
					}
				} else if monthCnt == 10 {
					if monthCnt <= currMonth && yearCnt == currYear {
						product.October += report.SoldItems[i].Quantity[j]
					} else if monthCnt > currMonth && yearCnt+1 == currYear {
						product.October += report.SoldItems[i].Quantity[j]
					}
				} else if monthCnt == 11 {
					if monthCnt <= currMonth && yearCnt == currYear {
						product.November += report.SoldItems[i].Quantity[j]
					} else if monthCnt > currMonth && yearCnt+1 == currYear {
						product.November += report.SoldItems[i].Quantity[j]
					}
				} else if monthCnt == 12 {
					if monthCnt <= currMonth && yearCnt == currYear {
						product.December += report.SoldItems[i].Quantity[j]
					} else if monthCnt > currMonth && yearCnt+1 == currYear {
						product.December += report.SoldItems[i].Quantity[j]
					}
				}
			}
			product.AvergaePrice = product.AvergaePrice / float64(product.TotalUnitSold)
			json.NewEncoder(w).Encode(product)
		}
	} else {
		fmt.Fprintln(w, "Stats Based On Year For the Last One Decade!!")
		for i := 0; i < len(report.SoldItems); i++ {
			// var product StatsBasedOnYear
			product := make(map[string]interface{})
			product["productName"] = report.SoldItems[i].ProductName
			product["productImage"] = report.SoldItems[i].ProductImage
			totalUnitSold := 0
			averagePrice := 0.0
			currYear := time.Now().Year()
			for k := currYear; k >= currYear-10; k-- {
				y := strconv.Itoa(k)
				product[y] = 0
			}
			for j := 0; j < len(report.SoldItems[i].Quantity); j++ {
				totalUnitSold += report.SoldItems[i].Quantity[j]
				averagePrice += float64(report.SoldItems[i].ProductPrice[j] * report.SoldItems[i].Quantity[j])

				yearCnt := report.SoldItems[i].DeliveryDate[j].Year()
				if yearCnt >= currYear-10 {
					year := strconv.Itoa(yearCnt)
					product[year] = product[year].(int) + report.SoldItems[i].Quantity[j]
				}
			}
			product["totalUnitSold"] = totalUnitSold
			averagePrice = averagePrice / float64(totalUnitSold)
			product["avergaePrice"] = averagePrice
			json.NewEncoder(w).Encode(product)
		}
	}
}

func ViewReportByProductId(w http.ResponseWriter, r *http.Request) {
	Generic.SetupResponse(&w, r)
	w.Header().Set("Content-Type", "application/json")

	//taking the sellerid and productId from the url
	var params = mux.Vars(r)
	ids := params["id"]
	sellerId, err := primitive.ObjectIDFromHex(ids)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	idx := params["productId"]
	productId, err := primitive.ObjectIDFromHex(idx)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//taking the statsBasedOn from body of postman
	data, err := ioutil.ReadAll(r.Body)
	asString := string(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	statsBasedOn := make(map[string]interface{})
	err = json.Unmarshal([]byte(asString), &statsBasedOn)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//findOne report using sellerId from report API
	var report Model.Report
	reportCollection := conn.ConnectDB("reports")
	reportFilter := bson.M{"sellerId": sellerId}
	errr := reportCollection.FindOne(context.TODO(), reportFilter).Decode(&report)
	if errr != nil {
		conn.GetError(errr, w)
		return
	}

	//looping through all the product to take report statsBasedOnMonth
	if statsBasedOn["statsBasedOn"].(string) == "month" || statsBasedOn["statsBasedOn"].(string) == "Month" {
		for i := 0; i < len(report.SoldItems); i++ {
			if report.SoldItems[i].ProductId == productId {
				fmt.Fprintln(w, "Stats Based On Month For the Last One Year for "+report.SoldItems[i].ProductName+"!!")
				var product StatsBasedOnMonth
				product.ProductName = report.SoldItems[i].ProductName
				product.ProductImage = report.SoldItems[i].ProductImage
				for j := 0; j < len(report.SoldItems[i].Quantity); j++ {
					product.TotalUnitSold += report.SoldItems[i].Quantity[j]
					product.AvergaePrice += float64(report.SoldItems[i].ProductPrice[j] * report.SoldItems[i].Quantity[j])

					monthCnt := report.SoldItems[i].DeliveryDate[j].Month()
					yearCnt := report.SoldItems[i].DeliveryDate[j].Year()
					currMonth := time.Now().Month()
					currYear := time.Now().Year()
					if monthCnt == 1 {
						if monthCnt <= currMonth && yearCnt == currYear {
							product.January += report.SoldItems[i].Quantity[j]
						} else if monthCnt > currMonth && yearCnt+1 == currYear {
							product.January += report.SoldItems[i].Quantity[j]
						}
					} else if monthCnt == 2 {
						if monthCnt <= currMonth && yearCnt == currYear {
							product.February += report.SoldItems[i].Quantity[j]
						} else if monthCnt > currMonth && yearCnt+1 == currYear {
							product.February += report.SoldItems[i].Quantity[j]
						}
					} else if monthCnt == 3 {
						if monthCnt <= currMonth && yearCnt == currYear {
							product.March += report.SoldItems[i].Quantity[j]
						} else if monthCnt > currMonth && yearCnt+1 == currYear {
							product.March += report.SoldItems[i].Quantity[j]
						}
					} else if monthCnt == 4 {
						if monthCnt <= currMonth && yearCnt == currYear {
							product.April += report.SoldItems[i].Quantity[j]
						} else if monthCnt > currMonth && yearCnt+1 == currYear {
							product.April += report.SoldItems[i].Quantity[j]
						}
					} else if monthCnt == 5 {
						if monthCnt <= currMonth && yearCnt == currYear {
							product.May += report.SoldItems[i].Quantity[j]
						} else if monthCnt > currMonth && yearCnt+1 == currYear {
							product.May += report.SoldItems[i].Quantity[j]
						}
					} else if monthCnt == 6 {
						if monthCnt <= currMonth && yearCnt == currYear {
							product.June += report.SoldItems[i].Quantity[j]
						} else if monthCnt > currMonth && yearCnt+1 == currYear {
							product.June += report.SoldItems[i].Quantity[j]
						}
					} else if monthCnt == 7 {
						if monthCnt <= currMonth && yearCnt == currYear {
							product.July += report.SoldItems[i].Quantity[j]
						} else if monthCnt > currMonth && yearCnt+1 == currYear {
							product.July += report.SoldItems[i].Quantity[j]
						}
					} else if monthCnt == 8 {
						if monthCnt <= currMonth && yearCnt == currYear {
							product.August += report.SoldItems[i].Quantity[j]
						} else if monthCnt > currMonth && yearCnt+1 == currYear {
							product.August += report.SoldItems[i].Quantity[j]
						}
					} else if monthCnt == 9 {
						if monthCnt <= currMonth && yearCnt == currYear {
							product.September += report.SoldItems[i].Quantity[j]
						} else if monthCnt > currMonth && yearCnt+1 == currYear {
							product.September += report.SoldItems[i].Quantity[j]
						}
					} else if monthCnt == 10 {
						if monthCnt <= currMonth && yearCnt == currYear {
							product.October += report.SoldItems[i].Quantity[j]
						} else if monthCnt > currMonth && yearCnt+1 == currYear {
							product.October += report.SoldItems[i].Quantity[j]
						}
					} else if monthCnt == 11 {
						if monthCnt <= currMonth && yearCnt == currYear {
							product.November += report.SoldItems[i].Quantity[j]
						} else if monthCnt > currMonth && yearCnt+1 == currYear {
							product.November += report.SoldItems[i].Quantity[j]
						}
					} else if monthCnt == 12 {
						if monthCnt <= currMonth && yearCnt == currYear {
							product.December += report.SoldItems[i].Quantity[j]
						} else if monthCnt > currMonth && yearCnt+1 == currYear {
							product.December += report.SoldItems[i].Quantity[j]
						}
					}
				}
				product.AvergaePrice = product.AvergaePrice / float64(product.TotalUnitSold)
				json.NewEncoder(w).Encode(product)
				break
			}
		}
	} else {
		for i := 0; i < len(report.SoldItems); i++ {
			if report.SoldItems[i].ProductId == productId {
				fmt.Fprintln(w, "Stats Based On Year For the Last One Decade for "+report.SoldItems[i].ProductName+"!!")
				product := make(map[string]interface{})
				product["productName"] = report.SoldItems[i].ProductName
				product["productImage"] = report.SoldItems[i].ProductImage
				totalUnitSold := 0
				averagePrice := 0.0
				currYear := time.Now().Year()
				for k := currYear; k >= currYear-10; k-- {
					y := strconv.Itoa(k)
					product[y] = 0
				}
				for j := 0; j < len(report.SoldItems[i].Quantity); j++ {
					totalUnitSold += report.SoldItems[i].Quantity[j]
					averagePrice += float64(report.SoldItems[i].ProductPrice[j] * report.SoldItems[i].Quantity[j])

					yearCnt := report.SoldItems[i].DeliveryDate[j].Year()
					if yearCnt >= currYear-10 {
						year := strconv.Itoa(yearCnt)
						product[year] = product[year].(int) + report.SoldItems[i].Quantity[j]
					}
				}
				product["totalUnitSold"] = totalUnitSold
				averagePrice = averagePrice / float64(totalUnitSold)
				product["avergaePrice"] = averagePrice
				json.NewEncoder(w).Encode(product)
				break
			}
		}
	}
}
