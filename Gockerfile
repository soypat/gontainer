# This is NOT a Dockerfile. This is LIKE a Dockerfile.
# Run these commands in the container through bash/shell
# i.e. `gontainer --chrt /fs run ash` and you'll have a working
# Python installation

# /etc/resolv.conf for internet access
echo nameserver 8.8.8.8 > /etc/resolv.conf
env PYTHONUNBUFFERED=1
# Install barebones Python
apk add --update --no-cache python3 && ln -sf python3 /usr/bin/python
# Install pip
python3 -m ensurepip
pip3 install --no-cache --upgrade pip setuptools
# Install GNU compiler collection so python can compile+build numpy and other packages
apk add --update alpine-sdk
apk add python3-dev
# Install packages
ARCHFLAGS=-Wno-error=unused-command-line-argument-hard-error-in-future pip install --upgrade numpy
ARCHFLAGS=-Wno-error=unused-command-line-argument-hard-error-in-future pip install --upgrade pandas
