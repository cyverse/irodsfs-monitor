package service

import (
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/cyverse/irodsfs-monitor/types"
	log "github.com/sirupsen/logrus"
)

const (
	DataLifeSpanDays = 7
)

// Storage is a storage object
type Storage struct {
	Instances     map[string]types.ReportInstance
	FileTransfers map[string][]types.ReportFileTransfer
	Mutex         sync.Mutex
}

// NewStorage creates a storage
func NewStorage() *Storage {
	return &Storage{
		Instances:     map[string]types.ReportInstance{},
		FileTransfers: map[string][]types.ReportFileTransfer{},
		Mutex:         sync.Mutex{},
	}
}

// Init initializes the storage
func (storage *Storage) Init() error {
	logger := log.WithFields(log.Fields{
		"package":  "service",
		"function": "Storage.Destroy",
	})

	logger.Info("Initializing the storage")

	return nil
}

// Destroy destroys the storage
func (storage *Storage) Destroy() {
	logger := log.WithFields(log.Fields{
		"package":  "service",
		"function": "Storage.Destroy",
	})

	logger.Info("Destroying the storage")
}

// ListInstances lists instances
func (storage *Storage) ListInstances() []types.ReportInstance {
	storage.Mutex.Lock()
	defer storage.Mutex.Unlock()

	result := []types.ReportInstance{}
	for _, v := range storage.Instances {
		result = append(result, v)
	}

	sort.SliceStable(result, func(i int, j int) bool {
		return result[i].CreationTime.Before(result[j].CreationTime)
	})

	return result
}

// GetInstance returns instance
func (storage *Storage) GetInstance(instanceID string) (types.ReportInstance, bool) {
	storage.Mutex.Lock()
	defer storage.Mutex.Unlock()

	if v, ok := storage.Instances[instanceID]; ok {
		return v, true
	}

	return types.ReportInstance{}, false
}

// AddInstance adds an instance
func (storage *Storage) AddInstance(instance types.ReportInstance) {
	// clear old
	storage.clearAWeekOld()

	storage.Mutex.Lock()
	defer storage.Mutex.Unlock()

	storage.Instances[instance.InstanceID] = instance
}

// UpdateInstanceLastActivityTime updates the instance's last activity time
func (storage *Storage) UpdateInstanceLastActivityTime(instanceID string) error {
	storage.Mutex.Lock()
	defer storage.Mutex.Unlock()

	if instance, ok := storage.Instances[instanceID]; ok {
		instance.LastActivityTime = time.Now().UTC()
		storage.Instances[instanceID] = instance
		return nil
	}

	return fmt.Errorf("unable to find an instance for ID %s", instanceID)
}

// TerminateInstance sets the instance terminated
func (storage *Storage) TerminateInstance(instanceID string) error {
	storage.Mutex.Lock()
	defer storage.Mutex.Unlock()

	if instance, ok := storage.Instances[instanceID]; ok {
		instance.Terminated = true
		instance.LastActivityTime = time.Now().UTC()
		instance.TerminationTime = time.Now().UTC()
		storage.Instances[instanceID] = instance
		return nil
	}

	return fmt.Errorf("unable to find an instance for ID %s", instanceID)
}

// ListInstances lists instances
func (storage *Storage) ListFileTransfers() []types.ReportFileTransfer {
	storage.Mutex.Lock()
	defer storage.Mutex.Unlock()

	result := []types.ReportFileTransfer{}
	for _, v := range storage.FileTransfers {
		result = append(result, v...)
	}

	return result
}

// ListInstances lists instances
func (storage *Storage) ListFileTransfersForInstance(instanceID string) []types.ReportFileTransfer {
	storage.Mutex.Lock()
	defer storage.Mutex.Unlock()

	if v, ok := storage.FileTransfers[instanceID]; ok {
		return v
	}

	return []types.ReportFileTransfer{}
}

// AddInstance adds an instance
func (storage *Storage) AddFileTransfer(transfer types.ReportFileTransfer) error {
	// clear old
	storage.clearAWeekOld()

	storage.Mutex.Lock()
	defer storage.Mutex.Unlock()

	if _, ok := storage.Instances[transfer.InstanceID]; ok {
		if existingList, ok2 := storage.FileTransfers[transfer.InstanceID]; ok2 {
			existingList = append(existingList, transfer)
			storage.FileTransfers[transfer.InstanceID] = existingList
		} else {
			storage.FileTransfers[transfer.InstanceID] = []types.ReportFileTransfer{transfer}
		}
		return nil
	}

	return fmt.Errorf("unable to find an instance for ID %s", transfer.InstanceID)
}

// CleanUp clears all instance and transfer data
func (storage *Storage) CleanUp() {
	logger := log.WithFields(log.Fields{
		"package":  "service",
		"function": "Storage.CleanUp",
	})

	storage.Mutex.Lock()
	defer storage.Mutex.Unlock()

	storage.Instances = map[string]types.ReportInstance{}
	storage.FileTransfers = map[string][]types.ReportFileTransfer{}

	logger.Info("Cleaned up storage")
}

// clearOld clears old instance and transfer data
func (storage *Storage) clearOld(daysOld int) {
	logger := log.WithFields(log.Fields{
		"package":  "service",
		"function": "Storage.clearOld",
	})

	storage.Mutex.Lock()
	defer storage.Mutex.Unlock()

	instanceIDToBeRemoved := []string{}
	lastWeek := time.Now().AddDate(0, 0, -1*daysOld)
	for instanceID, instance := range storage.Instances {
		if instance.CreationTime.Before(lastWeek) {
			// delete
			instanceIDToBeRemoved = append(instanceIDToBeRemoved, instanceID)
		}
	}

	for _, instanceID := range instanceIDToBeRemoved {
		delete(storage.FileTransfers, instanceID)
		delete(storage.Instances, instanceID)
	}

	logger.Infof("Cleaned up old data that are %d days old", daysOld)
}

// clearAWeekOld clears old instance and transfer data
func (storage *Storage) clearAWeekOld() {
	storage.clearOld(DataLifeSpanDays)
}
