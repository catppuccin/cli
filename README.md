<h3 align="center">
	<img src="https://raw.githubusercontent.com/catppuccin/catppuccin/main/assets/logos/exports/1544x1544_circle.png" width="100" alt="Logo"/><br/>
	<img src="https://raw.githubusercontent.com/catppuccin/catppuccin/main/assets/misc/transparent.png" height="30" width="0px"/>
	Catppuccin CLI
	<img src="https://raw.githubusercontent.com/catppuccin/catppuccin/main/assets/misc/transparent.png" height="30" width="0px"/>
</h3>

<p align="center">
	<a href="https://github.com/catppuccin/cli/stargazers"><img src="https://img.shields.io/github/stars/catppuccin/cli?colorA=363a4f&colorB=b7bdf8&style=for-the-badge"></a>
	<a href="https://github.com/catppuccin/cli/issues"><img src="https://img.shields.io/github/issues/catppuccin/cli?colorA=363a4f&colorB=f5a97f&style=for-the-badge"></a>
	<a href="https://github.com/catppuccin/cli/contributors"><img src="https://img.shields.io/github/contributors/catppuccin/cli?colorA=363a4f&colorB=a6da95&style=for-the-badge"></a>
</p>

<p align="center">

<img src="https://raw.githubusercontent.com/catppuccin/catppuccin/main/assets/misc/sample.png"/>
</p>

## Catppuccin CLI
A work-in-progress CLI for Catppuccin themes.

It allows you to:
- Install themes with one command
- Uninstall themes with one command
- Update themes with one command

## Installation
- You can download the executable for this project from the releases section. Download the release as per your OS. 

## Development 
- The foremost requirement to develop is to make sure that go version 1.19 is installed. 
- Development and contribution guidelines along with the future development plans have been added to the [wiki](https://github.com/catppuccin/cli/wiki/Contributing).

## Building
You can simply build the cli with `make`. This will automatically install the required dependencies and build the program. The outputed executable will be in `builds/`. You can also use GoReleaser to build it for all platforms.

## Docker image 
- The cli also has a Docker image. To build it, run `docker build --network=host -t ctp:latest .`.
- To run the built image, run the command `docker run -it --rm --net=host ctp:latest help`.

## TODO
- [ ] Hooks 
  - [x] Install hooks
  - [ ] Uninstall hooks
- [x] Command hooks: To execute shell scripts.  
- [x] Web hooks: To handle `xdg-open`, `open` or equivalent command on Windows. 
- [x] Rework remove function from scratch: Need to find a way to save the flavour user installs. 
- [x] Better error handling overall 
- [x] Use `gofmt` from now on. 
- [x] Rewrite the wiki for catppuccin/cli to make first contributions easier. 
- [x] Refactoring: 
  - [x] Move `cmd` to `internal`
  - [x] Move `main.go` to `cmd/ctp` => Reason: Check [#25](https://github.com/catppuccin/cli/issues/25)


&nbsp;

<p align="center">
	<img src="https://raw.githubusercontent.com/catppuccin/catppuccin/main/assets/footers/gray0_ctp_on_line.svg?sanitize=true" />
</p>

<p align="center">
	Copyright &copy; 2021-present <a href="https://github.com/catppuccin" target="_blank">Catppuccin Org</a>
</p>

<p align="center">
	<a href="https://github.com/catppuccin/catppuccin/blob/main/LICENSE"><img src="https://img.shields.io/static/v1.svg?style=for-the-badge&label=License&message=MIT&logoColor=d9e0ee&colorA=363a4f&colorB=b7bdf8"/></a>
</p>
