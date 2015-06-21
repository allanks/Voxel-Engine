package Graphics

import "github.com/go-gl/mathgl/mgl32"

type BufferBinder interface {
	BindBuffers(...[]float32)
}

type UniformBinder interface {
	BindUniforms(...[]float32)
}

type TextureBinder interface {
	BindTexture(uint32)
}

type InstanceRenderer interface {
	RenderInstances([]float32, int32)
}

type BufferCreator interface {
	CreateBuffers()
}

type UniformCreator interface {
	CreateUniforms()
}

type FragmentBinder interface {
	BindFragData()
}

type ProgramStarter interface {
	StartPrograms()
}

type ProjectionInitializer interface {
	BindProjection(float32, float32)
}

type ProjectionUpdater interface {
	UpdateProjection(states mgl32.Mat4)
}

type ProjectionController interface {
	ProjectionInitializer
	ProjectionUpdater
}

type OpenGLController interface {
	BufferCreator
	UniformCreator
	FragmentBinder
	BufferBinder
	UniformBinder
	TextureBinder
	InstanceRenderer
	ProgramStarter
	ProjectionController
}

type OpenGLInitializer interface {
	Init()
}

type ProgramCreator interface {
	NewProgram(string, string) uint32
}

type TextureCreator interface {
	CreateTexture(string) uint32
}

type OpenGLDepthToggle interface {
	DepthToggle(bool)
}

type OpenGLClear interface {
	Clear()
}

type OpenGLControl interface {
	OpenGLInitializer
	ProgramCreator
	TextureCreator
	OpenGLDepthToggle
	OpenGLClear
}

type OpenGLControlCreator interface {
	CreateControl() *OpenGLControl
}

type OpenGLControllerCreator interface {
	CreateControllers() []*OpenGLController
}

type OpenGLStarter interface {
	OpenGLControlCreator
	OpenGLControllerCreator
}
