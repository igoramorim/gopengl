#version 330 core

uniform sampler2D texture0;
uniform sampler2D texture1;

in vec2 TexCoord;

out vec4 FragColor;

void main() {
	// FragColor = texture(texture0, TexCoord);
	FragColor = mix(texture(texture0, TexCoord), texture(texture1, TexCoord), 0.2);
}
