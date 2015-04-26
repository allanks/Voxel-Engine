#version 450

uniform sampler2D tex;

in vec2 fragData;
smooth in vec3 vNormal;

out vec4 outputColor;

struct SimpleLight
{
	vec3 vColor;
	vec3 vDirection;
	float intensity;
};

uniform SimpleLight sun;

void main() {
	vec4 vTex = texture(tex, fragData);
	float diffuse = max(0.0, dot(normalize(vNormal), -sun.vDirection));
	outputColor = vTex*vec4(sun.vColor*(sun.intensity+diffuse),1.0);
}