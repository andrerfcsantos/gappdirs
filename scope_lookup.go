package gappdirs

// ScopeForPath returns the matching scope for path for app name and a bool indicating if a match was found.
//
// Directory comparison and lookup checks path candidates in local, then user, then system precedence.
// In other words, it assumes local scope for directory comparison.
//
// The path does not need to exist on disk.
func ScopeForPath(appName string, path string) (scope Scope, found bool) {
	ctx := topLevelScopedContext(appName, ScopeLocal)
	return scopedScopeForPath(ctx, path)
}

// ScopeForPath returns the matching scope for path for the Resolver app name and scope.
//
// It checks only the scope layers available for the resolver:
//   - local resolvers check local, user, then system
//   - user resolvers check user then system
//   - system resolvers check system only
//
// The path does not need to exist on disk.
func (r *Resolver) ScopeForPath(path string) (scope Scope, found bool) {
	return scopedScopeForPath(r.ctx, path)
}
