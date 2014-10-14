goplaceholder
=============
a small golang lib to generate placeholder images

## Usage
get it
```
$ go get github.com/michiwend/goplaceholder
```

simple example
```Go
placeholder, err := goplaceholder.Placeholder(
    "Lorem ipsum!",
    "/usr/share/fonts/TTF/DejaVuSans-Bold.ttf",
    color.RGBA{150, 150, 150, 255},
    color.RGBA{204, 204, 204, 255},
    400, 200)
```

results in

![example placeholder](example/lorem.png)
