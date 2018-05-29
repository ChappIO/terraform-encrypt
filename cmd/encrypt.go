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
	base64String := base64.StdEncoding.EncodeToString(message.CipherText)

	targetFile := outputFile

	if targetFile == "" {
		targetFile = sourceFile
	}

	file, err := os.Create(targetFile)
	if err != nil {
		fmt.Printf("Failed to write to '%s': %s", sourceFile, err.Error())
		os.Exit(1)
	}
	defer file.Close()

	// Write the base64 string in chunks
	written := 0
	for written < len(base64String) {
		start := written
		end := start + 80
		if end >= len(base64String) {
			end = len(base64String)
		}

		str := base64String[start:end]
		file.WriteString(str)
		file.WriteString("\n")
		written += end - start
	}
}
