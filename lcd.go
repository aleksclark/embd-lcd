// Package lcd provides a high-level interface for standard hd44780 LCDs
// LCDs may be connected via I2C, eventually GPIO will be supported
package lcd

import "fmt"
import "time"
import "github.com/aleksclark/portwriter"

type LCDCommand struct {
	data int
	hold int
}

type LCD struct {
	rsPin int
	rwpin int
	enPin int
	d4Pin int
	d5Pin int
	d6Pin int
	d7Pin int
	backlightPin int
	writer portwriter.PortWriter
	commChan chan LCDCommand
}


// NewSI2C returns a new LCD instance, using the SharedI2CWriter provided in portwriter.
// NewSI2C initializes the LCD before returning it.
func NewSI2C( rsPin, rwpin, enPin, d4Pin, d5Pin, d6Pin, d7Pin, backlightPin int, i2cPort, i2cAddr byte) *LCD {
	w := portwriter.NewSharedI2CWriter(i2cPort, i2cAddr)
	c := make(chan LCDCommand)
	l := LCD{rsPin, rwpin, enPin, d4Pin, d5Pin, d6Pin, d7Pin, backlightPin, w, c}
	l.Initialize()
	return &l
}

func (l *LCD) Initialize() {
	go lcdWriter(l);
	funcSet := 0
	funcSet = setBit(funcSet, l.d5Pin)
	funcSet = setBit(funcSet, l.backlightPin)
	// 4-bit mode must be clocked 3 times for it to take
	l.clockByte(funcSet, 5)
	l.clockByte(funcSet, 5)
	
	// l.clockByte(funcSet, 5)
	// l.clockByte(funcSet, 5)
	// function set 2 lines, 16 chars, 5x8
	l.SendCommand(40)
	// display on, cursor off
	l.SendCommand(12)
	// clear entire display, set to pos 0
	l.SendCommand(1)
	l.SendCommand(2)
	l.PrintRow1("")
	l.PrintRow2("")
}

func (l *LCD) PrintRow1(str string) {
	l.SendCommand(128)
	l.PrintText(str)
}

func (l *LCD) PrintRow2(str string) {
	l.SendCommand(168)
	l.PrintText(str)
}

func (l *LCD) PrintText(str string) {
	str = fmt.Sprintf("%-17s", str)
	for _, value := range str {
		l.SendData(int(value), false)
	}
}

func (l *LCD) SendCommand(data int) {
	l.SendData(data, true)
}

func (l *LCD) SendData(data int, command bool) {
	nib := 0
	if (!command) {
		nib = setBit(nib, l.rsPin)
	}

	nib = setBit(nib, l.backlightPin)
	// send bits 7-4
	nib = nib ^ ((data >> 4) << 4)
	l.clockByte(nib, 1)
	// clear bits 7-4 but leave remaining bits in place
	nib = nib &^ (255 << 4)
	//set bits 3-0 of data to bits 7-4 of nib
	nib = nib ^ (data << 4)
	l.clockByte(nib, 1)

}

// clockByte clocks in a single byte, holding the enPin high for the specified hold in milliseconds.
func (l *LCD) clockByte(data int, hold int) {
	msg := LCDCommand{data, hold}
	l.commChan <- msg
}

func lcdWriter(l *LCD) {
	for {
		msg, more := <-l.commChan
		if more {
			timedByteWrite(msg, l)
		} else {
			fmt.Println("all done!")
			return
		}
	}
}

func timedByteWrite(c LCDCommand, l *LCD) {

	enHigh := setBit(c.data, l.enPin)
	l.writer.WriteByte(byte(enHigh))

	enHigh = clearBit(c.data, l.enPin)
	l.writer.WriteByte(byte(enHigh))
	delayMilli(c.hold)

}

func printInt(data int) {
    int_string := ""
    for i := 7; i > -1; i-- {

        if hasBit(data, i) {
            int_string += "1"
        } else {
            int_string += "-"
        }
        
    }
    fmt.Println(int_string)
}

func hasBit(n, pos int) bool {
	val := n & (1 << uint(pos))
	return (val > 0)
}

func setBit(n, pos int) int {
    n |= 1<<uint(pos);
    return n
}

func clearBit(n, pos int) int {
    mask := ^(1 << uint(pos))
    n &= mask
    return n
}

func delayMilli(n int) {
	holdTime := time.Duration(n)*time.Millisecond
	time.Sleep(holdTime)
}