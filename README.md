# Sprite Sheet Viewer

A lightweight tool for viewing and inspecting sprite sheets with configurable grid sizes and margins.  


## Features

- Load PNG and JPEG sprite sheets
- Adjust grid size and margin settings in real-time
- Auto-grid sprite detection
- Scroll through large sprite sheets
- Visual sprite preview with row / column indicators


## Example
<div align="center">
  <video src="https://github.com/user-attachments/assets/6ef2a6b3-1f8c-4f40-8d90-df67bae71471" width="800" alt="Sprite Sheet Viewer Screenshot">
  </video>  
</div>  


## Configuration

The viewer provides two main settings:
- **Margin**: Space between sprites  (Limit: 10px)
- **Grid Size**: Size of each sprite cell (Limit: 64px)

## Running the Viewer

Download the latest version from the [Github Releases](https://github.com/ztkent/spritesheet-viewer/releases).

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
