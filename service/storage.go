package service

import (
	"fmt"
	"time"

	"github.com/cyverse/irodsfs-monitor/types"
	log "github.com/sirupsen/logrus"
)

// Storage is a storage object
type Storage struct {
	Instances     map[string]types.ReportInstance
	FileTransfers map[string][]types.ReportFileTransfer
}

// NewStorage creates a storage
func NewStorage() *Storage {
	return &Storage{
		Instances:     map[string]types.ReportInstance{},
		FileTransfers: map[string][]types.ReportFileTransfer{},
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
	result := []types.ReportInstance{}
	for _, v := range storage.Instances {
		result = append(result, v)
	}

	return result
}

// GetInstance returns instance
func (storage *Storage) GetInstance(instanceID string) (types.ReportInstance, bool) {
	if v, ok := storage.Instances[instanceID]; ok {
		return v, true
	}

	return types.ReportInstance{}, false
}

// AddInstance adds an instance
func (storage *Storage) AddInstance(instance types.ReportInstance) {
	storage.Instances[instance.InstanceID] = instance
}

// UpdateInstanceLastActivityTime updates the instance's last activity time
func (storage *Storage) UpdateInstanceLastActivityTime(instanceID string) error {
	if instance, ok := storage.Instances[instanceID]; ok {
		instance.LastActivityTime = time.Now().UTC()
		storage.Instances[instanceID] = instance
		return nil
	}

	return fmt.Errorf("unable to find an instance for ID %s", instanceID)
}

// TerminateInstance sets the instance terminated
func (storage *Storage) TerminateInstance(instanceID string) error {
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
	result := []types.ReportFileTransfer{}
	for _, v := range storage.FileTransfers {
		result = append(result, v...)
	}

	return result
}

// ListInstances lists instances
func (storage *Storage) ListFileTransfersForInstance(instanceID string) []types.ReportFileTransfer {
	if v, ok := storage.FileTransfers[instanceID]; ok {
		return v
	}

	return []types.ReportFileTransfer{}
}

// AddInstance adds an instance
func (storage *Storage) AddFileTransfer(transfer types.ReportFileTransfer) error {
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
