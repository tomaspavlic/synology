package main

import "fmt"

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	syn := NewSynologyCore("10.180.0.3", 5001)
	if err := syn.Login("tomas", "133789"); err != nil {
		panic(err)
	}

	cameras, err := syn.SurveillanceStationCameraList()
	fmt.Println(cameras, err)

	err = syn.SurveillanceStationCameraDisable(cameras)
	fmt.Println(err)
}
