package synology

import (
	"strconv"
	"strings"
)

// https://global.download.synology.com/download/Document/Software/DeveloperGuide/Package/SurveillanceStation/All/enu/Surveillance_Station_Web_API.pdf

const (
	surveillanceStationVersion = 9

	// SurveillanceStation api names
	surveillanceStationCameraApiName = "SYNO.SurveillanceStation.Camera"

	// SurveillanceStationCamera method names
	listMethod    = "List"
	disableMethod = "Disable"
	enableMethod  = "Enable"
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
	NewName string
	Id      int
	Ip      string
	Port    int
	Status  CameraStatus
}

type SurveillanceStationInfoCameraListResponse struct {
	Total   int
	Cameras []Camera
}

type SurveillanceStation struct {
	core      *SynologyCore
	cameraApi *Api
}

// SurveillanceStationCameraList get the list of all cameras.
func (s *SurveillanceStation) SurveillanceStationCameraList() ([]Camera, error) {
	response, err := s.core.makeRequest(s.cameraApi.Path, s.cameraApi.Name, listMethod, surveillanceStationVersion, nil)
	if err != nil {
		return nil, err
	}

	data, err := unmarshal[SurveillanceStationInfoCameraListResponse](response)
	return data.Cameras, err
}

// mapCameraIds maps camera slice into slice of camera ids.
func mapCameraIds(cameras []Camera) string {
	ids := make([]string, len(cameras))
	for i, camera := range cameras {
		ids[i] = strconv.Itoa(camera.Id)
	}

	return strings.Join(ids, ",")
}

// SurveillanceStationCameraDisable disables cameras.
func (s *SurveillanceStation) SurveillanceStationCameraDisable(cameras []Camera) error {
	response, err := s.core.makeRequest(
		s.cameraApi.Path,
		s.cameraApi.Name,
		disableMethod,
		surveillanceStationVersion,
		map[string]string{"idList": mapCameraIds(cameras)})

	if err != nil {
		return err
	}

	_, err = unmarshal[struct{}](response)

	return err
}

// SurveillanceStationCameraEnable enables cameras.
func (s *SurveillanceStation) SurveillanceStationCameraEnable(cameras []Camera) error {
	response, err := s.core.makeRequest(
		s.cameraApi.Path,
		s.cameraApi.Name,
		enableMethod,
		surveillanceStationVersion,
		map[string]string{"idList": mapCameraIds(cameras)})

	if err != nil {
		return err
	}

	_, err = unmarshal[struct{}](response)

	return err
}
