package errors

const (
	Conflict   Class = "conflict"          // Action cannot be performed
	Internal   Class = "internal"          // Internal error
	Invalid    Class = "invalid"           // Validation failed
	NotFound   Class = "not_found"         // Entity does not exist
	PermDenied Class = "permission_denied" // Does not have permission
	Other      Class = "other"             // Unclassified error
)
