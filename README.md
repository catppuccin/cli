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
- Clone the repository and switch to the dev branch. 
- Make all the changes to dev branch. To build the executable run `go build -o ctp`. 

## Note 
- You need to make sure that you set the environmental variable `$ORG_OVERRIDE` to `catppuccin-rfc` or the tool will search for `.catppuccin.yaml` in the `catppuccin` organisation which currently doesn't host the yaml files. This is a temporary measure to test the tool during its development.  

## TODO
- [ ] Hooks 
- [ ] Command hooks 
- [ ] Web hooks: To handle `xdg-open`, `open` or equivalent command on Windows. 
- [ ] Rework remove function from scratch: Need to find a way to save the flavour user installs. 
- [ ] Better error handling overall 
- [ ] Use to `gofmt` from now on. 
- [ ] Rewrite the wiki for catppuccin/cli to make first contributions easier. 

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
