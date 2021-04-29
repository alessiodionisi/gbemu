package gameboy

import (
	"image"

	"github.com/adnsio/gbemu/pkg/gameboy/cpu"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const ScreenWidth int = 160
const ScreenHeight int = 144

type Gameboy struct {
	log zerolog.Logger
	cpu *cpu.CPU

	IsRunning bool
}

func (g *Gameboy) Screen() *image.RGBA {
	rgba := image.NewRGBA(
		image.Rect(
			0, 0,
			ScreenWidth, ScreenHeight,
		),
	)

	return rgba
}

func (g *Gameboy) BackgroundMap() *image.RGBA {
	rgba := image.NewRGBA(
		image.Rect(
			0, 0,
			256, 256,
		),
	)

	return rgba
}

func (g *Gameboy) Resume() {

}

func (g *Gameboy) Pause() {

}

func (g *Gameboy) update() {
	for g.IsRunning {

	}

	g.log.Info().Msg("update exit")
}

func New() *Gameboy {
	return &Gameboy{
		log: log.With().Str("package", "gameboy").Logger(),
		cpu: cpu.New(),
	}
}
