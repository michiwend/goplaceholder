goplaceholder
=============
a small golang lib to generate placeholder images.

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

## Webservice
You can try it as a web service at
[placeholder.michiwend.com](http://placeholder.michiwend.com/400x300.png?text=lorem%20ipsum!).

The following requests are allowed:
* /800x600.png
* /800x600.png?text=foo
* /500.png
* /500.png?text=foo
