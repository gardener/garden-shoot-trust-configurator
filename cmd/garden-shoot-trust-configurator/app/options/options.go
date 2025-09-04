// SPDX-FileCopyrightText: SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

package options

import (
	"errors"
	"slices"
	"time"

	"github.com/spf13/pflag"
)

// Options contain the server options.
type Options struct {
	ResyncOptions ResyncOptions
}

// ResyncOptions holds options regarding the resync interval between reconciliations.
type ResyncOptions struct {
	Duration time.Duration
}

// AddFlags adds the [ResyncOptions] flags to the flagset.
func (o *ResyncOptions) AddFlags(fs *pflag.FlagSet) {
	fs.DurationVar(&o.Duration, "resync-period", time.Minute*30, "The period between reconciliations of cluster shoot trust configurator information.")
}

// Validate checks if options are valid.
func (o *ResyncOptions) Validate() []error {
	var errs []error
	if o.Duration <= 0 {
		errs = append(errs, errors.New("--resync-period must be positive"))
	}
	return errs
}

// ApplyTo applies the options to the configuration.
func (o *ResyncOptions) ApplyTo(c *ResyncConfig) error {
	c.Duration = o.Duration
	return nil
}

// ResyncConfig holds configurations regarding the resync interval between reconciliations.
type ResyncConfig struct {
	Duration time.Duration
}

// NewOptions return options with default values.
func NewOptions() *Options {
	opts := &Options{
		ResyncOptions: ResyncOptions{},
	}
	return opts
}

// AddFlags adds server options to flagset
func (o *Options) AddFlags(fs *pflag.FlagSet) {
	o.ResyncOptions.AddFlags(fs)
}

// ApplyTo applies the options to the configuration.
func (o *Options) ApplyTo(server *Config) error {
	return o.ResyncOptions.ApplyTo(&server.Resync)
}

// Validate checks if options are valid.
func (o *Options) Validate() []error {
	return slices.Concat(
		o.ResyncOptions.Validate(),
	)
}

// Config has all the context to run the shoot trust configurator server.
type Config struct {
	Resync ResyncConfig
}
