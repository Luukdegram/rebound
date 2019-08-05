#version 410 core
layout (location = 0) in vec3 position;
layout (location = 1) in vec2 textureCoords;
layout (location = 2) in vec3 normal;
layout (location = 3) in vec4 tangent;

out vec2 pass_textureCoords;
out vec3 surfaceNormal;
out vec3 lightVec;
out vec3 cameraVec;
out float visibility;

uniform mat4 transformMatrix;
uniform mat4 projectionMatrix;
uniform mat4 viewMatrix;
uniform vec3 lightPos;

uniform float useFakeLighting;

const float density = 0.07;
const float gradient = 1.5;

void main(void) {
    vec4 worldPos = transformMatrix * vec4(position, 1.0);
    vec4 positionRelativeToCam = viewMatrix * worldPos;
    gl_Position = projectionMatrix * positionRelativeToCam;
    pass_textureCoords = textureCoords;

    vec3 actualNormal = normal;
    if (useFakeLighting > 0.5) {
        actualNormal = vec3(0.0, 1.0, 0.0);
    }

    surfaceNormal = (transformMatrix * vec4(actualNormal, 0.0)).xyz;
    lightVec = lightPos - worldPos.xyz;
    cameraVec = (inverse(viewMatrix) * vec4(0,0,0,1)).xyz - worldPos.xyz;

    float dist = abs(positionRelativeToCam.z);
    visibility = exp(-pow(dist*density, gradient));
    visibility = clamp(visibility, 0.0, 1.0);
}