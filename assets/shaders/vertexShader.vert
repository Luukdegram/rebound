#version 410
in vec3 vp;

out vec3 colour;

void main() {
    gl_Position = vec4(vp, 1.0);
    colour = vec3(vp.x+0.5, 1.0, vp.y+0.5);
}