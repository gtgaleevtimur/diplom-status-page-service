package voice

import (
	"io/ioutil"
	"log"
	"os"
	"statusPage/internal/entities"
	"statusPage/internal/sms"
	"strconv"
	"strings"
	"sync"
)

func VoiceCallReader(path string, wg *sync.WaitGroup) []entities.VoiceCallData {
	out := make(chan []entities.VoiceCallData)
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(out)
		file, err := os.Open(path)
		if err != nil {
			log.Fatal("Cannot open voice file:", err)
		}
		defer file.Close()
		reader, err := ioutil.ReadAll(file)
		if err != nil {
			log.Fatal("Cannot read voice file:", err)
		}
		var resultIn []entities.VoiceCallData
		lines := strings.Split(string(reader), "\n")
		for _, value := range lines {
			splitVal := strings.Split(value, ";")
			if len(splitVal) == 8 {
				if checkVoiceCall(splitVal) {
					percentValue, _ := strconv.Atoi(splitVal[1])
					responseTime, _ := strconv.Atoi(splitVal[2])
					connVal, _ := strconv.ParseFloat(splitVal[4], 32)
					connVal32 := float32(connVal)
					ttfbVal, _ := strconv.Atoi(splitVal[5])
					voicePurVal, _ := strconv.Atoi(splitVal[6])
					medianVoicVal, _ := strconv.Atoi(splitVal[7])
					res := entities.VoiceCallData{
						Country:             splitVal[0],
						Bandwidth:           percentValue,
						ResponseTime:        responseTime,
						Provider:            splitVal[3],
						ConnectionStability: connVal32,
						TTFB:                ttfbVal,
						VoicePurity:         voicePurVal,
						MedianOfCallsTime:   medianVoicVal,
					}
					resultIn = append(resultIn, res)
				}
			}
		}
		out <- resultIn
	}()
	var result = <-out
	return result
}

func checkVoiceCall(value []string) bool {
	if value[0] == sms.CountryAlpha2()[value[0]] {
		percentValue, err := strconv.Atoi(value[1])
		if err != nil {
			log.Printf("Значение пропускной способности канала %v не соответсвует ожидаемому.", value)
			return false
		}
		if -1 < percentValue && percentValue < 101 {
			_, err := strconv.Atoi(value[2])
			if err == nil {
				providers := map[string]string{"TransparentCalls": "TransparentCalls", "E-Voice": "E-Voice", "JustPhone": "JustPhone"}
				if value[3] == providers[value[3]] {
					_, err := strconv.ParseFloat(value[4], 32)
					if err != nil {
						log.Printf("Значение стабильности соединения %v не соответсвует ожидаемому.", value)
						return false
					} else {
						_, err := strconv.Atoi(value[5])
						if err == nil {
							_, err := strconv.Atoi(value[6])
							if err == nil {
								_, err := strconv.Atoi(value[7])
								if err == nil {
									return true
								} else {
									log.Printf("Значение медианы звонка %v не соответсвует ожидаемому.", value)
									return false
								}
							} else {
								log.Printf("Значение чистоты связи %v не соответсвует ожидаемому.", value)
								return false
							}
						} else {
							log.Printf("Значение TTFB %v не соответсвует ожидаемому.", value)
							return false
						}
					}
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
