#version 330 core

layout (location = 0) in vec3 position;
layout (location = 1) in vec3 normal;
layout (location = 2) in vec2 texCoords;

uniform mat4 model;
uniform mat4 view;
uniform mat4 projection;

out vec3 FragPos;
out vec3 Normal;
out vec2 TexCoords;

void main() {
	FragPos = vec3(model * vec4(position, 1.0));

	// NOTE: Inversing matrices is a costly operation for shaders.
	// It should be done in the CPU.
	Normal = mat3(transpose(inverse(model))) * normal;

	// Note that we read the multiplication from right to left
	gl_Position = projection * view * vec4(FragPos, 1.0);

	TexCoords = texCoords;
}
