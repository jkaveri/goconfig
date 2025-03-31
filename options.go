package goconfig

// Option is a function type that modifies a Loader's configuration.
type Option func(*Loader)

// WithPrefix sets a prefix for all environment variables.
// The name of the environment variable will be prefix_fieldName.
// For example, with prefix "APP" and field "Host", the environment variable would be "APP_HOST".
func WithPrefix(prefix string) Option {
	return func(c *Loader) {
		c.prefix = prefix
	}
}

// WithSeparator sets the separator used for nested field names in environment variables.
// The default separator is "_". For example, with separator "." and nested field "DB.Host",
// the environment variable would be "DB.HOST".
func WithSeparator(sep string) Option {
	return func(c *Loader) {
		c.sep = sep
	}
}

// WithArraySeparator sets the separator used for array values in environment variables.
// The default separator is ",". For example, with separator ";" and array field "Numbers",
// the environment variable would be "1;2;3;4".
func WithArraySeparator(sep string) Option {
	return func(c *Loader) {
		c.arraySep = sep
	}
}

// WithKeyTransformer sets a custom function to transform field names into environment variable names.
// This allows for custom naming conventions beyond the default behavior.
// The transformer function receives the field name and returns the transformed name.
func WithKeyTransformer(transformer func(key string) string) Option {
	return func(c *Loader) {
		c.fieldNameTransformer = transformer
	}
}
