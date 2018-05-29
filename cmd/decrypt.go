package cmd

import (
	"github.com/spf13/cobra"
	"fmt"
	"io/ioutil"
	"github.com/ChappIO/terraform-encrypt/crypt"
	"os"
	"encoding/base64"
)

var decryptCmd = &cobra.Command{
	Use:   "decrypt [sourceFiles...]",
	Short: "Decrypt a file so it can be viewed",
	Args:  requireOneInputFile,
	Run: func(cmd *cobra.Command, args []string) {
		password := findPassword()

		for _, file := range args {
			decryptFile(file, password)
		}
	},
}

func decryptFile(sourceFile string, password string) {
	base64String, err := ioutil.ReadFile(sourceFile)
	if err != nil {
		fmt.Printf("Failed to open '%s': %s\n", sourceFile, err.Error())
		os.Exit(1)
	}
	ciphertext, err := base64.StdEncoding.DecodeString(string(base64String))
	if err != nil {
		fmt.Printf("Corrupted file '%s': %s\n", sourceFile, err.Error())
		os.Exit(1)
	}

	message := crypt.Message{
		CipherText: ciphertext,
	}
	err = message.Decrypt(password)
	if err != nil {
		fmt.Printf("Failed to decrypt %s: %s\n", sourceFile, err)
		os.Exit(1)
	}

	targetFile := outputFile

	if targetFile == "" {
		targetFile = sourceFile
	}

	file, err := os.Create(targetFile)
	if err != nil {
		fmt.Printf("Failed to write to '%s': %s\n", sourceFile, err.Error())
		os.Exit(1)
	}
	defer file.Close()
	file.Write(message.PlainText)
}
