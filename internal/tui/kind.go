package tui

type Kind int

const (
	SearchKind Kind = iota
	PreviewListKind
	ContentKind
	InfoKind
	PinListKind
)
