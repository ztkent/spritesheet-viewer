# Sprite Sheet Viewer

A lightweight tool for viewing and inspecting sprite sheets with configurable grid sizes and margins.

## Features

- Load PNG and JPEG sprite sheets
- Adjust grid size and margin settings in real-time
- Auto-grid sprite detection
- Scroll through large sprite sheets
- Visual sprite preview with row / column indicators

## Configuration

The viewer provides two main settings:
- **Margin**: Space between sprites  (Limit: 10px)
- **Grid Size**: Size of each sprite cell (Limit: 64px)

## Running the Viewer

Download the latest version from the [Github Releases]().

### Building from Source

Requires Go and the following dependencies:
```bash
go get github.com/gen2brain/raylib-go/raylib
```

To build and run:
```bash
go build
./view_spritesheet
```

## System Requirements
- macOS
