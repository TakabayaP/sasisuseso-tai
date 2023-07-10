# sasisuseso-tai
Interactive Sensing 2023

Made by Katsumi Kobayashi, Menma.

# How to use
## 1. Install
First, install gocv following the instruction on https://github.com/hybridgroup/gocv.

```bash
$ git clone https://github.com/hybridgorup/gocv.git
$ cd gocv
$ make install
```
If go run ./cmd/version/main.go does not work, try installing opencv via yay.

```bash
$ yay -S opencv
```

 You may also need to export PKG_CONFIG_PATH.

```bash
$ export PKG_CONFIG_PATH="/usr/lib/pkgconfig:$PKG_CONFIG_PATH"
```

If it still doesn't work,try settind ldconfig.
Add ```/usr/local/lib/opencv/``` to ```/etc/ld.so.conf.d/opencv.conf```, and run ```sudo ldconfig```.

If everything goes well, run ```go run .``` in the ```sasisuseso-tai/server``` directory.
This may take a while since it will also build the dependencies (I guess). Don't forget to plug in a camera.