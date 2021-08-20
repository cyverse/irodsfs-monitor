package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

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

// addHandlers adds web server handlers
func (svc *MonitorService) addHandlers() {
	svc.Router.HandleFunc("/instances", svc.addInstance).Methods("POST")
	svc.Router.HandleFunc("/instances", svc.listInstances).Methods("GET")
	svc.Router.HandleFunc("/instances/{instance_id}", svc.getInstance).Methods("GET")
	svc.Router.HandleFunc("/instances/{instance_id}", svc.terminateInstance).Methods("DELETE")

	svc.Router.HandleFunc("/transfers", svc.addTransfer).Methods("POST")
	svc.Router.HandleFunc("/transfers", svc.listTransfers).Methods("GET")
	svc.Router.HandleFunc("/transfers/{instance_id}", svc.listTransfersForInstance).Methods("GET")
	svc.Router.HandleFunc("/cleanup", svc.cleanUp).Methods("DELETE")
	svc.Router.HandleFunc("/cleanup/{days}", svc.cleanUpDaysOld).Methods("DELETE")
}

// Init initializes the service
func (svc *MonitorService) Init() error {
	return nil
}

// Start starts the service
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

func (svc *MonitorService) getClientIP(r *http.Request) string {
	addr := r.Header.Get("X-Real-Ip")
	if addr == "" {
		addr = r.Header.Get("X-Forwarded-For")
	}

	if addr == "" {
		addr = r.RemoteAddr
	}

	// erase port number
	addrs := strings.Split(addr, ":")
	if len(addrs) > 0 {
		addr = addrs[0]
	}
	return addr
}

func (svc *MonitorService) addInstance(w http.ResponseWriter, r *http.Request) {
	logger := log.WithFields(log.Fields{
		"package":  "service",
		"function": "MonitorService.addInstance",
	})

	logger.Infof("Page access request (%s) from %s to %s", r.Method, r.RemoteAddr, r.RequestURI)
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

	nowUTC := time.Now().UTC()

	instance.ClientHostIP = svc.getClientIP(r)
	if instance.CreationTime.IsZero() {
		instance.CreationTime = nowUTC
	}

	instance.Terminated = false
	instance.LastActivityTime = nowUTC

	svc.Storage.AddInstance(instance)
	w.WriteHeader(http.StatusAccepted)
}

func (svc *MonitorService) listInstances(w http.ResponseWriter, r *http.Request) {
	logger := log.WithFields(log.Fields{
		"package":  "service",
		"function": "MonitorService.listInstances",
	})

	logger.Infof("Page access request (%s) from %s to %s", r.Method, r.RemoteAddr, r.RequestURI)

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

	logger.Infof("Page access request (%s) from %s to %s", r.Method, r.RemoteAddr, r.RequestURI)

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

func (svc *MonitorService) terminateInstance(w http.ResponseWriter, r *http.Request) {
	logger := log.WithFields(log.Fields{
		"package":  "service",
		"function": "MonitorService.terminateInstance",
	})

	logger.Infof("Page access request (%s) from %s to %s", r.Method, r.RemoteAddr, r.RequestURI)

	varMap := mux.Vars(r)
	instanceID, ok := varMap["instance_id"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("instance_id is not given"))
		return
	}

	err := svc.Storage.TerminateInstance(instanceID)
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func (svc *MonitorService) addTransfer(w http.ResponseWriter, r *http.Request) {
	logger := log.WithFields(log.Fields{
		"package":  "service",
		"function": "MonitorService.addTransfer",
	})

	logger.Infof("Page access request (%s) from %s to %s", r.Method, r.RemoteAddr, r.RequestURI)
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

	err = svc.Storage.AddFileTransfer(transfer)
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	err = svc.Storage.UpdateInstanceLastActivityTime(transfer.InstanceID)
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func (svc *MonitorService) listTransfers(w http.ResponseWriter, r *http.Request) {
	logger := log.WithFields(log.Fields{
		"package":  "service",
		"function": "MonitorService.listTransfers",
	})

	logger.Infof("Page access request (%s) from %s to %s", r.Method, r.RemoteAddr, r.RequestURI)

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

	logger.Infof("Page access request (%s) from %s to %s", r.Method, r.RemoteAddr, r.RequestURI)

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

func (svc *MonitorService) cleanUp(w http.ResponseWriter, r *http.Request) {
	logger := log.WithFields(log.Fields{
		"package":  "service",
		"function": "MonitorService.cleanUp",
	})

	logger.Infof("Page access request (%s) from %s to %s", r.Method, r.RemoteAddr, r.RequestURI)

	svc.Storage.CleanUp()
	w.WriteHeader(http.StatusAccepted)
}

func (svc *MonitorService) cleanUpDaysOld(w http.ResponseWriter, r *http.Request) {
	logger := log.WithFields(log.Fields{
		"package":  "service",
		"function": "MonitorService.cleanUpDaysOld",
	})

	logger.Infof("Page access request (%s) from %s to %s", r.Method, r.RemoteAddr, r.RequestURI)

	varMap := mux.Vars(r)
	daysString, ok := varMap["days"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("days is not given"))
		return
	}

	days, err := strconv.Atoi(daysString)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("days is not number"))
		return
	}

	svc.Storage.clearOld(days)
	w.WriteHeader(http.StatusAccepted)
}
