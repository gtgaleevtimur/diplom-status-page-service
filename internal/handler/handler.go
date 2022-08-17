package handler

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"log"
	"net/http"
	"statusPage/internal/entities"
	"statusPage/internal/resultData"
)

func Build(router *chi.Mux, store *resultData.ResultDataStorage) {
	router.Use(middleware.Recoverer)

	controller := NewController(store)

	router.Get("/", controller.GetData)

}

type Controller struct {
	storage *resultData.ResultDataStorage
}

func NewController(storage *resultData.ResultDataStorage) *Controller {
	return &Controller{
		storage: storage,
	}
}

func (c *Controller) GetData(w http.ResponseWriter, r *http.Request) {
	var result entities.ResultT
	resultSetT := c.storage.GetResultData()
	checkFull := c.storage.IsFull()
	switch checkFull {
	case true:
		result.Status = true
		result.Data = resultSetT
	case false:
		result.Status = false
		result.Error = "Error on collect data"
	}
	res, err := json.Marshal(result)
	if err != nil {
		log.Printf("Ошибка преобразования ResultT в json: %v", err)
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Write(res)
}
