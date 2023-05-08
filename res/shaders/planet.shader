#shader vertex
#version 330

layout (location = 0) in vec3 aPos;
layout (location = 1) in vec3 aNormal;

out vec3 VertexPos;
out vec3 VertexNormal;
out vec3 FragPos;
out vec3 Normal;
out mat4 Model;

uniform mat4 model;
uniform mat4 view;
uniform mat4 projection;

void main() {
    VertexPos = aPos;
    VertexNormal = aNormal;
    Model = model;
    FragPos = vec3(model * vec4(aPos, 1.0));
    Normal = mat3(transpose(inverse(model))) * aNormal; 

    gl_Position = projection * view * vec4(FragPos, 1.0);
}

#shader fragment
#version 330

in vec3 VertexPos;
in vec3 VertexNormal;
in vec3 FragPos;
in vec3 Normal;
in mat4 Model;

out vec4 FragColor;

// Colors
uniform vec3 shoreColLow;
uniform vec3 shoreColHigh;
uniform vec3 flatColLow;
uniform vec3 flatColHigh;
uniform vec3 steepColLow;
uniform vec3 steepColHigh;

// Textures
uniform sampler2D mainTexture;
uniform sampler2D normalMap;
uniform float texScale;
uniform float nMapScale;

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
    float gloss = 1.5;

    vec3 lightToFrag = normalize(FragPos - lightPos);
    vec3 camToFrag = normalize(camPos - FragPos);
    vec3 reflection = reflect(lightToFrag, normal);

    // The reflection is calculated by the dot product of the reflected light on the surface
    // and the vector from the surface to the camera
    float specLight = pow(clamp(dot(camToFrag, reflection), 0.0, 1.0), gloss);
    return specLight * intensity;
}

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
    vec3 weight = vec3(pow(abs(Normal.x), 0.5), pow(abs(Normal.y), 0.5), pow(abs(Normal.z), 0.5));
    weight /= dot(weight, vec3(1));

    // Return final color
    return colX * weight.x + colY * weight.y + colZ * weight.z;
}

// Maps normal map to six sides of the model
vec3 triplanarNormal(vec3 pos, sampler2D normalMap) {
    // Calculate tex coords for sampling in three directions
    // fract ensures that all sampling is within [0.0, 1.0]
    vec2 uvX = vec2(fract(pos.z * nMapScale), fract(pos.y * nMapScale));
    vec2 uvY = vec2(fract(pos.x * nMapScale), fract(pos.z * nMapScale));
    vec2 uvZ = vec2(fract(pos.x * nMapScale), fract(pos.y * nMapScale));

    // Sample normalmap normals
    // Also convert color range [0.0, 1.0] to normal range [-1.0, 1.0]
    vec3 normalX = texture(normalMap, uvX).rgb * 2.0 - 1.0;
    vec3 normalY = texture(normalMap, uvY).rgb * 2.0 - 1.0;
    vec3 normalZ = texture(normalMap, uvZ).rgb * 2.0 - 1.0;

    // Calculate normal in every direction
    vec3 tnormalX = vec3(normalX.xy + VertexNormal.zy, normalX.z * VertexNormal.x);
    vec3 tnormalY = vec3(normalY.xy + VertexNormal.xz, normalY.z * VertexNormal.y);
    vec3 tnormalZ = vec3(normalZ.xy + VertexNormal.xy, normalZ.z * VertexNormal.z);

    // Calculate how much every tangent normal will contribute to final normal
    vec3 weight = vec3(pow(abs(VertexNormal.x), 1.5), pow(abs(VertexNormal.y), 1.5), pow(abs(VertexNormal.z), 1.5));
    weight /= dot(weight, vec3(1.0));

    // Calculate normal in model space
    vec3 modelNormal = normalize(tnormalX.zyx * weight.x + tnormalY.xzy * weight.y + tnormalZ.xyz * weight.z);
    // Return normal in world space
    return mat3(transpose(inverse(Model))) * modelNormal;
}

vec3 lerp(vec3 va, vec3 vb, float k) {
    k = clamp(k, 0.0, 1.0);
    return va * (1.0 - k) + vb * k;
}

// Calculates color based on model height
vec3 heightColor(vec3 pos) {
    vec3 col;

    float height = length(pos) - 1;
    float flatness = dot(normalize(VertexNormal), normalize(pos));

    col = shoreColLow;

    col = lerp(col, shoreColHigh, (height - 0.01) / (0.02 - 0.01));

    col = lerp(col, flatColLow, (height - 0.02) / (0.04 - 0.03));

    col = lerp(col, flatColHigh, (height - 0.04) / (0.05 - 0.04));

    // Color less flat areas as steep color
    col = lerp(col, steepColLow, (max(height - 0.02, 0) * (5.0)) * (0.9 - flatness) * (15.0));

    col = lerp(col, steepColLow, (height - 0.05) / (0.05));

    col = lerp(col, steepColHigh, (height - 0.12) * (50.0));

    return col;
}

void main() {
    // Partially sample texture
    vec3 texColor = vec3(0.7) + triplanarTexture(VertexPos, mainTexture) * 0.3;

    // Calculate normal to be used for lighting calculations
    vec3 lightingNormal = triplanarNormal(VertexPos, normalMap);

    vec3 heightColor = heightColor(VertexPos);

    // Ambient light: the natural light in space
    float ambientLight = 0.1;
    float diffuseLight = calculateDiffuseLight(lightingNormal);
    float specularLight = calculateSpecularLight(lightingNormal);

    // Phong shading combines the different lighting types
    float phong = ambientLight + diffuseLight + specularLight;

    FragColor = vec4(texColor * heightColor * phong * lightColor, 1.0);

    //FragColor = vec4(lightingNormal, 1.0);
}