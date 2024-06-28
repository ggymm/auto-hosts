buildGUI:
	@echo "Building GUI"
	@cd cmd && go build -o ../bin/AutoHostsGUI.exe

buildCLI:
	@echo "Building CLI"
	@cd cmd/cli && go build -o ../../bin/AutoHostsCLI.exe