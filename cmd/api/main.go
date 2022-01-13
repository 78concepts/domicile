package main

import (
	"78concepts.com/domicile/internal/broker"
	"78concepts.com/domicile/internal/database"
	"78concepts.com/domicile/internal/devices"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func NewServer(ctx context.Context, client broker.MqttClient, devicesService *devices.Service) *Server {
	return &Server{ctx: ctx, client: client, devicesService: devicesService}
}

type Server struct {
	ctx context.Context
	client broker.MqttClient
	devicesService *devices.Service
}

func (s *Server) Index(w http.ResponseWriter, r *http.Request){

	log.Println("Endpoint hit: /")

	html :=
		`<html>
			<head>
				<script>
					function on(groupId){
						window.location.replace('/groupOn?group=' + groupId);
					}
					function off(groupId){
						window.location.replace('/groupOff?group=' + groupId);
					}
				</script>
			</head>
			<body>
				<h1>Welcome to the machine</h1>
			`
	fmt.Fprintf(w, html)

	areas, err := s.devicesService.GetAreas(s.ctx)

	if err != nil {
		log.Println("Unable to list areas", err)
	}

	for _, area := range areas {
		fmt.Fprintf(w, "<strong>" + area.Name + "</strong><br /><br />")
		fmt.Fprintf(w, "<a href=\"reports/?area=%s&type=temperature\">Temperature reports</a> [<a href=\"graphs/?area=%s&type=temperature\">Graph</a>]<br />", area.Uuid.Get(), area.Uuid.Get())
		fmt.Fprintf(w, "<a href=\"reports/?area=%s&type=humidity\">Humidity reports</a> [<a href=\"graphs/?area=%s&type=humidity\">Graph</a>]<br />", area.Uuid.Get(), area.Uuid.Get())
		fmt.Fprintf(w, "<a href=\"reports/?area=%s&type=pressure\">Pressure reports</a> [<a href=\"graphs/?area=%s&type=pressure\">Graph</a>]<br />", area.Uuid.Get(), area.Uuid.Get())
		fmt.Fprintf(w, "<a href=\"reports/?area=%s&type=illuminance\">Illuminance reports</a> [<a href=\"graphs/?area=%s&type=illuminance\">Graph</a>]<br />", area.Uuid.Get(), area.Uuid.Get())
		fmt.Fprintf(w, "<br /><br />")
	}

	//fmt.Fprintf(w, "<br /><br />")
	groups, err := s.devicesService.GetGroups(s.ctx)

	fmt.Fprintf(w, "<strong>Groups</strong><br /><br />")

	for _, group := range groups {
		fmt.Fprintf(w, "%s -- <button onClick=\"on(%d)\">Toggle on</button>  -- <button onClick=\"off(%d)\">Toggle off</button><br />", group.FriendlyName, group.Id, group.Id)
	}

	fmt.Fprintf(w, "</body></html>")
}

func (s *Server) ListAllReports(w http.ResponseWriter, r *http.Request) {

	log.Println("Endpoint hit: /reports?area=" + r.URL.Query().Get("area") + "&type=" + r.URL.Query().Get("type"))

	//if !middleware.ValidateRequiredQueryParam(w, r, "group") || !middleware.ValidateValidSetQueryParam(w, r, "type", []string{devices.TemperatureReport, devices.HumidityReport, devices.PressureReport, devices.IlluminanceReport}) {
	//	return
	//}

	uuid, _ := uuid.FromString(r.URL.Query().Get("area"))
	area, err := s.devicesService.GetArea(s.ctx, uuid)

	if area == nil || err != nil {
		//middleware.NotFound(w, "group", r.URL.Query().Get("group"))
		return
	}

	var response string

	switch reportType := r.URL.Query().Get("type"); reportType {
		case "temperature":
			data, _ := s.devicesService.GetTemperatureReports(s.ctx, area.Id)
			json, _ := json.Marshal(data)
			response = string(json)
		case "humidity":
			data, _ := s.devicesService.GetHumidityReports(s.ctx, area.Id)
			json, _ := json.Marshal(data)
			response = string(json)
	}

	fmt.Fprintf(w, string(response))

}

func (s *Server) ListAllGroups(w http.ResponseWriter, r *http.Request) {

	log.Println("Endpoint hit: /groups")

	//if !middleware.ValidateRequiredQueryParam(w, r, "group") || !middleware.ValidateValidSetQueryParam(w, r, "type", []string{devices.TemperatureReport, devices.HumidityReport, devices.PressureReport, devices.IlluminanceReport}) {
	//	return
	//}

	groups, err := s.devicesService.GetGroups(s.ctx)

	if err != nil {
		//middleware.NotFound(w, "group", r.URL.Query().Get("group"))
		return
	}

	var response string

	json, _ := json.Marshal(groups)
	response = string(json)

	fmt.Fprintf(w, string(response))

}

func (s *Server) GraphReports(w http.ResponseWriter, r *http.Request) {

	log.Println("Endpoint hit: /graphs?area=" + r.URL.Query().Get("area") + "&type=" + r.URL.Query().Get("type"))

	//if !middleware.ValidateRequiredQueryParam(w, r, "group") || !middleware.ValidateValidSetQueryParam(w, r, "type", []string{devices.TemperatureReport, devices.HumidityReport, devices.PressureReport, devices.IlluminanceReport}) {
	//	return
	//}

	uuid, _ := uuid.FromString(r.URL.Query().Get("area"))
	area, err := s.devicesService.GetArea(s.ctx, uuid)

	if area == nil || err != nil {
		//middleware.NotFound(w, "group", r.URL.Query().Get("group"))
		return
	}

	data, _ := s.devicesService.GetTemperatureReports(s.ctx, area.Id)

	var cleanedData []devices.TemperatureReport
	var previousDate time.Time
	for _, report := range data {
		if report.Date.Sub(previousDate).Seconds() >= 60 {
			cleanedData = append(cleanedData, report)
		}
		previousDate = report.Date
	}

	var values []string
	for _, report := range cleanedData {
		s := fmt.Sprintf("%v", report.Value)
		values = append(values, s)
	}

	location, err := time.LoadLocation("Australia/Melbourne")

	if err != nil {
		panic(err)
	}

	var dates []string
	for _, report := range cleanedData {
		dates = append(dates, "'" + report.Date.In(location).Format("Jan 2 15:04:05") + "'")
	}


	script :=
		`<html>
			<head>
				<script src='https://cdn.jsdelivr.net/npm/chart.js@3.2.1/dist/chart.min.js'></script>
			</head>
			<body>
				<div class="chart-container" style="position: relative; max-width:1024px; margin: auto; ">
					<canvas id="myChart" width="1024" height="768"></canvas>
				</div>
				<script>
					var ctx = document.getElementById('myChart').getContext('2d');
					var myChart = new Chart(ctx, {
						type: 'line',
						data: {
							labels: [` + strings.Join(dates, ", ") + `],
							datasets: [{
								label: '` + area.Name + " - " + r.URL.Query().Get("type") + `',
							  	tension: 0.4,
								data: [` + strings.Join(values, ", ") + `],
								borderColor: 'rgb(75, 192, 255)',
								borderWidth: 4
							}]
						},
						options: {
							responsive: false,
							scales: {
								y: {
									suggestedMin: 10,
									suggestedMax: 32
								}
							}
						}
					});
				</script>
			</body>
		</html>`

	fmt.Fprintf(w, script);
}

func (s *Server) TurnGroupOn(w http.ResponseWriter, r *http.Request) {

	log.Println("Endpoint hit: /groupOn?group=" + r.URL.Query().Get("group"))

	//if !middleware.ValidateRequiredQueryParam(w, r, "group") || !middleware.ValidateValidSetQueryParam(w, r, "type", []string{devices.TemperatureReport, devices.HumidityReport, devices.PressureReport, devices.IlluminanceReport}) {
	//	return
	//}

	id, _ := strconv.ParseUint(r.URL.Query().Get("group"), 10, 64)
	group, err := s.devicesService.GetGroup(s.ctx, id)

	if group == nil || err != nil {
		//middleware.NotFound(w, "group", r.URL.Query().Get("group"))
		return
	}

	devices.TurnGroupOn(s.client, group)

	http.Redirect(w, r, "/", 302)
}

func (s *Server) TurnGroupOff(w http.ResponseWriter, r *http.Request) {

	log.Println("Endpoint hit: /groupOff?group=" + r.URL.Query().Get("group"))

	//if !middleware.ValidateRequiredQueryParam(w, r, "group") || !middleware.ValidateValidSetQueryParam(w, r, "type", []string{devices.TemperatureReport, devices.HumidityReport, devices.PressureReport, devices.IlluminanceReport}) {
	//	return
	//}

	id, _ := strconv.ParseUint(r.URL.Query().Get("group"), 10, 64)
	group, err := s.devicesService.GetGroup(s.ctx, id)

	if group == nil || err != nil {
		//middleware.NotFound(w, "group", r.URL.Query().Get("group"))
		return
	}

	devices.TurnGroupOff(s.client, group)

	http.Redirect(w, r, "/", 302)
}


func handleRequests(server *Server) {

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", server.Index)
	router.HandleFunc("/reports", server.ListAllReports)
	router.HandleFunc("/graphs", server.GraphReports)
	router.HandleFunc("/groups", server.ListAllGroups)
	router.HandleFunc("/groupOn", server.TurnGroupOn)
	router.HandleFunc("/groupOff", server.TurnGroupOff)

	log.Fatal(http.ListenAndServe(":8080", router))
}

func main() {

	// Connect to the database
	dbPool:= database.NewPGXPool()

	ctx, ctxCancel := context.WithCancel(context.Background())

	var client = broker.NewMqttClient(ctx, ctxCancel)

	devicesService:= devices.NewService(&devices.PostgresRepository{Postgres: dbPool})

	server:= NewServer(ctx, client, devicesService)

	handleRequests(server)
}

