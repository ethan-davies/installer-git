# Installer Template

This is a template for creating a installer for your application. Follow these instructions to set up and customize the installer for your own use.

## Setup

1. Clone or download this repository to your local machine.

2. Open the `main.go` file in your preferred text editor.

3. Locate the following section of code in `main.go`:

```go
const gitHubAPIRepoURL = "https://api.github.com/repos/ethan-davies/myapp"
const gitHubRepoURL = "https://github.com/ethan-davies/myapp"
const name = "myapp"
const installFileName = "myapp"

const windowsBinaryName = "myapp"
const linuxBinaryName = "myapp"
const darwinBinaryName = "myapp-macos"
```
Change these variables to your liking. Following the instructions.

## Build
After you have finished setting up your installer you can build it. 

If you are on windows use the following command:
```bash
./build.bat
```
This should build the installer into a `bin` folder

Otherwise, if you are on MacOS or Linux run 
```bash
chmod +x build.sh
./your_script.sh
```