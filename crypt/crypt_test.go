package crypt

import (
	"testing"
)

func TestEncryptedMessageCanBeDecryptedWithSameKey(t *testing.T) {
	inputs := []string{
		"Sup?!",
		"This is some secret file which is longer than the block size",
	}

	for _, input := range inputs {
		message := Message{
			PlainText: []byte(input),
		}

		message.Encrypt("this is a key as random as it is")

		message.PlainText = nil

		err := message.Decrypt("this is a key as random as it is")

		output := string(message.PlainText)

		if output != input {
			t.Errorf("Expected decrypt(encrypt()) to return '%s' but got '%s'", input, output)
		}
		if err != nil {
			t.Errorf("An unexpected error was returned: %s", err)
		}
	}
}

func TestEncryptedMessageCannotBeDecryptedWithArbitraryKey(t *testing.T) {
	inputs := []string{
		"Sup?!",
		"This is some secret file which is longer than the block size",
	}

	for _, input := range inputs {

		message := Message{
			PlainText: []byte(input),
		}

		message.Encrypt("this is a key as random as it is")

		message.PlainText = nil

		err := message.Decrypt("this message has a different key")

		output := string(message.PlainText)

		if output == input {
			t.Errorf("The message was decrypted using an arbitrary key")
		}

		if err == nil {
			t.Errorf("The expected error was not returned")
		}
	}
}
