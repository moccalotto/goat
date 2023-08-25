#version 460 core

out vec4 fragColor;

// in vec4 vColor;
in vec2 vTexCoord;

uniform sampler2D uniTexture;
uniform float uniColorMix;
uniform vec4 uniColor;
uniform vec4 uniSubTexPos;

void main() {
  vec2 tmp = vec2(
      mix(uniSubTexPos.x, uniSubTexPos.z, vTexCoord.x),		// mix == lerp
      mix(uniSubTexPos.y, uniSubTexPos.w, vTexCoord.y)		// mix == lerp
      );
  fragColor = mix(texture(uniTexture, tmp), uniColor, uniColorMix);
}
