package data

// ReleaseContext returns a newly allocated call frame after its result was
// extracted. Escaped closure, generator and reference frames stay alive;
// concrete pooled contexts enforce that lifetime rule themselves.
func ReleaseContext(ctx Context) {
	if pooled, ok := ctx.(interface{ ReleasePooled() }); ok {
		pooled.ReleasePooled()
	}
}
