package cmd

import (
	"github.com/spf13/cobra"
	"fmt"
	"io/ioutil"
	"github.com/ChappIO/terraform-encrypt/crypt"
	"os"
	"encoding/base64"
)

var encryptCmd = &cobra.Command{
	Use:   "encrypt [sourceFiles...]",
	Short: "Encrypt a file for safe storage",
	Args:  requireOneInputFile,
	Run: func(cmd *cobra.Command, args []string) {
		password := findPassword()

		for _, file := range args {
			encryptFile(file, password)
		}
	},
}

func encryptFile(sourceFile string, password string) {
	plaintext, err := ioutil.ReadFile(sourceFile)
	if err != nil {
		fmt.Printf("Failed to open '%s': %s", sourceFile, err.Error())
		os.Exit(1)
	}

	message := crypt.Message{
		PlainText: plaintext,
	}
	message.Encrypt(password)
	base64String := wrapString(base64.StdEncoding.EncodeToString(message.CipherText), 80)

	targetFile := outputFile

	if targetFile == "" {
		targetFile = sourceFile
	}

	if targetFile == "-" {
		fmt.Println(base64String);
	} else {
		file, err := os.Create(targetFile)
		if err != nil {
			fmt.Printf("Failed to write to '%s': %s", sourceFile, err.Error())
			os.Exit(1)
		}
		defer file.Close()
		file.WriteString(base64String);
	}
}

func wrapString(input string, width int) string {
	wrappedString := "";
	// Write the base64 string in chunks
	written := 0
	for written < len(input) {
		start := written
		end := start + width
		if end >= len(input) {
			end = len(input)
		}

		str := input[start:end]
		wrappedString = wrappedString + str + "\n";
		written += end - start
	}

	return wrappedString;
}
