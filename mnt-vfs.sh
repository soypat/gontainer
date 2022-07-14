#!/bin/bash
# Patricio Whittingslow. Courtesy of http://souptonuts.sourceforge.net/quota_tutorial.html
TYPE="ext3" #Type can be ext2, ntfs, ext4 you name it
MNTDIR="/mnt/$1"
NAME="$1.${TYPE}"

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
mkdir -p $MNTDIR
BLOCKCOUNT=`expr $2 \* 1000 \* 1000 / 512`
echo "vfs of size ${2}MB"
# We first make the image of the quota size
dd if=/dev/zero of=./$NAME count=${BLOCKCOUNT}
# NExt we make the filesystem on the image
/sbin/mkfs -t $TYPE -q ./$NAME -F

# append a line to fstab
#This will make the filesystem always available on reboot, plus it's easier to mount and unmout when testing. 
echo "$DIR/$NAME  $MNTDIR   $TYPE rw,loop,usrquota,grpquota  0 0" &>> /etc/fstab

mount $MNTDIR
