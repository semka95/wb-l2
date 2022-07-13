package pattern

import (
	"fmt"
)

/*
	Реализовать паттерн «фасад».
Объяснить применимость паттерна, его плюсы и минусы,а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Facade_pattern
*/

type filamentType string

const (
	PETG filamentType = "PETG"
	PLA               = "PLA"
	ABS               = "ABS"
)

func newFilamentType(fType string) (filamentType, error) {
	switch fType {
	case "PETG":
		return PETG, nil
	case "PLA":
		return PLA, nil
	case "ABS":
		return ABS, nil
	}
	return "", fmt.Errorf("not allowed '%s' filament", fType)
}

type filament struct {
	fType filamentType
}

func newFilament(fType string) (*filament, error) {
	t, err := newFilamentType(fType)
	if err != nil {
		return nil, err
	}
	return &filament{
		fType: t,
	}, nil
}

func (f *filament) load() {
	fmt.Printf("loading %s filament\n", f.fType)
}

func (f *filament) unload() {
	fmt.Println("unloading filament")
}

type heatbed struct {
	temperature int
}

func newHeatbed(temperature int) *heatbed {
	return &heatbed{
		temperature: temperature,
	}
}

func (h *heatbed) setTemperature(temperature int) error {
	if temperature < 0 || temperature > 100 {
		return fmt.Errorf("%d temperature not allowed", temperature)
	}

	h.temperature = temperature
	return nil
}

func (h *heatbed) heatUp() {
	fmt.Printf("heating up to %d degrees\n", h.temperature)
}

type hotend struct {
	temperature int
}

func newHotend(temperature int) *hotend {
	return &hotend{
		temperature: temperature,
	}
}

func (h *hotend) setHotend(temperature int) error {
	if temperature < 0 || temperature > 300 {
		return fmt.Errorf("%d temperature not allowed", temperature)
	}

	h.temperature = temperature
	return nil
}

func (h *hotend) heatUp() {
	fmt.Printf("heating up to %d degrees\n", h.temperature)
}

type printerFacade struct {
	filament *filament
	heatbed  *heatbed
	hotend   *hotend
}

func newPrinterFacade(fType string, heatbedTemp int, hotendTemp int) (*printerFacade, error) {
	f, err := newFilament(fType)
	if err != nil {
		return nil, err
	}
	return &printerFacade{
		filament: f,
		heatbed:  newHeatbed(heatbedTemp),
		hotend:   newHotend(hotendTemp),
	}, nil
}

func (p *printerFacade) preparePrinting() {
	fmt.Println("unloaing old filament")
	p.filament.unload()
	p.filament.load()

	fmt.Println("heating bed")
	p.heatbed.heatUp()

	fmt.Println("heating hotend")
	p.hotend.heatUp()
}
