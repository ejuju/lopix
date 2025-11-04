# TODO

- Implement Parser to be able to factor/unify parsing frame and animation
- Factor "ReadGrid" routine from Frame.ReadFrom and Animation.ReadFrom
- Rename Frame.ReadFrom to func Frame.ParseFrom(r io.Reader) (err error)
- Web editor
- Add unit tests

Ideas:
- Support rendering frames and animations to terminal (2ch per cell for square)
