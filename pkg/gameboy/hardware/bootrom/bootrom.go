package bootrom

const (
	Start = 0x0000
	End   = 0x00ff
	Size  = End - Start + 1
)

type Bootrom struct {
	Enabled bool
	Data    [Size]uint8
}

func NewBootrom() *Bootrom {
	return &Bootrom{}
}

func (b *Bootrom) Load(data []uint8) {
	for i := Start; i <= End; i++ {
		b.Data[i] = data[i]
	}
}

func (b *Bootrom) Read(addr uint16) uint8 {
	return b.Data[addr]
}
