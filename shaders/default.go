package shaders

const (
	//FragmentShader is a default fragmentation shader.
	//It supports lights, skycolour, reflectivity, shinedamper.
	FragmentShader = `
	#version 410 core
	in vec2 pass_textureCoords;
	in vec3 surfaceNormal;
	in vec3 lightVec;
	in vec3 cameraVec;
	in float visibility;
	
	out vec4 frag_colour;
	
	uniform sampler2D textureSampler;
	uniform sampler2D backgroundTexture;
	uniform vec3 lightColour;
	uniform float shineDamper;
	uniform float reflectivity;
	uniform vec3 skyColour;
	
	void main() {
		vec3 unitNormal = normalize(surfaceNormal);
		vec3 unitLightVector = normalize(lightVec);
	
		float nDot = dot(unitNormal, unitLightVector);
		float brightness = max(nDot, 0.2);
		vec3 diffuse = brightness * lightColour;
	
		vec3 unitCameraVector = normalize(cameraVec);
		vec3 lightDirection = -unitLightVector;
		vec3 reflectedLightDirection = reflect(lightDirection, unitNormal);
	
		float specularFactor = dot(reflectedLightDirection, unitCameraVector);
		specularFactor = max(specularFactor, 0.0);
		float dampedFactor = pow(specularFactor, shineDamper);
		vec3 finalSpecular = dampedFactor * reflectivity * lightColour;
	
		vec4 textureColor = texture(textureSampler, pass_textureCoords);
		if (textureColor.a < 0.5) {
			discard;
		}
	
		frag_colour = vec4(diffuse, 1.0) * textureColor + vec4(finalSpecular, 1.0);
		frag_colour = mix(vec4(skyColour, 1.0), frag_colour, visibility);
	}
	` + "\x00"

	//VertexShader is a default vertex shader.
	//It allows visibility turn on/off, a camera and fake lighting
	VertexShader = `
	#version 410 core
	layout (location = 0) in vec3 position;
	layout (location = 1) in vec2 textureCoords;
	layout (location = 2) in vec2 textureCoords2;
	layout (location = 3) in vec3 normal;
	layout (location = 4) in vec4 tangents;
	layout (location = 5) in vec4 color;
	layout (location = 6) in vec4 joints;
	layout (location = 7) in vec4 weights;
	
	out vec2 pass_textureCoords;
	out vec3 surfaceNormal;
	out vec3 lightVec;
	out vec3 cameraVec;
	out float visibility;
	
	uniform mat4 transformMatrix;
	uniform mat4 projectionMatrix;
	uniform mat4 viewMatrix;
	uniform vec3 lightPos;
	
	uniform float useFakeLighting;
	
	const float density = 0.07;
	const float gradient = 1.5;
	
	void main(void) {
		vec4 worldPos = transformMatrix * vec4(position, 1.0);
		vec4 positionRelativeToCam = viewMatrix * worldPos;
		gl_Position = projectionMatrix * positionRelativeToCam;
		pass_textureCoords = textureCoords;
	
		vec3 actualNormal = normal;
		if (useFakeLighting > 0.5) {
			actualNormal = vec3(0.0, 1.0, 0.0);
		}
	
		surfaceNormal = (transformMatrix * vec4(actualNormal, 0.0)).xyz;
		lightVec = lightPos - worldPos.xyz;
		cameraVec = (inverse(viewMatrix) * vec4(0,0,0,1)).xyz - worldPos.xyz;
	
		float dist = abs(positionRelativeToCam.z);
		visibility = exp(-pow(dist*density, gradient));
		visibility = clamp(visibility, 0.0, 1.0);
	}` + "\x00"
)
