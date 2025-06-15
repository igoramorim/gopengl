H# What
Converting [my studies about OpenGL done in C++](https://github.com/igoramorim/learn-opengl) to Go.

The idea is to have small scenes that show how to do _something_ using OpenGL.

# How to run

To show the possible scenes.
````
$ go run cmd/cli/main.go
````

To execute a scene.
````
$ go run cmd/cli/main.go ${scene}
````

## Note

I used [assimp-go](https://github.com/bloeys/assimp-go) to load 3D models in some scenes.

Not sure why but to make it work on **MacOS** I had to build [assimp](https://github.com/assimp/assimp/blob/master/Build.md) from source. Copied the `assimp/bin/libassimp.5.dylib` generated to `/usr/local/bin/libassimp.5.dylib`.

Now to run a scene I have to:
````
go build cmd/cli/main.go && install_name_tool -add_rpath /usr/local/bin main && ./main {scene}
````

If there is no changes in the source code between running two scenes there is no need to execute the `build` and `install_name_tool` commands. In fact, trying to run will return an error saying that the the binary is already in the `LC_RPATH`.

# Scenes

## triangle
![](/images/triangle.png)

## shaders_uniforms
![](/images/shaders_uniforms.png)

## textures
![](/images/textures.png)

## transformations
![](/images/transformations.png)

## coordinate_system
![](/images/coordinate_system.png)

## cube
![](/images/cube.png)

## camera
no preview

## light_colors
![](/images/light_colors.png)

## basic_light
![](/images/basic_light.png)

## materials
![](/images/materials.png)

## light_maps
![](/images/light_maps.png)

## directional_light
![](/images/directional_light.png)

## point_light
![](/images/point_light.png)

## spotlight
![](/images/spotlight.png)

## model_loading
![](/images/model_loading.png)

## depth_testing
![](/images/depth_testing.png)

## stencil_testing
![](/images/stencil_testing.png)
