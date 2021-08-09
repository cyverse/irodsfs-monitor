package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/cyverse/irodsfs-monitor/types"
	"github.com/rs/xid"
	log "github.com/sirupsen/logrus"
)

// APIClient is a struct that holds connection information of a API client
type APIClient struct {
	APIRootURL string
}

// NewAPIClient creates a new API client
func NewAPIClient(apiRootURL string) *APIClient {
	return &APIClient{
		APIRootURL: apiRootURL,
	}
}

func (client *APIClient) makeAPIURL(apiPath string) string {
	u := client.APIRootURL
	if !strings.HasSuffix(u, "/") {
		// does not end with '/'
		u = u + "/"
	}

	if strings.HasPrefix(apiPath, "/") {
		return u + apiPath[1:]
	}

	return u + apiPath
}

// AddInstance registers an instance
func (client *APIClient) AddInstance(instance *types.ReportInstance) (string, error) {
	logger := log.WithFields(log.Fields{
		"package":  "client",
		"function": "APIClient.AddInstance",
	})

	if len(instance.ClientHostname) == 0 {
		hostname, err := os.Hostname()
		if err == nil {
			instance.ClientHostname = hostname
		}
	}

	if len(instance.InstanceID) == 0 {
		// generate an id
		instance.InstanceID = xid.New().String()
	}

	if instance.CreationTime.IsZero() {
		instance.CreationTime = time.Now().UTC()
	}

	JSONBytes, err := json.Marshal(instance)
	if err != nil {
		logger.Error(err)
		return "", err
	}

	url := client.makeAPIURL("/instances")
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		logger.Error(err)
		return "", err
	}

	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Body = ioutil.NopCloser(bytes.NewReader(JSONBytes))
	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		logger.Error(err)
		return "", err
	}

	if resp.StatusCode != http.StatusAccepted && resp.StatusCode != http.StatusOK {
		logger.Error(fmt.Sprintf("service error returned - %s", resp.Status))
		return "", fmt.Errorf("service error returned - %s", resp.Status)
	}

	return instance.InstanceID, nil
}

// ListInstances lists instances registered
func (client *APIClient) ListInstances(instance *types.ReportInstance) ([]types.ReportInstance, error) {
	logger := log.WithFields(log.Fields{
		"package":  "client",
		"function": "APIClient.ListInstances",
	})

	url := client.makeAPIURL("/instances")
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	if resp.StatusCode != http.StatusAccepted && resp.StatusCode != http.StatusOK {
		logger.Error(fmt.Sprintf("service error returned - %s", resp.Status))
		return nil, fmt.Errorf("service error returned - %s", resp.Status)
	}

	var instances []types.ReportInstance
	responseJSON, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	err = json.Unmarshal(responseJSON, &instances)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return instances, nil
}

// GetInstance returns an instance registered
func (client *APIClient) GetInstance(instanceID string) (types.ReportInstance, error) {
	logger := log.WithFields(log.Fields{
		"package":  "client",
		"function": "APIClient.GetInstance",
	})

	url := client.makeAPIURL(fmt.Sprintf("/instances/%s", instanceID))
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logger.Error(err)
		return types.ReportInstance{}, err
	}

	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		logger.Error(err)
		return types.ReportInstance{}, err
	}

	if resp.StatusCode != http.StatusAccepted && resp.StatusCode != http.StatusOK {
		logger.Error(fmt.Sprintf("service error returned - %s", resp.Status))
		return types.ReportInstance{}, fmt.Errorf("service error returned - %s", resp.Status)
	}

	var instance types.ReportInstance
	responseJSON, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error(err)
		return types.ReportInstance{}, err
	}

	err = json.Unmarshal(responseJSON, &instance)
	if err != nil {
		logger.Error(err)
		return types.ReportInstance{}, err
	}

	return instance, nil
}

// TerminateInstance sets the instance terminated
func (client *APIClient) TerminateInstance(instanceID string) error {
	logger := log.WithFields(log.Fields{
		"package":  "client",
		"function": "APIClient.TerminateInstance",
	})

	url := client.makeAPIURL(fmt.Sprintf("/instances/%s", instanceID))
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		logger.Error(err)
		return err
	}

	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		logger.Error(err)
		return err
	}

	if resp.StatusCode != http.StatusAccepted && resp.StatusCode != http.StatusOK {
		logger.Error(fmt.Sprintf("service error returned - %s", resp.Status))
		return fmt.Errorf("service error returned - %s", resp.Status)
	}

	return nil
}

// AddFileTransfer adds a file transfer
func (client *APIClient) AddFileTransfer(transfer *types.ReportFileTransfer) error {
	logger := log.WithFields(log.Fields{
		"package":  "client",
		"function": "APIClient.AddFileTransfer",
	})

	if len(transfer.InstanceID) == 0 {
		return fmt.Errorf("invalid instance id")
	}

	JSONBytes, err := json.Marshal(transfer)
	if err != nil {
		logger.Error(err)
		return err
	}

	url := client.makeAPIURL("/transfers")
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		logger.Error(err)
		return err
	}

	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	httpClient := &http.Client{}
	req.Body = ioutil.NopCloser(bytes.NewReader(JSONBytes))
	resp, err := httpClient.Do(req)
	if err != nil {
		logger.Error(err)
		return err
	}

	if resp.StatusCode != http.StatusAccepted && resp.StatusCode != http.StatusOK {
		logger.Error(fmt.Sprintf("service error returned - %s", resp.Status))
		return fmt.Errorf("service error returned - %s", resp.Status)
	}

	return nil
}

// ListFileTransfers lists all file transfers
func (client *APIClient) ListFileTransfers() ([]types.ReportFileTransfer, error) {
	logger := log.WithFields(log.Fields{
		"package":  "client",
		"function": "APIClient.ListFileTransfers",
	})

	url := client.makeAPIURL("/transfers")
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	if resp.StatusCode != http.StatusAccepted && resp.StatusCode != http.StatusOK {
		logger.Error(fmt.Sprintf("service error returned - %s", resp.Status))
		return nil, fmt.Errorf("service error returned - %s", resp.Status)
	}

	var transfers []types.ReportFileTransfer
	responseJSON, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	err = json.Unmarshal(responseJSON, &transfers)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return transfers, nil
}

// ListFileTransfersForInstance lists all file transfers
func (client *APIClient) ListFileTransfersForInstance(instanceID string) ([]types.ReportFileTransfer, error) {
	logger := log.WithFields(log.Fields{
		"package":  "client",
		"function": "APIClient.ListFileTransfersForInstance",
	})

	url := client.makeAPIURL(fmt.Sprintf("/transfers/%s", instanceID))
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	if resp.StatusCode != http.StatusAccepted && resp.StatusCode != http.StatusOK {
		logger.Error(fmt.Sprintf("service error returned - %s", resp.Status))
		return nil, fmt.Errorf("service error returned - %s", resp.Status)
	}

	var transfers []types.ReportFileTransfer
	responseJSON, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	err = json.Unmarshal(responseJSON, &transfers)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return transfers, nil
}
