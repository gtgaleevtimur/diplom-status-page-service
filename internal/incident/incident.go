package incident

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"statusPage/internal/entities"
	"sync"
)

func GetResultIncidentData(wg *sync.WaitGroup) []entities.IncidentData {
	out := make(chan []entities.IncidentData)
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(out)

		data := incidentDataCollection()

		sort.Slice(data, func(i, j int) bool {
			return data[i].Status < data[j].Status
		})

		out <- data
	}()
	var result = <-out
	return result
}

func incidentDataCollection() []entities.IncidentData {
	var result []entities.IncidentData
	resp, err := http.Get("http://127.0.0.1:8383/accendent")
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
			return []entities.IncidentData{}
		}
	} else {
		return []entities.IncidentData{}
	}
	return result
}
