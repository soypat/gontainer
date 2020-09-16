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

**VFS mount usage:**

There is a file included with this repo called `mnt-vfs.sh` which
is a script that creates a size-limited filesystem and mounts it.

Usage:
```shell script
chmod +x ./mnt-vfs.sh
./mnt-vfs.sh my-fs-name 100 # Will create a 100MB vfs with the name `my-fs-name`
```

### Creating a quota limited virtual filesystem
For this example we will create a 2GB image mounted at `/mnt/my-vfs`. I assume you have a alpine
linux filesystem at `/alpine-fs` and are running everything as root.

```shell script
wget http://urlTo.thefile/mnt-vfs.sh
chmod +x ./mnt-vfs.sh
./mnt-vfs.sh my-vfs 2000
cp -r /alpine-fs /mnt/my-vfs
```

Now you have a virtual filesystem ready with a 2GB limit! You can
run `gontainer` on it!

```shell script
gontainer --chrt /mnt/my-vfs run sh
```





