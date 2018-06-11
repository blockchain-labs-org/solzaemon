package token

type Pos int

type File struct {
	base int
	size int

	lines [][]int
}

func NewFile() *File {
	return &File{lines: [][]int{[]int{}}}
}

func (f *File) AddLine(offset int) {
	f.lines = append(f.lines, []int{})
}

func (f *File) AddCharacter(offset int) {
	lx := len(f.lines) - 1
	f.lines[lx] = append(f.lines[lx], offset)
}

func (f *File) Line(offset int) int {
	for l, ld := range f.lines {
		for _, cd := range ld {
			if offset == cd {
				return l + 1
			}
		}
	}
	return 0
}

func (f *File) Character(offset int) int {
	for _, ld := range f.lines {
		for c, cd := range ld {
			if offset == cd {
				return c + 1
			}
		}
	}
	return 0
}

func (f *File) Offset(line, character int) int {
	for l, ld := range f.lines {
		if l+1 != line {
			continue
		}
		for c, cd := range ld {
			if c+1 == character {
				return cd
			}
		}
	}
	return 0
}
