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
uniform float camNear;
uniform float camFar;

uniform vec3 lightPos;
uniform vec3 planetOrigin;

uniform mat4 viewMatrix;
uniform mat4 projMatrix;

uniform sampler2D colorTexture;
uniform sampler2D depthTexture;

uniform sampler2D oceanNormalMap;

// Create uniform block of up to 10 planets
layout(std140) uniform PlanetPositions {
    vec4 planetPositions[10];
};

float atmosphereScale = 0.5;

// The solution for atmosphere scattering is based on Sebastian Lagues implementation
// in this video on YouTube: https://www.youtube.com/watch?v=DxfEbulyFcY
vec3 scatteringCoefficients = 6.0 * vec3(pow(400.0 / 700.0, 4.0), 
                                         pow(400.0 / 530.0, 4.0), 
                                         pow(400.0 / 440.0, 4.0));

// Calculate how much of a ray from the camera intersects with a sphere (from video above)
// Returns a vector with the distance to the sphere and the travelled distance through it
vec2 raySphereIntersection(vec3 rayPosition, vec3 rayDirection, vec3 sphereOrigin, float sphereRadius) {
    vec3 offset = rayPosition - sphereOrigin;
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

// Get thickness of the atmosphere at the point in world space, 
// around the planet of the planetData (from video)
float densityAtPoint(vec3 point, vec4 planetData) {
    float heightAboveSurface = length(point - planetData.xyz) - planetData.w;
    float height01 = heightAboveSurface / (planetData.w + atmosphereScale - planetData.w);
    float localDensity = exp(-height01) * (1 - height01);

    return localDensity;
}

// Get the optical depth along the ray from rayOrigin in direction of rayDir
float opticalDepth(vec3 rayOrigin, vec3 rayDir, float rayLength, vec4 planetData) {
    vec3 point = rayOrigin;
    float numOpticalDepthPoints = 10.0;
    float stepSize = rayLength / (numOpticalDepthPoints - 1);
    float opticalDepth = 0.0;

    // Get total optical depth by adding the depth of each depth point
    for (int i = 0; i < numOpticalDepthPoints; i++) {
        float localDensity = densityAtPoint(point, planetData);
        opticalDepth += localDensity * stepSize;
        point += rayDir * stepSize;
    }

    return opticalDepth;
}

// Calculate atmosphere scattering along the ray from rayOrigin in direction of
// rayDir, through the planet of planetData (from video above)
vec3 scattering(vec3 rayOrigin, vec3 rayDir, float rayLength, vec3 originalColor, vec4 planetData) {
    // Get scattering on 15 points along the ray
    vec3 scatteringPoint = rayOrigin;
    float scatteringPoints = 15.0;
    float stepSize = rayLength / (scatteringPoints - 1);
    vec3 totalScattering = vec3(0.0);
    float viewRayOpticalDepth = 0.0;

    for (int i = 0; i < scatteringPoints; i++) {
        vec3 sunDir = lightPos - scatteringPoint;
        vec2 sunRayLength = raySphereIntersection(scatteringPoint, sunDir, planetData.xyz, planetData.w + atmosphereScale);
        float sunRayOpticalDepth = opticalDepth(scatteringPoint, sunDir, sunRayLength.y, planetData);
        viewRayOpticalDepth = opticalDepth(scatteringPoint, -rayDir, stepSize * i, planetData);
        
        // Calculate the light reaching each point multiplied by the wavelength coefficients
        vec3 transmittance = vec3(exp(-(sunRayOpticalDepth + viewRayOpticalDepth) * scatteringCoefficients));
        float localDensity = densityAtPoint(scatteringPoint, planetData);

        totalScattering += localDensity * transmittance * stepSize * scatteringCoefficients;
        scatteringPoint += rayDir * stepSize;
    }
    
    float originalColorTransmittance = exp(-viewRayOpticalDepth);
    return originalColor * originalColorTransmittance + totalScattering;
}

// Applies a glow effect using the same technique as the scattering function,
// emitting the wavelengths and transmittance to get a glowing fire effect
vec3 shine(vec3 rayOrigin, vec3 rayDir, float rayLength, vec3 originalColor, vec4 planetData) {
    // Calculates scattering on 5 points along the ray, starting at the ray origin
    vec3 scatteringPoint = rayOrigin;
    float scatteringPoints = 5.0;
    float stepSize = rayLength / (scatteringPoints - 1);
    vec3 totalScattering = vec3(0.0);

    for (int i = 0; i < scatteringPoints; i++) {
        float localDensity = densityAtPoint(scatteringPoint, planetData);

        totalScattering += localDensity * stepSize;
        scatteringPoint += rayDir * stepSize;
    }

    // Multiply by an orange vector to get a fire-y color
    return originalColor+ totalScattering * vec3(1.0, 0.7, 0.0);
}

// Maps normal map to six sides of the model
vec3 triplanarNormal(vec3 pos, vec3 normal, sampler2D normalMap) {
    float nMapScale = 2.0;

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
    vec3 tnormalX = vec3(normalX.xy + normal.zy, normalX.z * normal.x);
    vec3 tnormalY = vec3(normalY.xy + normal.xz, normalY.z * normal.y);
    vec3 tnormalZ = vec3(normalZ.xy + normal.xy, normalZ.z * normal.z);

    // Calculate how much every tangent normal will contribute to final normal
    vec3 weight = vec3(pow(abs(normal.x), 1.5), pow(abs(normal.y), 1.5), pow(abs(normal.z), 1.5));
    weight /= dot(weight, vec3(1.0));

    return normalize(tnormalX.zyx * weight.x + tnormalY.xzy * weight.y + tnormalZ.xyz * weight.z);
}

void main() {
    // Get the base color from the color texture
    vec4 finalColor = texture(colorTexture, texCoords);

    // Get the world coordinates of the current fragment/pixel
    mat4 inverseViewProjMatrix = inverse(projMatrix * viewMatrix);
    vec4 clipCoord = vec4(texCoords * 2.0 - 1.0, 0.0, 1.0);
    vec4 worldCoord = inverseViewProjMatrix * clipCoord;
    worldCoord /= worldCoord.w;

    // Cast a ray from the camera in direction towards the current fragment
    vec3 fragRay = normalize(worldCoord.xyz - camPos);

    // Get the depth from a custom depth texture
    float depth = texture(depthTexture, texCoords).r * (camFar - camNear);

    // Apply post processing effects for each planet
    for (int i = 0; i < 10; i++) {
        // The ocean effect is using the method presented by Sebastian Lague in 
        // this video on YouTube: https://youtu.be/lctXaT9pxA0
        vec2 oceanIntersection = raySphereIntersection(worldCoord.xyz, fragRay, planetPositions[i].xyz, planetPositions[i].w);
        float distToOcean = oceanIntersection.x;
        float distThroughOcean = oceanIntersection.y;
        float oceanViewDepth = min(distThroughOcean, depth - distToOcean);

        if (oceanViewDepth > 0.0 && i != 0) {
            // Calculate diffuse lighting
            vec3 surfaceFragPos = worldCoord.xyz + fragRay * distToOcean;
            vec3 normal = normalize(surfaceFragPos - planetPositions[i].xyz);
            normal = triplanarNormal(normal * planetPositions[i].w, normal, oceanNormalMap);

            float diffuseLight = clamp(dot(normal, -surfaceFragPos), 0.0, 0.7);

            // Calculate specular lighting
            vec3 lightToFrag = normalize(surfaceFragPos - planetPositions[0].xyz);
            vec3 camToFrag = normalize(surfaceFragPos - camPos);
            vec3 reflection = reflect(lightToFrag, normal);

            float specularValue = clamp(dot(reflection, -camToFrag), 0.0, 1.0);
            float specularLight = pow(specularValue, 32) * 5;

            // Apply water colors with phong shading
            finalColor = vec4(0.31, 0.25, 0.71, 0.5) * vec4(vec3(specularLight + diffuseLight + 0.3), 1.0);
        }

        // Get ray interaction with atmospheres
        vec2 atmosphereIntersection = raySphereIntersection(worldCoord.xyz, fragRay, planetPositions[i].xyz, planetPositions[i].w + atmosphereScale);

        float distToAtmosphere = atmosphereIntersection.x;
        float distThroughAtmosphere = min(atmosphereIntersection.y, depth - distToAtmosphere);

        if (distThroughAtmosphere > 0.0) {
            vec3 point = worldCoord.xyz + fragRay * (distToAtmosphere);
            vec3 light = vec3(0.0);

            // Apply glow effect if sun or atmosphere if planet
            if (i == 0) {
                light = shine(point, fragRay, distThroughAtmosphere, finalColor.xyz, planetPositions[i]);
            } else {
                light = scattering(point, fragRay, distThroughAtmosphere, finalColor.xyz, planetPositions[i]);
            }

            finalColor = vec4(light, 0);
        }
    }

    FragColor = finalColor;
}