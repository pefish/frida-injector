package main

import (
	"context"
	"encoding/json"

	"github.com/frida/frida-go/frida"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx           context.Context
	deviceManager *frida.DeviceManager
	attachSession *frida.Session
	script        *frida.Script
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		deviceManager: frida.NewDeviceManager(),
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) ListDevices() ([]string, error) {
	devices, err := a.deviceManager.EnumerateDevices()
	if err != nil {
		return nil, err
	}
	deviceIds := make([]string, 0)
	for _, d := range devices {
		deviceIds = append(deviceIds, d.ID())
	}
	return deviceIds, nil
}

type ProcessInfo struct {
	Name     string `json:"name"`
	Id       int    `json:"id"`
	ParentId int    `json:"parent_id"`
	User     string `json:"user"`
	Path     string `json:"path"`
}

func (a *App) ListProcessesOfDevice(deviceId string) ([]*ProcessInfo, error) {
	device, err := a.deviceManager.DeviceByID(deviceId)
	if err != nil {
		return nil, err
	}
	processes, err := device.EnumerateProcesses(frida.ScopeFull)
	if err != nil {
		return nil, err
	}
	processInfos := make([]*ProcessInfo, 0)
	for _, p := range processes {
		pInfo := ProcessInfo{
			Name: p.Name(),
			Id:   p.PID(),
		}
		if ppid, ok := p.Params()["ppid"]; ok {
			pInfo.ParentId = int(ppid.(int64))
		}
		if user, ok := p.Params()["user"]; ok {
			pInfo.User = user.(string)
		}
		if path, ok := p.Params()["path"]; ok {
			pInfo.Path = path.(string)
		}
		processInfos = append(processInfos, &pInfo)
	}
	return processInfos, nil
}

func (a *App) DetachProcess() error {
	if a.attachSession != nil && !a.attachSession.IsDetached() {
		return a.attachSession.Detach()
	}
	return nil
}

func (a *App) AttachProcess(deviceId string, pid int) error {
	device, err := a.deviceManager.DeviceByID(deviceId)
	if err != nil {
		return err
	}

	session, err := device.Attach(pid, nil)
	if err != nil {
		return err
	}
	a.attachSession = session

	return nil
}

type Msg struct {
	Type string `json:"type"`

	Description  string `json:"description,omitempty"`
	Stack        string `json:"stack,omitempty"`
	FileName     string `json:"fileName,omitempty"`
	LineNumber   uint64 `json:"lineNumber,omitempty"`
	ColumnNumber uint64 `json:"columnNumber,omitempty"`

	Level   string `json:"level,omitempty"`
	Payload string `json:"payload,omitempty"`
}

func (a *App) CancelScript() error {
	if a.script != nil && !a.script.IsDestroyed() {
		return a.script.Unload()
	}
	return nil
}

func (a *App) InjectScript(scriptStr string) error {
	script, err := a.attachSession.CreateScript(scriptStr)
	if err != nil {
		return err
	}
	a.script = script

	script.On("message", func(msgStr string) {
		var msg Msg
		err := json.Unmarshal([]byte(msgStr), &msg)
		if err != nil {
			runtime.EventsEmit(a.ctx, "log-parse-error", msgStr)
			return
		}
		switch msg.Type {
		case "error":
			runtime.EventsEmit(a.ctx, "error", msg)
		case "log":
			runtime.EventsEmit(a.ctx, "log", msg)
		}

	})
	err = script.Load()
	if err != nil {
		return err
	}

	return nil
}
