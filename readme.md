# LoPix: Generate pixel art from code!

To install:
```
go install github.com/ejuju/lopix@latest
```

Usage:
```
lopix {"png" or "gif"} {input file path} {output file path} {scaling factor}
```

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
