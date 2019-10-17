#!/bin/sh

# vi:nu:et:sts=4 ts=4 sw=4

echo "Testing the main package..."
if go test -v ./cmd/*.go ; then
	:
else
	echo "ERROR - Main package testing failed!"
	exit 4
fi

echo "Testing the support packages..."
if go test -v ./pkg/... ; then
	:
else
	echo "ERROR - Support packages testing failed!"
	exit 4
fi


