#version 450

uniform sampler2D tex;

in vec2 fragData;
in vec4 sunlight;

out vec4 outputColor;

void main() {
	vec4 vTex = texture(tex, fragData);
	outputColor = vTex*sunlight;
}