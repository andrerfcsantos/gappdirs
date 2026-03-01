package gappdirs

import (
	"io/fs"
	"path/filepath"
	"reflect"
	"testing"
)

func TestScopedConstructorsApplyOptions(t *testing.T) {
	t.Run("missing app name defaults", func(t *testing.T) {
		r := NewUserResolver("")
		if r.ctx.appName != "unnamed_app" {
			t.Fatalf("default app name mismatch: want %q, got %q", "unnamed_app", r.ctx.appName)
		}
	})

	t.Run("sanitizes app name", func(t *testing.T) {
		wd := t.TempDir()
		r := newScopedResolver(
			"A/B Demo",
			ScopeLocal,
			[]ResolverOption{WithLocalDir(wd)},
			func(string, category) ([]string, error) { return []string{filepath.Join(wd, "u")}, nil },
			func(string, category) ([]string, error) { return []string{filepath.Join(wd, "s")}, nil },
		)
		if r.ctx.appName != "A_B_Demo" {
			t.Fatalf("sanitized app name mismatch: want %q, got %q", "A_B_Demo", r.ctx.appName)
		}

		got := r.DataDir()
		want := filepath.Join(wd, ".A_B_Demo", "data")
		if got != want {
			t.Fatalf("local data dir mismatch: want %q, got %q", want, got)
		}
	})

	t.Run("nil option is ignored", func(t *testing.T) {
		r := NewUserResolver("demo", nil)
		if r.ctx.scope != ScopeUser {
			t.Fatalf("scope mismatch: want %s, got %s", ScopeUser, r.ctx.scope)
		}
	})

	t.Run("default permissions are applied", func(t *testing.T) {
		wd := t.TempDir()
		r := newScopedResolver(
			"demo",
			ScopeSystem,
			[]ResolverOption{WithLocalDir(wd)},
			func(string, category) ([]string, error) { return []string{filepath.Join(wd, "u")}, nil },
			func(string, category) ([]string, error) { return []string{filepath.Join(wd, "s")}, nil },
		)
		if r.ctx.scope != ScopeSystem {
			t.Fatalf("scope mismatch: want %s, got %s", ScopeSystem, r.ctx.scope)
		}
		if r.ctx.defaultDirPerm != DefaultDirPerm {
			t.Fatalf("default dir permission mismatch: want %o, got %o", DefaultDirPerm, r.ctx.defaultDirPerm)
		}
	})

	t.Run("working dir options append in order and dedupe", func(t *testing.T) {
		wdA := t.TempDir()
		wdB := t.TempDir()
		wdC := t.TempDir()
		r := newScopedResolver(
			"demo",
			ScopeLocal,
			[]ResolverOption{
				WithLocalDir(wdA),
				WithLocalDirs("", wdB, wdA, "  ", wdC, wdB),
				WithDefaultDirPerm(0o700),
				WithDefaultDirPerm(0o750),
			},
			func(string, category) ([]string, error) { return []string{filepath.Join(wdB, "u")}, nil },
			func(string, category) ([]string, error) { return []string{filepath.Join(wdB, "s")}, nil },
		)
		wantWorkingDirs := []string{wdA, wdB, wdC}
		if !reflect.DeepEqual(r.ctx.workingDirs, wantWorkingDirs) {
			t.Fatalf("working dirs mismatch:\nwant: %#v\ngot:  %#v", wantWorkingDirs, r.ctx.workingDirs)
		}
		if r.ctx.defaultDirPerm != 0o750 {
			t.Fatalf("permission mismatch: want %o, got %o", 0o750, r.ctx.defaultDirPerm)
		}
	})
}

func TestScopedConstructorsSetExpectedScope(t *testing.T) {
	type scopedCtor func(string, ...ResolverOption) *Resolver

	for _, tc := range []struct {
		name      string
		ctor      scopedCtor
		wantScope Scope
	}{
		{name: "user", ctor: NewUserResolver, wantScope: ScopeUser},
		{name: "system", ctor: NewSystemResolver, wantScope: ScopeSystem},
		{name: "local", ctor: NewLocalResolver, wantScope: ScopeLocal},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r := tc.ctor("demo", WithLocalDir(t.TempDir()))
			if r.ctx.scope != tc.wantScope {
				t.Fatalf("scope mismatch: want %s, got %s", tc.wantScope, r.ctx.scope)
			}
		})
	}
}

func TestSanitizeAppName(t *testing.T) {
	for _, tc := range []struct {
		name  string
		input string
		want  string
	}{
		{name: "mixed case and spaces", input: "My App", want: "My_App"},
		{name: "single dot", input: ".", want: "."},
		{name: "double dots", input: "..", want: ".."},
		{name: "forward slash", input: "my/app", want: "my_app"},
		{name: "backslash", input: `my\app`, want: "my_app"},
		{name: "preserves repeated underscores", input: "My   App", want: "My___App"},
		{name: "empty defaults", input: "   ", want: "unnamed_app"},
		{name: "trim outer spaces", input: "  Demo App  ", want: "Demo_App"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := sanitizeAppName(tc.input)
			if got != tc.want {
				t.Fatalf("sanitized name mismatch: want %q, got %q", tc.want, got)
			}
		})
	}
}

func TestOptionFallbacks(t *testing.T) {
	wd := t.TempDir()
	userFn := func(string, category) ([]string, error) { return []string{filepath.Join(wd, "u")}, nil }
	systemFn := func(string, category) ([]string, error) { return []string{filepath.Join(wd, "s")}, nil }

	t.Run("empty working dir option is ignored", func(t *testing.T) {
		r := newScopedResolver("demo", ScopeUser, []ResolverOption{WithLocalDir("  ")}, userFn, systemFn)
		if len(r.ctx.workingDirs) != 0 {
			t.Fatalf("working dirs should remain empty when option is ignored, got %#v", r.ctx.workingDirs)
		}
	})

	t.Run("invalid permission falls back to default", func(t *testing.T) {
		r := newScopedResolver("demo", ScopeUser, []ResolverOption{WithDefaultDirPerm(0), WithLocalDir(wd)}, userFn, systemFn)
		if r.ctx.defaultDirPerm != DefaultDirPerm {
			t.Fatalf("permission fallback mismatch: want %o, got %o", DefaultDirPerm, r.ctx.defaultDirPerm)
		}

		r = newScopedResolver("demo", ScopeUser, []ResolverOption{WithDefaultDirPerm(fs.ModeDir | 0o755), WithLocalDir(wd)}, userFn, systemFn)
		if r.ctx.defaultDirPerm != DefaultDirPerm {
			t.Fatalf("permission fallback mismatch for mode bits: want %o, got %o", DefaultDirPerm, r.ctx.defaultDirPerm)
		}
	})
}

func TestDirsScopeOrdering(t *testing.T) {
	wd := t.TempDir()
	userData := filepath.Join(wd, "user", "data")
	systemData := filepath.Join(wd, "system", "data")
	userConfig := filepath.Join(wd, "user", "config")
	systemConfig := filepath.Join(wd, "system", "config")

	user := map[category][]string{
		categoryData:   {userData},
		categoryConfig: {userConfig},
	}
	system := map[category][]string{
		categoryData:   {systemData},
		categoryConfig: {systemConfig},
	}

	localWorkingDirA := t.TempDir()
	localWorkingDirB := t.TempDir()
	rLocal := mustResolverWithLocalDirs(t, ScopeLocal, []string{localWorkingDirA, localWorkingDirB}, user, system)
	gotLocal := rLocal.DataDirs()
	wantLocal := []string{
		filepath.Join(localWorkingDirA, ".demo", "data"),
		filepath.Join(localWorkingDirB, ".demo", "data"),
		userData,
		systemData,
	}
	if !reflect.DeepEqual(gotLocal, wantLocal) {
		t.Fatalf("local dirs mismatch:\nwant: %#v\ngot:  %#v", wantLocal, gotLocal)
	}

	rUser := mustResolver(t, ScopeUser, wd, user, system)
	gotUser := rUser.DataDirs()
	wantUser := []string{userData, systemData}
	if !reflect.DeepEqual(gotUser, wantUser) {
		t.Fatalf("user dirs mismatch:\nwant: %#v\ngot:  %#v", wantUser, gotUser)
	}

	rSystem := mustResolver(t, ScopeSystem, wd, user, system)
	gotSystem := rSystem.ConfigDirs()
	wantSystem := []string{systemConfig}
	if !reflect.DeepEqual(gotSystem, wantSystem) {
		t.Fatalf("system dirs mismatch:\nwant: %#v\ngot:  %#v", wantSystem, gotSystem)
	}
}

func TestDirMethods(t *testing.T) {
	wd := t.TempDir()
	userData := filepath.Join(wd, "user", "data")
	systemData := filepath.Join(wd, "system", "data")

	r := mustResolver(t, ScopeUser, wd,
		map[category][]string{categoryData: {userData}},
		map[category][]string{categoryData: {systemData}},
	)

	got := r.DataDir()
	if got != userData {
		t.Fatalf("most relevant data dir mismatch: want %q, got %q", userData, got)
	}

	gotDataDirs := r.DataDirs()
	if !reflect.DeepEqual(gotDataDirs, []string{userData, systemData}) {
		t.Fatalf("unexpected data dirs: %#v", gotDataDirs)
	}
}

func TestResolverScopeForPath(t *testing.T) {
	wd := t.TempDir()
	userData := filepath.Join(wd, "user", "data")
	systemData := filepath.Join(wd, "system", "data")

	userFn := func(string, category) ([]string, error) {
		return []string{userData}, nil
	}
	systemFn := func(string, category) ([]string, error) {
		return []string{systemData}, nil
	}

	localResolver := newScopedResolver("demo", ScopeLocal, []ResolverOption{WithLocalDir(wd)}, userFn, systemFn)
	userResolver := newScopedResolver("demo", ScopeUser, []ResolverOption{WithLocalDir(wd)}, userFn, systemFn)
	systemResolver := newScopedResolver("demo", ScopeSystem, []ResolverOption{WithLocalDir(wd)}, userFn, systemFn)

	localPath := filepath.Join(wd, ".demo", "data", "nested", "item.db")
	userPath := filepath.Join(userData, "nested", "user.db")
	systemPath := filepath.Join(systemData, "nested", "system.db")

	scope, ok := localResolver.ScopeForPath(localPath)
	if !ok || scope != ScopeLocal {
		t.Fatalf("local scope mismatch: want (%s, true), got (%s, %t)", ScopeLocal, scope, ok)
	}

	_, ok = userResolver.ScopeForPath(localPath)
	if ok {
		t.Fatalf("user resolver should not match local path %q", localPath)
	}

	scope, ok = userResolver.ScopeForPath(userPath)
	if !ok || scope != ScopeUser {
		t.Fatalf("user scope mismatch: want (%s, true), got (%s, %t)", ScopeUser, scope, ok)
	}

	_, ok = systemResolver.ScopeForPath(userPath)
	if ok {
		t.Fatalf("system resolver should not match user path %q", userPath)
	}

	scope, ok = systemResolver.ScopeForPath(systemPath)
	if !ok || scope != ScopeSystem {
		t.Fatalf("system scope mismatch: want (%s, true), got (%s, %t)", ScopeSystem, scope, ok)
	}

	sanitizedResolver := newScopedResolver("A/B Demo", ScopeLocal, []ResolverOption{WithLocalDir(wd)}, userFn, systemFn)
	sanitizedPath := filepath.Join(wd, ".A_B_Demo", "cache", "cache.db")
	scope, ok = sanitizedResolver.ScopeForPath(sanitizedPath)
	if !ok || scope != ScopeLocal {
		t.Fatalf("sanitized local scope mismatch: want (%s, true), got (%s, %t)", ScopeLocal, scope, ok)
	}
}

func mustResolver(t *testing.T, scope Scope, workingDir string, user map[category][]string, system map[category][]string, extraOpts ...ResolverOption) *Resolver {
	return mustResolverWithLocalDirs(t, scope, []string{workingDir}, user, system, extraOpts...)
}

func mustResolverWithLocalDirs(t *testing.T, scope Scope, workingDirs []string, user map[category][]string, system map[category][]string, extraOpts ...ResolverOption) *Resolver {
	t.Helper()

	opts := []ResolverOption{WithLocalDirs(workingDirs...)}
	opts = append(opts, extraOpts...)

	r := newScopedResolver(
		"demo",
		scope,
		opts,
		func(_ string, cat category) ([]string, error) {
			return append([]string(nil), user[cat]...), nil
		},
		func(_ string, cat category) ([]string, error) {
			return append([]string(nil), system[cat]...), nil
		},
	)

	return r
}
