#version 460 core

uniform float u_rotAngle;
// in vec3 gl_FragCoord;
out vec4 fragColor;

void main() {
  fragColor = vec4(vec3(0.99 * (sin(u_rotAngle*3) * 0.5 + 0.5),
                    0.25 * (cos(u_rotAngle + gl_FragCoord.y) * 0.5 + 0.5),
                    0.99 * (sin(u_rotAngle + gl_FragCoord.x) * 0.5 + 0.5)),
               1);
}