package hue

import (
	"bytes"
	"encoding/json"
)

// Light - encapsulates the controls for a specific philips hue light
type Light struct {
	ID     string
	Name   string
	bridge *Bridge
}

// LightState ...
type LightState struct {
	Hue       int       `json:"hue"`
	On        bool      `json:"on"`
	Effect    string    `json:"effect"`
	Alert     string    `json:"alert"`
	Bri       int       `json:"bri"`
	Sat       int       `json:"sat"`
	Ct        int       `json:"ct"`
	Xy        []float32 `json:"xy"`
	Reachable bool      `json:"reachable"`
	ColorMode string    `json:"colormode"`
}

// SetLightState ...
type SetLightState struct {
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

// LightAttributes ...
type LightAttributes struct {
	State           LightState        `json:"state"`
	Type            string            `json:"type"`
	Name            string            `json:"name"`
	ModelID         string            `json:"modelid"`
	SoftwareVersion string            `json:"swversion"`
	PointSymbol     map[string]string `json:"pointsymbol"`
}

// GetLightAttributes - retrieves light attributes and state as per
// http://developers.meethue.com/1_lightsapi.html#14_get_light_attributes_and_state
func (l *Light) GetLightAttributes() (*LightAttributes, error) {
	response, err := l.bridge.get("/lights/" + l.ID)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	result := new(LightAttributes)
	err = json.NewDecoder(response.Body).Decode(&result)
	return result, err
}

// SetName - sets the name of a light as per
// http://developers.meethue.com/1_lightsapi.html#15_set_light_attributes_rename
func (l *Light) SetName(newName string) ([]Result, error) {
	params := map[string]string{"name": newName}
	data, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	response, err := l.bridge.put("/lights/"+l.ID, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var results []Result
	err = json.NewDecoder(response.Body).Decode(&results)
	return results, err
}

// On - a convenience method to turn on a light and set its effect to "none"
func (l *Light) On() ([]Result, error) {
	state := SetLightState{
		On:     true,
		Effect: "none",
	}
	return l.SetState(state)
}

// Off - a convenience method to turn off a light
func (l *Light) Off() ([]Result, error) {
	state := SetLightState{On: false}
	return l.SetState(state)
}

// ColorLoop - a convenience method to turn on a light and have it begin
// a colorloop effect
func (l *Light) ColorLoop() ([]Result, error) {
	state := SetLightState{
		On:     true,
		Effect: "colorloop",
	}
	return l.SetState(state)
}

// SetState sets the state of a light as per
// http://developers.meethue.com/1_lightsapi.html#16_set_light_state
func (l *Light) SetState(state SetLightState) ([]Result, error) {
	data, err := json.Marshal(state)
	if err != nil {
		return nil, err
	}

	response, err := l.bridge.put("/lights/"+l.ID+"/state", bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var results []Result
	err = json.NewDecoder(response.Body).Decode(&results)
	return results, err
}
