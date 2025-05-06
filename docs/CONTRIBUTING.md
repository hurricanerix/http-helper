# Contributing

Pull requests are welcome, however as this is a side project, I have limited time to vet/test new
features or review code.  I will do my best to be timely, but I can make no guarantees.

## Running from docker

```shell
$ docker build -t http-helper/hs .
$ docker run -p 8000:8000 -v ./testdata:/data http-helper/hs
```
