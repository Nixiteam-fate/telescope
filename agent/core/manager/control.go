package manager

import (
	"os"
)

func HandleOsSignal(osSignal chan os.Signal) {
	<-osSignal

	HandleSignalDirect()
}

func HandleSignalDirect(){
	os.Exit(0)
}
