# gontainer
###  a lightweight containerization tool

**Requirements:**
* Linux
* A linux filesystem if you wish to containerize 
to a linux environment. You can get one from [Alpine Linux](https://alpinelinux.org/downloads/)


**Usage:**

How to open a ash window, which is like bash (Alpine linux fs)
```shell script
gontainer -chrt "/home/alpine-fs" run ash 
```

Example on how to run python (must be installed beforehand, see [Gockerfile](Gockerfile))
```shell script
gontainer -chrt "/home/alpine-fs" run python3 
```

## Quickstart

Requirements:
* Gontainer binary
* [`mnt-vfs.sh`](./mnt-vfs.sh) shell script
* Linux filesystem image (alpine available at [Alpine Linux](https://alpinelinux.org/downloads/))
* Root privilidges for all commands

1. Create a blank virtual file system and mount using the `mnt-vfs.sh` script.
    ```sh
    cd /dev/external-hard-drive
    sudo mnt-vfs.sh VFS 200
    ```
    This command will create a directory for the VFS at `/mnt/VFS` and create a 
    blank image for the filesystem of type `ext3` in the working directory called `VFS.ext3`.
    The size of the image will be 200MB, trying to store more than that on the filesystem will be impossible. Finally it mounts the newly created `VFS.ext3` filesystem image to `/mnt/VFS` so
    that the filesystem is ready to use. It also adds an mounting line to `/etc/fstab` file so that
    when mounting the filesystem correct parameters are used.

2. Copy linux filesystem image to newly created VFS. For this example the filesystem
is located at `$HOME/alpine-fs`.
    ```sh
    sudo cp -r $HOME/alpine-fs/* /mnt/VFS/
    ```
    Check if copy was succesfull:
    ```sh
    # ls /mnt/upsya-vfs
    alpine-fs  bin  dev  etc  home  lib  lost+found  media  mnt  opt  proc  root  run  sbin  srv  sys  tmp  usr  var
    ```
3. You may now use gontainer on the image. Open a shell with `ash`
    ```sh
    sudo gontainer -chrt /mnt/my-vfs run sh
    ```
    Look at the [`Gockerfile`](./Gockerfile) for typical first commands on getting a
    python installation up and running.

To **uninstall** virtual filesystem delete the line in your `/etc/fstab` folder corresponding
to the virtual filesystem (only one line with second value set to mount directory).
Finally delete the mount directory.