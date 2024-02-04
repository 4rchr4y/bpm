package require

import (
	"io"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
)

func TestNoArgs(t *testing.T) {
	t.Run("Command should pass without arguments when no arguments are expected", func(t *testing.T) {
		cmd := createTestCommand(NoArgs, []string{})
		err := cmd.Execute()
		require.NoError(t, err, "Expected no error when no arguments are passed")
	})

	t.Run("Command should fail with arguments when no arguments are expected", func(t *testing.T) {
		cmd := createTestCommand(NoArgs, []string{"one"})
		err := cmd.Execute()
		require.Error(t, err, "Expected an error when arguments are passed")
		require.Contains(t, err.Error(), `"root" accepts no arguments`, "Error message should indicate that no arguments are accepted")
	})
}

func TestExactArgs(t *testing.T) {
	t.Run("Command should pass with one argument when exactly one argument is expected", func(t *testing.T) {
		cmd := createTestCommand(ExactArgs(1), []string{"one"})
		err := cmd.Execute()
		require.NoError(t, err, "Expected no error when one argument is passed and one is expected")
	})

	t.Run("Command should fail without arguments when exactly one argument is expected", func(t *testing.T) {
		cmd := createTestCommand(ExactArgs(1), []string{})
		err := cmd.Execute()
		require.Error(t, err, "Expected an error when no arguments are passed but one is expected")
		require.Contains(t, err.Error(), `"root" requires 1 argument`, "Error message should indicate that one argument is required")
	})

	t.Run("Command should fail when exactly two arguments are expected but not provided", func(t *testing.T) {
		cmd := createTestCommand(ExactArgs(2), []string{"one"})
		err := cmd.Execute()
		require.Error(t, err, "Expected an error when less than two arguments are passed")
		require.Contains(t, err.Error(), `"root" requires 2 arguments`, "Error message should indicate that two arguments are required")
	})
}

func TestMaximumNArgs(t *testing.T) {
	t.Run("Command should pass with fewer arguments than maximum allowed", func(t *testing.T) {
		cmd := createTestCommand(MaximumNArgs(2), []string{"one"})
		err := cmd.Execute()
		require.NoError(t, err, "Expected no error when the number of arguments is less than the maximum allowed")
	})

	t.Run("Command should pass with exactly maximum allowed arguments", func(t *testing.T) {
		cmd := createTestCommand(MaximumNArgs(2), []string{"one", "two"})
		err := cmd.Execute()
		require.NoError(t, err, "Expected no error when the number of arguments equals the maximum allowed")
	})

	t.Run("Command should fail with more arguments than maximum allowed", func(t *testing.T) {
		cmd := createTestCommand(MaximumNArgs(2), []string{"one", "two", "three"})
		err := cmd.Execute()
		require.Error(t, err, "Expected an error when the number of arguments exceeds the maximum allowed")
		require.Contains(t, err.Error(), `"root" accepts at most 2 arguments`, "Error message should indicate that the maximum number of arguments is exceeded")
	})
}

func TestMinimumNArgs(t *testing.T) {
	t.Run("Command should fail with fewer arguments than minimum required", func(t *testing.T) {
		cmd := createTestCommand(MinimumNArgs(2), []string{"one"})
		err := cmd.Execute()
		require.Error(t, err, "Expected an error when the number of arguments is less than the minimum required")
		require.Contains(t, err.Error(), `"root" requires at least 2 arguments`, "Error message should indicate that the minimum number of arguments is not met")
	})

	t.Run("Command should pass with exactly minimum required arguments", func(t *testing.T) {
		cmd := createTestCommand(MinimumNArgs(2), []string{"one", "two"})
		err := cmd.Execute()
		require.NoError(t, err, "Expected no error when the number of arguments equals the minimum required")
	})

	t.Run("Command should pass with more than minimum required arguments", func(t *testing.T) {
		cmd := createTestCommand(MinimumNArgs(2), []string{"one", "two", "three"})
		err := cmd.Execute()
		require.NoError(t, err, "Expected no error when the number of arguments exceeds the minimum required")
	})
}

func createTestCommand(validateFunc cobra.PositionalArgs, args []string) *cobra.Command {
	cmd := &cobra.Command{
		Use:  "root",
		Run:  func(*cobra.Command, []string) {},
		Args: validateFunc,
	}
	cmd.SetArgs(args)
	cmd.SetOutput(io.Discard)
	return cmd
}
