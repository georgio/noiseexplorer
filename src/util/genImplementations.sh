#!/bin/bash
echo "[NoiseExplorer] Generating implementations..."
cd ../
for pattern in ../patterns/*.noise; do
	node noiseExplorer \
		--generate=go --pattern=$pattern
	node noiseExplorer \
		--generate=rs --pattern=$pattern
	node noiseExplorer \
		--generate=wasm --pattern=$pattern
done

cd ../implementations/go
go get -d ./...
