#version 410 core
layout (location = 0) in vec3 position;
layout (location = 1) in vec2 textureCoords;

out vec2 pass_textureCoords;
out vec4 vertexColor;

void main() {
    gl_Position = vec4(position, 1.0);
    vertexColor = vec4(0.5, 0.0, 0.0, 1.0);
    pass_textureCoords = textureCoords;
}