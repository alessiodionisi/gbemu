package main

import (
	"os"

	"github.com/adnsio/gbemu/pkg/gameboy"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/veandco/go-sdl2/sdl"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out: os.Stderr,
	})

	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		log.Fatal().Err(err).Send()
	}
	defer sdl.Quit()

	screenScale := int32(3)

	window, err := sdl.CreateWindow(
		"gbemu",
		sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED,
		int32(gameboy.ScreenWidth)*screenScale, int32(gameboy.ScreenHeight)*screenScale,
		sdl.WINDOW_SHOWN|sdl.WINDOW_ALLOW_HIGHDPI,
	)
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	defer renderer.Destroy()

	gameboyTexture, err := renderer.CreateTexture(
		uint32(sdl.PIXELFORMAT_RGBA32),
		sdl.TEXTUREACCESS_STREAMING,
		int32(gameboy.ScreenWidth), int32(gameboy.ScreenHeight),
	)
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	defer gameboyTexture.Destroy()

	// bootromPath := "assets/roms/dmg-boot.bin"

	gameboy := gameboy.New()

	// run gameboy cycles
	// go gameboy.Update()

	// loop
	running := true
	for running {
		// pool sdl events
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch typedEvent := event.(type) {
			case *sdl.QuitEvent:
				log.Info().Msg("quit event")
				running = false
			case *sdl.KeyboardEvent:
				keyCode := typedEvent.Keysym.Sym

				if typedEvent.State == sdl.RELEASED && keyCode == sdl.K_ESCAPE {
					log.Info().Msg("escape key released")
				}
			}
		}

		// clear
		if err := renderer.Clear(); err != nil {
			log.Fatal().Err(err).Send()
		}

		// TODO: draw the system UI

		// get gameboy screen
		screen := gameboy.Screen()

		// draw gameboy texture
		if err := gameboyTexture.Update(nil, screen.Pix, screen.Stride); err != nil {
			log.Fatal().Err(err).Send()
		}

		// render gameboy texture
		// TODO: if gameboy is paused draw this texture on *sdl.Rect
		if err := renderer.Copy(gameboyTexture, nil, nil); err != nil {
			log.Fatal().Err(err).Send()
		}

		// present the render
		renderer.Present()

		// some delay
		sdl.Delay(1)
	}

	// bye bye sdl
	sdl.Quit()
}
