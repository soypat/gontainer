# gontainer
###  a lightweight containerization tool

**Requirements:**
* Linux
* A linux filesystem if you wish to containerize 
to a linux environment. You can get one from [Alpine Linux](https://alpinelinux.org/downloads/)


**Usage:**

How to open a shell window (Alpine linux fs)
```shell script
gontainer --chrt "/home/alpine-fs" run sh 
```

How to open a ash window, which is like bash (Alpine linux fs)
```shell script
gontainer --chrt "/home/alpine-fs" run ash 
```

Example on how to run python (must be installed beforehand, see [Gockerfile](Gockerfile))
```shell script
gontainer --chrt "/home/alpine-fs" run python3 
```