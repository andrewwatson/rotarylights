package main

import (
	"fmt"
	"github.com/stianeikeland/go-rpio"
	"os"
	"os/signal"
	"time"
	 "github.com/gizak/termui" 
)

var (
	led1, led2, led3, led4, led5, led6, led7, led8 rpio.Pin
	rot_sw, rot_dt, rot_clk                        rpio.Pin

	pinSlots []rpio.Pin
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

func main() {

	if err := rpio.Open(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer rpio.Close()

   err := termui.Init()
   if err != nil {
      panic(err)
   }

   defer termui.Close()

	killChannel := make(chan bool)

	termui.Handle("/sys/kbd/q", func(termui.Event) {
		termui.StopLoop()
		renderLED(0)
		killChannel <- true
	})

	setup()
	renderLED(0)

	renderCUI(0, 0, 0, 0)

   go rotary(killChannel)

	termui.Loop()

}

func renderCUI(clk, dt, lastCLK rpio.State, value int) {

	strs := []string{
		time.Now().Format(time.RFC822),
		fmt.Sprintf("DT PIN: %d", dt),
		fmt.Sprintf("CLK PIN: %d", clk),
		fmt.Sprintf("LAST CLK PIN: %d", lastCLK),
		fmt.Sprintf("Value: %d", value),
	}	

	ls := termui.NewList()
	ls.Items = strs
	ls.ItemFgColor = termui.ColorYellow
	ls.BorderLabel = "Current Values"
	ls.Height = 10 
	ls.Width = 50
	ls.Y = 0

	termui.Render(ls)

}

func rotary(killChannel chan bool) {

	encoderVal := 0

	var lastCLK rpio.State

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)


	lastCLK = rpio.High

	for {

		select {
		case <- killChannel:
			fmt.Println("shutdown signal")
			renderLED(0)
			return
		default:
		}

		dt := rot_dt.Read()
		clk := rot_dt.Read()

		if dt == rpio.High {

			if lastCLK == rpio.Low && clk == rpio.High {
				encoderVal += 1
			}

			if lastCLK == rpio.High && clk == rpio.Low {
				encoderVal -= 1
			}

		}

		if lastCLK != clk {
			lastCLK = clk
		}

		if rot_sw.Read() == rpio.Low {
			// switch done been mashed down
			// fmt.Println("push that button down")
			encoderVal = 0
		}

		// blink(led8, 2)
		// fmt.Printf("Encoder Value: %d\n", encoderVal)
		renderLED(encoderVal)
		renderCUI(dt, clk, lastCLK, encoderVal)
		time.Sleep(5 * time.Millisecond)
	}
}

func renderLED(current int) {

	if current > 8 {
		current = 8
	}

	// fmt.Printf("Current Value: %d", current)

	for i := 0; i < 8; i++ {

		thePin := pinSlots[i]
		thePin.Low()

		if i < current {
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

	pinSlots = []rpio.Pin{led1, led2, led3, led4, led5, led6, led7, led8}

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

	renderLED(0)
	time.Sleep(delay)
	renderLED(1)
	time.Sleep(delay)
	renderLED(2)
	time.Sleep(delay)
	renderLED(3)
	time.Sleep(delay)
	renderLED(4)
	time.Sleep(delay)
	renderLED(5)
	time.Sleep(delay)
	renderLED(6)
	time.Sleep(delay)
	renderLED(7)
	time.Sleep(delay)
	renderLED(8)
	time.Sleep(250 * time.Millisecond)
	renderLED(0)
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
