package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
	"os"
	"syscall"
)

func init() {
	for _, cmd := range []*cobra.Command{encryptCmd, decryptCmd} {
		rootCmd.AddCommand(cmd)
		cmd.Flags().StringVarP(&outputFile, "output", "o", "", "The target file location. Can only be used if a single file is passed. Specify '-' to output to stdout.")
		cmd.Flags().StringVarP(&vaultPassword, "password", "p", "", "The vault password. This defaults to the value of environment variable `VAULT_PASSWORD`.")
		cmd.Flags().BoolVarP(&confirmPassword, "confirm-password", "c", false, "Confirm the vault password when prompting.")
	}
}

var outputFile = ""
var vaultPassword = ""
var confirmPassword = false

var rootCmd = &cobra.Command{
	Use:   "terraform-encrypt",
	Short: "Manage secret files as code using terraform",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func requireOneInputFile(_ *cobra.Command, args []string) error {
	if outputFile != "" {
		if len(args) == 0 {
			return errors.New("provide one input file")
		}
		if len(args) > 1 {
			return errors.New("only one input file can be provided when using the --output option")
		}
	}
	if len(args) < 1 {
		return errors.New("provide at least one input file")
	}
	return nil
}

func findPassword() string {
	// 1. Command Line Option
	result := vaultPassword
	if result != "" {
		return result
	}

	// 2. Environment Variable
	result = os.Getenv("VAULT_PASSWORD")
	if result != "" {
		return result
	}

	// 3. Prompt
	for result == "" {
		prompted := promptPassword("Vault Password")

		// Confirm Password
		if confirmPassword && prompted != promptPassword("Confirm Password") {
			fmt.Println("Passwords do not match, please try again.")
			continue
		}

		result = prompted
	}
	return result
}

func promptPassword(prompt string) string {
	fmt.Print(prompt + ": ")
	password, err := terminal.ReadPassword(syscall.Stdin)
	if err != nil {
		panic(err)
	}
	fmt.Println()
	return string(password)
}
