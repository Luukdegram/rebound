package rebound

import (
	"fmt"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/luukdegram/rebound/internal/thread"
)

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
	
		float dist = length(positionRelativeToCam.xyz);
		visibility = exp(-pow((dist*density), gradient));
		visibility = clamp(visibility, 0.0, 1.0);
	}` + "\x00"

	newVShader = `
	#version 410 core
	layout (location = 0) in vec3 position;
	layout (location = 1) in vec2 textureCoords;
	layout (location = 2) in vec2 textureCoords2;
	layout (location = 3) in vec3 normal;
	layout (location = 4) in vec4 tangents;
	layout (location = 5) in vec4 color;
	layout (location = 6) in vec4 joints;
	layout (location = 7) in vec4 weights;

	out vec3 FragPos;  
	out vec3 Normal;
	out vec2 TexCoords;
	
	uniform mat4 model;
	uniform mat4 projection;
	uniform mat4 view;
	
	void main(void) {
		// Calculate fragment position
		FragPos = vec3(model * vec4(position, 1.0));
		Normal = normal;

		// Pass our texture coords
		TexCoords = textureCoords;
		
		// calculate the vector position from the 3D world to 2D view
		gl_Position = projection * view * vec4(FragPos, 1.0);
	}` + "\x00"

	newFShader = `
	#version 410 core
	out vec4 FragColor;

	in vec2 TexCoords;
	in vec3 Normal;
	in vec3 FragPos;

	struct Material {
		sampler2D diffuse;
		vec3 specular;
		float shininess;
	};

	struct DirLight {
		vec3 direction;
		vec3 ambient;
		vec3 diffuse;
		vec3 specular;
	};

	struct PointLight {
		vec3 position;
    
		float constant;
		float linear;
		float quadratic;  

		vec3 ambient;
		vec3 diffuse;
		vec3 specular;
	};
	
	uniform vec3 lightColour;
	uniform vec3 viewPos;
	uniform Material material;
	uniform DirLight light;
	uniform int amountLights;
	
	#define NR_POINT_LIGHTS 4  
	uniform PointLight pointLights[NR_POINT_LIGHTS];

	vec3 CalcDirLight(DirLight light, vec3 normal, vec3 viewDir);
	vec3 CalcPointLight(PointLight light, vec3 normal, vec3 fragPos, vec3 viewDir);

	void main()
	{
		vec4 texColour = texture(material.diffuse, TexCoords);
		if (texColour.a < 0.5) {
			discard;
		}

		//Calculate the normals
		vec3 norm = normalize(Normal);
		
		// Calculate view direction
		vec3 viewDir = normalize(viewPos - FragPos);
		
		// Calculate directional light 
		vec3 light = CalcDirLight(light, norm, viewDir);

		int size = amountLights;
		if (size > 0) 
		{
			if (size > NR_POINT_LIGHTS)
				size = NR_POINT_LIGHTS;
			// Calculate point lights
			for(int i = 0; i < size; i++) {
				light += CalcPointLight(pointLights[i], norm, FragPos, viewDir);
			}
		}

		// Set the final result pixel
		FragColor = vec4(light, 1.0);
	}

	vec3 CalcDirLight(DirLight light, vec3 normal, vec3 viewDir)
	{
		vec3 lightDir = normalize(-light.direction);
		// diffuse shading
		float diff = max(dot(normal, lightDir), 0.0);
		// specular shading
		vec3 reflectDir = reflect(-lightDir, normal);
		float spec = pow(max(dot(viewDir, reflectDir), 0.0), material.shininess);
		// combine results
		vec3 ambient  = light.ambient  * vec3(texture(material.diffuse, TexCoords));
		vec3 diffuse  = light.diffuse  * diff * vec3(texture(material.diffuse, TexCoords));
		vec3 specular = light.specular * (spec * material.specular);
		return (ambient + diffuse + specular);
	}

	vec3 CalcPointLight(PointLight light, vec3 normal, vec3 fragPos, vec3 viewDir)
	{
		vec3 lightDir = normalize(light.position - fragPos);
		// diffuse shading
		float diff = max(dot(normal, lightDir), 0.0);
		// specular shading
		vec3 reflectDir = reflect(-lightDir, normal);
		float spec = pow(max(dot(viewDir, reflectDir), 0.0), material.shininess);
		// attenuation
		float distance    = length(light.position - fragPos);
		float attenuation = 1.0 / (light.constant + light.linear * distance + 
					light.quadratic * (distance * distance));    
		// combine results
		vec3 ambient  = light.ambient  * vec3(texture(material.diffuse, TexCoords));
		vec3 diffuse  = light.diffuse  * diff * vec3(texture(material.diffuse, TexCoords));
		vec3 specular = light.specular * (spec * material.specular);
		ambient  *= attenuation;
		diffuse  *= attenuation;
		specular *= attenuation;
		return (ambient + diffuse + specular);
	}
	` + "\x00"
)

var shaderIds []uint32

// Shader contains the logic to render a shader
type Shader interface {
	// Setup runs at the beginning of the renderer's update() function, before any entities are being rendered.
	// The camera is provided to retreive its view and projection matrixes
	Setup(Camera)
	// Render runs while rendering each entity, the corresponding RenderComponent is provided in the Render runction.
	// Within this function you can set entity specific shader options
	Render(RenderComponent)
	// ID returns the shader's ID, this is generated by the NewShader function
	ID() uint32
}

// BasicShader is the default shader part of the Rebound engine.
type BasicShader struct {
	id          uint32
	SceneLight  *Light
	PointLights []PointLight
}

//NewShader returns a new ShaderComponent by compiling the given vertexShader and fragmentShader
//Returns an error if any of the shaders could not be compiled
func NewShader(vertexShader, fragmentShader string) (id uint32, err error) {
	err = thread.CallErr(func() error {
		vID, err := compileShader(vertexShader, gl.VERTEX_SHADER)
		if err != nil {
			return err
		}

		shaderIds = append(shaderIds, vID)

		fID, err := compileShader(fragmentShader+"\x00", gl.FRAGMENT_SHADER)
		if err != nil {
			return err
		}

		shaderIds = append(shaderIds, fID)

		id = gl.CreateProgram()
		gl.AttachShader(id, vID)
		gl.AttachShader(id, fID)

		gl.LinkProgram(id)
		gl.ValidateProgram(id)

		gl.DetachShader(id, vID)
		gl.DetachShader(id, fID)

		return nil
	})

	return
}

// NewBasicShader creates a default shader, provided by the Rebound engine
func NewBasicShader() (*BasicShader, error) {
	id, err := NewShader(newVShader, newFShader)
	if err != nil {
		return nil, err
	}

	bs := &BasicShader{
		id: id,
		SceneLight: &Light{
			Direction: [3]float32{-0.2, -1.0, -0.3},
		},
	}

	return bs, nil
}

// ID returns the shader id generated by opengl
func (bs *BasicShader) ID() uint32 {
	return bs.id
}

// Setup loads variables into the shader pre-entity rendering
func (bs *BasicShader) Setup(c Camera) {
	if bs.SceneLight != nil {
		LoadVec3(bs, "light.ambient", [3]float32{0.2, 0.2, 0.2})
		LoadVec3(bs, "light.diffuse", [3]float32{0.5, 0.5, 0.5})
		LoadVec3(bs, "light.specular", [3]float32{1, 1, 1})
		LoadVec3(bs, "light.direction", bs.SceneLight.Direction)
	}

	LoadInt(bs, "amountLights", len(bs.PointLights))

	for index, l := range bs.PointLights {
		prefix := fmt.Sprintf("pointLights[%v].", index)
		LoadVec3(bs, prefix+"position", l.Position)
		LoadVec3(bs, prefix+"ambient", l.Ambient)
		LoadVec3(bs, prefix+"diffuse", l.Diffuse)
		LoadVec3(bs, prefix+"specular", l.Specular)
		LoadFloat(bs, prefix+"constant", l.Constant)
		LoadFloat(bs, prefix+"linear", l.Linear)
		LoadFloat(bs, prefix+"quadratic", l.Quadratic)
	}

	LoadVec3(bs, "viewPos", c.Position)
	LoadMat(bs, "projection", c.ProjectMat)
	LoadMat(bs, "view", NewViewMatrix(c))
}

// Render loads variables into the shader based on current RenderComponent
func (bs *BasicShader) Render(rc RenderComponent) {
	// Set material
	LoadVec3(bs, "material.specular", [3]float32{0.5, 0.5, 0.5})
	LoadFloat(bs, "material.shininess", 64)

	tmMat := NewTransformationMatrix(rc.Position, rc.Rotation, rc.Scale)
	LoadMat(bs, "model", tmMat)
}

//GetUniformLocation returns the location of the uniform given, returning the OpenGL id as an int32
func GetUniformLocation(s Shader, name string) int32 {
	return gl.GetUniformLocation(s.ID(), gl.Str(name+"\x00"))
}

//LoadFloat loads a uniform float into the shader
func LoadFloat(s Shader, name string, value float32) {
	loc := GetUniformLocation(s, name)
	gl.Uniform1f(loc, value)
}

// LoadInt loads an integer into the shader
func LoadInt(s Shader, name string, value int) {
	loc := GetUniformLocation(s, name)
	gl.Uniform1i(loc, int32(value))
}

//LoadVec3 loads a uniform Vector into the shader
func LoadVec3(s Shader, name string, value [3]float32) {
	loc := GetUniformLocation(s, name)
	gl.Uniform3f(loc, value[0], value[1], value[2])
}

// LoadVec4 loads a uniform Vector with 4 elements into the shader
func LoadVec4(s Shader, name string, value [4]float32) {
	loc := GetUniformLocation(s, name)
	gl.Uniform4f(loc, value[0], value[1], value[2], value[3])
}

//LoadBool loads a boolean into the shader
func LoadBool(s Shader, name string, value bool) {
	loc := GetUniformLocation(s, name)
	var float float32
	if value {
		float = 1
	}
	gl.Uniform1f(loc, float)
}

//LoadMat loads a matrix into the shader
func LoadMat(s Shader, name string, value [16]float32) {
	loc := GetUniformLocation(s, name)
	gl.UniformMatrix4fv(loc, 1, false, &value[0])
}

//startHader starts the shader program
func startHader(s Shader) {
	gl.UseProgram(s.ID())
}

//stopShader stops the current shader program
func stopShader() {
	gl.UseProgram(0)
}

//CleanUpShaders deletes the program
func CleanUpShaders() {
	thread.Call(func() {
		for _, id := range shaderIds {
			gl.DeleteProgram(id)
		}

		for _, id := range shaderIds {
			gl.DeleteShader(id)
		}
	})
}

func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}
	return shader, nil
}
