package gappdirs

import (
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
)

func TestNewResolverValidatesAndAppliesOptions(t *testing.T) {
	t.Run("missing app name", func(t *testing.T) {
		_, err := NewResolver("")
		if err == nil {
			t.Fatal("expected error for empty app name")
		}
	})

	t.Run("sanitizes app name", func(t *testing.T) {
		wd := t.TempDir()
		r, err := newResolver(
			"A/B Demo",
			[]Option{WithScope(ScopeLocal), WithWorkingDir(wd)},
			func(string, category) ([]string, error) { return []string{filepath.Join(wd, "u")}, nil },
			func(string, category) ([]string, error) { return []string{filepath.Join(wd, "s")}, nil },
		)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if r.ctx.appName != "a_b_demo" {
			t.Fatalf("sanitized app name mismatch: want %q, got %q", "a_b_demo", r.ctx.appName)
		}

		got, err := r.DataDir()
		if err != nil {
			t.Fatalf("unexpected error from DataDir: %v", err)
		}
		want := filepath.Join(wd, ".a_b_demo", "data")
		if got != want {
			t.Fatalf("local data dir mismatch: want %q, got %q", want, got)
		}
	})

	t.Run("nil option", func(t *testing.T) {
		_, err := NewResolver("demo", nil)
		if err == nil {
			t.Fatal("expected error for nil option")
		}
	})

	t.Run("default scope and permissions", func(t *testing.T) {
		wd := t.TempDir()
		r, err := newResolver(
			"demo",
			[]Option{WithWorkingDir(wd)},
			func(string, category) ([]string, error) { return []string{filepath.Join(wd, "u")}, nil },
			func(string, category) ([]string, error) { return []string{filepath.Join(wd, "s")}, nil },
		)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if r.ctx.scope != ScopeUser {
			t.Fatalf("default scope mismatch: want %s, got %s", ScopeUser, r.ctx.scope)
		}
		if r.ctx.defaultDirPerm != defaultDirPerm {
			t.Fatalf("default dir permission mismatch: want %o, got %o", defaultDirPerm, r.ctx.defaultDirPerm)
		}
	})

	t.Run("duplicate options use last value", func(t *testing.T) {
		wd := t.TempDir()
		r, err := newResolver(
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
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if r.ctx.scope != ScopeLocal {
			t.Fatalf("scope mismatch: want %s, got %s", ScopeLocal, r.ctx.scope)
		}
		if r.ctx.defaultDirPerm != 0o750 {
			t.Fatalf("permission mismatch: want %o, got %o", 0o750, r.ctx.defaultDirPerm)
		}
	})
}

func TestScopedConstructorsForceScope(t *testing.T) {
	type scopedCtor func(string, ...Option) (Resolver, error)

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
			r, err := tc.ctor("demo", WithScope(tc.conflicting), WithWorkingDir(t.TempDir()))
			if err != nil {
				t.Fatalf("constructor error: %v", err)
			}
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
	r, err := NewResolver("demo", WithWorkingDir(t.TempDir()))
	if err != nil {
		t.Fatalf("constructor error: %v", err)
	}
	if _, ok := r.(*resolver); !ok {
		t.Fatalf("expected private resolver implementation, got %T", r)
	}
}

func TestSanitizeAppName(t *testing.T) {
	for _, tc := range []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{name: "mixed case and spaces", input: "My App", want: "my_app"},
		{name: "single dot", input: ".", want: "_"},
		{name: "double dots", input: "..", want: "__"},
		{name: "forward slash", input: "my/app", want: "my_app"},
		{name: "backslash", input: `my\app`, want: "my_app"},
		{name: "collapse underscores", input: "My   App", want: "my_app"},
		{name: "empty", input: "   ", wantErr: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := sanitizeAppName(tc.input)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected error for %q", tc.input)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tc.want {
				t.Fatalf("sanitized name mismatch: want %q, got %q", tc.want, got)
			}
		})
	}
}

func TestOptionValidation(t *testing.T) {
	wd := t.TempDir()
	userFn := func(string, category) ([]string, error) { return []string{filepath.Join(wd, "u")}, nil }
	systemFn := func(string, category) ([]string, error) { return []string{filepath.Join(wd, "s")}, nil }

	for _, tc := range []struct {
		name string
		opts []Option
	}{
		{name: "invalid scope", opts: []Option{WithScope(Scope(-1)), WithWorkingDir(wd)}},
		{name: "empty working dir", opts: []Option{WithWorkingDir("  ")}},
		{name: "zero default permission", opts: []Option{WithDefaultDirPerm(0), WithWorkingDir(wd)}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, err := newResolver("demo", tc.opts, userFn, systemFn)
			if err == nil {
				t.Fatal("expected option validation error")
			}
		})
	}
}

func TestDirsNilResolver(t *testing.T) {
	var r *resolver
	if _, err := r.DataDirs(); err == nil {
		t.Fatal("expected error for nil resolver")
	}
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
	gotLocal, err := rLocal.DataDirs()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	wantLocal := []string{
		filepath.Join(wd, ".demo", "data"),
		userData,
		systemData,
	}
	if !reflect.DeepEqual(gotLocal, wantLocal) {
		t.Fatalf("local dirs mismatch:\nwant: %#v\ngot:  %#v", wantLocal, gotLocal)
	}

	rUser := mustResolver(t, ScopeUser, wd, user, system)
	gotUser, err := rUser.DataDirs()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	wantUser := []string{userData, systemData}
	if !reflect.DeepEqual(gotUser, wantUser) {
		t.Fatalf("user dirs mismatch:\nwant: %#v\ngot:  %#v", wantUser, gotUser)
	}

	rSystem := mustResolver(t, ScopeSystem, wd, user, system)
	gotSystem, err := rSystem.ConfigDirs()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
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

	got, err := r.DataDir()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != userData {
		t.Fatalf("most relevant data dir mismatch: want %q, got %q", userData, got)
	}

	gotDataDirs, err := r.DataDirs()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(gotDataDirs, []string{userData, systemData}) {
		t.Fatalf("unexpected data dirs: %#v", gotDataDirs)
	}
}

func TestDirsDedupeAndNormalization(t *testing.T) {
	wd := t.TempDir()
	dupBase := filepath.Join(wd, "shared")

	userData := []string{dupBase, filepath.Join(dupBase, ".")}
	systemData := []string{dupBase}
	if runtime.GOOS == "windows" {
		systemData = []string{filepath.Join(wd, "SHARED")}
	}

	r := mustResolver(t, ScopeUser, wd,
		map[category][]string{categoryData: userData},
		map[category][]string{categoryData: systemData},
	)

	got, err := r.DataDirs()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if runtime.GOOS == "windows" {
		if len(got) != 1 {
			t.Fatalf("expected 1 deduped path on windows, got %#v", got)
		}
		if filepath.Clean(got[0]) != filepath.Clean(dupBase) {
			t.Fatalf("deduped path mismatch: want %q, got %q", dupBase, got[0])
		}
		return
	}

	if !reflect.DeepEqual(got, []string{dupBase}) {
		t.Fatalf("dedupe mismatch:\nwant: %#v\ngot:  %#v", []string{dupBase}, got)
	}
}

func mustResolver(t *testing.T, scope Scope, workingDir string, user map[category][]string, system map[category][]string, extraOpts ...Option) *resolver {
	t.Helper()

	opts := []Option{WithScope(scope), WithWorkingDir(workingDir)}
	opts = append(opts, extraOpts...)

	r, err := newResolver(
		"demo",
		opts,
		func(_ string, cat category) ([]string, error) {
			return append([]string(nil), user[cat]...), nil
		},
		func(_ string, cat category) ([]string, error) {
			return append([]string(nil), system[cat]...), nil
		},
	)
	if err != nil {
		t.Fatalf("newResolver error: %v", err)
	}

	return r
}
