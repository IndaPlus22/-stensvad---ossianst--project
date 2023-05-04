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
uniform vec3 lightPos;
uniform vec3 planetOrigin;

uniform float planetRadius;
uniform float atmosphereScale;

uniform float far;
uniform float near;

uniform mat4 viewMatrix;
uniform mat4 projMatrix;

uniform sampler2D colorTexture;
uniform sampler2D depthTexture;

// The solution for atmosphere scattering is based on Sebastian Lagues implementation
// in this video on YouTube: https://www.youtube.com/watch?v=DxfEbulyFcY
float scatteringStrength = 6.0;
vec3 wavelengths = vec3(7.0, 5.3, 4.4);
vec3 scatteringCoefficients = vec3(pow(4.0 / wavelengths.x, 4.0), pow(4.0 / wavelengths.y, 4.0), pow(4.0 / wavelengths.z, 4.0)) * scatteringStrength;

// Calculate how much of a ray from the camera intersects with a sphere
// Returns a vector with the distance to the sphere and the travelled distance through it
vec2 raySphereIntersection(vec3 worldCoord, vec3 rayDirection, vec3 sphereOrigin, float sphereRadius) {
    vec3 offset = worldCoord - sphereOrigin;
    float a = 1.0;
    float b = 2.0 * dot(offset, rayDirection);
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

// TODO: Städa, gör mer readable, skriv kommentarer, ta bort onödig boilerplate
float densityAtPoint(vec3 point) {
    float heightAboveSurface = length(point - planetOrigin) - planetRadius;
    float height01 = heightAboveSurface / (atmosphereScale - planetRadius);
    float localDensity = exp(-height01) * (1 - height01);

    return localDensity;
}

// TODO: Städa, gör mer readable, skriv kommentarer, ta bort onödig boilerplate
float opticalDepth(vec3 rayOrigin, vec3 rayDir, float rayLength) {
    float numOpticalDepthPoints = 10.0;
    vec3 point = rayOrigin;
    float stepSize = rayLength / (numOpticalDepthPoints - 1);
    float opticalDepth = 0.0;

    for (int i = 0; i < numOpticalDepthPoints; i++) {
        float localDensity = densityAtPoint(point);
        opticalDepth += localDensity * stepSize;
        point += rayDir * stepSize;
    }

    return opticalDepth;
}

// TODO: Städa, gör mer readable, skriv kommentarer, ta bort onödig boilerplate
vec3 rayleighScattering(vec3 rayOrigin, vec3 rayDir, float rayLength, vec3 originalColor) {
    vec3 scatteringPoint = rayOrigin;
    float scatteringPoints = 15.0;
    float stepSize = rayLength / (scatteringPoints - 1);
    vec3 totalScattering = vec3(0.0);
    float viewRayOpticalDepth = 0.0;

    for (int i = 0; i < scatteringPoints; i++) {
        vec3 sunDir = lightPos - scatteringPoint;
        vec2 sunRayLength = raySphereIntersection(scatteringPoint, sunDir, planetOrigin, atmosphereScale);
        float sunRayOpticalDepth = opticalDepth(scatteringPoint, sunDir, sunRayLength.y);
        viewRayOpticalDepth = opticalDepth(scatteringPoint, -rayDir, stepSize * i);
        
        vec3 transmittance = vec3(exp(-(sunRayOpticalDepth + viewRayOpticalDepth) * scatteringCoefficients));
        float localDensity = densityAtPoint(scatteringPoint);

        totalScattering += localDensity * transmittance * stepSize * scatteringCoefficients;
        scatteringPoint += rayDir * stepSize;
    }
    float originalColorTransmittance = exp(-viewRayOpticalDepth);
    return originalColor * originalColorTransmittance + totalScattering;
}

float LinearizeDepth(float z) {
    return (z - near) / (far - near);
}

void main() {
    vec4 finalColor = texture(colorTexture, texCoords);

    // Get the world coordinates of the current fragment/pixel
    mat4 inverseViewProjMatrix = inverse(projMatrix * viewMatrix);
    vec4 clipCoord = vec4(texCoords * 2.0 - 1.0, 0.0, 1.0);
    vec4 worldCoord = inverseViewProjMatrix * clipCoord;
    worldCoord /= worldCoord.w;

    // A ray from the current fragment in direction of the camera relative to the planet
    vec3 fragRay = worldCoord.xyz + (planetOrigin - camPos);
    
    float depth = LinearizeDepth(texture(depthTexture, texCoords).r);
    depth *= length(fragRay);

    fragRay = normalize(fragRay);

    vec2 intersection = raySphereIntersection(worldCoord.xyz, fragRay, planetOrigin, atmosphereScale);

    float distToAtmosphere = intersection.x;
    float distThroughAtmosphere = min(intersection.y, depth - distToAtmosphere);

    if (distThroughAtmosphere > 0.0) {
        vec3 point = worldCoord.xyz + fragRay * (distToAtmosphere);
        vec3 light = rayleighScattering(point, fragRay, distThroughAtmosphere, finalColor.xyz);
        finalColor = vec4(light, 0);
    }

    FragColor = finalColor;
}