package main

import (
	"fmt"
	"github.com/stianeikeland/go-rpio"
	"os"
	"time"
)

var (
	led1, led2, led3, led4, led5, led6, led7, led8 rpio.Pin
	rot_sw, rot_dt, rot_clk                        rpio.Pin
)

const (
	PIN_LED_1 = 14
	PIN_LED_2 = 15
	PIN_LED_3 = 16
	PIN_LED_4 = 17
	PIN_LED_5 = 18
	PIN_LED_6 = 19
	PIN_LED_7 = 20
	PIN_LED_8 = 21

	PIN_ROT_SW  = 4
	PIN_ROT_DT  = 5
	PIN_ROT_CLK = 26
)

func init() {

}

func main() {

	if err := rpio.Open(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer rpio.Close()
	setup()
	rotary()

}

func blink(pin rpio.Pin, times int) {

	for i := 0; i < times; i++ {
		pin.Low()
		time.Sleep(20 * time.Millisecond)
		pin.High()
		time.Sleep(20 * time.Millisecond)
		pin.Low()
	}
}

func getEncoderDiff() int {

//	oldA := rpio.High
//	oldB := rpio.High

//	result := 0

	dt := rot_dt.Read()
	clk := rot_dt.Read()

	// newA := rot_dt.Read()
	// newB := rot_clk.Read()

	if dt == rpio.High {

		if clk == rpio.High {
			return 1
		} else {
			return -1
		}
	}

	return 0
/*
	if (newA != oldA || newB != oldB) {

		fmt.Println("change")
		// something has changed
		if oldA == rpio.High && newA == rpio.Low {
			result = (int(oldB)*2 - 1)
		}
	}

	oldA = newA
	oldB = newB
	return result
*/

}

func rotary() {

	encoderVal := 0

	for {
		encoderVal = encoderVal + getEncoderDiff()

		if (rot_sw.Read() == rpio.Low) {
			// switch done been mashed down
			// fmt.Println("push that button down")
			encoderVal = 0
		}

		// blink(led8, 2)
		fmt.Printf("Encoder Value: %d\n", encoderVal)
		render(encoderVal)
		time.Sleep(10 * time.Millisecond)
	}
}

func render(current int) {

	if (current > 7) {
		current = 7
	}

	for i := 0; i < 8; i++ {
		
		thePin := rpio.Pin(i)
		thePin.Low()

		if (i < current) {
			thePin.High()
		}
	}

}

func setup() {

	led1 = rpio.Pin(PIN_LED_1)
	led2 = rpio.Pin(PIN_LED_2)
	led3 = rpio.Pin(PIN_LED_3)
	led4 = rpio.Pin(PIN_LED_4)
	led5 = rpio.Pin(PIN_LED_5)
	led6 = rpio.Pin(PIN_LED_6)
	led7 = rpio.Pin(PIN_LED_7)
	led8 = rpio.Pin(PIN_LED_8)

	rot_sw = rpio.Pin(PIN_ROT_SW)
	rot_dt = rpio.Pin(PIN_ROT_DT)
	rot_clk = rpio.Pin(PIN_ROT_CLK)

	led1.Output()
	led2.Output()
	led3.Output()
	led4.Output()
	led5.Output()
	led6.Output()
	led7.Output()
	led8.Output()

	rot_sw.Output()
	rot_sw.High()
	rot_sw.Input()

	rot_dt.Input()
	rot_clk.Input()

	delay := 50 * time.Millisecond

	led1.High()
	time.Sleep(delay)
	led2.High()
	time.Sleep(delay)
	led3.High()
	time.Sleep(delay)
	led4.High()
	time.Sleep(delay)
	led5.High()
	time.Sleep(delay)
	led6.High()
	time.Sleep(delay)
	led7.High()
	time.Sleep(delay)
	led8.High()
	time.Sleep(delay)

	led1.Low()
	time.Sleep(delay)
	led2.Low()
	time.Sleep(delay)
	led3.Low()
	time.Sleep(delay)
	led4.Low()
	time.Sleep(delay)
	led5.Low()
	time.Sleep(delay)
	led6.Low()
	time.Sleep(delay)
	led7.Low()
	time.Sleep(delay)
	led8.Low()
}
