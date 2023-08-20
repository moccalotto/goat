#version 460 core

in vec3 iVert;
in vec4 iColor;
in vec2 iTexCoord;

out vec4 vColor;
out vec2 vTexCoord;

uniform mat3 uTransformation;

void main() {
  vColor = iColor;
  vTexCoord = iTexCoord;

  gl_Position = vec4(
    uTransformation * iVert,
    1.0
  );
}