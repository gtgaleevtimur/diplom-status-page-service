package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const minResponseTime = 30
const maxResponseTime = 2000

const minConnectionStability = 600
const maxConnectionStability = 1000

const minVoicePurity = 0
const maxVoicePurity = 92

const minVoiceCallMedian = 3
const maxVoiceCallMedian = 60

const minTTFB = 2
const maxTTFB = 980

const minBandwidth = 0
const maxBandwidth = 100

const minEmailDeliveryTime = 0
const maxEmailDeliveryTime = 600

const smsFilename = "sms.data"
const mmsApiUrl = "http://localhost:8282/mms" // to params
const voiceFilename = "voice.data"
const emailFilename = "email.data"
const billingFilename = "billing.data"
const supportApiUrl = "http://localhost:8282/support"
const accendentListFilename = "accendents.data"

var firstSMSRowForCorrupt int
var secondSMSRowForCorrupt int

var firstVoiceRowForCorrupt int
var secondVoiceRowForCorrupt int

var firstEmailRowForCorrupt int
var secondEmailRowForCorrupt int

var MMSCollection []MMSItem
var SupportCollection []SupportItem
var AccendentCollection []AccendentItem

type MMSItem struct {
	Country      string `json:"country"`
	Provider     string `json:"provider"`
	Bandwidth    string `json:"bandwidth"`
	ResponseTime string `json:"response_time"`
}

type SupportItem struct {
	Topic         string `json:"topic"`
	ActiveTickets int    `json:"active_tickets"`
}

type AccendentItem struct {
	Topic  string `json:"topic"`
	Status string `json:"status"`
}

const accendentStatusActive = "active"
const accendentStatusClosed = "closed"

var AccendentTopics = []string{
	"sms delivery in EU",
	"MMS connection stability",
	"Voice call connection purity",
	"Checkout page is down",
	"Support overload",
	"Buy phone number not working in US",
	"API Slow latency",
}

func init() {
	rand.Seed(time.Now().UnixNano())

	firstSMSRowForCorrupt = rand.Intn(70)
	fmt.Printf("First sms row for currupt %d\n", firstSMSRowForCorrupt+1)

	secondSMSRowForCorrupt = rand.Intn(90)
	fmt.Printf("Second sms row for currupt %d\n", secondSMSRowForCorrupt+1)

	firstVoiceRowForCorrupt = rand.Intn(70)
	fmt.Printf("First Voice row for currupt %d\n", firstVoiceRowForCorrupt+1)

	secondVoiceRowForCorrupt = rand.Intn(90)
	fmt.Printf("Second Voice row for currupt %d\n", secondVoiceRowForCorrupt+1)

	firstEmailRowForCorrupt = rand.Intn(70)
	fmt.Printf("First Email row for currupt %d\n", firstEmailRowForCorrupt+1)

	secondEmailRowForCorrupt = rand.Intn(90)
	fmt.Printf("Second Email row for currupt %d\n", secondEmailRowForCorrupt+1)
}

func main() {
	shuffleSmsData()

	MMSCollection = shuffleMMSData()

	shuffleVoiceData()
	shuffleEmailData()
	shuffleBillingData()

	SupportCollection = shuffleSupportData()
	AccendentCollection = shuffleAccendentData()

	listenAndServeHTTP()
}

func shuffleSmsData() {
	var data string
	for i, country := range getCountriesList() {
		row := strings.Join([]string{
			country,
			getRandomBandwidthInString(),
			getRandomResponseTimeInString(),
			getSmsProviderByCountry(country),
		}, ";") + "\n"

		if i == firstSMSRowForCorrupt || i == secondSMSRowForCorrupt {
			row = strings.Replace(row, ";", "", rand.Intn(4))
			row = strings.Replace(row, "R", "", rand.Intn(3))
			row = strings.Replace(row, "C", "", rand.Intn(3))

			fmt.Println("sms row corrupted")
		}

		data += row
	}

	err := ioutil.WriteFile(getFilapathByFilename(smsFilename), []byte(data), 0644)
	if err != nil {
		fmt.Printf("Error in write sms data: %s", err.Error())
	}
}

func shuffleMMSData() []MMSItem {
	data := make([]MMSItem, 0)
	for _, country := range getCountriesList() {
		data = append(
			data,
			MMSItem{
				Country:      country,
				Provider:     getMMSProviderByCountry(country),
				Bandwidth:    getRandomBandwidthInString(),
				ResponseTime: getRandomResponseTimeInString(),
			},
		)
	}

	return data
}

func shuffleVoiceData() {
	var data string
	for i, country := range getCountriesList() {
		row := strings.Join([]string{
			country,
			getRandomBandwidthInString(),
			getRandomResponseTimeInString(),
			getVoiceCallProviderByCountry(country),
			getRandomConnectionStability(),
			getRandomTTFB(),
			getRandomVoicePurity(),
			getRandomMedianOfCallsTime(),
		}, ";") + "\n"

		if i == firstVoiceRowForCorrupt || i == secondVoiceRowForCorrupt {
			row = strings.Replace(row, ";", "", rand.Intn(4))
			row = strings.Replace(row, "R", "", rand.Intn(3))
			row = strings.Replace(row, "C", "", rand.Intn(3))

			fmt.Println("Voice row corrupted")
		}

		data += row
	}

	err := ioutil.WriteFile(getFilapathByFilename(voiceFilename), []byte(data), 0644)
	if err != nil {
		fmt.Printf("Error in write sms data: %s", err.Error())
	}
}

func shuffleEmailData() {
	var data string
	providersList := getEmailProvidersList()
	i := 0
	for _, country := range getCountriesList() {
		for _, provider := range providersList {
			row := strings.Join([]string{
				country,
				provider,
				getRandomEmailDeliveryTime(),
			}, ";") + "\n"

			if i == firstEmailRowForCorrupt || i == secondEmailRowForCorrupt {
				row = strings.Replace(row, ";", "", rand.Intn(4))
				row = strings.Replace(row, "A", "", rand.Intn(3))
				row = strings.Replace(row, "a", "", rand.Intn(3))
				row = strings.Replace(row, "O", "", rand.Intn(3))
				row = strings.Replace(row, "o", "", rand.Intn(3))
				row = strings.Replace(row, "M", "", rand.Intn(3))
				row = strings.Replace(row, "m", "", rand.Intn(3))
				row = strings.Replace(row, "P", "", rand.Intn(3))
				row = strings.Replace(row, "p", "", rand.Intn(3))

				fmt.Println("Email row corrupted")
			}

			data += row
			i++
		}
	}

	err := ioutil.WriteFile(getFilapathByFilename(emailFilename), []byte(data), 0644)
	if err != nil {
		fmt.Printf("Error in write email data: %s", err.Error())
	}
}

func shuffleBillingData() {
	data := ""
	for i := 0; i < 6; i++ {
		value := getRandomIntBetweenValues(0, 150)
		if value > 50 {
			value = 1
		} else {
			value = 0
		}

		data = data + fmt.Sprintf("%d", value)
		// create customer
		// purchase
		// payout
		// recurring
		// fraud control
		// checkout page
	}

	err := ioutil.WriteFile(getFilapathByFilename(billingFilename), []byte(data), 0644)
	if err != nil {
		fmt.Printf("Error in write sms data: %s", err.Error())
	}
}

func shuffleSupportData() []SupportItem {
	data := make([]SupportItem, 0)
	for _, topic := range getSupportTopicsList() {
		data = append(data, SupportItem{Topic: topic, ActiveTickets: getRandomSupportTickets()})
	}

	return data
}

func shuffleAccendentData() []AccendentItem {
	collection := make([]AccendentItem, 0)
	status := ""
	for _, topic := range AccendentTopics {
		if getRandomIntBetweenValues(0, 1) == 1 {
			status = accendentStatusActive
		} else {
			status = accendentStatusClosed
		}

		collection = append(collection, AccendentItem{Topic: topic, Status: status})
	}

	return collection
}

func getCountriesList() []string {
	return []string{"RU", "US", "GB", "FR", "BL", "AT", "BG", "DK", "CA", "ES", "CH", "TR", "PE", "NZ", "MC"}
}

func getSmsProviderByCountry(country string) string {
	smsProviderMap := map[string]string{
		"RU": "Topolo",
		"US": "Rond",
		"GB": "Topolo",
		"FR": "Topolo",
		"BL": "Kildy",
		"AT": "Topolo",
		"BG": "Rond",
		"DK": "Topolo",
		"CA": "Rond",
		"ES": "Topolo",
		"CH": "Topolo",
		"TR": "Rond",
		"PE": "Topolo",
		"NZ": "Kildy",
		"MC": "Kildy",
	}

	return smsProviderMap[country]
}

func getMMSProviderByCountry(country string) string {
	smsProviderMap := map[string]string{
		"RU": "Topolo",
		"US": "Rond",
		"GB": "Topolo",
		"FR": "Topolo",
		"BL": "Kildy",
		"AT": "Topolo",
		"BG": "Rond",
		"DK": "Topolo",
		"CA": "Rond",
		"ES": "Topolo",
		"CH": "Topolo",
		"TR": "Rond",
		"PE": "Topolo",
		"NZ": "Kildy",
		"MC": "Kildy",
	}

	return smsProviderMap[country]
}

func getVoiceCallProviderByCountry(country string) string {
	voiceProviderMap := map[string]string{
		"RU": "TransparentCalls",
		"US": "E-Voice",
		"GB": "TransparentCalls",
		"FR": "TransparentCalls",
		"BL": "E-Voice",
		"AT": "TransparentCalls",
		"BG": "E-Voice",
		"DK": "JustPhone",
		"CA": "JustPhone",
		"ES": "E-Voice",
		"CH": "JustPhone",
		"TR": "TransparentCalls",
		"PE": "JustPhone",
		"NZ": "JustPhone",
		"MC": "E-Voice",
	}

	return voiceProviderMap[country]
}

func getEmailProvidersList() []string {
	return []string{
		"Gmail",
		"Yahoo",
		"Hotmail",
		"MSN",
		"Orange",
		"Comcast",
		"AOL",
		"Live",
		"RediffMail",
		"GMX",
		"Protonmail",
		"Yandex",
		"Mail.ru",
	}
}

func getSupportTopicsList() []string {
	return []string{
		"sms",
		"MMS",
		"Email",
		"Billing",
		"Create account",
		"API",
		"Marketing",
		"Privacy",
		"GDPR",
		"Other",
	}
}

func getRandomSupportTickets() int {
	return getRandomIntBetweenValues(0, 8)
}

func getFilapathByFilename(filename string) string {
	return "" + filename
}

func getRandomBandwidthInString() string {
	return strconv.Itoa(getRandomIntBetweenValues(minBandwidth, maxBandwidth))
}

func getRandomResponseTimeInString() string {
	return strconv.Itoa(getRandomIntBetweenValues(minResponseTime, maxResponseTime))
}

func getRandomConnectionStability() string {
	stability := getRandomIntBetweenValues(minConnectionStability, maxConnectionStability)

	return fmt.Sprintf("%.2f", float32(stability)/1000)
}

func getRandomTTFB() string {
	return strconv.Itoa(getRandomIntBetweenValues(minTTFB, maxTTFB))
}

func getRandomVoicePurity() string {
	return strconv.Itoa(getRandomIntBetweenValues(minVoicePurity, maxVoicePurity))
}

func getRandomMedianOfCallsTime() string {
	return strconv.Itoa(getRandomIntBetweenValues(minVoiceCallMedian, maxVoiceCallMedian))
}

func getRandomEmailDeliveryTime() string {
	return strconv.Itoa(getRandomIntBetweenValues(minEmailDeliveryTime, maxEmailDeliveryTime))
}

func getRandomIntBetweenValues(min int, max int) int {
	return rand.Intn(max-min) + min
}

func listenAndServeHTTP() {
	router := mux.NewRouter()

	router.HandleFunc("/mms", handleMMS)
	router.HandleFunc("/support", handleSupport)
	router.HandleFunc("/accendent", handleAccendent)
	router.HandleFunc("/test", handleTest).Methods("GET", "OPTIONS")

	http.ListenAndServe("127.0.0.1:8383", router)
}

func handleMMS(w http.ResponseWriter, r *http.Request) {
	response(w, r, MMSCollection)
}

func handleSupport(w http.ResponseWriter, r *http.Request) {
	response(w, r, SupportCollection)
}

func handleAccendent(w http.ResponseWriter, r *http.Request) {
	response(w, r, AccendentCollection)
}

func handleTest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Write([]byte("{\n  \"status\": true,\n  \"data\": {\n    \"sms\": [\n      [\n        {\n          \"country\": \"Canada\",\n          \"bandwidth\": \"12\",\n          \"response_time\": \"67\",\n          \"provider\": \"Rond\"\n        },\n        {\n          \"country\": \"Great Britain\",\n          \"bandwidth\": \"98\",\n          \"response_time\": \"593\",\n          \"provider\": \"Kildy\"\n        },\n        {\n          \"country\": \"Russian Federation\",\n          \"bandwidth\": \"77\",\n          \"response_time\": \"1734\",\n          \"provider\": \"Topolo\"\n        }\n      ],\n      [\n        {\n          \"country\": \"Great Britain\",\n          \"bandwidth\": \"98\",\n          \"response_time\": \"593\",\n          \"provider\": \"Kildy\"\n        },\n        {\n          \"country\": \"Canada\",\n          \"bandwidth\": \"12\",\n          \"response_time\": \"67\",\n          \"provider\": \"Rond\"\n        },\n        {\n          \"country\": \"Russian Federation\",\n          \"bandwidth\": \"77\",\n          \"response_time\": \"1734\",\n          \"provider\": \"Topolo\"\n        }\n      ]\n    ],\n    \"mms\": [\n      [\n        {\n          \"country\": \"Great Britain\",\n          \"bandwidth\": \"98\",\n          \"response_time\": \"593\",\n          \"provider\": \"Kildy\"\n        },\n        {\n          \"country\": \"Canada\",\n          \"bandwidth\": \"12\",\n          \"response_time\": \"67\",\n          \"provider\": \"Rond\"\n        },\n        {\n          \"country\": \"Russian Federation\",\n          \"bandwidth\": \"77\",\n          \"response_time\": \"1734\",\n          \"provider\": \"Topolo\"\n        }\n      ],\n      [\n        {\n          \"country\": \"Canada\",\n          \"bandwidth\": \"12\",\n          \"response_time\": \"67\",\n          \"provider\": \"Rond\"\n        },\n        {\n          \"country\": \"Great Britain\",\n          \"bandwidth\": \"98\",\n          \"response_time\": \"593\",\n          \"provider\": \"Kildy\"\n        },\n        {\n          \"country\": \"Russian Federation\",\n          \"bandwidth\": \"77\",\n          \"response_time\": \"1734\",\n          \"provider\": \"Topolo\"\n        }\n      ]\n    ],\n    \"voice_call\": [\n      {\n        \"country\": \"US\",\n        \"bandwidth\": \"53\",\n        \"response_time\": \"321\",\n        \"provider\": \"TransparentCalls\",\n        \"connection_stability\": 0.72,\n        \"ttfb\": 442,\n        \"voice_purity\": 20,\n        \"median_of_call_time\": 5\n      },\n      {\n        \"country\": \"US\",\n        \"bandwidth\": \"53\",\n        \"response_time\": \"321\",\n        \"provider\": \"TransparentCalls\",\n        \"connection_stability\": 0.72,\n        \"ttfb\": 442,\n        \"voice_purity\": 20,\n        \"median_of_call_time\": 5\n      },\n      {\n        \"country\": \"US\",\n        \"bandwidth\": \"53\",\n        \"response_time\": \"321\",\n        \"provider\": \"E-Voice\",\n        \"connection_stability\": 0.72,\n        \"ttfb\": 442,\n        \"voice_purity\": 20,\n        \"median_of_call_time\": 5\n      },\n      {\n        \"country\": \"US\",\n        \"bandwidth\": \"53\",\n        \"response_time\": \"321\",\n        \"provider\": \"E-Voice\",\n        \"connection_stability\": 0.72,\n        \"ttfb\": 442,\n        \"voice_purity\": 20,\n        \"median_of_call_time\": 5\n      }\n    ],\n    \"email\": [\n      [\n        {\n          \"country\": \"RU\",\n          \"provider\": \"Gmail\",\n          \"delivery_time\": 195\n        },\n        {\n          \"country\": \"RU\",\n          \"provider\": \"Gmail\",\n          \"delivery_time\": 393\n        },\n        {\n          \"country\": \"RU\",\n          \"provider\": \"Gmail\",\n          \"delivery_time\": 393\n        }\n      ],\n      [\n        {\n          \"country\": \"RU\",\n          \"provider\": \"Gmail\",\n          \"delivery_time\": 393\n        },\n        {\n          \"country\": \"RU\",\n          \"provider\": \"Gmail\",\n          \"delivery_time\": 393\n        },\n        {\n          \"country\": \"RU\",\n          \"provider\": \"Gmail\",\n          \"delivery_time\": 393\n        }\n      ]\n    ],\n    \"billing\": {\n      \"create_customer\": true,\n      \"purchase\": true,\n      \"payout\": true,\n      \"recurring\": false,\n      \"fraud_control\": true,\n      \"checkout_page\": false\n    },\n    \"support\": [\n      3,\n      62\n    ],\n    \"incident\": [\n      {\"topic\":  \"Topic 1\", \"status\": \"active\"},\n      {\"topic\":  \"Topic 2\", \"status\": \"active\"},\n      {\"topic\":  \"Topic 3\", \"status\": \"closed\"},\n      {\"topic\":  \"Topic 4\", \"status\": \"closed\"}\n    ]\n  },\n  \"error\": \"\"\n}"))
}

func response(w http.ResponseWriter, r *http.Request, responseStruct interface{}) {
	response, _ := json.Marshal(responseStruct)

	w.Write(response)
}
