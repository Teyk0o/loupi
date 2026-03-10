package utils

import (
	"testing"
)

const testKey = "0102030405060708091011121314151617181920212223242526272829303132"

func TestEncryptDecryptRoundTrip(t *testing.T) {
	enc, err := NewEncryptor(testKey)
	if err != nil {
		t.Fatalf("NewEncryptor: %v", err)
	}

	original := "sensitive health data: stomach pain, nausea"
	ciphertext, err := enc.Encrypt(original)
	if err != nil {
		t.Fatalf("Encrypt: %v", err)
	}

	if ciphertext == original {
		t.Fatal("ciphertext should differ from plaintext")
	}

	decrypted, err := enc.Decrypt(ciphertext)
	if err != nil {
		t.Fatalf("Decrypt: %v", err)
	}

	if decrypted != original {
		t.Errorf("expected %q, got %q", original, decrypted)
	}
}

func TestEncryptEmptyString(t *testing.T) {
	enc, err := NewEncryptor(testKey)
	if err != nil {
		t.Fatalf("NewEncryptor: %v", err)
	}

	ciphertext, err := enc.Encrypt("")
	if err != nil {
		t.Fatalf("Encrypt empty: %v", err)
	}
	if ciphertext != "" {
		t.Errorf("expected empty string, got %q", ciphertext)
	}

	decrypted, err := enc.Decrypt("")
	if err != nil {
		t.Fatalf("Decrypt empty: %v", err)
	}
	if decrypted != "" {
		t.Errorf("expected empty string, got %q", decrypted)
	}
}

func TestDecryptWrongKey(t *testing.T) {
	enc1, _ := NewEncryptor(testKey)
	enc2, _ := NewEncryptor("ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff")

	ciphertext, err := enc1.Encrypt("secret data")
	if err != nil {
		t.Fatalf("Encrypt: %v", err)
	}

	_, err = enc2.Decrypt(ciphertext)
	if err == nil {
		t.Fatal("expected error when decrypting with wrong key")
	}
}

func TestDecryptTamperedData(t *testing.T) {
	enc, _ := NewEncryptor(testKey)

	ciphertext, err := enc.Encrypt("secret data")
	if err != nil {
		t.Fatalf("Encrypt: %v", err)
	}

	// Tamper with the ciphertext
	tampered := ciphertext[:len(ciphertext)-2] + "AA"
	_, err = enc.Decrypt(tampered)
	if err == nil {
		t.Fatal("expected error when decrypting tampered data")
	}
}

func TestEncryptDecryptJSON(t *testing.T) {
	enc, _ := NewEncryptor(testKey)

	type Symptom struct {
		Type     string `json:"type"`
		Severity int    `json:"severity"`
	}

	original := []Symptom{
		{Type: "headache", Severity: 3},
		{Type: "nausea", Severity: 5},
	}

	ciphertext, err := enc.EncryptJSON(original)
	if err != nil {
		t.Fatalf("EncryptJSON: %v", err)
	}

	var decrypted []Symptom
	if err := enc.DecryptJSON(ciphertext, &decrypted); err != nil {
		t.Fatalf("DecryptJSON: %v", err)
	}

	if len(decrypted) != len(original) {
		t.Fatalf("expected %d items, got %d", len(original), len(decrypted))
	}
	for i, s := range decrypted {
		if s.Type != original[i].Type || s.Severity != original[i].Severity {
			t.Errorf("item %d: expected %+v, got %+v", i, original[i], s)
		}
	}
}

func TestNewEncryptorInvalidKey(t *testing.T) {
	_, err := NewEncryptor("not-hex")
	if err == nil {
		t.Fatal("expected error for non-hex key")
	}

	_, err = NewEncryptor("0102030405060708")
	if err == nil {
		t.Fatal("expected error for short key")
	}
}

func TestEncryptProducesDifferentCiphertexts(t *testing.T) {
	enc, _ := NewEncryptor(testKey)

	ct1, _ := enc.Encrypt("same plaintext")
	ct2, _ := enc.Encrypt("same plaintext")

	if ct1 == ct2 {
		t.Fatal("encrypting the same plaintext should produce different ciphertexts (random nonce)")
	}
}
