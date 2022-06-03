package aldlstruct

import (
	"fmt"
)

var b1 uint8
var ct map[uint8]int = map[uint8]int{
	255: -40,
	251: -30,
	250: -25,
	247: -20,
	245: -15,
	241: -10,
	237: -5,
	231: 0,
	225: 5,
	218: 10,
	209: 15,
	199: 20,
	189: 25,
	177: 30,
	165: 35,
	152: 40,
	139: 45,
	126: 50,
	114: 55,
	102: 60,
	92:  65,
	81:  70,
	72:  75,
	64:  80,
	56:  85,
	50:  90,
	44:  95,
	39:  100,
	34:  105,
	30:  110,
	26:  115,
	23:  120,
	21:  125,
	18:  130,
	16:  135,
	14:  140,
	13:  145,
	12:  150,
	0:   200,
}

func new136() aldler {
	res := aldl{
		FrameLength: 23,
		flags: map[uint8]string{
			0:  "", //"MW2",
			11: "", //"MALFUNCTION FLG1",
			12: "", //"MALFUNCTION FLG2",
			13: "", //"MALFUNCTION FLG3",
			14: "", //"AIR FUEL MODE",
			16: "", //"MCU2IO",
		},
		dataStruct: map[uint8]string{
			1:  "",
			2:  "EPROM ID",
			3:  "IAC PRESENT POSITION",
			4:  "COOLANT TEMPERATURE",
			5:  "VEHICLE SPEED",
			6:  "MANIFOLD ABSOLUTE PRESSURE",
			7:  "ENGINE RPM",
			8:  "TPS VOLTS",
			9:  "CLOSED LOOP INTEGRATOR",
			10: "OXYGEN SENSOR",
			15: "VOLTAGE",
			17: "KNOCK SENSOR COUNTER",
			18: "BASE PULSE COURSE CORRECTION(BLM)",
			19: "O2 CROSS COUNTS",
			20: "FUEL PUMP POWER SWITCH VOLTAGE",
			21: "DESIRED IDLE RPM",
			22: "RESCALED TPS",
		},
		parseFunctions: []func(uint8, uint8) string{
			mw2,
			getEPROMUID1,
			getEPROMUID2,
			getIAC,
			getCoolantTemp,
			mph,
			mainfoldPressure,
			rpm,
			tps,
			toHex,
			oxygen,
			malf1,
			malf2,
			malf3,
			airFuelMode,
			voltage,
			mcu2io,
			knockSensor,
			blm,
			aldlxntr,
			fuelPumpVoltage,
			desiredIdle,
			rescTPS,
		},
	}
	return res
}

func (a aldl) GetFrameLength() uint8 {
	return a.FrameLength
}

func (a aldl) GetDataStruct() map[uint8]string {
	return a.dataStruct
}

func (a aldl) GetFlags() map[uint8]string {
	return a.flags
}

func (a aldl) TranslateByte(pos, data uint8) string {
	return a.parseFunctions[pos](pos, data)
}

func getEPROMUID1(pos, data uint8) string {
	b1 = data
	return ""
}

func getEPROMUID2(pos, data uint8) string {
	inte := int(b1)*256 + int(data)
	return fmt.Sprintf("%d (%d+%d)", inte, b1, data)
}

func getIAC(pos, data uint8) string {
	return fmt.Sprintf("%d counts", data)
}

func knockSensor(pos, data uint8) string {
	return getIAC(pos, data)
}

func aldlxntr(pos, data uint8) string {
	return getIAC(pos, data)
}

func getCoolantTemp(pos, data uint8) string {
	return fmt.Sprintf("%d C", ct[data])
}

func mph(pos, data uint8) string {
	return fmt.Sprintf("%d MPH", data)
}

func blm(pos, data uint8) string {
	return fmt.Sprintf("%d", data)
}

func mainfoldPressure(pos, data uint8) string {
	return fmt.Sprintf("%.4f Volts", float32(data)/51)
}

func rpm(pos, data uint8) string {
	return fmt.Sprintf("%d RPM", int(data)*25)
}

func tps(pos, data uint8) string {
	return mainfoldPressure(pos, data)
}

func oxygen(pos, data uint8) string {
	return fmt.Sprintf("%.2f mVolts", float32(data)*4.42)
}

func fuelPumpVoltage(pos, data uint8) string {
	return fmt.Sprintf("%.2f Volts", float32(data)*0.1)
}

func voltage(pos, data uint8) string {
	return fmt.Sprintf("%.2f mVolts", float32(data)*0.1)
}

func desiredIdle(pos, data uint8) string {
	return fmt.Sprintf("%.2f RPM", float32(data)*12.5)
}

func rescTPS(pos, data uint8) string {
	return fmt.Sprintf("%.2f %%", float32(data)/2.55)
}

func mw2(pos, data uint8) string {

	dec := []string{
		"ROAD SPEED PULSE OCCURRED(6.25 MSEC CHECK)",
		"ESC 43B READY FOR SECOND P.E.",
		"REFERENCE PULSE OCCURRED(6.25 MSEC CHECK)",
		"DIAGNOSTIC SWITCH IN FACTORY TEST POSITION",
		"DIAGNOSTIC SWITCH IN DIAGNOSTIC POSITION",
		"DIAGNOSTIC SWITCH IN ALDL POSITION",
		"HIGH BATTERY VOLTAGE-DISABLE MCU SOLENOID DISCRTS",
		"",
	}

	res := ""
	for i, k := range getBits(data) {
		if k {
			res = fmt.Sprintf("%s%s\r\n", res, dec[i])
		}
	}

	return res
}

func malf1(pos, data uint8) string {

	dec := []string{
		"CODE 24  VEHICLE SPEED SENSOR",
		"CODE 23  not used",
		"CODE 22  THROTTLE POSITION SENSOR LOW",
		"CODE 21  THROTTLE POSITION SENSOR HIGH",
		"CODE 15  COOLANT SENSOR LOW TEMPERATURE",
		"CODE 14  COOLANT SENSOR HIGH TEMPERATURE",
		"CODE 13  OXYGEN SENSOR",
		"CODE 12  NO REFERENCE PULSES(ENGINE NOT RUNNING)",
	}

	res := ""
	for i, k := range getBits(data) {
		if k {
			res = fmt.Sprintf("%s%s\r\n", res, dec[i])
		}
	}

	return res
}

func malf2(pos, data uint8) string {

	dec := []string{
		"CODE 42  EST MONITOR ERROR",
		"CODE 41  not used",
		"CODE 35  not used",
		"CODE 34  MAP SENSOR LOW",
		"CODE 33  MAP SENSOR HIGH",
		"CODE 32  EGR FAILURE",
		"CODE 31  not used",
		"CODE 25  not used",
	}

	res := "" // fmt.Sprintf("%b, %v\r\n", dec, getBits(data))
	for i, k := range getBits(data) {
		if k {
			res = fmt.Sprintf("%s%s\r\n", res, dec[i])
		}
	}

	return res
}

func malf3(pos, data uint8) string {

	dec := []string{
		"CODE 55  ADU ERROR",
		"1 = MALF CODE 54  FUEL PUMP RELAY MALFUNCTION",
		"CODE 53  not used",
		"CODE 52  CAL-PAK MISSING",
		"CODE 51  PROM ERROR",
		"CODE 45  OXYGEN SENSOR RICH",
		"CODE 44  OXYGEN SENSOR LEAN",
		"CODE 43  ESC FAILURE",
	}

	res := ""
	for i, k := range getBits(data) {
		if k {
			res = fmt.Sprintf("%s%s\r\n", res, dec[i])
		}
	}

	return res
}

func airFuelMode(pos, data uint8) string {

	dec := []string{
		"X",
		"LEARN CONTROL ENABLE",
		"LOW BATTERY",
		"4-3 DOWNSHIFT FLAG FOR TCC UNLOCK",
		"ASYNCHRONOUS PULSE FLAG (AP FLAG)",
		"OLD HIGH GEAR FLAG, 0= HIGH GEAR LAST TIME",
		"REACH",
		"CLOSED LOOP FLAG (1=CLOSED LOOP, 0=OPEN LOOP)",
	}

	res := ""
	for i, k := range getBits(data) {
		switch i {
		case 5:
			if k {
				res = fmt.Sprintf("%sOLD HIGH GEAR\r\n", res)
			}
			break
		case 6:
			if k {
				res = fmt.Sprintf("%sREACH\r\n", res)
			} else {
				res = fmt.Sprintf("%sLEAN\r\n", res)
			}
			break
		case 7:
			if k {
				res = fmt.Sprintf("%sCLOSED LOOP\r\n", res)
			} else {
				res = fmt.Sprintf("%sOPEN LOOP\r\n", res)
			}
			break
		default:
			if k {
				res = fmt.Sprintf("%s%s\r\n", res, dec[i])
			}
			break
		}
	}

	return res
}

func mcu2io(pos, data uint8) string {

	dec := []string{
		"FAN #1 RELAY ENGAGED",
		"FAN #2 RELAY ENGAGED",
		"A/C ENABLED, SOLENOID ON",
		"TCC CONVERTER LOCKED, SOLENOID ON",
		"PARK/NEUTRAL",
		"HIGH GEAR ACTIVE",
		"HEATED WINDSHIELD SWITCH ON",
		"AIR CONDITIONER",
	}

	res := ""
	for i, k := range getBits(data) {
		if k {
			res = fmt.Sprintf("%s%s\r\n", res, dec[i])
		}
	}

	return res
}
