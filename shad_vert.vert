#version 460 core

in vec2 iVert;
in vec4 iColor;
in vec2 iTexCoord;

out vec4 ioColor;
out vec2 ioTexCoord;

void main() {
  ioColor = iColor;
  ioTexCoord = iTexCoord;
  gl_Position = vec4(iVert * 0.7, 1, 1);
}