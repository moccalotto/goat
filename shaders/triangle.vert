#version 460 core

in vec3 iVert;

layout (location=0) uniform float scale;

void main() {
  gl_Position = vec4(iVert * scale, 1);
}