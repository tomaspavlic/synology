package main

import (
	"strconv"
	"strings"
)

// https://global.download.synology.com/download/Document/Software/DeveloperGuide/Package/SurveillanceStation/All/enu/Surveillance_Station_Web_API.pdf

const (
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

func (s *SynologyCore) SurveillanceStationCameraList() ([]Camera, error) {
	info, err := s.Find(surveillanceStationCameraApiName)
	if err != nil {
		return nil, err
	}

	response, err := s.makeRequest(info.Path, info.Name, listMethod, 9, nil)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	data, err := readResponse[SurveillanceStationInfoCameraListResponse](response)
	return data.Cameras, err
}

func mapCameraIds(cameras []Camera) string {
	ids := make([]string, len(cameras))
	for i, camera := range cameras {
		ids[i] = strconv.Itoa(camera.Id)
	}

	return strings.Join(ids, ",")
}

func (s *SynologyCore) SurveillanceStationCameraDisable(cameras []Camera) error {
	info, err := s.Find(surveillanceStationCameraApiName)
	if err != nil {
		return err
	}

	response, err := s.makeRequest(info.Path, info.Name, disableMethod, 9, map[string]string{"idList": mapCameraIds(cameras)})
	if err != nil {
		return err
	}
	defer response.Body.Close()

	_, err = readResponse[struct{}](response)

	return err
}

func (s *SynologyCore) SurveillanceStationCameraEnable(cameras []Camera) error {
	info, err := s.Find(surveillanceStationCameraApiName)
	if err != nil {
		return err
	}

	response, err := s.makeRequest(info.Path, info.Name, enableMethod, 9, map[string]string{"idList": mapCameraIds(cameras)})
	if err != nil {
		return err
	}
	defer response.Body.Close()

	_, err = readResponse[struct{}](response)

	return err
}
