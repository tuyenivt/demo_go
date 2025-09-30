package main

import "concurrency/app"

func main() {
	app.SimpleGoroutine("World")

	app.SimpleWaitGroup()

	app.SimpleChannel()

	app.ReadWriteChannel("World")

	app.MultiChannel()

	app.RaceConditionMutex()

	app.RaceConditionAtomic()

	app.RaceConditionNewCond()

	app.SimpleContextCancel()

	app.SimpleContextTimeout()

	app.SimpleFanInFanOut()

	app.SimpleSyncPool()
}
