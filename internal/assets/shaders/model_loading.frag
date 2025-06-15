#version 330 core

struct Light {
	vec3 position;
	vec3 ambient;
	vec3 diffuse;
	vec3 specular;
	float constant;
	float linear;
	float quadratic;
};

in vec3 Normal;
in vec3 FragPos;
in vec2 TexCoords;

uniform sampler2D texture_diffuse1;
uniform sampler2D texture_specular1;

uniform Light light;
uniform vec3 viewPos;

out vec4 FragColor;

void main() {
	// ambient
	vec3 ambient = light.ambient * vec3(texture(texture_diffuse1, TexCoords));

	// diffuse
	vec3 norm = normalize(Normal);
	vec3 lightDir = normalize(light.position - FragPos);
	float diff = max(dot(norm, lightDir), 0.0);
	vec3 diffuse = light.diffuse * diff * vec3(texture(texture_diffuse1, TexCoords));

	// specular
	vec3 viewDir = normalize(viewPos - FragPos);
	vec3 reflectDir = reflect(-lightDir, norm);
	float shine = 32.0;
	float spec = pow(max(dot(viewDir, reflectDir), 0.0), shine);
	vec3 specTexel = texture(texture_specular1, TexCoords).rgb;
	vec3 specular = light.specular * spec * specTexel;

	// attenuation
	float distance = length(light.position - FragPos);
	float attenuation = 1.0 / (light.constant + light.linear * distance + light.quadratic * (distance * distance));

	ambient *= attenuation;
	diffuse *= attenuation;
	specular *= attenuation;

	vec3 result = ambient + diffuse + specular;
	FragColor = vec4(result, 1.0);

	// FragColor = texture(texture_specular1, TexCoords);
	// FragColor = vec4(0.0, TexCoords.x, TexCoords.y, 1.0);
}
