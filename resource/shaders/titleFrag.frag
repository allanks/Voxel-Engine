#version 450

uniform sampler2D tex;
uniform vec4 vColor;

in vec2 fragData;

out vec4 outputColor;

void main() {
	vec4 vTexColor = texture2D(tex, fragData);
	outputColor = vTexColor*vColor;
}