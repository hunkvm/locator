package validation

import (
	"fmt"
	"net"
	"slices"
	"strconv"
	"time"
)

type Validatable interface {
	Validate() []Error
}

func ValidateAddress(addr string) error {
	_, p, err := net.SplitHostPort(addr)
	if err != nil {
		return err
	}
	port, err := strconv.Atoi(p)
	if err != nil || port < 1 || port > 65535 {
		format := "invalid port number %s in address"
		return fmt.Errorf(format, p)
	}
	return nil
}

func ValidateOneOf[T comparable](value T, allowed []T) error {
	if slices.Contains(allowed, value) {
		return nil
	}

	return fmt.Errorf("must be one of %v", allowed)
}

func ValidateDuration(d time.Duration, min, max *time.Duration) error {
	if min != nil && *min == 0 && d < 0 {
		return fmt.Errorf("duration cannot be negative")
	} else if min != nil && *min > 0 && d < *min {
		return fmt.Errorf("duration is less than minimum %s", *min)
	}

	if max != nil && *max > 0 && d > *max {
		return fmt.Errorf("duration exceeds maximum %s", *max)
	}
	return nil
}

func ValidateInteger[T ~int | ~uint | ~int64 | ~uint64](i T, min, max *T) error {
	if min != nil && *min == 0 && i < 0 {
		return fmt.Errorf("integer cannot be negative")
	} else if min != nil && *min > 0 && i < *min {
		return fmt.Errorf("integer is less than minimum %d", *min)
	}

	if max != nil && *max > 0 && i > *max {
		return fmt.Errorf("integer exceeds maximum %d", *max)
	}
	return nil
}
