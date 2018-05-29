package cmd

import (
	"github.com/spf13/cobra"
	"fmt"
	"os"
	"errors"
	"golang.org/x/crypto/ssh/terminal"
	"syscall"
)

func init() {
	for _, cmd := range []*cobra.Command{encryptCmd, decryptCmd} {
		rootCmd.AddCommand(cmd)
		cmd.Flags().StringVarP(&outputFile, "output", "o", "", "The target file location. Can only be used if a single file is passed.")
		cmd.Flags().StringVarP(&vaultPassword, "password", "p", "", "The vault password")
	}
}

var outputFile = ""
var vaultPassword = ""

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

func requireOneInputFile(cmd *cobra.Command, args []string) error {
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
		fmt.Print("Vault Password: ")
		passwd, err := terminal.ReadPassword(int(syscall.Stdin))
		if err != nil {
			panic(err)
		}
		result = string(passwd)
		fmt.Println()
	}
	return result
}
