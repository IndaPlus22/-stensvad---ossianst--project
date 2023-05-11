// Shader taken from https://learnopengl.com/Advanced-OpenGL/Cubemaps
#shader vertex
#version 330 core
layout (location = 0) in vec3 aPos;

out vec3 TexCoords;

uniform mat4 projection;
uniform mat4 view;

void main()
{
    TexCoords = aPos;
    vec4 pos = projection * view * vec4(aPos, 1.0);
    gl_Position = pos.xyww;
}  

#shader fragment
#version 330 core
layout(location = 0) out vec4 FragColor;
layout(location = 1) out vec4 DepthColor;

in vec3 TexCoords;

uniform samplerCube skybox;

void main()
{    
    FragColor = texture(skybox, TexCoords);
    DepthColor.r = 1.0;
}