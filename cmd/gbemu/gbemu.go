package main

import (
	"flag"
	"fmt"
	"github.com/adnsio/gbemu/internal/renderer"
	"github.com/adnsio/gbemu/pkg/gameboy"
	"io/ioutil"
	"os"
)

func loadFileData(path string) []uint8 {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}

	return data
}

func main() {
	var bootromPath, cartridgePath string
	var debugWindows bool
	//var maxFramesPerSecond int

	flag.StringVar(&bootromPath, "bootrom", "assets/bios/dmg_boot.bin", "bootrom path")
	flag.StringVar(&cartridgePath, "cartridge", "assets/roms/tetris.gb", "cartridge path")
	flag.BoolVar(&debugWindows, "debug-windows", true, "enabled debug windows")
	//flag.IntVar(&maxFramesPerSecond, "max-fps", 60, "max frames per second")

	flag.Usage = func() {
		fmt.Fprint(flag.CommandLine.Output(), "Usage: gbemu [options]\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	gbCfg := gameboy.Config{}

	// testing
	//bootromPath = ""
	//cartridgePath = "assets/test_roms/cpu_instrs.gb"
	//cartridgePath = "assets/test_roms/cpu_instrs-individual/06-ld r,r.gb"

	if bootromPath != "" {
		gbCfg.Bootrom = loadFileData(bootromPath)
	}

	if cartridgePath != "" {
		gbCfg.Cartridge = loadFileData(cartridgePath)
	}

	gb := gameboy.NewGameBoy(gbCfg)

	rdr := renderer.NewRenderer(renderer.Config{
		GameBoy:      gb,
		DebugWindows: debugWindows,
	})

	rdr.Run()
}
