#shader vertex
#version 330

layout (location = 0) in vec3 aPos;
layout (location = 1) in vec3 aNormal;

out vec3 VertexPos;
out vec3 vertexNormal;
out vec3 Normal;

uniform mat4 model;
uniform mat4 view;
uniform mat4 projection;

void main() {
    VertexPos = aPos;
    vertexNormal = aNormal;
    vec3 FragPos = vec3(model * vec4(aPos, 1.0));
    Normal = mat3(transpose(inverse(model))) * aNormal; 

    gl_Position = projection * view * vec4(FragPos, 1.0);
}

#shader fragment
#version 330

in vec3 VertexPos;
in vec3 vertexNormal;
in vec3 Normal;

out vec4 FragColor;

// Textures
uniform sampler2D mainTexture;
uniform float texScale;

// Maps a texture to six sides of the model
vec3 triplanarTexture(vec3 pos, sampler2D tex) {
    // Calculate tex coords for sampling in three directions
    // fract ensures that all sampling is within [0.0, 1.0]
    vec2 uvX = vec2(fract(pos.z * texScale), fract(pos.y * texScale));
    vec2 uvY = vec2(fract(pos.x * texScale), fract(pos.z * texScale));
    vec2 uvZ = vec2(fract(pos.x * texScale), fract(pos.y * texScale));

    // Sample texture colors
    vec3 colX = vec3(texture(tex, uvX));
    vec3 colY = vec3(texture(tex, uvY));
    vec3 colZ = vec3(texture(tex, uvZ));

    // Calculate how much every color will contribute to final color
    vec3 weight = vec3(pow(abs(Normal.x), 1.0), pow(abs(Normal.y), 1.0), pow(abs(Normal.z), 1.0));
    weight /= dot(weight, vec3(1));

    // Return final color
    return colX * weight.x + colY * weight.y + colZ * weight.z;
}

void main() {
    vec3 texColor = triplanarTexture(VertexPos, mainTexture);

    FragColor = vec4(texColor, 1.0);
}