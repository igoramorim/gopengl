#version 330 core

in vec2 TexCoords;

uniform sampler2D texture0;

out vec4 FragColor;

float near = 0.1;
float far = 100.0;

float linearizeDepth(float depth) {
	float z = depth * 2.0 - 1.0; // back to NDC
	return (2.0 * near * far) / (far + near - z * (far - near));
}

void main() {
	// FragColor = texture(texture0, TexCoords);
	// FragColor = vec4(vec3(gl_FragCoord.z), 1.0);

	float depth = linearizeDepth(gl_FragCoord.z) / far;
	FragColor = vec4(vec3(depth), 1.0);
}
