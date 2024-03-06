package cost

import (
	"driver/elevator"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

func Cost(
	hall_requests [][2]bool,
	localElevator elevator.Elevator,
	externalElevators map[string]elevator.ElevatorState) [][2]bool { //REMEMBER TO CHANGE TYPES HERE

	input := elevator.HRAInput{
		HallRequests: hall_requests,
		ElevatorState: map[string]elevator.ElevatorState{
			"aaa_self": elevator.ElevatorState{
				Behavior:    strings.ToLower(elevator.ElevBehaviourToString(localElevator.Behaviour)[3:]),
				Floor:       localElevator.Floor,
				Direction:   strings.ToLower(elevator.ElevDirnToString(localElevator.Dirn)),
				CabRequests: elevator.GetCabRequests(localElevator),
			},
		},
	}

	for key, value := range externalElevators {
		input.ElevatorState[key] = value
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
