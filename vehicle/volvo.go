package vehicle

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

	"github.com/andig/evcc/api"
	"github.com/andig/evcc/provider"
	"github.com/andig/evcc/util"
)

// "customeraccounts"
// "attributes"
// "status"
// "position"
// "trips"

const (
	volvoAPI = "https://vocapi.wirelesscar.net/customerapi/rest/v3.0/"
)

type volvoDynamicResponse interface{}

type volvoVehiclesResponse struct {
	AccountVehicleRelations []struct {
	} `json:"accountVehicleRelations"`
}

// Volvo is an api.Vehicle implementation for Volvo cars
type Volvo struct {
	*embed
	*util.HTTPHelper
	user, password, vin string
	token               string
	tokenValid          time.Time
	chargeStateG        func() (float64, error)
}

func init() {
	registry.Add("volvo", NewVolvoFromConfig)
}

// NewVolvoFromConfig creates a new vehicle
func NewVolvoFromConfig(other map[string]interface{}) (api.Vehicle, error) {
	cc := struct {
		Title               string
		Capacity            int64
		User, Password, VIN string
		Cache               time.Duration
	}{}
	if err := util.DecodeOther(other, &cc); err != nil {
		return nil, err
	}

	log := util.NewLogger("volvo")

	v := &Volvo{
		embed:      &embed{cc.Title, cc.Capacity},
		HTTPHelper: util.NewHTTPHelper(log),
		user:       cc.User,
		password:   cc.Password,
		vin:        cc.VIN,
	}

	if cc.VIN == "" {
		vehicles, err := v.vehicles()
		_ = vehicles
		_ = err
		// if err != nil {
		// 	return nil, fmt.Errorf("cannot get vehicles: %v", err)
		// }

		// if len(vehicles) != 1 {
		// 	return nil, fmt.Errorf("cannot find vehicle: %v", vehicles)
		// }

		// v.vin = vehicles[0].Vin
		// log.DEBUG.Printf("found vehicle: %v", v.vin)
	}

	v.chargeStateG = provider.NewCached(v.chargeState, cc.Cache).FloatGetter()

	return v, nil
}

func (v *Volvo) request(uri string) (*http.Request, error) {
	req, err := http.NewRequest(http.MethodGet, volvoAPI+uri, nil)
	if err == nil {
		basicAuth := base64.StdEncoding.EncodeToString([]byte(v.user + ":" + v.password))
		req.Header.Set("Authorization", fmt.Sprintf("Basic %s", basicAuth))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Device-Id", "Device")
		req.Header.Set("X-OS-Type", "Android")
		req.Header.Set("X-Originator-Type", "App")
		req.Header.Set("X-OS-Version", "22")
	}

	return req, nil
}

// vehicles implements returns the list of user vehicles
func (v *Volvo) vehicles() (volvoVehiclesResponse, error) {
	var resp volvoVehiclesResponse

	req, err := v.request("customeraccounts")
	if err == nil {
		var b []byte
		b, err = v.RequestJSON(req, &resp)
		println(string(b))
	}

	return resp, err
}

// chargeState implements the Vehicle.ChargeState interface
func (v *Volvo) chargeState() (float64, error) {
	var resp volvoDynamicResponse

	req, err := v.request("status")
	if err != nil {
		return 0, err
	}

	_, err = v.RequestJSON(req, &resp)
	return 0, err
}

// ChargeState implements the Vehicle.ChargeState interface
func (v *Volvo) ChargeState() (float64, error) {
	return v.chargeStateG()
}
