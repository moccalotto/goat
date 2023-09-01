#version 460 core

in vec3 iVert;
in vec2 iTexCoord;

out vec2 vTexCoord;

uniform mat3 uniTransformation;

void main() {
  vTexCoord = iTexCoord;

  gl_Position = vec4(
    uniTransformation * iVert,
    1.0
  );
}
