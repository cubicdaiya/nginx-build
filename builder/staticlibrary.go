package builder

type StaticLibrary struct {
	Name    string
	Version string
	Option  string
}

func MakeStaticLibrary(builder *Builder) StaticLibrary {
	return StaticLibrary{
		Name:    builder.name(),
		Version: builder.Version,
		Option:  builder.option()}
}
