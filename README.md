# embd-lcd
lcd driver for embd

Currently only supports 16x2 char LCDs, over the shared I2C interface in portwriter ( https://github.com/aleksclark/portwriter ). No fancyness, just write up to 16 chars to any given row:

```Go
package main
import "fmt"
import "time"
import "github.com/aleksclark/embd-lcd"

func main() {
	fmt.Println("Starting")
	lcd1 := lcd.NewSI2C(0,1,2,4,5,6,7,3, 0, 0x26)
	
	t := time.Now().Local()
	lcd1.PrintRow1("Time: " + t.Format("15:04:05"))
	lcd1.PrintRow2("Hello World!")

	for {
		t = time.Now().Local()
		lcd3.PrintRow1("Time:  " + t.Format("15:04:05"))
		delayMilli(1000)
	}

}

func delayMilli(n int) {
	holdTime := time.Duration(n)*time.Millisecond
	time.Sleep(holdTime)
}
```
