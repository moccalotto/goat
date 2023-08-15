#version 460 core

in vec2 in_vert;

// u_scale contains:    [ widthFactor,  heightFactor ]
uniform vec2 u_scale;

// // u_rotation contains: [ sin(R) , cos(R) ]  where R is the rotation angle
// uniform vec2 u_rotation;

uniform float u_rotAngle;

void main() {
  vec2 u_rotation = vec2(sin(u_rotAngle), cos(u_rotAngle));
  vec2 in_vert = in_vert * u_scale;
  vec2 v_rotated = vec2(in_vert.x * u_rotation.y + in_vert.y * u_rotation.x,
                        in_vert.y * u_rotation.y - in_vert.x * u_rotation.x);

  gl_Position = vec4(v_rotated, 0.0, 1.0);
}