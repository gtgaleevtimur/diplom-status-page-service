package support

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"statusPage/internal/entities"
	"sync"
)

func GetResultSupportData(wg *sync.WaitGroup) []int {
	out := make(chan []int)
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(out)
		sup := supportRequest()
		totalTicket := 0
		var load int
		for _, value := range sup {
			totalTicket += value.ActiveTickets
		}

		if totalTicket <= 9 {
			load = 1
		} else if totalTicket <= 16 {
			load = 2
		} else {
			load = 3
		}
		timeToRequest := float64(60) / float64(18)
		averageTimeRequest := float64(load) * timeToRequest
		result := []int{load, int(averageTimeRequest)}

		out <- result
	}()
	var result = <-out
	return result

}

func supportRequest() []entities.SupportData {
	var result []entities.SupportData
	resp, err := http.Get("http://127.0.0.1:8383/support")
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
			return []entities.SupportData{}
		}
	} else {
		return []entities.SupportData{}
	}
	return result
}
