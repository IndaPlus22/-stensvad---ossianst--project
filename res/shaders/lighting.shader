#shader vertex
#version 330

layout (location = 0) in vec3 aPos;
layout (location = 1) in vec3 aNormal;

out vec3 VertexPos;
out vec3 FragPos;
out vec3 Normal;

uniform mat4 model;
uniform mat4 view;
uniform mat4 projection;

void main() {
    VertexPos = aPos;
    FragPos = vec3(model * vec4(aPos, 1.0));
    Normal = mat3(transpose(inverse(model))) * aNormal; 

    gl_Position = projection * view * vec4(FragPos, 1.0);
}

#shader fragment
#version 330

in vec3 VertexPos;
in vec3 FragPos;
in vec3 Normal;

out vec4 FragColor;

uniform sampler2D mainTexture;
uniform float texScale;
uniform sampler2D normalMap;
uniform float nmapScale;

uniform vec3 camPos;
uniform vec3 lightPos;
uniform vec3 lightColor;

// Diffuse light lights a surface in relation to its angle to the light source
float calculateDiffuseLight(vec3 normal) {
    vec3 lightDirection = normalize(FragPos - lightPos);
    return clamp(dot(normal, -lightDirection), 0.0, 0.9);
}

// Specular light is the refletion on glossy areas
float calculateSpecularLight(vec3 normal) {
    // The intensity of the glow and how much light is reflected
    float intensity = 0.3;
    int gloss = 2;

    vec3 lightToFrag = normalize(FragPos - lightPos);
    vec3 camToFrag = normalize(camPos - FragPos);
    vec3 reflection = reflect(lightToFrag, normal);

    // The reflection is calculated by the dot product of the reflected light on the surface
    // and the vector from the surface to the camera
    float specLight = pow(clamp(dot(camToFrag, reflection), 0.0, 1.0), gloss);
    return specLight * intensity;
}

vec3 triplanarTexture(vec3 pos, sampler2D tex) {
    vec2 uvX = vec2(fract(pos.z * texScale), fract(pos.y * texScale));
    vec2 uvY = vec2(fract(pos.x * texScale), fract(pos.z * texScale));
    vec2 uvZ = vec2(fract(pos.x * texScale), fract(pos.y * texScale));

    vec3 colX = vec3(texture2D(tex, uvX));
    vec3 colY = vec3(texture2D(tex, uvY));
    vec3 colZ = vec3(texture2D(tex, uvZ));

    vec3 weight = vec3(pow(abs(Normal.x), 2), pow(abs(Normal.y), 2), pow(abs(Normal.z), 2));

    weight /= dot(weight, vec3(1));

    return colX * weight.x + colY * weight.y + colZ * weight.z;
}

vec3 triplanarNormal(vec3 pos, vec3 surfaceNormal, sampler2D normalMap) {
    vec2 uvX = vec2(fract(pos.z * nmapScale), fract(pos.y * nmapScale));
    vec2 uvY = vec2(fract(pos.x * nmapScale), fract(pos.z * nmapScale));
    vec2 uvZ = vec2(fract(pos.x * nmapScale), fract(pos.y * nmapScale));

    vec3 tnormalX = vec3(texture2D(normalMap, uvX));
    vec3 tnormalY = vec3(texture2D(normalMap, uvY));
    vec3 tnormalZ = vec3(texture2D(normalMap, uvZ));

    tnormalX = vec3(tnormalX.xy + surfaceNormal.zy, tnormalX.z * surfaceNormal.x);
    tnormalY = vec3(tnormalY.xy + surfaceNormal.xz, tnormalY.z * surfaceNormal.y);
    tnormalZ = vec3(tnormalZ.xy + surfaceNormal.xy, tnormalZ.z * surfaceNormal.z);

    vec3 weight = vec3(pow(abs(surfaceNormal.x), 2), pow(abs(surfaceNormal.y), 2), pow(abs(surfaceNormal.z), 2));
    weight /= dot(weight, vec3(1));

    return normalize(tnormalX.zyx * weight.x + tnormalY.xzy * weight.y + tnormalZ.xyz * weight.z);
}

void main() {
    // Triplanar mapping
    vec3 texColor = vec3(0.5) + triplanarTexture(VertexPos, mainTexture) * 0.3;

    vec3 lightingNormal = triplanarNormal(VertexPos, Normal, normalMap);

    // Ambient light: the natural light in space
    float ambientLight = 0.1;
    float diffuseLight = calculateDiffuseLight(lightingNormal);
    float specularLight = calculateSpecularLight(lightingNormal);

    // Phong shading combines the different lighting types
    float phong = ambientLight + diffuseLight + specularLight;

    // When adding textures, use texture2D() to get color value and multiply with the phong shading for the final FragColor
    FragColor = vec4(texColor * phong * lightColor, 1.0);
}