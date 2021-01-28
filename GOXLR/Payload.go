package GOXLR

type Payload struct {
}

func (p Payload) Routing(action ActionType, input Input, output Output) []byte {
	return []byte(`{
			"action": "com.tchelicon.goxlr.routingtable",
			"event": "keyUp",
			"payload": {
				"settings":{
					"RoutingAction": "` + string(action) + `",
					"RoutingInput": "` + string(input) + `", 
					"RoutingOutput": "` + string(output) + `"
				}
			}
		}`)
}
