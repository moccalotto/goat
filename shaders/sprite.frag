#version 460 core

out vec4 fragColor;

// in vec4 vColor;
in vec2 vTexCoord;

uniform sampler2D uniTexture; // the texture to use
uniform vec4 uniColor;        // the colot to use
uniform float uniColorMix;    // how much of the output color comes from uniColor
uniform vec4 uniSubTexPos;    // which part of the texture do we want to use

void main() {
  vec2 tmp =
      vec2(mix(uniSubTexPos.x, uniSubTexPos.z, vTexCoord.x), // mix == lerp
           mix(uniSubTexPos.y, uniSubTexPos.w, vTexCoord.y)  // mix == lerp
      );
  fragColor = mix(texture(uniTexture, tmp), uniColor, uniColorMix);
}
