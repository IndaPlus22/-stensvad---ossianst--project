#shader vertex
#version 330

layout (location = 0) in vec2 aPos;
layout (location = 1) in vec2 aTexCoords;

out vec2 texCoords;

void main() {
    gl_Position = vec4(aPos.x, aPos.y, 0.0, 1.0);
    texCoords = aTexCoords;
}

#shader fragment
#version 330

out vec4 FragColor;
in vec2 texCoords;

uniform vec3 camDir;
uniform vec3 camPos;
uniform mat4 viewMatrix;
uniform mat4 projMatrix;
uniform sampler2D screenTexture;
//uniform sampler2D depthTexture;

// Calculate how much of a ray from the camera intersects with a sphere (the atmosphere)
// Returns a vector with the distance to the sphere and the travelled distance through it
vec2 raySphereIntersection(vec3 camPosition, vec3 cameraDirection, vec3 sphereOrigin, float sphereRadius) {
    vec3 offset = camPosition - sphereOrigin;
    float a = 1.0;
    float b = 2.0 * dot(offset, cameraDirection);
    float c = dot(offset, offset) - sphereRadius * sphereRadius;
    float d = b * b - 4.0 * a * c;

    if (d > 0.0) {
        float s = sqrt(d);
        float distToSphereNear = max(0.0, (-b - s) / (2.0 * a));
        float distToSphereFar = (-b + s) / (2.0 * a);

        if (distToSphereFar >= 0.0) {
            return vec2(distToSphereNear, distToSphereFar - distToSphereNear);
        }
    }

    return vec2(1000.0, 0.0);
}

void main() {
    float sphereRadius = 1.3;
    vec4 originalColor = texture(screenTexture, texCoords);

    mat4 inverseViewProjMatrix = inverse(projMatrix * viewMatrix);
    vec4 clipCoord = vec4(texCoords * 2.0 - 1.0, 0.0, 1.0);
    vec4 worldCoord = inverseViewProjMatrix * clipCoord;
    worldCoord /= worldCoord.w;

    vec3 fragRay = normalize(worldCoord.xyz + (vec3(0.0) - camPos));

    vec2 hitInfo = raySphereIntersection(worldCoord.xyz, fragRay, vec3(0.0), sphereRadius);

    float distToAtmosphere = hitInfo.x;
    float distThroughAtmosphere = hitInfo.y;

    FragColor = originalColor + vec4(distThroughAtmosphere / (2.0 * sphereRadius)) * vec4(0.4, 0.7, 1.0, 0.5);
}
