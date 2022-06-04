package synology

import "strconv"

// https://global.download.synology.com/download/Document/Software/DeveloperGuide/Package/SurveillanceStation/All/enu/Surveillance_Station_Web_API.pdf

const (
	surveillanceStationVersion = 9

	// SurveillanceStation api names
	surveillanceStationCameraApiName = "SYNO.SurveillanceStation.Camera"

	// SurveillanceStationCamera method names
	surveillanceStationCameraListMethod    = "List"
	surveillanceStationCameraDisableMethod = "Disable"
	surveillanceStationCameraEnableMethod  = "Enable"
)

type CameraStatus uint8

const (
	Normal CameraStatus = iota + 1
	Deleted
	Disconnected
	Unavailable
	Ready
	Inaccessible
	Disabled
	Unrecognized
	Setting
	ServerDisconnected
	Migrating
	Others
	StorageRemoved
	Stopping
	ConnectHistFailed
	Unauthorized
	RTSPError
	NoVideo
)

type Camera struct {
	core *SynologyCore

	NewName string
	Id      int
	Ip      string
	Port    int
	Status  CameraStatus
}

type surveillanceStationInfoCameraListResponse struct {
	Total   int
	Cameras []*Camera
}

type SurveillanceStation struct {
	*SynologyCore
}

// List get the list of all cameras.
func (s *SurveillanceStation) List() ([]*Camera, error) {
	response, err := s.SynologyCore.makeRequest(
		entryPath,
		surveillanceStationCameraApiName,
		surveillanceStationCameraListMethod,
		surveillanceStationVersion,
		nil,
	)

	if err != nil {
		return nil, err
	}

	data, err := unmarshal[surveillanceStationInfoCameraListResponse](response)

	for _, camera := range data.Cameras {
		camera.core = s.SynologyCore
	}

	return data.Cameras, err
}

func (s *Camera) changeCamerasState(method string) error {
	response, err := s.core.makeRequest(
		entryPath,
		surveillanceStationCameraApiName,
		method,
		surveillanceStationVersion,
		parameters{"idList": strconv.Itoa(s.Id)},
	)

	if err != nil {
		return err
	}

	_, err = unmarshal[struct{}](response)

	return err
}

// Disable disables cameras.
func (s *Camera) Disable() error {
	return s.changeCamerasState(surveillanceStationCameraDisableMethod)
}

// Enable enables cameras.
func (s *Camera) Enable() error {
	return s.changeCamerasState(surveillanceStationCameraEnableMethod)
}
