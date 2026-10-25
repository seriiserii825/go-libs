package menu

import "testing"

func TestBreadcrumbPath(t *testing.T) {
	if got := BreadcrumbPath(); got != "" {
		t.Fatalf("empty stack: got %q, want \"\"", got)
	}

	popMain := PushBreadcrumb("Main Menu")
	if got, want := BreadcrumbPath(), "Main Menu"; got != want {
		t.Errorf("got %q, want %q", got, want)
	}

	popPage := PushBreadcrumb("Page checks")
	if got, want := BreadcrumbPath(), "Main Menu > Page checks"; got != want {
		t.Errorf("got %q, want %q", got, want)
	}

	popSEO := PushBreadcrumb("SEO")
	if got, want := BreadcrumbPath(), "Main Menu > Page checks > SEO"; got != want {
		t.Errorf("got %q, want %q", got, want)
	}

	popSEO()
	if got, want := BreadcrumbPath(), "Main Menu > Page checks"; got != want {
		t.Errorf("after pop: got %q, want %q", got, want)
	}

	popPage()
	popMain()
	if got := BreadcrumbPath(); got != "" {
		t.Errorf("after popping everything: got %q, want \"\"", got)
	}
}
