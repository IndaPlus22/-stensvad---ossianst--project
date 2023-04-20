#shader vertex
#version 330

layout (location = 0) in vec3 aPos;
layout (location = 1) in vec3 aNormal;

out vec3 FragPos;
out vec3 Normal;

uniform mat4 model;
uniform mat4 view;
uniform mat4 projection;

void main()
{
    FragPos = vec3(model * vec4(aPos, 1.0));
    Normal = mat3(transpose(inverse(model))) * aNormal; 

    gl_Position = projection * view * vec4(FragPos, 1.0);
}

#shader fragment
#version 330

in vec3 FragPos;
in vec3 Normal;

out vec4 FragColor;

uniform vec3 camPos;
uniform vec3 lightPos;
uniform vec3 lightColor;

// Diffuse light lights a surface in relation to its angle to the light source
float calculateDiffuseLight() {
    vec3 lightDirection = normalize(FragPos - lightPos);
    return clamp(dot(Normal, -lightDirection), 0.0, 0.9);
}

// Specular light is the refletion on glossy areas
float calculateSpecularLight() {
    // The intensity of the glow and how much light is reflected
    float intensity = 0.5;
    int gloss = 16;

    vec3 lightToFrag = normalize(FragPos - lightPos);
    vec3 camToFrag = normalize(camPos - FragPos);
    vec3 reflection = reflect(lightToFrag, Normal);

    // The reflection is calculated by the dot product of the reflected light on the surface
    // and the vector from the surface to the camera
    float specLight = pow(clamp(dot(camToFrag, reflection), 0.0, 1.0), gloss);
    return specLight * intensity;
}

void main()
{
    // Ambient light: the natural light in space
    vec3 ambientLight = vec3(0.1, 0.1, 0.1);
    float diffuseLight = calculateDiffuseLight();
    float specularLight = calculateSpecularLight();

    // Phong shading combines the different lighting types
    vec3 phong = ambientLight + diffuseLight + specularLight;

    // When adding textures, use texture2D() to get color value and multiply with the phong shading for the final FragColor
    FragColor = vec4(phong * lightColor, 1.0);
}