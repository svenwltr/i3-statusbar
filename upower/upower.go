package upower

import (
	"fmt"
	"github.com/godbus/dbus"
	"strings"
	"time"
)

type Paths []dbus.ObjectPath

func (d *Paths) Batteries() *Paths {
	const (
		PREFIX = "/org/freedesktop/UPower/devices/battery"
	)

	var sub = make(Paths, 0)
	for _, path := range *d {
		if strings.HasPrefix(string(path), PREFIX) {
			sub = append(sub, path)
		}
	}

	return &sub

}

type Upower struct {
	conn *dbus.Conn
}

func New() (*Upower, error) {
	var u *Upower = new(Upower)
	var err error

	u.conn, err = dbus.SystemBus()
	if err != nil {
		return nil, err
	}

	return u, nil

}

func (u *Upower) Enumerate() (*Paths, error) {
	const (
		OBJECT_PATH = "/org/freedesktop/UPower"
		DEST        = "org.freedesktop.UPower"
		INTERFACE   = "org.freedesktop.UPower.EnumerateDevices"
	)
	var devices = new([]dbus.ObjectPath)
	var err error

	err = u.conn.Object(DEST, OBJECT_PATH).Call(INTERFACE, 0).Store(devices)
	if err != nil {
		return nil, err
	}

	var d = Paths(*devices)
	return &d, nil

}

type Device struct {
	obj dbus.BusObject
}

func (u *Upower) Details(path dbus.ObjectPath) *Device {
	const (
		DEST = "org.freedesktop.UPower"
	)

	var dev = new(Device)
	dev.obj = u.conn.Object(DEST, path)
	return dev

}

func (d *Device) getProperty(key string) (dbus.Variant, error) {
	const PREFIX = "org.freedesktop.UPower.Device."
	return d.obj.GetProperty(PREFIX + key)

}

func (d *Device) getPropertyFloat64(key string) (float64, error) {
	var err error
	var variant dbus.Variant

	variant, err = d.getProperty(key)
	if err != nil {
		return 0., err
	}

	switch v := variant.Value().(type) {
	case float64:
		return v, nil
	default:
		return 0., fmt.Errorf("Invalid type.")
	}
}

func (d *Device) getPropertyInt64(key string) (int64, error) {
	var err error
	var variant dbus.Variant

	variant, err = d.getProperty(key)
	if err != nil {
		return 0., err
	}

	switch v := variant.Value().(type) {
	case int64:
		return v, nil
	default:
		return 0., fmt.Errorf("Invalid type.")
	}
}

func (d *Device) getPropertyUint32(key string) (uint32, error) {
	var err error
	var variant dbus.Variant

	variant, err = d.getProperty(key)
	if err != nil {
		return 0., err
	}

	switch v := variant.Value().(type) {
	case uint32:
		return v, nil
	default:
		return 0., fmt.Errorf("Invalid type.")
	}
}

func (d *Device) GetPercentage() (int, error) {
	f, err := d.getPropertyFloat64("Percentage")
	if err != nil {
		return 0, err
	}

	return int(f), nil

}

func (d *Device) GetTimeToFull() (time.Duration, error) {
	i, err := d.getPropertyInt64("TimeToFull")
	if err != nil {
		return 0, err
	}

	return time.Duration(int(i) * int(time.Second)), nil

}

func (d *Device) GetTimeToEmpty() (time.Duration, error) {
	i, err := d.getPropertyInt64("TimeToEmpty")
	if err != nil {
		return 0, err
	}

	return time.Duration(int(i) * int(time.Second)), nil

}

type State uint32

const (
	UNKNOWN           State = iota
	CHARGING          State = iota
	DISCHARGING       State = iota
	EMPTY             State = iota
	FULLY_CHARGED     State = iota
	PENDING_CHARGE    State = iota
	PENDING_DISCHARGE State = iota
)

func (d *Device) GetState() (State, error) {
	i, err := d.getPropertyUint32("State")
	if err != nil {
		return 0, err
	}

	return State(i), nil

}

func (d *Device) GetStateText() (string, error) {
	var state State
	var err error

	state, err = d.GetState()
	if err != nil {
		return "", err
	}

	switch state {
	case CHARGING:
		return "⚡", nil
	case FULLY_CHARGED:
		return "⚡", nil
	case DISCHARGING:
		return "🔋", nil
	default:
		return "?", nil
	}

}
