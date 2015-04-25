#version 450

uniform sampler2D tex;

in vec2 fragData;

out vec4 outputColor;

void main() {
	outputColor = texture(tex, fragData);
}