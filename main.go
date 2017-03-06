package main

import (
	"fmt"
	"github.com/gizak/termui"
	"github.com/stianeikeland/go-rpio"
	"os"
	"os/signal"
	"time"
)

var (
	led1, led2, led3, led4, led5, led6, led7, led8 rpio.Pin
	rot_c, rot_a, rot_b                        rpio.Pin

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

	PIN_ROT_A  = 26
	PIN_ROT_B  = 5
	PIN_ROT_C  = 4
)

func main() {

	if err := rpio.Open(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer rpio.Close()
	killChannel := make(chan bool)

	sigChannel := make(chan os.Signal, 1)
	signal.Notify(sigChannel, os.Interrupt)

	/*
		err := termui.Init()
		if err != nil {
			panic(err)
		}

		defer termui.Close()


		termui.Handle("/sys/kbd/q", func(termui.Event) {
			termui.StopLoop()
			renderLED(0)
			killChannel <- true
		})
	*/

	setup()
	renderLED(0)

	// renderCUI(0, 0, 0, 0, "")

	messages := make(chan string)
	go rotary(killChannel, messages)

	var msg string
	for {
		select {
		case msg = <-messages:
			fmt.Sprintf("MESSAGE: %s\n", msg)
		case <- sigChannel:
			renderLED(0)
			os.Exit(0)
		default:
		}

	}

	//	termui.Loop()

}

func renderCUI(clk, dt, lastCLK rpio.State, value int, message string) {

	strs := []string{
		time.Now().Format(time.Stamp),
		fmt.Sprintf("DT PIN: %d", dt),
		fmt.Sprintf("CLK PIN: %d", clk),
		fmt.Sprintf("LAST CLK PIN: %d", lastCLK),
		fmt.Sprintf("Value: %d", value),
		fmt.Sprintf("Event: %s", message),
	}

	ls := termui.NewList()
	ls.Items = strs
	ls.ItemFgColor = termui.ColorYellow
	ls.BorderLabel = "Current Values"
	ls.Height = 20
	ls.Width = 150
	ls.Y = 2

	termui.Render(ls)

}

func rotary(killChannel chan bool, messages chan string) {

	encoderVal := 0

	var lastA rpio.State

	lastA = rot_a.Read()

	message := ""

	for {

		a := rot_a.Read()
		b := rot_b.Read()

		if (a != lastA) {

			if (b != a) {
				encoderVal += 1
				message = "increment"

			} else {
				encoderVal -= 1
				message = "decrement"
			}

			renderLED(encoderVal)
			//renderCUI(dt, clk, lastCLK, encoderVal, message)

		}

		if rot_c.Read() == rpio.Low {
			encoderVal = 0
			renderLED(encoderVal)
		}

		select {
		case <-killChannel:
			fmt.Println("shutdown signal")
			renderLED(0)
			return
		case messages <- message:
		default:
		}

		lastA = a

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

	rot_a = rpio.Pin(PIN_ROT_A)
	rot_b = rpio.Pin(PIN_ROT_B)
	rot_c = rpio.Pin(PIN_ROT_C)

	led1.Output()
	led2.Output()
	led3.Output()
	led4.Output()
	led5.Output()
	led6.Output()
	led7.Output()
	led8.Output()

	rot_a.Output()
	rot_a.High()
	rot_a.Input()

	rot_b.Input()
	rot_c.Input()

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
