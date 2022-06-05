package configure

type Options struct {
	Values map[string]OptionValue
	Bools  map[string]OptionBool
}

type OptionValue struct {
	Name  string
	Desc  string
	Value *string
}

type OptionBool struct {
	Name    string
	Desc    string
	Enabled *bool
}

func MakeArgsBool() map[string]OptionBool {
	return make(map[string]OptionBool)
}

func MakeArgsString() map[string]OptionValue {
	return make(map[string]OptionValue)
}
