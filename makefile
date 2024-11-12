# Check if we're running on Windows
ifeq ($(OS),Windows_NT)
    MAKEOS := windows
else
    MAKEOS := unix
endif

build:
	@echo "Running tests..."
	@go test || exit
	@printf "\n"

	make $(MAKEOS)
windows:
	SHELL := powershell.exe
	@printf "Building for $(shell @ver)"

unix:
	@# Build the package and put it
	@# in your local bin folder
	@printf "Building for $(shell uname | tr '[:upper:]' '[:lower:]').\n"
	go build -o bin/gurl .
	mv bin/gurl ~/.local/bin

	@printf "\n"
	@gurl