package dependencies

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type Dependencies struct {
	// Making each field an interface rather than a struct
	// would would allow us to create mocks for each dependency,
	// which would make testing easier. Something to consider.
	db *SensorDB
}

func NewDependencies() (*Dependencies, error) {
	db, err := initDB()
	if err != nil {
		return nil, err
	}

	return &Dependencies{db: db}, nil
}

type SensorListResponse struct {
	// Fields need to be exported in order to be JSON encoded
	Success  bool     `json:"success"`
	Sensors  []Sensor `json:"sensor"`
	ErrorMsg string   `json:"error_msg"`
}

func (d *Dependencies) ListSensorsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Debug("Running ListSensorsHandler")
		resp := SensorListResponse{}

		writeResp := func(errMsg string) {
			if errMsg != "" {
				resp.ErrorMsg = errMsg
				log.Error(errMsg)
			}

			respJson, err := json.Marshal(resp)
			if err != nil {
				log.Errorf("Internal Server Error - could not marshal JSON for response: %s", err.Error())
				w.WriteHeader(500)
				return
			}
			fmt.Fprint(w, string(respJson))
		}

		list, err := d.db.ListSensors()
		if err != nil {
			writeResp(fmt.Sprintf("Could not retrieve list of sensors: %s", err.Error()))
			return
		}

		resp = SensorListResponse{
			Success:  true,
			Sensors:  list,
			ErrorMsg: "",
		}

		writeResp("")
	}
}

type SensorResponse struct {
	// Fields need to be exported in order to be JSON encoded
	Success  bool    `json:"success"`
	Sensor   *Sensor `json:"sensor"`
	ErrorMsg string  `json:"error_msg"`
}

func (d *Dependencies) InsertSensorHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Debug("Running InsertSensorsHandler")
		resp := SensorResponse{}

		writeResp := func(errMsg string) {
			if errMsg != "" {
				resp.ErrorMsg = errMsg
				log.Error(errMsg)
			}

			respJson, err := json.Marshal(resp)
			if err != nil {
				log.Errorf("Internal Server Error - could not marshal JSON for response: %s", err.Error())
				w.WriteHeader(500)
				return
			}
			fmt.Fprint(w, string(respJson))
		}

		var sensor Sensor
		b, err := io.ReadAll(r.Body)
		if err != nil {
			writeResp(fmt.Sprintf("Could not read request body: %s", err.Error()))
			return
		}

		err = json.Unmarshal(b, &sensor)
		if err != nil {
			writeResp(fmt.Sprintf("Could not parse JSON body: %s", err.Error()))
			return
		}

		recvSensor, err := d.db.InsertSensor(sensor)
		if err != nil {
			writeResp(fmt.Sprintf("Could not insert new sensor: %s", err.Error()))
			return
		}

		resp = SensorResponse{
			Success:  true,
			Sensor:   recvSensor,
			ErrorMsg: "",
		}
		writeResp("")
	}
}

func (d *Dependencies) GetSensorByNameHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Debug("Running GetSensorByNameHandler")
		resp := SensorResponse{}

		writeResp := func(errMsg string) {
			if errMsg != "" {
				resp.ErrorMsg = errMsg
				log.Error(errMsg)
			}

			respJson, err := json.Marshal(resp)
			if err != nil {
				log.Errorf("Internal Server Error - could not marshal JSON for response: %s", err.Error())
				w.WriteHeader(500)
				return
			}
			fmt.Fprint(w, string(respJson))
		}

		// If this route is hit, `ok` should always return true,
		// so this check is probably unnecessary.
		name, ok := mux.Vars(r)["name"]
		if !ok {
			writeResp("Could not fetch sensor metadata: no sensor name provided")
			return
		}

		sensor, err := d.db.GetSensorByName(name)
		if err != nil {
			writeResp(fmt.Sprintf("Could not fetch sensor metadata: %s", err.Error()))
			return
		}

		resp = SensorResponse{
			Success:  true,
			Sensor:   sensor,
			ErrorMsg: "",
		}
		writeResp("")
	}
}

func (d *Dependencies) UpdateSensorHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Debug("Running UpdateSensorHandler")
		resp := SensorResponse{}

		writeResp := func(errMsg string) {
			if errMsg != "" {
				resp.ErrorMsg = errMsg
				log.Error(errMsg)
			}

			respJson, err := json.Marshal(resp)
			if err != nil {
				log.Errorf("Internal Server Error - could not marshal JSON for response: %s", err.Error())
				w.WriteHeader(500)
				return
			}
			fmt.Fprint(w, string(respJson))
		}

		var sensor Sensor
		b, err := io.ReadAll(r.Body)
		if err != nil {
			writeResp(fmt.Sprintf("Could not read request body: %s", err.Error()))
			return
		}

		err = json.Unmarshal(b, &sensor)
		if err != nil {
			writeResp(fmt.Sprintf("Could not parse JSON body: %s", err.Error()))
			return
		}

		name, ok := mux.Vars(r)["name"]
		if !ok {
			writeResp("Could not fetch sensor metadata: no sensor name provided")
			return
		}

		recvSensor, err := d.db.UpdateSensor(name, &sensor)
		if err != nil {
			writeResp(fmt.Sprintf("Could not fetch sensor metadata: %s", err.Error()))
			return
		}

		resp = SensorResponse{
			Success:  true,
			Sensor:   recvSensor,
			ErrorMsg: "",
		}
		writeResp("")
	}
}

func (d *Dependencies) DeleteSensorHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Debug("Running DeleteSensorHandler")
		resp := SensorResponse{}

		writeResp := func(errMsg string) {
			if errMsg != "" {
				resp.ErrorMsg = errMsg
				log.Error(errMsg)
			}

			respJson, err := json.Marshal(resp)
			if err != nil {
				log.Errorf("Internal Server Error - could not marshal JSON for response: %s", err.Error())
				w.WriteHeader(500)
				return
			}
			fmt.Fprint(w, string(respJson))
		}

		name, ok := mux.Vars(r)["name"]
		if !ok {
			writeResp("Could not fetch sensor metadata: no sensor name provided")
			return
		}

		recvSensor, err := d.db.DeleteSensor(name)
		if err != nil {
			writeResp(fmt.Sprintf("Could not fetch sensor metadata: %s", err.Error()))
			return
		}

		resp = SensorResponse{
			Success:  true,
			Sensor:   recvSensor,
			ErrorMsg: "",
		}

		writeResp("")
	}
}

func (d *Dependencies) NearestSensorHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Debug("Running NearestSensorHandler")
		resp := SensorResponse{}

		writeResp := func(errMsg string) {
			if errMsg != "" {
				resp.ErrorMsg = errMsg
				log.Error(errMsg)
			}

			respJson, err := json.Marshal(resp)
			if err != nil {
				log.Errorf("Internal Server Error - could not marshal JSON for response: %s", err.Error())
				w.WriteHeader(500)
				return
			}
			fmt.Fprint(w, string(respJson))
		}

		b, err := io.ReadAll(r.Body)
		if err != nil {
			writeResp(fmt.Sprintf("Could not read request body: %s", err.Error()))
			return
		}

		var point Point
		err = json.Unmarshal(b, &point)
		if err != nil {
			writeResp(fmt.Sprintf("Could not parse JSON body: %s", err.Error()))
			return
		}

		list, err := d.db.ListSensors()
		if err != nil {
			writeResp(fmt.Sprintf("Could not retrieve list of sensors: %s", err.Error()))
			return
		}

		nearest := nearestSensorToPoint(point, list)
		if nearest == nil {
			writeResp("List of known sensors is empty")
			return
		}

		resp = SensorResponse{
			Success:  true,
			Sensor:   nearest,
			ErrorMsg: "",
		}
		writeResp("")
	}
}
