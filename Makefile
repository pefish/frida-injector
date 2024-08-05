
name = "Frida Injector"

mac-arm64:
	wails build -clean -platform darwin/arm64
	create-dmg ./build/bin/$(name).dmg ./build/bin/$(name).app
	rm -rf ./build/bin/$(name).app

mac-amd64:
	wails build -clean -platform darwin/amd64
	create-dmg ./build/bin/$(name).dmg ./build/bin/$(name).app
	rm -rf ./build/bin/$(name).app


windows-amd64:
	wails build -nsis -platform windows/amd64

windows-arm64:
	wails build -nsis -platform windows/arm64
