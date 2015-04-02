# Voxel-Engine
As it stands this is a generic voxel graphics engine with a very small memory footprint ~150MB for roughly 400k cubes
The goal is to reduce this further
Engine created using the go-gl graphics package

Overview
The final intent is a 3D voxel based game

You will need a local MongoDB running, install on your local machine and then run mongod from a terminal
The Go code will setup the other necessary data

To run execute these commands 
"go build github.com/allanks/Voxel-Engine/src/main" 
"go run src/main/main.go"