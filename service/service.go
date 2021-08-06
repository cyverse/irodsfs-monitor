package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/cyverse/irodsfs-monitor/types"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

// MonitorService is a service object
type MonitorService struct {
	Config    *Config
	WebServer *http.Server
	Router    *mux.Router
	Storage   *Storage
}

// NewMonitorService creates a new monitor service
func NewMonitorService(config *Config) *MonitorService {

	webServerRouter := mux.NewRouter()
	webServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.ServicePort),
		Handler: webServerRouter,
	}

	service := &MonitorService{
		Config:    config,
		WebServer: webServer,
		Router:    webServerRouter,
		Storage:   NewStorage(),
	}

	service.addHandlers()

	return service
}

// AddHandlers adds web server handlers
func (svc *MonitorService) addHandlers() {
	svc.Router.HandleFunc("/instances", svc.addInstance).Methods("POST")
	svc.Router.HandleFunc("/instances", svc.listInstances).Methods("GET")
	svc.Router.HandleFunc("/instances/{instance_id}", svc.getInstance).Methods("GET")

	svc.Router.HandleFunc("/transfers", svc.addTransfer).Methods("POST")
	svc.Router.HandleFunc("/transfers", svc.listTransfers).Methods("GET")
	svc.Router.HandleFunc("/transfers/{instance_id}", svc.listTransfersForInstance).Methods("GET")
}

func (svc *MonitorService) Init() error {
	return nil
}

func (svc *MonitorService) Start() error {
	logger := log.WithFields(log.Fields{
		"package":  "service",
		"function": "MonitorService.Start",
	})

	logger.Info("Starting the iRODS FUSE Lite Monitoring service")

	err := svc.WebServer.ListenAndServe()
	if err != nil {
		logger.Error(err)
		return err
	}

	// should not return
	return nil
}

// Destroy destroys the service
func (svc *MonitorService) Destroy() {
	logger := log.WithFields(log.Fields{
		"package":  "service",
		"function": "MonitorService.Destroy",
	})

	logger.Info("Destroying the iRODS FUSE Lite Monitoring service")

	err := svc.WebServer.Close()
	if err != nil {
		logger.Error(err)
	}
}

func (svc *MonitorService) addInstance(w http.ResponseWriter, r *http.Request) {
	logger := log.WithFields(log.Fields{
		"package":  "service",
		"function": "MonitorService.addInstance",
	})

	logger.Infof("Page access request from %s to %s", r.RemoteAddr, r.RequestURI)
	var instance types.ReportInstance

	requestJSON, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	err = json.Unmarshal(requestJSON, &instance)
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	svc.Storage.AddInstance(instance)
	w.WriteHeader(http.StatusAccepted)
}

func (svc *MonitorService) listInstances(w http.ResponseWriter, r *http.Request) {
	logger := log.WithFields(log.Fields{
		"package":  "service",
		"function": "MonitorService.listInstances",
	})

	logger.Infof("Page access request from %s to %s", r.RemoteAddr, r.RequestURI)

	instances := svc.Storage.ListInstances()
	responseJSON, err := json.Marshal(instances)
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_, err = w.Write(responseJSON)
	if err != nil {
		logger.Error(err)
		return
	}
}

func (svc *MonitorService) getInstance(w http.ResponseWriter, r *http.Request) {
	logger := log.WithFields(log.Fields{
		"package":  "service",
		"function": "MonitorService.getInstance",
	})

	logger.Infof("Page access request from %s to %s", r.RemoteAddr, r.RequestURI)

	varMap := mux.Vars(r)
	instanceID, ok := varMap["instance_id"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("instance_id is not given"))
		return
	}

	instance, ok := svc.Storage.GetInstance(instanceID)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	responseJSON, err := json.Marshal(instance)
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_, err = w.Write(responseJSON)
	if err != nil {
		logger.Error(err)
		return
	}
}

func (svc *MonitorService) addTransfer(w http.ResponseWriter, r *http.Request) {
	logger := log.WithFields(log.Fields{
		"package":  "service",
		"function": "MonitorService.addTransfer",
	})

	logger.Infof("Page access request from %s to %s", r.RemoteAddr, r.RequestURI)
	var transfer types.ReportFileTransfer

	requestJSON, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	err = json.Unmarshal(requestJSON, &transfer)
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	svc.Storage.AddFileTransfer(transfer)
	w.WriteHeader(http.StatusAccepted)
}

func (svc *MonitorService) listTransfers(w http.ResponseWriter, r *http.Request) {
	logger := log.WithFields(log.Fields{
		"package":  "service",
		"function": "MonitorService.listTransfers",
	})

	logger.Infof("Page access request from %s to %s", r.RemoteAddr, r.RequestURI)

	transfers := svc.Storage.ListFileTransfers()
	responseJSON, err := json.Marshal(transfers)
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_, err = w.Write(responseJSON)
	if err != nil {
		logger.Error(err)
		return
	}
}

func (svc *MonitorService) listTransfersForInstance(w http.ResponseWriter, r *http.Request) {
	logger := log.WithFields(log.Fields{
		"package":  "service",
		"function": "MonitorService.listTransfersForInstance",
	})

	logger.Infof("Page access request from %s to %s", r.RemoteAddr, r.RequestURI)

	varMap := mux.Vars(r)
	instanceID, ok := varMap["instance_id"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("instance_id is not given"))
		return
	}

	transfers := svc.Storage.ListFileTransfersForInstance(instanceID)
	responseJSON, err := json.Marshal(transfers)
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_, err = w.Write(responseJSON)
	if err != nil {
		logger.Error(err)
		return
	}
}
