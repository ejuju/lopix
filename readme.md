# Lopix: Pixel Art Toolkit

Render pixel art (PNG/GIF) from code!

## Quick start

To install:
```
go install github.com/ejuju/lopix@latest
```

Usage:
```
lopix {"png" or "gif"} {input_lopix_file} {output_file_path} {scaling_factor}
```

## Lopix text format

Lopix relies on a simple code format for declaring grids of colored pixels.
In the end, all we need to know to render an image is the dimensions of the pixel grid (width and height),
and which colors to assign to each pixel.

In order to do so, you start by declaring the width and height at the beggining of the Lopix file, for example:
```
5x6
```
This declares that the width of the grid is 5 pixels and its height is 6 pixels.

Then on the next lines we declare a color palette to use
(a maximum of 16 colors can be specified: one per line, in hexadecimal RGBA format):
```
#ff0000
#00ff00
#0000ff
```

Then you can assign colors to each pixel, using index of the color to use in the palette defined above
(ranging from `0` for the first one, `f` for the sixteenth color):
```
00000
01220
01220
01220
01110
00000
```

NB: You must leave a blank line after the color palette and terminate each grid row with a trailing line feed.

Which gives us the following file all together:
```
5x6
#ff0000
#00ff00
#0000ff

00000
01220
01220
01220
01110
00000
```

You can now render this file as PNG using our CLI or our Go library.

## Examples

### Generate a PNG (using CLI)

To generate a PNG: first, define a frame (in a `.lopix` file), for example:
```
16x16
#d6d6d6
#ff4000
#242424

0000000000000000
0000000220000000
0000000200000000
0011111111111100
0011111111111100
0011111111111100
0011121111211100
0011121111211100
0011122112211100
0011111111111100
0011212121121100
0011222222221100
0011112121211100
0011111111111100
0000000000000000
0000000000000000
```

And then render it:
```
lopix png src.lopix out.png 20
```

Which results in:  
![PNG of a pumpkin](/examples/demo-0/demo.png)

### Generate a GIF (using CLI)

To generate a GIF: first, define an animation (in a `.lopix` file), for example:
```
2*30
16x16
#d6d6d6
#ff4000
#242424

0000000000000000
0000000220000000
0000000200000000
0011111111111100
0011111111111100
0011111111111100
0011121111211100
0011121111211100
0011122112211100
0011111111111100
0011212121121100
0011222222221100
0011112121211100
0011111111111100
0000000000000000
0000000000000000

0000000000000000
0000000000000000
0000000220000000
0000000200000000
0011111111111100
0011111111111100
0011111111111100
0011121111211100
0011121111211100
0011122112211100
0011111111111100
0011212121121100
0011222222221100
0011112121211100
0011111111111100
0000000000000000
```

And then render it:
```
lopix gif src.lopix out.gif 20
```

Which results in:  
![GIF of a pumpkin](/examples/demo-1/demo.gif)


### Use as Go library

Generate a PNG:
```go
v := &lopix.Frame{}
err := v.ParseFrom(f)
if err != nil {
    panic(err)
}
err = v.EncodePNG(b, 1080/v.W()) // Scale to from 16x16 to 1080x1080p.
if err != nil {
    panic(err)
}
```

Generate a GIF:
```go
v := &lopix.Animation{}
err := v.ParseFrom(f)
if err != nil {
    panic(err)
}
err = v.EncodeGIF(b, 1080/v.W()) // Scale to from 16x16 to 1080x1080p.
if err != nil {
    panic(err)
}
```
