package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/godbus/dbus"
	"github.com/svenwltr/i3-statusbar/upower"
)

const (
	DATE_LAYOUT = "2. Jan 2006, 15:04:05"
	COLOR_LABEL = "#999999"
)

func main() {
	fmt.Println(`{ "version": 1 }`)
	fmt.Println(`[`)
	fmt.Println(`[]`)

	for {
		printLine(getLine())
		time.Sleep(time.Second)
	}
}

func printLine(line StatusLine) {
	bytes, err := json.Marshal(line.Lines)
	if err != nil {
		panic(err)
	}

	fmt.Print(`,`)
	fmt.Print(string(bytes))
	fmt.Println()
}

func getLine() StatusLine {
	var line StatusLine

	//line.AddLabel("Uptime: ")
	//line.Add().SetFullText(getUptime())
	line.Add().SetFullText(time.Now().Format(DATE_LAYOUT))

	line.Prepend(getPowerState())

	return line
}

func getPowerState() *StatusLine {
	const (
		DISCHARGING = "🔋"
		CHARGING    = "🔌"
	)

	var u *upower.Upower
	var p *upower.Paths
	var err error

	u, err = upower.New()
	if err != nil {
		log.Fatal(err)
	}

	p, err = u.Enumerate()
	if err != nil {
		log.Fatal(err)
	}

	var line StatusLine

	for _, path := range *p.Batteries() {
		line.AddLabel("Battery:")
		var bat = u.Details(path)
		var percentage int
		var timeToFull, timeToEmpty time.Duration
		var text, state string
		var err error

		percentage, err = bat.GetPercentage()
		if err != nil {
			panic(err)
		}

		state, err = bat.GetStateText()
		if err != nil {
			panic(err)
		}

		timeToFull, err = bat.GetTimeToFull()
		if err != nil {
			panic(err)
		}

		timeToEmpty, err = bat.GetTimeToEmpty()
		if err != nil {
			panic(err)
		}

		text = fmt.Sprintf("%d%% %s", percentage,
			state)

		if timeToFull != 0 {
			text += " " + timeToFull.String()
		}

		if timeToEmpty != 0 {
			text += " " + timeToEmpty.String()
		}

		line.Add().SetFullText(text)

	}

	return &line
}

func getPowerState_() *StatusLine {
	const (
		BAT_PREFIX    = "/org/freedesktop/UPower/devices/battery"
		LIST_METHOD   = "org.freedesktop.UPower.EnumerateDevices"
		DETAIL_METHOD = "org.freedesktop.DBus.Properties.GetAll"
		DETAIL_ARG1   = "org.freedesktop.UPower.Device"
	)

	var line StatusLine

	conn, err := dbus.SystemBus()
	if err != nil {
		log.Fatal(err)
	}
	//defer conn.Close()

	var devices = new([]dbus.ObjectPath)

	err = conn.
		Object("org.freedesktop.UPower", "/org/freedesktop/UPower").
		Call("org.freedesktop.UPower.EnumerateDevices", 0).
		Store(devices)
	if err != nil {
		log.Fatal(err)
	}

	for _, path := range *devices {
		if !strings.HasPrefix(string(path), BAT_PREFIX) {
			continue
		}

		var device = conn.
			Object("org.freedesktop.UPower", path)

		timeToFull, err := device.GetProperty("org.freedesktop.UPower.Device.TimeToFull")
		if err != nil {
			log.Fatal(err)
		}

		//timeToEmpty, err := device.GetProperty("org.freedesktop.UPower.Device.TimeToEmpty")
		if err != nil {
			log.Fatal(err)
		}

		percentage, err := device.GetProperty("org.freedesktop.UPower.Device.Percentage")
		if err != nil {
			log.Fatal(err)
		}

		line.AddLabel("Power:")

		switch v := percentage.Value().(type) {
		case float64:
			line.Add().
				SetFullText(fmt.Sprintf("%d%%", int64(v)))
		default:
			panic("Invalid type.")
		}

		switch v := timeToFull.Value().(type) {
		case int64:
			var left = strings.TrimRight(time.Duration(int(v/60)*int(time.Minute)).String(), "0s")
			line.Add().
				SetFullText(fmt.Sprintf("%s", left))
		default:
			panic("Invalid type.")
		}

	}

	return &line
}

func getUptime() string {
	bytes, err := ioutil.ReadFile("/proc/uptime")

	if err != nil {
		return fmt.Sprint(err)
	}

	secondstr := strings.Split(string(bytes), " ")[0]
	seconds, err := strconv.ParseFloat(secondstr, 64)
	duration := time.Duration(int(seconds) * int(time.Second))

	return duration.String()
}
