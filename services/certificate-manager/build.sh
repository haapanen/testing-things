#!/bin/sh

go build -gcflags="all=-N -l" -o /certificate-manager
