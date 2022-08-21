package email

import (
	"io/ioutil"
	"log"
	"os"
	"sort"
	"statusPage/internal/entities"
	"statusPage/internal/sms"
	"strconv"
	"strings"
	"sync"
)

//Функция обработки результирующих данных о системе эмейл,возвращающая топ 3 самых быстрых и топ 3 самых медленных провайдеров

func GetResultEmailData(path string, wg *sync.WaitGroup) map[string][][]entities.EmailData {
	out := make(chan map[string][][]entities.EmailData)
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(out)

		result := make(map[string][][]entities.EmailData)
		data := emailDataReader(path)
		countrySlice := []string{}
		for _, value := range data {
			countrySlice = append(countrySlice, value.Country)
		}
		countrySlice = removeDuplicates(countrySlice)
		for _, valueCountry := range countrySlice {
			var fastProviders []entities.EmailData
			var slowProviders []entities.EmailData
			for _, value := range data {
				if value.Country == valueCountry {
					fastProviders = append(fastProviders, value)
					slowProviders = append(slowProviders, value)
				}
			}
			country := sms.CountryFromAlpa(valueCountry)
			sort.Slice(fastProviders, func(i, j int) bool {
				return fastProviders[i].DeliveryTime < fastProviders[j].DeliveryTime
			})
			sort.Slice(slowProviders, func(i, j int) bool {
				return slowProviders[i].DeliveryTime > slowProviders[j].DeliveryTime
			})
			result[country] = [][]entities.EmailData{
				fastProviders[:3],
				slowProviders[:3],
			}
		}
		out <- result
	}()
	var result = <-out
	return result
}

//Функция для удаления дубликатов стран,итерируемся по срезу и создаем другой срез уже без дубликатов.

func removeDuplicates(slice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, value := range slice {
		if _, ok := keys[value]; !ok {
			keys[value] = true
			list = append(list, value)
		}
	}
	return list
}

//Функция сбора данных эмейл системы.

func emailDataReader(path string) []entities.EmailData {
	var result []entities.EmailData
	file, err := os.Open(path)
	if err != nil {
		log.Fatal("Cannot open email file:", err)
	}
	defer file.Close()

	reader, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal("Cannot read email file:", err)
	}
	lines := strings.Split(string(reader), "\n")
	for _, value := range lines {
		splitVal := strings.Split(value, ";")
		if len(splitVal) == 3 {
			if checkEmailData(splitVal) {
				deliveryValue, _ := strconv.Atoi(splitVal[2])
				res := entities.EmailData{
					Country:      splitVal[0],
					Provider:     splitVal[1],
					DeliveryTime: deliveryValue,
				}
				result = append(result, res)
			}
		}
	}
	return result
}

//Функция проверки на валидность полученных данных.

func checkEmailData(value []string) bool {
	if value[0] == sms.CountryAlpha2()[value[0]] {
		providers := map[string]string{"Gmail": "Gmail", "Yahoo": "Yahoo", "Hotmail": "Hotmail", "MSN": "MSN", "Orange": "Orange", "Comcast": "Comcast", "AOL": "AOL", "Live": "Live", "RediffMail": "RediffMail", "GMX": "GMX", "Protonmail": "Protonmail", "Yandex": "Yandex", "Mail.ru": "Mail.ru"}
		if value[1] == providers[value[1]] {
			_, err := strconv.Atoi(value[2])
			if err != nil {
				log.Printf("Значение ответа доставки письма ms %v не соответсвует ожидаемому.", value)
				return false
			} else {
				return true
			}
		} else {
			log.Printf("Значение провайдера %v не соответсвует ожидаемому.", value)
			return false
		}
	} else {
		log.Printf("Значение страны alpha-2 %v не соответсвует ожидаемому.", value)
		return false
	}
}
