#version 460 core

out vec4 fragColor;

in vec4 ioColor;
in vec2 ioTexCoord;

uniform sampler2D uniTexture;

void main() {
  fragColor = mix(texture(uniTexture, ioTexCoord), ioColor, 0.6);
  /*
  fragColor = vec4(
    0.5 + 0.5 * sin(ioTexCoord.y),
    0.5 + 0.5 * sin(ioTexCoord.y),
    // 0.5 + 0.5 * sin(length(ioColor.rgb)),
    0,
    ioColor.a
  );
  */
}