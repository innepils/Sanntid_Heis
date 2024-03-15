package backup

import (
	"src/config"
	"src/elevator"
	"src/elevator_io"
	"encoding/gob"
	"fmt"
	"os"
)

func SaveBackupToFile(filename string, allRequests [config.N_FLOORS][config.N_BUTTONS]elevator.RequestType) {
	var cabRequests [config.N_FLOORS]bool
	for floor := range allRequests {
		if allRequests[floor][config.N_BUTTONS-1] == elevator.ConfirmedRequest {
			cabRequests[floor] = true
		} else {
			cabRequests[floor] = false
		}
	}
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	err = encoder.Encode(cabRequests)
	if err != nil {
		return
	}
}

func LoadBackupFromFile(filename string, ch_buttonPressed chan elevator_io.ButtonEvent) {
	var data [config.N_FLOORS]bool

	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error decoding data from backup")
	}
	defer file.Close()

	decoder := gob.NewDecoder(file)
	err = decoder.Decode(&data)
	if err != nil {
		fmt.Println("Error decoding data from backup")
	}

	for i, element := range data {
		if element {
			ch_buttonPressed <- elevator_io.ButtonEvent{BtnFloor: i, BtnType: elevator_io.BT_Cab}
		}
	}
}
