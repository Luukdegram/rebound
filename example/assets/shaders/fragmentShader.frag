#version 410 core
in vec2 pass_textureCoords;
in vec3 surfaceNormal;
in vec3 lightVec;
in vec3 cameraVec;
in float visibility;

out vec4 frag_colour;

uniform sampler2D textureSampler;
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