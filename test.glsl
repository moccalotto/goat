#version 460 core

out vec4 fragColor;

in vec4 ioColor;

uniform sampler2D uniTexture;

void main() {
   fragColor = mix(texture(uniTexture, TexCoord)) * ioColor, 0.5;
    color = mix(texture(ourTexture0, TexCoord), texture(ourTexture1, TexCoord) * vec4(ourColor, 1.0f), 0.5);
  fragColor = ioColor;
}
