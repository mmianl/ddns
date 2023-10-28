package internal

import (
	"testing"
)

// Test NewStaticIPAddressProvider() method of StaticIPAddressProvider
func TestNewStaticIPAddressProvider(t *testing.T) {
	// Create a new StaticIPAddressProvider
	provider := NewStaticIPAddressProvider(defaultStaticIPAddressProviderConfig)

	// Verify that the returned StaticIPAddressProvider has the same static ip address as the one passed via configuration
	want := "127.0.0.1"
	got := provider.address

	if got != want {
		t.Errorf("got %s, wanted %s", got, want)
	}
}

// Test GetIPAddress() method of StaticIPAddressProvider
func TestStaticIPAddressProviderGetIPAddress(t *testing.T) {
	// Create a new StaticIPAddressProvider
	provider := NewStaticIPAddressProvider(defaultStaticIPAddressProviderConfig)

	// Call GetIPAddress() method
	got, err := provider.GetIPAddress()

	// Verify that the returned ip address is the same as the one passed via configuration
	want := "127.0.0.1"

	if err != nil {
		t.Errorf("got %s, wanted %s", err, want)
	}

	if *got != want {
		t.Errorf("got %s, wanted %s", *got, want)
	}
}
