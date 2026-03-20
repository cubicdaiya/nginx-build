package builder

type StaticLibrary struct {
	Name                 string
	Version              string
	Option               string
	EnablesHTTPSSLModule bool
}

func MakeStaticLibrary(builder *Builder) StaticLibrary {
	enablesHTTPSSLModule := false
	switch builder.Component {
	case ComponentOpenSSL, ComponentLibreSSL, ComponentCustomSSL:
		enablesHTTPSSLModule = true
	}

	return StaticLibrary{
		Name:                 builder.name(),
		Version:              builder.Version,
		Option:               builder.option(),
		EnablesHTTPSSLModule: enablesHTTPSSLModule,
	}
}
