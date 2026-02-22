package gappdirs

import (
	"io/fs"
	"path/filepath"
	"reflect"
	"testing"
)

func TestNewResolverAppliesOptions(t *testing.T) {
	t.Run("missing app name defaults", func(t *testing.T) {
		r := NewResolver("")
		impl, ok := r.(*resolver)
		if !ok {
			t.Fatalf("unexpected resolver implementation type: %T", r)
		}
		if impl.ctx.appName != "unnamed_app" {
			t.Fatalf("default app name mismatch: want %q, got %q", "unnamed_app", impl.ctx.appName)
		}
	})

	t.Run("sanitizes app name", func(t *testing.T) {
		wd := t.TempDir()
		r := newResolver(
			"A/B Demo",
			[]Option{WithScope(ScopeLocal), WithWorkingDir(wd)},
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
		r := NewResolver("demo", nil)
		impl, ok := r.(*resolver)
		if !ok {
			t.Fatalf("unexpected resolver implementation type: %T", r)
		}
		if impl.ctx.scope != ScopeUser {
			t.Fatalf("default scope mismatch: want %s, got %s", ScopeUser, impl.ctx.scope)
		}
	})

	t.Run("default scope and permissions", func(t *testing.T) {
		wd := t.TempDir()
		r := newResolver(
			"demo",
			[]Option{WithWorkingDir(wd)},
			func(string, category) ([]string, error) { return []string{filepath.Join(wd, "u")}, nil },
			func(string, category) ([]string, error) { return []string{filepath.Join(wd, "s")}, nil },
		)
		if r.ctx.scope != ScopeUser {
			t.Fatalf("default scope mismatch: want %s, got %s", ScopeUser, r.ctx.scope)
		}
		if r.ctx.defaultDirPerm != defaultDirPerm {
			t.Fatalf("default dir permission mismatch: want %o, got %o", defaultDirPerm, r.ctx.defaultDirPerm)
		}
	})

	t.Run("duplicate options use last value", func(t *testing.T) {
		wd := t.TempDir()
		r := newResolver(
			"demo",
			[]Option{
				WithWorkingDir(wd),
				WithScope(ScopeSystem),
				WithScope(ScopeLocal),
				WithDefaultDirPerm(0o700),
				WithDefaultDirPerm(0o750),
			},
			func(string, category) ([]string, error) { return []string{filepath.Join(wd, "u")}, nil },
			func(string, category) ([]string, error) { return []string{filepath.Join(wd, "s")}, nil },
		)
		if r.ctx.scope != ScopeLocal {
			t.Fatalf("scope mismatch: want %s, got %s", ScopeLocal, r.ctx.scope)
		}
		if r.ctx.defaultDirPerm != 0o750 {
			t.Fatalf("permission mismatch: want %o, got %o", 0o750, r.ctx.defaultDirPerm)
		}
	})
}

func TestScopedConstructorsForceScope(t *testing.T) {
	type scopedCtor func(string, ...Option) Resolver

	for _, tc := range []struct {
		name        string
		ctor        scopedCtor
		conflicting Scope
		wantScope   Scope
	}{
		{name: "user", ctor: NewUserResolver, conflicting: ScopeLocal, wantScope: ScopeUser},
		{name: "system", ctor: NewSystemResolver, conflicting: ScopeLocal, wantScope: ScopeSystem},
		{name: "local", ctor: NewLocalResolver, conflicting: ScopeSystem, wantScope: ScopeLocal},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r := tc.ctor("demo", WithScope(tc.conflicting), WithWorkingDir(t.TempDir()))
			impl, ok := r.(*resolver)
			if !ok {
				t.Fatalf("unexpected resolver implementation type: %T", r)
			}
			if impl.ctx.scope != tc.wantScope {
				t.Fatalf("scope mismatch: want %s, got %s", tc.wantScope, impl.ctx.scope)
			}
		})
	}
}

func TestConstructorsReturnResolverInterface(t *testing.T) {
	r := NewResolver("demo", WithWorkingDir(t.TempDir()))
	if _, ok := r.(*resolver); !ok {
		t.Fatalf("expected private resolver implementation, got %T", r)
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

	t.Run("invalid scope defaults to user", func(t *testing.T) {
		r := newResolver("demo", []Option{WithScope(Scope(-1)), WithWorkingDir(wd)}, userFn, systemFn)
		if r.ctx.scope != ScopeUser {
			t.Fatalf("scope fallback mismatch: want %s, got %s", ScopeUser, r.ctx.scope)
		}
	})

	t.Run("empty working dir option is ignored", func(t *testing.T) {
		r := newResolver("demo", []Option{WithWorkingDir("  ")}, userFn, systemFn)
		if r.ctx.workingDir != "" {
			t.Fatalf("working dir should remain empty when option is ignored, got %q", r.ctx.workingDir)
		}
	})

	t.Run("invalid permission falls back to default", func(t *testing.T) {
		r := newResolver("demo", []Option{WithDefaultDirPerm(0), WithWorkingDir(wd)}, userFn, systemFn)
		if r.ctx.defaultDirPerm != defaultDirPerm {
			t.Fatalf("permission fallback mismatch: want %o, got %o", defaultDirPerm, r.ctx.defaultDirPerm)
		}

		r = newResolver("demo", []Option{WithDefaultDirPerm(fs.ModeDir | 0o755), WithWorkingDir(wd)}, userFn, systemFn)
		if r.ctx.defaultDirPerm != defaultDirPerm {
			t.Fatalf("permission fallback mismatch for mode bits: want %o, got %o", defaultDirPerm, r.ctx.defaultDirPerm)
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

	rLocal := mustResolver(t, ScopeLocal, wd, user, system)
	gotLocal := rLocal.DataDirs()
	wantLocal := []string{
		filepath.Join(wd, ".demo", "data"),
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

func mustResolver(t *testing.T, scope Scope, workingDir string, user map[category][]string, system map[category][]string, extraOpts ...Option) *resolver {
	t.Helper()

	opts := []Option{WithScope(scope), WithWorkingDir(workingDir)}
	opts = append(opts, extraOpts...)

	r := newResolver(
		"demo",
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
