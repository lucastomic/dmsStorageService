package contextypes

// CTXRequestIDKey is a type used exclusively as a key in context.Context for storing and retrieving
// the unique request ID associated with an HTTP request. This key ensures type safety and reduces
// the likelihood of key collisions when using context values.
type CTXRequestIDKey struct{}

// ContextPathVarKey is a type used as a context key for storing and retrieving path variables.
type ContextPathVarKey string
