// SPDX-FileCopyrightText: SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

package options

import (
	"errors"
	"net"
	"slices"
	"strconv"
	"time"

	"github.com/spf13/pflag"
)

// Options contain the server options.
type Options struct {
	ResyncOptions  ResyncOptions
	ServingOptions ServingOptions
}

// ServingOptions are options applied to the shoot trust configurator server.
type ServingOptions struct {
	Address string
	Port    uint
}

// AddFlags adds server options to flagset
func (o *ServingOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.Address, "address", "", "The IP address that the server will listen on. If unspecified all interfaces will be used.")
	fs.UintVar(&o.Port, "port", 10443, "The port that the server will listen on.")
}

// Validate checks if options are valid.
// Returns errors if the port is out of range or address is invalid.
func (o *ServingOptions) Validate() []error {
	var errs []error
	if o.Port == 0 || o.Port > 65535 {
		errs = append(errs, errors.New("--port must be between 1 and 65535"))
	}
	if o.Address != "" && net.ParseIP(o.Address) == nil {
		errs = append(errs, errors.New("--address must be a valid IP address"))
	}
	return errs
}

// ApplyTo applies the options to the configuration.
func (o *ServingOptions) ApplyTo(c *ServingConfig) error {
	c.Address = net.JoinHostPort(o.Address, strconv.FormatUint(uint64(o.Port), 10))
	return nil
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
		ResyncOptions:  ResyncOptions{},
		ServingOptions: ServingOptions{},
	}
	return opts
}

// AddFlags adds server options to flagset
func (o *Options) AddFlags(fs *pflag.FlagSet) {
	o.ServingOptions.AddFlags(fs)
	o.ResyncOptions.AddFlags(fs)
}

// ApplyTo applies the options to the configuration.
func (o *Options) ApplyTo(server *Config) error {
	if err := o.ResyncOptions.ApplyTo(&server.Resync); err != nil {
		return err
	}

	return o.ServingOptions.ApplyTo(&server.Serving)
}

// Validate checks if options are valid.
func (o *Options) Validate() []error {
	return slices.Concat(
		o.ResyncOptions.Validate(),
		o.ServingOptions.Validate(),
	)
}

// Config has all the context to run the shoot trust configurator server.
type Config struct {
	Resync  ResyncConfig
	Serving ServingConfig
}

// ServingConfig has the context to run an http server.
type ServingConfig struct {
	Address string
}
