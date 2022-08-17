package mms

import (
	"encoding/json"
	"github.com/jinzhu/copier"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"statusPage/internal/entities"
	"statusPage/internal/sms"
	"strconv"
	"sync"
)

func GetResultMMSData(wg *sync.WaitGroup) [][]entities.MMSData {
	out := make(chan [][]entities.MMSData)
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(out)
		first := mmsRequest()
		second := []entities.MMSData{}
		err := copier.Copy(&second, &first)
		if err != nil {
			log.Print(err)
		}

		sort.Slice(first, func(i, j int) bool {
			return first[i].Provider < first[j].Provider
		})

		sort.Slice(second, func(i, j int) bool {
			return second[i].Country < second[j].Country
		})
		result := [][]entities.MMSData{
			first, second,
		}
		out <- result
	}()
	var result = <-out
	return result
}

func mmsRequest() []entities.MMSData {
	var result []entities.MMSData
	resp, err := http.Get("http://127.0.0.1:8383/mms")
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println(err)
		}
		err = json.Unmarshal(body, &result)
		if err != nil {
			log.Println(err)
			return []entities.MMSData{}
		}
		i := 0
		for _, value := range result {
			if checkMMS(value) {
				value.Country = sms.CountryFromAlpa(value.Country)
				result[i] = value
				i++
			}
		}
		result = result[:i]
		return result
	} else {
		return []entities.MMSData{}
	}
}

func checkMMS(value entities.MMSData) bool {
	if value.Country == sms.CountryAlpha2()[value.Country] {
		percentValue, err := strconv.Atoi(value.Bandwidth)
		if err != nil {
			log.Printf("Значение пропускной способности канала %v не соответсвует ожидаемому.", value)
			return false
		}
		if -1 < percentValue && percentValue < 101 {
			_, err := strconv.Atoi(value.ResponseTime)
			if err == nil {
				providers := map[string]string{"Topolo": "Topolo", "Rond": "Rond", "Kildy": "Kildy"}
				if value.Provider == providers[value.Provider] {
					return true
				} else {
					log.Printf("Значение провайдера %v не соответсвует ожидаемому.", value)
					return false
				}
			} else {
				log.Printf("Значение ответа в ms %v не соответсвует ожидаемому.", value)
				return false
			}
		} else {
			log.Printf("Значение пропускной способности канала %v не соответсвует ожидаемому.", value)
			return false
		}
	} else {
		log.Printf("Значение страны alpha-2 %v не соответсвует ожидаемому.", value)
		return false
	}
}
