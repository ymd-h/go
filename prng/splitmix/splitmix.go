package splitmix

type (
	SplitMix32 struct {
		x uint32
	}

	SplitMix64 struct {
		x uint64
	}
)

func NewSplitMix32(seed uint32) *SplitMix32 {
	return &SplitMix32{x: seed}
}

func (s *SplitMix32) Next() uint32 {
	s.x += 0x9e3779b9
	z := s.x
	z = (z ^ (z >> 16)) * 0x85ebca6b
	z = (z ^ (z >> 13)) * 0xc2b2ae35
	return z ^ (z >> 16)
}

func NewSplitMix64(seed uint64) *SplitMix64 {
	return &SplitMix64{x: seed}
}

func (s *SplitMix64) Next() uint64 {
	s.x += 0x9e3779b97f4a7c15
	z := s.x
	z = (z ^ (z >> 30)) * 0xbf58476d1ce4e5b9
	z = (z ^ (z >> 27)) * 0x94d049bb133111eb
	return z ^ (z >> 31)
}
