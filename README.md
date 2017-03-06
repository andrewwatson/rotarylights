# Rotary Lights

This program reads the output of a rotary encoder and powers a line of LEDs according
to the relative value.  Turn the knob clockwise and the number goes up!

## Schematics

Coming Soon!

## UI
I used TermUI [https://github.com/gizak/termui] to build a simple interface to show
the state of the pins.  When you're in a tight loop with only a 5ms delay you can't
just spew data out to the console.
