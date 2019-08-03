#version 410 core
layout (location = 0) in vec3 position;
layout (location = 1) in vec2 textureCoords;
layout (location = 2) in vec3 normal;
layout (location = 3) in vec4 tangent;

out vec2 pass_textureCoords;
out vec3 surfaceNormal;
out vec3 lightVec;
out vec3 cameraVec;

uniform mat4 transformMatrix;
uniform mat4 projectionMatrix;
uniform mat4 viewMatrix;
uniform vec3 lightPos;

void main() {
    vec4 worldPos = transformMatrix * vec4(position, 1.0);
    gl_Position = projectionMatrix * viewMatrix * worldPos;
    pass_textureCoords = textureCoords;

    surfaceNormal = (transformMatrix * vec4(normal, 0.0)).xyz;
    lightVec = lightPos - worldPos.xyz;
    cameraVec = (inverse(viewMatrix) * vec4(0,0,0,1)).xyz - worldPos.xyz;
}