# HTTP Helper

A simple HTTP server for development and testing, similar to `python -m http.server` with a few more features.

## Running from docker

```shell
$ docker build -t http-helper/hs .
$ docker run -p 8000:8000 -v ./testdata:/data http-helper/hs
```

## License

Unless otherwise noted, all source code in this repo falls under a [MIT License](LICENSE).

The [Gopher image](testdata/images/gopher.png) was created by Renee French and uses a Creative Commons Attribution 4.0 licensed License.

The [hello.go](testdata/hello.go) was taken from the [Go Playground](https://go.dev/play/), and probably is copyright 2009 The Go Authors.

The [Water Bottle](testdata/WaterBottle.glb) was taken from [glTF Sample Models](https://github.com/KhronosGroup/glTF-Sample-Models/tree/main/2.0/WaterBottle) and is CC0.