#version 330 core

in vec2 TexCoords;

uniform sampler2D texture0;

out vec4 FragColor;

void main() {
	FragColor = texture(texture0, TexCoords);
}
