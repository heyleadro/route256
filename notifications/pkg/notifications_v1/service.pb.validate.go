// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: service.proto

package notifications_v1

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"net/mail"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	"google.golang.org/protobuf/types/known/anypb"
)

// ensure the imports are used
var (
	_ = bytes.MinRead
	_ = errors.New("")
	_ = fmt.Print
	_ = utf8.UTFMax
	_ = (*regexp.Regexp)(nil)
	_ = (*strings.Reader)(nil)
	_ = net.IPv4len
	_ = time.Duration(0)
	_ = (*url.URL)(nil)
	_ = (*mail.Address)(nil)
	_ = anypb.Any{}
	_ = sort.Sort
)

// Validate checks the field values on GetUserHistoryRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *GetUserHistoryRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on GetUserHistoryRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// GetUserHistoryRequestMultiError, or nil if none found.
func (m *GetUserHistoryRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *GetUserHistoryRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if m.GetUser() < 1 {
		err := GetUserHistoryRequestValidationError{
			field:  "User",
			reason: "value must be greater than or equal to 1",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	// no validation rules for Period

	if len(errors) > 0 {
		return GetUserHistoryRequestMultiError(errors)
	}

	return nil
}

// GetUserHistoryRequestMultiError is an error wrapping multiple validation
// errors returned by GetUserHistoryRequest.ValidateAll() if the designated
// constraints aren't met.
type GetUserHistoryRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m GetUserHistoryRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m GetUserHistoryRequestMultiError) AllErrors() []error { return m }

// GetUserHistoryRequestValidationError is the validation error returned by
// GetUserHistoryRequest.Validate if the designated constraints aren't met.
type GetUserHistoryRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e GetUserHistoryRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e GetUserHistoryRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e GetUserHistoryRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e GetUserHistoryRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e GetUserHistoryRequestValidationError) ErrorName() string {
	return "GetUserHistoryRequestValidationError"
}

// Error satisfies the builtin error interface
func (e GetUserHistoryRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sGetUserHistoryRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = GetUserHistoryRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = GetUserHistoryRequestValidationError{}

// Validate checks the field values on GetUserHistoryResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *GetUserHistoryResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on GetUserHistoryResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// GetUserHistoryResponseMultiError, or nil if none found.
func (m *GetUserHistoryResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *GetUserHistoryResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	for idx, item := range m.GetHis() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, GetUserHistoryResponseValidationError{
						field:  fmt.Sprintf("His[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, GetUserHistoryResponseValidationError{
						field:  fmt.Sprintf("His[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return GetUserHistoryResponseValidationError{
					field:  fmt.Sprintf("His[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if len(errors) > 0 {
		return GetUserHistoryResponseMultiError(errors)
	}

	return nil
}

// GetUserHistoryResponseMultiError is an error wrapping multiple validation
// errors returned by GetUserHistoryResponse.ValidateAll() if the designated
// constraints aren't met.
type GetUserHistoryResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m GetUserHistoryResponseMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m GetUserHistoryResponseMultiError) AllErrors() []error { return m }

// GetUserHistoryResponseValidationError is the validation error returned by
// GetUserHistoryResponse.Validate if the designated constraints aren't met.
type GetUserHistoryResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e GetUserHistoryResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e GetUserHistoryResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e GetUserHistoryResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e GetUserHistoryResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e GetUserHistoryResponseValidationError) ErrorName() string {
	return "GetUserHistoryResponseValidationError"
}

// Error satisfies the builtin error interface
func (e GetUserHistoryResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sGetUserHistoryResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = GetUserHistoryResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = GetUserHistoryResponseValidationError{}

// Validate checks the field values on HisItem with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *HisItem) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on HisItem with the rules defined in the
// proto definition for this message. If any rules are violated, the result is
// a list of violation errors wrapped in HisItemMultiError, or nil if none found.
func (m *HisItem) ValidateAll() error {
	return m.validate(true)
}

func (m *HisItem) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Order

	// no validation rules for Time

	if len(errors) > 0 {
		return HisItemMultiError(errors)
	}

	return nil
}

// HisItemMultiError is an error wrapping multiple validation errors returned
// by HisItem.ValidateAll() if the designated constraints aren't met.
type HisItemMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m HisItemMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m HisItemMultiError) AllErrors() []error { return m }

// HisItemValidationError is the validation error returned by HisItem.Validate
// if the designated constraints aren't met.
type HisItemValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e HisItemValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e HisItemValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e HisItemValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e HisItemValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e HisItemValidationError) ErrorName() string { return "HisItemValidationError" }

// Error satisfies the builtin error interface
func (e HisItemValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sHisItem.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = HisItemValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = HisItemValidationError{}