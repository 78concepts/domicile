package api

import (
	"78concepts.com/domicile/internal/service"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gofrs/uuid"
	"log"
	"net/http"
)

func NewReportsApi(ctx context.Context, reportsService *service.ReportsService, areasService *service.AreasService) *ReportsApi {
	return &ReportsApi{ctx: ctx, reportsService: reportsService, areasService: areasService}
}

type ReportsApi struct {
	ctx context.Context
	reportsService *service.ReportsService
	areasService *service.AreasService
}

func (a *ReportsApi) ListReports(w http.ResponseWriter, r *http.Request) {

	log.Println("Endpoint hit: /reports?area=" + r.URL.Query().Get("area") + "&type=" + r.URL.Query().Get("type"))

	uuid, _ := uuid.FromString(r.URL.Query().Get("area"))
	area, err := a.areasService.GetArea(a.ctx, uuid)

	if area == nil || err != nil {
		json, _ := json.Marshal(map[string]interface{} {"status": 404, "error": "Area not found"})
		w.WriteHeader(404)
		fmt.Fprintf(w, string(json))
		return
	}

	var response string

	switch reportType := r.URL.Query().Get("type"); reportType {
		case "temperature":
			data, _ := a.reportsService.GetTemperatureReports(a.ctx, area.Id)
			json, _ := json.Marshal(data)
			response = string(json)
		case "humidity":
			data, _ := a.reportsService.GetHumidityReports(a.ctx, area.Id)
			json, _ := json.Marshal(data)
			response = string(json)
		case "pressure":
			data, _ := a.reportsService.GetPressureReports(a.ctx, area.Id)
			json, _ := json.Marshal(data)
			response = string(json)
		case "illuminance":
			data, _ := a.reportsService.GetIlluminanceReports(a.ctx, area.Id)
			log.Println(data)
			json, _ := json.Marshal(data)
			response = string(json)
		default:
			json, _ := json.Marshal(map[string]interface{} {"status": 400, "error": "Invalid report type"})
			w.WriteHeader(400)
			fmt.Fprintf(w, string(json))
			return
	}

	fmt.Fprintf(w, response)

}

