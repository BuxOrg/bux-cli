package cmd

import "errors"

// ErrModelIsNil is returned when a model is nil
var ErrModelIsNil = errors.New("model is nil")

// ErrModeIsRequired is returned when a mode is required
var ErrModeIsRequired = errors.New("mode is required")

// ErrUnknownMode is returned when a mode is unknown
var ErrUnknownMode = errors.New("unknown mode")

// ErrServerModeIsNotImplemented is returned when a server mode is not implemented
var ErrServerModeIsNotImplemented = errors.New("server mode is not implemented")

// ErrFailedToLoadBux is returned when bux fails to load
var ErrFailedToLoadBux = errors.New("failed to load bux")

// ErrXpubIsRequired is returned when a xpub is required
var ErrXpubIsRequired = errors.New("xpub is required")

// ErrXpubIDIsRequired is returned when a xpub id is required
var ErrXpubIDIsRequired = errors.New("xpub id is required")

// ErrXprivIsRequired is returned when a xpriv is required
var ErrXprivIsRequired = errors.New("xpriv is required")

// ErrUnknownSubcommand is returned when a subcommand is unknown
var ErrUnknownSubcommand = errors.New("unknown subcommand")

// ErrXpubOrXpubIDIsRequired is returned when a xpub or xpub id is required
var ErrXpubOrXpubIDIsRequired = errors.New("xpub or xpub id is required")

// ErrNoXpubsFound is returned when no xpubs are found
var ErrNoXpubsFound = errors.New("no xpubs found")
