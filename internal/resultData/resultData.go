package resultData

import (
	"reflect"
	"statusPage/internal/billing"
	"statusPage/internal/email"
	"statusPage/internal/entities"
	"statusPage/internal/incident"
	"statusPage/internal/mms"
	"statusPage/internal/sms"
	"statusPage/internal/support"
	"statusPage/internal/voice"
	"sync"
	"time"
)

type ResultDataStorage struct {
	Time    time.Time
	Storage entities.ResultSetT
	sync.Mutex
}

func NewStorage() *ResultDataStorage {
	return &ResultDataStorage{
		Time: time.Now().Add(-31 * time.Second),
	}
}

func (r *ResultDataStorage) GetResultData() entities.ResultSetT {
	r.Lock()
	defer r.Unlock()
	t := time.Now()
	difference := t.Sub(r.Time)
	if difference > time.Second*30 {
		var wg sync.WaitGroup
		sms := sms.GetResultSMSData("../simulator/sms.data", &wg)
		mms := mms.GetResultMMSData(&wg)
		voice := voice.VoiceCallReader("../simulator/voice.data", &wg)
		billing := billing.BillingDataReader("../simulator/billing.data", &wg)
		support := support.GetResultSupportData(&wg)
		incident := incident.GetResultIncidentData(&wg)
		email := email.GetResultEmailData("../simulator/email.data", &wg)

		wg.Wait()

		result := entities.ResultSetT{
			SMS:       sms,
			MMS:       mms,
			VoiceCall: voice,
			Email:     email,
			Billing:   billing,
			Support:   support,
			Incidents: incident,
		}
		r.Storage = result
		r.Time = time.Now()
		return r.Storage
	} else {
		return r.Storage
	}
}

func (r *ResultDataStorage) IsFull() bool {
	r.Lock()
	defer r.Unlock()
	v := reflect.ValueOf(r.Storage)
	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).IsZero() == true {
			return false
		}
	}
	return true
}
