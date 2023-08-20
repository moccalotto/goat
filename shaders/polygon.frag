#version 460 core

out vec4 fragColor;

in vec4 vColor;
in vec2 vTexCoord;

uniform sampler2D uniTexture;

void main() {
  fragColor = mix(texture(uniTexture, vTexCoord), vColor, 0.59);
}