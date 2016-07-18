// Copyright 2016 NetApp, Inc. All Rights Reserved.

package storage_drivers

import (
	"encoding/json"
	"fmt"

	"github.com/netapp/netappdvp/apis/sfapi"
)

// CurrentDriverVersion is the expected version in the config file
const CurrentDriverVersion = 1

// DriverVersion is the actual release version number
const DriverVersion = "1.2.1"

// CommonStorageDriverConfig holds settings in common across all StorageDrivers
type CommonStorageDriverConfig struct {
	Version           int             `json:"version"`
	StorageDriverName string          `json:"storageDriverName"`
	Debug             bool            `json:"debug"`
	DisableDelete     bool            `json:"disableDelete"`
	StoragePrefixRaw  json.RawMessage `json:"storagePrefix,string"`
}

// ValidateCommonSettings attempts to "partially" decode the JSON into just the settings in CommonStorageDriverConfig
func ValidateCommonSettings(configJSON string) (*CommonStorageDriverConfig, error) {
	config := &CommonStorageDriverConfig{}

	// decode configJSON into config object
	err := json.Unmarshal([]byte(configJSON), &config)
	if err != nil {
		return nil, fmt.Errorf("Cannot decode json configuration error: %v", err)
	}

	// load storage drivers and validate the one specified actually exists
	if config.StorageDriverName == "" {
		return nil, fmt.Errorf("Missing storage driver name in configuration file")
	}

	// validate config file version information
	if config.Version != CurrentDriverVersion {
		return nil, fmt.Errorf("Unexpected config file version;  found %v expected %v", config.Version, CurrentDriverVersion)
	}

	return config, nil
}

// OntapStorageDriverConfig holds settings for OntapStorageDrivers
type OntapStorageDriverConfig struct {
	CommonStorageDriverConfig        // embedded types replicate all fields
	ManagementLIF             string `json:"managementLIF"`
	DataLIF                   string `json:"dataLIF"`
	IgroupName                string `json:"igroupName"`
	SVM                       string `json:"svm"`
	Username                  string `json:"username"`
	Password                  string `json:"password"`
	Aggregate                 string `json:"aggregate"`
}

// ESeriesStorageDriverConfig holds settings for ESeriesStorageDriver
type ESeriesStorageDriverConfig struct {
	CommonStorageDriverConfig

	//Web Proxy Services Info
	WebProxyHostname  string `json:"webProxyHostname"`
	WebProxyPort      string `json:"webProxyPort"`      // optional
	WebProxyUseHTTP   bool   `json:"webProxyUseHTTP"`   // optional
	WebProxyVerifyTLS bool   `json:"webProxyVerifyTLS"` // optional
	Username          string `json:"username"`          //rw
	Password          string `json:"password"`          //rw

	//Array Info
	ControllerA     string `json:"controllerA"`
	ControllerB     string `json:"controllerB"`
	PasswordArray   string `json:"passwordArray"`   //optional
	ArrayRegistered bool   `json:"arrayRegistered"` //optional

	//Host Networking
	HostDataIP string `json:"hostData_IP"` //for iSCSI can be either port if multipathing is setup
}

// SolidfireStorageDriverConfig holds settings for SolidfireStorageDrivers
type SolidfireStorageDriverConfig struct {
	CommonStorageDriverConfig // embedded types replicate all fields
	TenantName                string
	EndPoint                  string
	DefaultVolSz              int64 //Default volume size in GiB
	SVIP                      string
	InitiatorIFace            string //iface to use of iSCSI initiator
	Types                     *[]sfapi.VolType
}

// Drivers is a map of driver names -> object
var Drivers = make(map[string]StorageDriver)

// StorageDriver provides a common interface for storage related operations
type StorageDriver interface {
	Name() string
	Initialize(string) error
	Validate() error
	Create(name string, opts map[string]string) error
	Destroy(name string) error
	Attach(name, mountpoint string, opts map[string]string) error
	Detach(name, mountpoint string) error
	DefaultStoragePrefix() string
}
