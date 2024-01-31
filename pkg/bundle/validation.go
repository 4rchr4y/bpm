package bundle

import validation "github.com/go-ozzo/ozzo-validation/v4"

func (b *Bundle) Validate() error {
	return validation.ValidateStruct(
		b,
		validation.Field(b.BundleFile, validation.Required),
		validation.Field(b.LockFile, validation.Required),
		validation.Field(b.RegoFiles, validation.Required),
	)
}
