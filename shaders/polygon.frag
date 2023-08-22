#version 460 core

out vec4 fragColor;

// in vec4 vColor;
in vec2 vTexCoord;

uniform sampler2D uniTexture;
uniform float uniColorMix;
uniform vec4 uniColor;
uniform bool uniWireframe;

void main() {

  if (uniWireframe) {
    fragColor = vec4(1, 1, 1, 1);
    return;
  }

  if (uniColorMix > 0.999) {
    fragColor = uniColor;
    return;
  }

  if (uniColorMix < 0.001) {
    fragColor = texture(uniTexture, vTexCoord);
    return;
  }

  fragColor = mix(texture(uniTexture, vTexCoord), uniColor, uniColorMix);
}