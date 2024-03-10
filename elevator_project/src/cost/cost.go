package cost

import (
	"driver/config"
	"driver/elevator"
	"encoding/json"
	"fmt"
	"os/exec"
)

func Cost(
	hallRequests [config.N_FLOORS][config.N_BUTTONS - 1]bool,
	localElevator map[string]elevator.ElevatorState,
	externalElevators map[string]elevator.ElevatorState) [][2]bool { //REMEMBER TO CHANGE TYPES HERE

	input := elevator.HRAInput{
		HallRequests:  hallRequests,
		ElevatorState: localElevator,
	}

	for key, value := range externalElevators {
		input.ElevatorState[key] = value
	}

	jsonBytes, err := json.Marshal(input)
	if err != nil {
		fmt.Println("json.Marshal error: ", err)
		//die?
	}

	ret, err := exec.Command("./hall_request_assigner/hall_request_assigner", "-i", string(jsonBytes)).CombinedOutput()
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
