// A sprite sheet viewer application that allows users to load and
// view sprite sheets with configurable grid size and margin settings.
// It supports dynamic reloading and provides a graphical interface
// for viewing individual sprites within the sheet.

package main

import (
	"fmt"
	"os/exec"
	"runtime"
	"sort"
	"strconv"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/ztkent/beam/resources"
)

// UIState holds the application state and configuration.
type UIState struct {
	showFileDialog bool
	margin         int32
	gridSize       int32
	currentFile    string
	rm             *resources.ResourceManager
	sheet          *resources.SpriteSheet
	spriteNames    []string
	scrollOffset   float32
	loadError      string
	debugInfo      string
}

type Config struct {
	displaySize    int32
	padding        int32
	startX         int32
	startY         int32
	viewportHeight int32
	headerHeight   int32
}

// reload attempts to load or reload the current sprite sheet with the specified
// margin and grid size settings. It updates the internal state with any errors
// or debug information.
func (s *UIState) reload() {
	if s.currentFile == "" {
		return
	}

	newSprites := []resources.Resource{
		{
			Name:        "spritesheet",
			Path:        s.currentFile,
			IsSheet:     true,
			SheetMargin: int32(s.margin),
			GridSize:    int32(s.gridSize),
		},
	}

	if s.rm != nil {
		s.rm.Close()
	}

	s.rm = resources.NewResourceManagerWithGlobal(newSprites, nil)
	if s.rm == nil {
		s.loadError = "Failed to create resource manager"
		return
	}

	if len(s.rm.Scenes) == 0 || len(s.rm.Scenes[0].SpriteSheets) == 0 {
		s.loadError = "No sprites found in sheet"
		return
	}

	s.sheet = s.rm.Scenes[0].SpriteSheets[0]
	if s.sheet.Texture.ID == 0 {
		s.loadError = "Invalid texture"
		return
	}

	s.updateSpriteNames()
	s.debugInfo = fmt.Sprintf("Loaded %d sprites", len(s.spriteNames))
	s.loadError = ""
}

// updateSpriteNames refreshes the sorted list of sprite names from the current sheet.
func (s *UIState) updateSpriteNames() {
	s.spriteNames = nil
	for name := range s.sheet.Sprites {
		s.spriteNames = append(s.spriteNames, name)
	}
	sort.Slice(s.spriteNames, func(i, j int) bool {
		return naturalSort(s.spriteNames[i], s.spriteNames[j])
	})
}

func initConfig() Config {
	return Config{
		displaySize:    32,
		padding:        10,
		startX:         50,
		startY:         80,
		viewportHeight: 500,
		headerHeight:   40,
	}
}

func initUI() *UIState {
	rl.InitWindow(800, 600, "Sprite Sheet Viewer")
	rl.SetTargetFPS(60)

	return &UIState{
		margin:   1,
		gridSize: 16,
	}
}

// handleInput processes keyboard and mouse input events.
func (s *UIState) handleInput(showSettings *bool) {
	if rl.IsKeyPressed(rl.KeyO) && (rl.IsKeyDown(rl.KeyLeftControl) || rl.IsKeyDown(rl.KeyRightControl)) {
		if file := openFileDialog(); file != "" {
			s.currentFile = file
			s.reload()
		}
	}
	if rl.IsKeyPressed(rl.KeyEscape) && *showSettings {
		*showSettings = false
	}
	s.scrollOffset -= rl.GetMouseWheelMove() * 30
}

// handleScrolling manages scroll state based on content height and viewport
func (s *UIState) handleScrolling(contentHeight float32, viewportHeight int32) {
	maxScroll := float32(0)
	if contentHeight > float32(viewportHeight) {
		maxScroll = contentHeight - float32(viewportHeight)
	}
	if s.scrollOffset < 0 {
		s.scrollOffset = 0
	}
	if s.scrollOffset > maxScroll {
		s.scrollOffset = maxScroll
	}
}

// renderSprites draws all visible sprites from the sprite sheet.
func (s *UIState) renderSprites(cfg Config) {
	if s.sheet == nil || s.sheet.Texture.ID == 0 {
		if s.loadError == "" {
			rl.DrawText("No spritesheet loaded. Press 'Open File' to select one.", 50, cfg.startY, 20, rl.Gray)
		}
		return
	}

	spritesPerRow := (800 - cfg.startX*2) / (cfg.displaySize + cfg.padding)
	totalRows := len(s.spriteNames) / int(spritesPerRow)
	if len(s.spriteNames)%int(spritesPerRow) != 0 {
		totalRows++
	}

	contentHeight := float32(cfg.startY) + float32(totalRows*(int(cfg.displaySize)+int(cfg.padding)+20))
	s.handleScrolling(contentHeight, cfg.viewportHeight)

	x, y := cfg.startX, cfg.startY
	for _, name := range s.spriteNames {
		yPos := float32(y) - s.scrollOffset

		if yPos+float32(cfg.displaySize) < 0 || yPos > float32(600) {
			x += cfg.displaySize + cfg.padding
			if x > 700 {
				x = cfg.startX
				y += cfg.displaySize + cfg.padding + 20
			}
			continue
		}

		rect := s.sheet.Sprites[name]

		source := rl.Rectangle{
			X:      float32(rect.X),
			Y:      float32(rect.Y),
			Width:  float32(rect.Width),
			Height: float32(rect.Height),
		}

		dest := rl.Rectangle{
			X:      float32(x),
			Y:      yPos,
			Width:  float32(cfg.displaySize),
			Height: float32(cfg.displaySize),
		}
		rl.DrawTexturePro(s.sheet.Texture, source, dest, rl.Vector2{}, 0, rl.White)

		rl.DrawRectangleLinesEx(dest, 1, rl.Gray)
		rl.DrawText(name, int32(x), int32(int32(yPos)+cfg.displaySize+2), 10, rl.DarkGray)

		x += cfg.displaySize + cfg.padding
		if x > 700 {
			x = cfg.startX
			y += cfg.displaySize + cfg.padding + 20
		}
	}

	if contentHeight > float32(cfg.viewportHeight) {
		if s.scrollOffset > 0 {
			rl.DrawTriangle(
				rl.Vector2{X: 780, Y: 50},
				rl.Vector2{X: 790, Y: 60},
				rl.Vector2{X: 770, Y: 60},
				rl.Gray)
		}
		if s.scrollOffset < contentHeight-float32(cfg.viewportHeight) {
			rl.DrawTriangle(
				rl.Vector2{X: 780, Y: float32(cfg.viewportHeight + cfg.startY - 10)},
				rl.Vector2{X: 770, Y: float32(cfg.viewportHeight + cfg.startY - 20)},
				rl.Vector2{X: 790, Y: float32(cfg.viewportHeight + cfg.startY - 20)},
				rl.Gray)
		}
	}
}

// renderUI draws the application interface including header, buttons, and settings panel.
func (s *UIState) renderUI(cfg Config, showSettings *bool) {
	rl.DrawRectangle(0, 0, 800, cfg.headerHeight, rl.RayWhite)
	rl.DrawLine(0, cfg.headerHeight, 800, cfg.headerHeight, rl.LightGray)
	rl.DrawText("Sprite Sheet Viewer", 10, 10, 20, rl.Black)

	if drawButton(rl.Rectangle{X: 600, Y: 8, Width: 80, Height: 25}, "Settings") {
		*showSettings = !*showSettings
	}

	if drawButton(rl.Rectangle{X: 690, Y: 8, Width: 80, Height: 25}, "Open File") {
		if file := openFileDialog(); file != "" {
			s.currentFile = file
			s.reload()
		}
	}

	if s.debugInfo != "" {
		rl.DrawText(s.debugInfo, 450, 15, 10, rl.DarkGray)
	}

	if s.loadError != "" {
		rl.DrawText(s.loadError, 50, cfg.startY, 20, rl.Red)
	}

	if *showSettings {
		panelHeight := int32(90)
		panelWidth := int32(300)

		settingsRect := rl.Rectangle{X: 400 - float32(panelWidth/2), Y: float32(cfg.headerHeight + 5)}

		rl.DrawRectangle(
			int32(settingsRect.X),
			int32(settingsRect.Y),
			panelWidth,
			panelHeight,
			rl.ColorAlpha(rl.LightGray, 0.95),
		)

		rl.DrawRectangleLinesEx(
			rl.Rectangle{
				X:      settingsRect.X,
				Y:      settingsRect.Y,
				Width:  float32(panelWidth),
				Height: float32(panelHeight),
			},
			1,
			rl.Black,
		)

		oldMargin := s.margin
		oldGridSize := s.gridSize

		titleText := "Settings"
		titleWidth := rl.MeasureText(titleText, 15)
		rl.DrawText(titleText,
			int32(settingsRect.X+float32(panelWidth/2)-float32(titleWidth)/2),
			int32(settingsRect.Y+5),
			15,
			rl.Black)

		inputWidth := float32(60)
		inputHeight := float32(20)
		spacing := float32(40)

		totalWidth := inputWidth*2 + spacing
		startX := settingsRect.X + (float32(panelWidth)-totalWidth)/2

		marginInput := rl.Rectangle{
			X:      startX,
			Y:      settingsRect.Y + 45,
			Width:  inputWidth,
			Height: inputHeight,
		}

		gridInput := rl.Rectangle{
			X:      startX + inputWidth + spacing,
			Y:      settingsRect.Y + 45,
			Width:  inputWidth,
			Height: inputHeight,
		}

		s.margin = drawInputField(marginInput, "Margin", s.margin, 0, 10)
		s.gridSize = drawInputField(gridInput, "Grid Size", s.gridSize, 1, 64)

		helpText := "Use Up/Down keys when selected"
		helpWidth := rl.MeasureText(helpText, 10)
		helpX := settingsRect.X + float32(panelWidth/2) - float32(helpWidth)/2
		rl.DrawText(helpText, int32(helpX), int32(marginInput.Y+30), 10, rl.DarkGray)

		if oldMargin != s.margin || oldGridSize != s.gridSize {
			s.reload()
		}
	}
}

func main() {
	cfg := initConfig()
	state := initUI()
	defer rl.CloseWindow()
	defer state.rm.Close()

	showSettings := false
	rl.SetExitKey(0)

	for !rl.WindowShouldClose() {
		state.handleInput(&showSettings)

		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)

		state.renderSprites(cfg)
		state.renderUI(cfg, &showSettings)

		rl.EndDrawing()
	}
}

func openFileDialog() string {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("osascript", "-e", `POSIX path of (choose file with prompt "Choose a sprite sheet:" of type {"png","jpg","jpeg"})`)
	case "linux":
		cmd = exec.Command("zenity", "--file-selection", "--file-filter=Images (*.png *.jpg *.jpeg)")
	default:
		return ""
	}

	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(output))
}

func drawButton(bounds rl.Rectangle, text string) bool {
	mousePoint := rl.GetMousePosition()
	btnState := rl.ColorAlpha(rl.Gray, 0.6)
	isHovered := rl.CheckCollisionPointRec(mousePoint, bounds)
	isClicked := isHovered && rl.IsMouseButtonPressed(rl.MouseLeftButton)

	if isHovered {
		btnState = rl.ColorAlpha(rl.DarkGray, 0.6)
	}

	rl.DrawRectangleRec(bounds, btnState)
	rl.DrawText(text, int32(bounds.X+bounds.Width/2-float32(rl.MeasureText(text, 10))/2),
		int32(bounds.Y+bounds.Height/2-5), 10, rl.Black)

	return isClicked
}

func drawInputField(bounds rl.Rectangle, label string, value int32, min, max int32) int32 {
	rl.DrawText(label, int32(bounds.X), int32(bounds.Y-15), 10, rl.Black)
	rl.DrawRectangleRec(bounds, rl.White)
	rl.DrawRectangleLinesEx(bounds, 1, rl.Gray)

	valueText := strconv.Itoa(int(value))
	textX := int32(bounds.X + 5)
	textY := int32(bounds.Y + bounds.Height/2 - 5)
	rl.DrawText(valueText, textX, textY, 10, rl.Black)

	mousePoint := rl.GetMousePosition()
	if rl.CheckCollisionPointRec(mousePoint, bounds) {
		if rl.IsKeyPressed(rl.KeyUp) {
			value = int32(rl.Clamp(float32(value+1), float32(min), float32(max)))
		} else if rl.IsKeyPressed(rl.KeyDown) {
			value = int32(rl.Clamp(float32(value-1), float32(min), float32(max)))
		}
	}

	return value
}

func naturalSort(a, b string) bool {
	aParts := strings.Split(a, "_")
	bParts := strings.Split(b, "_")

	minLen := len(aParts)
	if len(bParts) < minLen {
		minLen = len(bParts)
	}

	for i := 0; i < minLen; i++ {
		aNum, aErr := strconv.Atoi(aParts[i])
		bNum, bErr := strconv.Atoi(bParts[i])

		if aErr == nil && bErr == nil {
			if aNum != bNum {
				return aNum < bNum
			}
		} else {
			if aParts[i] != bParts[i] {
				return aParts[i] < bParts[i]
			}
		}
	}
	return len(aParts) < len(bParts)
}
