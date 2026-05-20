package main

import (
	"embed"
	_ "embed"
	"fmt"
	"log"
	"runtime"
	"time"

	"github.com/BenedictKing/ccx/desktop/internal/backend"
	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
	"github.com/wailsapp/wails/v3/pkg/services/dock"
)

//go:embed all:frontend/dist
var assets embed.FS

func init() {
	application.RegisterEvent[string]("desktop:show-tab")
	application.RegisterEvent[string]("desktop:tray-error")
}

func main() {
	manager := backend.NewManager(backend.Options{})
	desktopService := NewDesktopService(manager)
	dockService := dock.New()

	app := application.New(application.Options{
		Name:        "CCX Desktop",
		Description: "CCX desktop shell and core service supervisor",
		Services: []application.Service{
			application.NewService(desktopService),
			application.NewService(dockService),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: false,
		},
	})
	desktopService.setApp(app)

	mainWindow := app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:     "CCX Desktop",
		Width:     1180,
		Height:    820,
		MinWidth:  960,
		MinHeight: 640,
		Mac: application.MacWindow{
			InvisibleTitleBarHeight: 50,
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarHiddenInset,
		},
		BackgroundColour: application.NewRGB(18, 24, 38),
		URL:              "/",
	})
	desktopService.setMainWindow(mainWindow)

	var mainWindowCentered bool
	showMainWindow := func(withFocus bool) {
		if !mainWindowCentered {
			mainWindow.Center()
			mainWindowCentered = true
		}
		if mainWindow.IsMinimised() {
			mainWindow.UnMinimise()
		}
		mainWindow.Show()
		if withFocus {
			if runtime.GOOS == "windows" {
				mainWindow.SetAlwaysOnTop(true)
				mainWindow.Focus()
				go func() {
					time.Sleep(150 * time.Millisecond)
					mainWindow.SetAlwaysOnTop(false)
				}()
			} else {
				mainWindow.Focus()
			}
		}
	}

	mainWindow.RegisterHook(events.Common.WindowClosing, func(e *application.WindowEvent) {
		mainWindow.Hide()
		e.Cancel()
	})

	app.Event.OnApplicationEvent(events.Mac.ApplicationShouldHandleReopen, func(event *application.ApplicationEvent) {
		showMainWindow(true)
	})

	app.OnShutdown(func() {
		desktopService.Shutdown()
	})

	tray := app.SystemTray.New()
	tray.SetTooltip("CCX Desktop")
	if icon, err := assets.ReadFile("frontend/dist/wails.png"); err == nil && len(icon) > 0 {
		tray.SetTemplateIcon(icon)
	}

	trayAction := func(label string, fn func() error) {
		go func() {
			if err := fn(); err != nil {
				log.Printf("[Desktop-Tray] %s 失败: %v", label, err)
				app.Event.Emit("desktop:tray-error", fmt.Sprintf("%s 失败: %v", label, err))
			}
		}()
	}

	buildTrayMenu := func(running bool, autostartEnabled bool) *application.Menu {
		menu := application.NewMenu()

		menu.Add("打开 CCX Web UI").OnClick(func(ctx *application.Context) {
			trayAction("打开 CCX Web UI", desktopService.ShowWebUITab)
		})
		menu.Add("显示状态页").OnClick(func(ctx *application.Context) {
			showMainWindow(true)
			app.Event.Emit("desktop:show-tab", "status")
		})
		menu.Add("显示 Agent 配置").OnClick(func(ctx *application.Context) {
			showMainWindow(true)
			app.Event.Emit("desktop:show-tab", "agent")
		})

		menu.AddSeparator()

		startItem := menu.Add("启动服务")
		startItem.OnClick(func(ctx *application.Context) {
			trayAction("启动服务", desktopService.StartService)
		})
		startItem.SetHidden(running)

		stopItem := menu.Add("停止服务")
		stopItem.OnClick(func(ctx *application.Context) {
			trayAction("停止服务", desktopService.StopService)
		})
		stopItem.SetHidden(!running)

		restartItem := menu.Add("重启服务")
		restartItem.OnClick(func(ctx *application.Context) {
			trayAction("重启服务", desktopService.RestartService)
		})
		restartItem.SetHidden(!running)

		menu.Add("在浏览器中打开").OnClick(func(ctx *application.Context) {
			trayAction("在浏览器中打开", desktopService.OpenWebUIInBrowser)
		})

		menu.AddSeparator()

		autostartItem := menu.AddCheckbox("开机自启", autostartEnabled)
		autostartItem.OnClick(func(ctx *application.Context) {
			newState := !autostartItem.Checked()
			if err := desktopService.SetAutostart(newState); err != nil {
				log.Printf("[Desktop-Tray] 切换开机自启失败: %v", err)
				app.Event.Emit("desktop:tray-error", fmt.Sprintf("切换开机自启失败: %v", err))
			}
		})

		menu.AddSeparator()

		menu.Add("退出").OnClick(func(ctx *application.Context) {
			app.Quit()
		})

		return menu
	}

	// 初始化托盘菜单
	initialStatus := desktopService.GetStatus()
	initialAutostart, _ := app.Autostart.IsEnabled()
	tray.SetMenu(buildTrayMenu(initialStatus.Running, initialAutostart))

	// 状态变化时动态刷新菜单
	go func() {
		ticker := time.NewTicker(3 * time.Second)
		defer ticker.Stop()
		lastRunning := initialStatus.Running
		lastAutostart := initialAutostart
		for range ticker.C {
			st := desktopService.GetStatus()
			asEnabled, _ := app.Autostart.IsEnabled()
			if st.Running != lastRunning || asEnabled != lastAutostart {
				lastRunning = st.Running
				lastAutostart = asEnabled
				tray.SetMenu(buildTrayMenu(st.Running, asEnabled))
			}
		}
	}()

	tray.AttachWindow(mainWindow)

	showMainWindow(false)

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
