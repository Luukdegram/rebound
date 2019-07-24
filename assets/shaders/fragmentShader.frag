#version 410 core
in vec2 pass_textureCoords;
in vec4 vertexColor
out vec4 frag_colour;
uniform sampler2D textureSampler;
void main() {
    frag_colour = vertexColor;
}