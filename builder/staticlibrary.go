package builder

type StaticLibrary struct {
	Name    string
	Version string
	Option  string
}

func MakeStaticLibrary(builder *Builder) StaticLibrary {
	name := builder.name()
	if builder.Component == ComponentOpenSSLQuic {
		name = "openssl-quic"
	}
	return StaticLibrary{
		Name:    name,
		Version: builder.Version,
		Option:  builder.option()}
}
