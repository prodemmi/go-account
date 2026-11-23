package pages

import "embed"

// CSS holds the static stylesheets served under /static/css. Embedding them
// keeps the admin/auth pages styled no matter what the working directory is.
//
//go:embed css/*
var CSS embed.FS

// Fonts holds the bundled IRANSansX family used by the admin panel and by apps
// that select it for their auth page. They are served under /static/fonts.
//
//go:embed fonts/*
var Fonts embed.FS
