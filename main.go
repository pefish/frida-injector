package main

import (
	"embed"
	"fmt"

	runtime1 "runtime"

	"github.com/tidwall/gjson"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/menu/keys"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

//go:embed all:frontend/dist
var assets embed.FS

//go:embed wails.json
var wailsJSON string

func main() {
	app := NewApp()

	version := gjson.Get(wailsJSON, "info.productVersion").String()
	name := gjson.Get(wailsJSON, "info.productName").String()

	AppMenu := menu.NewMenu()
	if runtime1.GOOS == "darwin" {
		mainSubmenu := AppMenu.AddSubmenu(name)
		mainSubmenu.AddText("About me", nil, func(cd *menu.CallbackData) {
			runtime.MessageDialog(app.ctx, runtime.MessageDialogOptions{
				Type:          runtime.InfoDialog,
				Title:         "About Me",
				Message:       fmt.Sprintf("version v%s", version),
				DefaultButton: "Close",
			})
		})
		mainSubmenu.AddSeparator()
		mainSubmenu.AddText("Quit", keys.CmdOrCtrl("q"), func(_ *menu.CallbackData) {
			runtime.Quit(app.ctx)
		})
		AppMenu.Append(menu.EditMenu())
	}

	err := wails.Run(&options.App{
		Title:  "Frida Injector",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
		DisableResize: true,
		Menu:          AppMenu,
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
