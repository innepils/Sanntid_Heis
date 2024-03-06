package cost

import (
	"driver/config"
	"driver/elevator"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

type HRAElevState struct {
	Behavior    string `json:"behaviour"`
	Floor       int    `json:"floor"`
	Direction   string `json:"direction"`
	CabRequests []bool `json:"cabRequests"`
}

type HRAInput struct {
	HallRequests [config.N_FLOORS][config.N_BUTTONS - 1]bool `json:"hallRequests"`
	States       map[string]HRAElevState                     `json:"states"`
}

func Cost(hall_requests [config.N_FLOORS][config.N_BUTTONS - 1]bool, localelevator elevator.Elevator, extern_elevators map[string]HRAElevState) [][2]bool { //REMEMBER TO CHANGE TYPES HERE

	input := HRAInput{
		HallRequests: hall_requests,
		States: map[string]HRAElevState{
			"self": HRAElevState{
				Behavior:    strings.ToLower(elevator.ElevBehaviourToString(localelevator.Behaviour)[3:]),
				Floor:       localelevator.Floor,
				Direction:   strings.ToLower(elevator.ElevDirnToString(localelevator.Dirn)),
				CabRequests: elevator.GetCabRequests(localelevator),
			},
		},
	}

	for key, value := range extern_elevators {
		input.States[key] = value
	}

	fmt.Println(input)

	jsonBytes, err := json.Marshal(input)
	if err != nil {
		fmt.Println("json.Marshal error: ", err)
		//die?
	}

	ret, err := exec.Command("./hall_request_assigner/hall_request_assigner.exe", "-i", string(jsonBytes)).CombinedOutput()
	if err != nil {
		fmt.Println("exec.Command error: ", err)
		fmt.Println(string(ret))
		//die?
	}

	output := new(map[string][][2]bool)
	err = json.Unmarshal(ret, &output)
	if err != nil {
		fmt.Println("json.Unmarshal error: ", err)
		//die?
	}

	return (*output)["self"]
}
