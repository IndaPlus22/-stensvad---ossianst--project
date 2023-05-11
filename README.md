# stensvad-ossianst-melvinbe-project
## Planets! 🌎
This project is a 3D renderer and planet generator made with [Go](https://go.dev/dl/ "Golang download") and [go.gl](https://github.com/go-gl/gl "go.gl page").

## Running the program

### Compile and run from the source code

To compile the program you need a working Go environment and a cgo compiler (like gcc) installed. Open project folder and run the following commands to download the required libraries:
* go get github.com/go-gl/gl/v4.1-core/gl
* go get github.com/go-gl/glfw/v3.3/glfw
* go get github.com/go-gl/mathgl/mgl32

It is recommended to have the latest GPU drivers installed.

On some computers there are problems getting all the required installations. In that case we recommend going through the official installation process given in the [go-gl docs](https://github.com/go-gl/gl "go-gl docs page")..

After the installation use the command go run . in the src folder to compile and run the program.


### Running the executable

**macOS:**
Download the program and run the src file in the src folder.

**Windows:**
Download the program and run the src.exe file in the src folder.

