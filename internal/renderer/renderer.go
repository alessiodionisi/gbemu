package renderer

import (
	"fmt"
	"image"
	"runtime"
	"time"

	"github.com/adnsio/gbemu/pkg/gameboy"
	"github.com/adnsio/gbemu/pkg/gameboy/bits"
	"github.com/adnsio/gbemu/pkg/gameboy/hardware/display"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	//MainWindowWidth        = 160
	//MainWindowHeight       = 144
	MainWindowScale = 1
	//TilesWindowWidth       = 8 * 16
	//TilesWindowHeight      = 8 * 16
	//TilesWindowScale       = 1
	BackgroundWindowWidth  = 8 * 32
	BackgroundWindowHeight = 8 * 32
	BackgroundWindowScale  = 1
)

type Config struct {
	DebugWindows bool
	GameBoy      *gameboy.GameBoy
}

type Renderer struct {
	MainWindow            *sdl.Window
	MainRenderer          *sdl.Renderer
	MainTexture           *sdl.Texture
	BackgroundWindow      *sdl.Window
	BackgroundRenderer    *sdl.Renderer
	BackgroundTexture     *sdl.Texture
	BackgroundImage       *image.RGBA
	IsDebugWindowsEnabled bool
	GameBoy               *gameboy.GameBoy
}

func NewRenderer(cfg Config) *Renderer {
	rdr := &Renderer{
		IsDebugWindowsEnabled: cfg.DebugWindows,
		GameBoy:               cfg.GameBoy,
		BackgroundImage:       image.NewRGBA(image.Rect(0, 0, BackgroundWindowWidth, BackgroundWindowHeight)),
	}

	return rdr
}

func (rdr *Renderer) CreateBackgroundWindow() {
	var err error

	rdr.BackgroundWindow, err = sdl.CreateWindow("gbemu - background", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, int32(BackgroundWindowWidth*BackgroundWindowScale), int32(BackgroundWindowHeight*BackgroundWindowScale), sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}

	rdr.BackgroundRenderer, err = sdl.CreateRenderer(rdr.BackgroundWindow, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}

	rdr.BackgroundTexture, err = rdr.BackgroundRenderer.CreateTexture(uint32(sdl.PIXELFORMAT_RGBA32), sdl.TEXTUREACCESS_STREAMING, int32(BackgroundWindowWidth), int32(BackgroundWindowHeight))
	if err != nil {
		panic(err)
	}
}

func (rdr *Renderer) UpdateBackgroundWindow() {
	bgMapSelect := 0
	bgTileSelect := 1

	for y := 0; y < BackgroundWindowHeight; y++ {
		row := y / 8

		for x := 0; x < BackgroundWindowWidth; x++ {
			col := x / 8

			var rawBgTileNum int
			// Bit 3 - BG Tile Map Display Select     (0=9800-9BFF, 1=9C00-9FFF)
			bgMapAddr := row*32 + col
			if bgMapSelect == 1 {
				rawBgTileNum = int(rdr.GameBoy.Hardware.Display.WindowMap[bgMapAddr])
			} else {
				rawBgTileNum = int(rdr.GameBoy.Hardware.Display.BackgroundMap[bgMapAddr])
			}

			var bgTileNum int
			// Bit 4 - BG & Window Tile Data Select   (0=8800-97FF, 1=8000-8FFF)
			if bgTileSelect == 1 {
				bgTileNum = rawBgTileNum
			} else {
				bgTileNum = 128 + (rawBgTileNum + 128)
			}

			bgTileAddr := uint16(bgTileNum * 16)

			line := (y % 8) * 2
			data1Addr := bgTileAddr + uint16(line)
			data2Addr := bgTileAddr + uint16(line+1)
			data1 := rdr.GameBoy.Hardware.Display.TileDataBank0[data1Addr]
			data2 := rdr.GameBoy.Hardware.Display.TileDataBank0[data2Addr]

			bit := (uint8(x)%8 - 7) * 0xff
			colorVal := (bits.Get(data2, bit) << 1) | bits.Get(data1, bit)
			color := rdr.GameBoy.Hardware.Display.ReadBackgroundPalette(colorVal)
			rdr.BackgroundImage.SetRGBA(x, y, rdr.GameBoy.Hardware.Display.ShadesOfGray[color])
		}
	}

	err := rdr.BackgroundTexture.Update(nil, rdr.BackgroundImage.Pix, rdr.BackgroundImage.Stride)
	if err != nil {
		panic(err)
	}

	err = rdr.BackgroundRenderer.Clear()
	if err != nil {
		panic(err)
	}

	err = rdr.BackgroundRenderer.Copy(rdr.BackgroundTexture, nil, nil)
	if err != nil {
		panic(err)
	}

	rdr.BackgroundRenderer.Present()
}

func (rdr *Renderer) CreateMainWindow() {
	var err error

	rdr.MainWindow, err = sdl.CreateWindow("gbemu", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, int32(display.Width*MainWindowScale), int32(display.Height*MainWindowScale), sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}

	rdr.MainRenderer, err = sdl.CreateRenderer(rdr.MainWindow, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}

	rdr.MainTexture, err = rdr.MainRenderer.CreateTexture(uint32(sdl.PIXELFORMAT_RGBA32), sdl.TEXTUREACCESS_STREAMING, int32(display.Width), int32(display.Height))
	if err != nil {
		panic(err)
	}
}

func (rdr *Renderer) UpdateMainWindow() {
	err := rdr.MainTexture.Update(nil, rdr.GameBoy.Hardware.Display.Image.Pix, rdr.GameBoy.Hardware.Display.Image.Stride)
	if err != nil {
		panic(err)
	}

	err = rdr.MainRenderer.Clear()
	if err != nil {
		panic(err)
	}

	err = rdr.MainRenderer.Copy(rdr.MainTexture, nil, nil)
	if err != nil {
		panic(err)
	}

	rdr.MainRenderer.Present()
}

func (rdr *Renderer) Run() {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	err := sdl.Init(sdl.INIT_VIDEO | sdl.INIT_AUDIO)
	if err != nil {
		panic(err)
	}
	defer sdl.Quit()

	rdr.CreateMainWindow()
	defer rdr.MainWindow.Destroy()
	defer rdr.MainRenderer.Destroy()
	defer rdr.MainTexture.Destroy()

	rdr.CreateBackgroundWindow()
	defer rdr.BackgroundWindow.Destroy()
	defer rdr.BackgroundRenderer.Destroy()
	defer rdr.BackgroundTexture.Destroy()

	ticker := time.NewTicker(time.Second / 60)
	running := true

	for range ticker.C {
		if !running {
			break
		}

		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				running = false
				break
			}
		}

		rdr.GameBoy.RunFrame()

		rdr.UpdateMainWindow()
		rdr.UpdateBackgroundWindow()

		if rdr.GameBoy.Hardware.Cartrdige.Title != "" {
			rdr.MainWindow.SetTitle(fmt.Sprintf("gbemu - %s", rdr.GameBoy.Hardware.Cartrdige.Title))
		}
	}
}
