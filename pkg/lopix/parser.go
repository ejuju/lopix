package lopix

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type Parser struct {
	bufr *bufio.Reader
	line int // Line number (index+1).
}

func NewParser(r io.Reader) *Parser {
	return &Parser{bufio.NewReader(r), 1}
}

type ParserError struct {
	line int
	err  error
}

func (err *ParserError) Line() int     { return err.line }
func (err *ParserError) Error() string { return fmt.Sprintf("at line %d: %s", err.line, err.err) }
func (err *ParserError) Unwrap() error { return err.err }

func (p *Parser) errf(f string, args ...any) *ParserError {
	return &ParserError{p.line, fmt.Errorf(f, args...)}
}

// Reads a line (removing any leading and trailing spaces).
func (p *Parser) ReadLine() (v string, err error) {
	line, err := p.bufr.ReadString('\n')
	if err != nil {
		return "", p.errf("read line: %w", err)
	}
	p.line++
	line = strings.TrimSpace(line)
	return line, nil
}

func (p *Parser) ParseAnimationInfo() (numFrames, delay int, err error) {
	line, err := p.ReadLine()
	if err != nil {
		return 0, 0, err
	}

	parts := strings.Split(line, "*")
	if len(parts) != 2 {
		return 0, 0, p.errf("invalid animation info line: %q", line)
	}

	numFrames, err = strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, p.errf("invalid number of frames: %w", err)
	} else if numFrames <= 0 || numFrames > MaxAnimationFrames {
		return 0, 0, p.errf("invalid number of frames: %d", numFrames)
	}

	frameDelay, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, p.errf("invalid frame delay: %w", err)
	} else if frameDelay <= 0 || frameDelay > MaxFrameDelay {
		return 0, 0, p.errf("invalid frame delay: %d", frameDelay)
	}

	return numFrames, frameDelay, nil
}

func (p *Parser) ParseDimensions() (w, h int, err error) {
	line, err := p.ReadLine()
	if err != nil {
		return 0, 0, err
	}

	parts := strings.Split(line, "x")
	if len(parts) != 2 {
		return 0, 0, p.errf("invalid dimension info line: %q", line)
	}

	w, err = strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, p.errf("invalid width: %w", err)
	} else if w <= 0 || w > MaxWidth {
		return 0, 0, p.errf("invalid width: %d", w)
	}

	h, err = strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, p.errf("invalid height: %w", err)
	} else if h <= 0 || h > MaxHeight {
		return 0, 0, p.errf("invalid height: %d", h)
	}

	return w, h, nil
}

func (p *Parser) ParsePalette() (v Palette, err error) {
	palette := Palette{}
	for i := 0; true; i++ {
		line, err := p.ReadLine()
		if err != nil {
			return v, err
		} else if line == "" {
			break // Reached end of palette (= blank line).
		} else if i > 15 {
			return v, p.errf("too many lines in palette")
		}
		palette[i], err = HexColor(line)
		if err != nil {
			return v, p.errf("invalid color: %w", err)
		}
	}
	return palette, nil
}

func (p *Parser) ParseGrid(w, h int) (grid []byte, err error) {
	grid = make([]byte, w*h)

	for y := range h {
		line, err := p.ReadLine()
		if err != nil {
			return nil, err
		}

		if line == "" {
			return nil, p.errf("blank line inside grid")
		} else if len(line) != w {
			return nil, p.errf("line size and width mismatch (%d vs %d)", len(line), w)
		}

		for x := range w {
			cell := hexToU4(line[x])
			grid[y*w+(x%w)] = cell
		}
	}

	return grid, nil
}

func (p *Parser) ReadBlankLine() (err error) {
	line, err := p.ReadLine()
	if err != nil {
		return err
	} else if line != "" {
		return p.errf("unexpected non-blank line: %q", line)
	}
	return nil
}
