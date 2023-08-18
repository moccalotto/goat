#version 460 core

in vec3 iVert;

void main() {
  gl_Position = vec4(iVert, 1);
}