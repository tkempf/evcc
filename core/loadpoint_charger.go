package core

import (
	"fmt"
	"time"
)

// SyncEnabled synchronizes charger settings to expected state
func (lp *LoadPoint) SyncEnabled() {
	enabled, err := lp.charger.Enabled()
	if err == nil && enabled != lp.enabled {
		lp.log.WARN.Printf("sync enabled state to %s", status[lp.enabled])
		err = lp.charger.Enable(lp.enabled)
	}

	if err != nil {
		lp.log.ERROR.Printf("charge controller error: %v", err)
	}
}

// chargerEnable switches charging on or off. Minimum cycle duration is guaranteed.
func (lp *LoadPoint) chargerEnable(enable bool) error {
	if remaining := (lp.GuardDuration - lp.clock.Since(lp.guardUpdated)).Truncate(time.Second); remaining > 0 {
		lp.log.DEBUG.Printf("charger %s - contactor delay %v", status[enable], remaining)
		return nil
	}

	if lp.enabled != enable {
		if err := lp.charger.Enable(enable); err != nil {
			return fmt.Errorf("charge controller error: %v", err)
		}

		lp.enabled = enable // cache
		lp.log.INFO.Printf("charger %s", status[enable])
		lp.guardUpdated = lp.clock.Now()
	} else {
		lp.log.DEBUG.Printf("charger %s", status[enable])
	}

	// if not enabled, current will be reduced to 0 in handler
	lp.bus.Publish(evChargeCurrent, lp.MinCurrent)

	return nil
}

// setTargetCurrent guards setting current against changing to identical value
// and violating MaxCurrent
func (lp *LoadPoint) setTargetCurrent(targetCurrent int64) error {
	target := clamp(targetCurrent, lp.MinCurrent, lp.MaxCurrent)

	if lp.targetCurrent != target {
		lp.log.DEBUG.Printf("set charge current: %dA", target)
		if err := lp.charger.MaxCurrent(target); err != nil {
			return fmt.Errorf("charge controller error: %v", err)
		}

		lp.targetCurrent = target // cache
	}

	// if not enabled, current will be reduced to 0 in handler
	lp.bus.Publish(evChargeCurrent, target)

	return nil
}

// Ramp performs ramping charger current up and down where targetCurrent=0
// signals disabled state
func (lp *LoadPoint) Ramp(targetCurrent int64, force ...bool) error {
	// reset guard updated
	if len(force) == 1 && force[0] {
		lp.guardUpdated = time.Time{}
	}

	// if targetCurrent == 0 disable
	if targetCurrent == 0 {
		return lp.chargerEnable(false)
	}

	// else set targetCurrent and optionally enable
	err := lp.setTargetCurrent(targetCurrent)
	if err == nil && !lp.enabled {
		err = lp.chargerEnable(true)
	}

	return err
}
