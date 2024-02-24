package cost

import (
	"encoding/json"
	"fmt"
	"os/exec"
)

type HRAElevState struct {
	Behavior    string `json:"behaviour"`
	Floor       int    `json:"floor"`
	Direction   string `json:"direction"`
	CabRequests []bool `json:"cabRequests"`
}

type HRAInput struct {
	HallRequests [][2]bool               `json:"hallRequests"`
	States       map[string]HRAElevState `json:"states"`
}

func Cost(hall_requests [][2]bool, elevator Elevator, extern_elevators map[string]HRAElevState) [][2]bool { //REMEMBER TO CHANGE TYPES HERE

	input := HRAInput{
		HallRequests: hall_requests,
		States: map[string]HRAElevState{
			"self": HRAElevState{
				Behavior:    elevator.Behaviour,
				Floor:       elevator.Floor,
				Direction:   elevator.Direction,
				CabRequests: elevator.CabRequests,
			},
		},
		extern_elevators,
	}

	jsonBytes, err := json.Marshal(input)
	if err != nil {
		fmt.Println("json.Marshal error: ", err)
		//die?
	}

	ret, err := exec.Command("../hall_request_assigner/hall_request_assigner", "-i", string(jsonBytes)).CombinedOutput()
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

	for _, value := range *output {
		return value
	}

	fmt.Println("Cost function terminated without output.")
	//die?
}
