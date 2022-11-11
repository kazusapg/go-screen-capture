package main

import (
	"fmt"
	"image"
	"log"
	"strconv"

	"github.com/go-vgo/robotgo"
	hook "github.com/robotn/gohook"
)

func main() {
	for {
		area, err := getArea()
		if err != nil {
			log.Println(err)
			continue
		}
		captureScreen(area)
		fmt.Print("Done.\n\n")
	}
}

func getArea() (image.Rectangle, error) {
	fmt.Println("Get the upper left position of the area to be captured by left clicking with the mouse.")
	leftUpperPoint := getLeftClickPos()
	fmt.Println("Capture left upper position is ", leftUpperPoint)
	// Coordinate adjustment considering display magnification
	leftUpperPoint = image.Point{
		X: applyScreenRatio(leftUpperPoint.X),
		Y: applyScreenRatio(leftUpperPoint.Y),
	}
	fmt.Println("Get the lower right  position of the area to be captured by left clicking with the mouse.")
	lowerRightPoint := getLeftClickPos()
	// Coordinate adjustment considering display magnification
	lowerRightPoint = image.Point{
		X: applyScreenRatio(lowerRightPoint.X),
		Y: applyScreenRatio(lowerRightPoint.Y),
	}
	fmt.Println(lowerRightPoint)
	fmt.Println("Capture lower right position is ", lowerRightPoint)

	rect := image.Rectangle{
		Min: leftUpperPoint,
		Max: lowerRightPoint,
	}

	p := rect.Size()

	if p.X < 0 || p.Y < 0 {
		return image.Rectangle{}, fmt.Errorf("the coordinates of the lower right corner of the capture area must be lower and more right than left upper position")
	}

	return rect, nil
}

func getLeftClickPos() image.Point {
	x, y := 0, 0
	robotgo.EventHook(hook.MouseDown, []string{}, func(e hook.Event) {
		x, y = robotgo.GetMousePos()
		robotgo.EventEnd()
	})

	s := robotgo.EventStart()
	<-robotgo.EventProcess(s)

	return image.Point{X: x, Y: y}
}

func applyScreenRatio(n int) int {
	scaleX, _ := robotgo.GetScaleSize()
	screenX, _ := robotgo.GetScreenSize()
	ratio := float64(scaleX) / float64(screenX)
	return int(float64(n) * ratio)
}

func captureScreen(area image.Rectangle) {
	fmt.Println("Please press ctrl + q to quit capture.")
	hook.Register(hook.KeyDown, []string{"q", "ctrl"}, func(e hook.Event) {
		fmt.Println("Stop.")
		hook.End()
	})

	imgNum := 1
	fmt.Println("To capture screen, Please Enter.")
	hook.Register(hook.KeyDown, []string{"enter"}, func(e hook.Event) {
		img := robotgo.CaptureImg(area.Min.X, area.Min.Y, area.Size().X, area.Size().Y)
		savePath := strconv.Itoa(imgNum) + ".png"
		robotgo.Save(img, savePath)
		fmt.Println(savePath, " was captured.")
		imgNum++
	})

	s := hook.Start()
	<-hook.Process(s)
}
