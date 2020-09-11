# /etc/resolv.conf for internet access
nameserver 8.8.8.8

# Install python/pip
ENV PYTHONUNBUFFERED=1
RUN apk add --update --no-cache python3 && ln -sf python3 /usr/bin/python
RUN python3 -m ensurepip
RUN pip3 install --no-cache --upgrade pip setuptools

# Install build tools
apk add --update alpine-sdk

# Install python compiler thingy
apk add python3-dev

# Pip install package
ARCHFLAGS=-Wno-error=unused-command-line-argument-hard-error-in-future pip install --upgrade numpy
