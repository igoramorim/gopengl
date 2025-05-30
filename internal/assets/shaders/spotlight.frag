#version 330 core

in vec3 Normal;
in vec3 FragPos;
in vec2 TexCoords;

struct Material {
	sampler2D diffuse;
	sampler2D specular;
	float shininess;
};

struct Light {
	vec3 position;
	vec3 direction;
	float cutOff;
	float outerCutOff;
	vec3 ambient;
	vec3 diffuse;
	vec3 specular;
	float constant;
	float linear;
	float quadratic;
};

uniform Material material;
uniform Light light;
uniform vec3 viewPos;

out vec4 FragColor;

void main() {
	// ambient
	vec3 ambient = light.ambient * vec3(texture(material.diffuse, TexCoords));

	// diffuse
	vec3 norm = normalize(Normal);
	vec3 lightDir = normalize(light.position - FragPos);
	float diff = max(dot(norm, lightDir), 0.0);
	vec3 diffuse = light.diffuse * diff * vec3(texture(material.diffuse, TexCoords));

	// specular
	vec3 viewDir = normalize(viewPos - FragPos);
	vec3 reflectDir = reflect(-lightDir, norm);
	float spec = pow(max(dot(viewDir, reflectDir), 0.0), material.shininess);
	vec3 specTexel = texture(material.specular, TexCoords).rgb;
	vec3 specular = light.specular * spec * specTexel;

	// spotlight
	float theta = dot(lightDir, normalize(-light.direction));
	float epsilon = (light.cutOff - light.outerCutOff);
	// Clamp is used to make sure the value is not negative. A negative value for example would make
	// everything outside the spotlight totally dark
	float intensity = clamp((theta - light.outerCutOff) / epsilon, 0.0, 1.0);
	diffuse *= intensity;
	specular *= intensity;

	// attenuation
	float distance = length(light.position - FragPos);
	float attenuation = 1.0 / (light.constant + light.linear * distance + light.quadratic * (distance * distance));
	ambient *= attenuation;
	diffuse *= attenuation;
	specular *= attenuation;

	vec3 result = ambient + diffuse + specular;
	FragColor = vec4(result, 1.0);
}
