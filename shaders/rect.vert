#version 460 core

in vec3 iVert;

uniform mat3 uniTransformation;

void main() {
  gl_Position = vec4(uniTransformation * iVert, 1.0);
}
