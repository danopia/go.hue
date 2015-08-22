package hue

import (
	"bytes"
	"encoding/json"
)

// GroupState ...
type GroupState struct {
	// On/Off state of the light. On=true, Off=false
	On bool `json:"on,omitempty"`

	// The brightness value to set the light to. Brightness is a
	// scale from 0 (the minimum the light is capable of) to 255
	// (the maximum).
	//
	// NOTE: Brightness of 0 is not off.
	Bri uint8 `json:"bri,omitempty"`

	// The hue value to set light to. The hue value is a wrapping
	// value between 0 and 65535. Both 0 and 65535 are red, 25500
	// is green and 46920 is blue.
	Hue uint16 `json:"hue,omitempty"`

	// Saturation of the light. 255 is the most saturated (colored)
	// and 0 is the least saturated (white).
	Sat uint8 `json:"sat,omitempty"`

	// The X and Y coordinates of a color in CIE color space.
	//
	// The first entry is the X coordinate, and the second entry
	// is the Y coordinate. Both X and Y must be between 0 and 1.
	//
	// If the specified coordinates are not in the CIE color space,
	// the closest color to the coordinates will be chosen.
	Xy []float64 `json:"xy,omitempty"`

	// The Mired Color temperature of the light. 2012 connected lights
	// are capable of 153 (6500K) to 500 (2000K).
	Ct uint16 `json:"ct,omitempty"`

	// The alert effect, which is a temporary change to the bulb’s state,
	// and has one of the following values:
	//
	// 'none'    – The light does not perform an alert effect.
	// 'select'  – The light performs one breathe cycle.
	// 'lselect' – The light performs breathe cycles for 15 seconds,
	//             or until an Alert of 'none' command is received.
	//
	// NOTE: This contains the last alert sent to the light and not it's
	// current state, i.e. after the breathe cycle has finished, the bridge
	// does not reset the alert to 'none'
	Alert string `json:"alert,omitempty"`

	// The dynamic effect of the light, currently 'none' and 'colorloop' are supported.
	// Other values will generate an error of type 7.
	//
	// Setting the effect to 'colorloop' will cycle through all hues using the current
	// brightness and saturation settings.
	Effect string `json:"effect,omitempty"`

	// The duration of the transition from the light's current state to the new state.
	// This is given as a multiple of 100ms and defaults to 400ms.
	//
	// Example: setting to `10` will make the transition last 1 second.
	TransitionTime uint16 `json:"transitiontime,omitempty"`

	// The scene identifier if the scene you wish to recall (optional)
	Scene string `json:"scene,omitempty"`
}

// SetGroupState sets the state of a group:
// http://www.developers.meethue.com/documentation/groups-api#25_set_group_state
func (b *Bridge) SetGroupState(groupID string, state GroupState) ([]Result, error) {
	data, err := json.Marshal(state)
	if err != nil {
		return nil, err
	}

	response, err := b.put("/groups/"+groupID+"/action", bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var results []Result
	err = json.NewDecoder(response.Body).Decode(&results)
	return results, err
}
