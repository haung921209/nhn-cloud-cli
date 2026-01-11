package interactive

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"golang.org/x/term"
)

// Option represents a selectable option in interactive mode
type Option struct {
	Value       string                 // Actual value to use
	Display     string                 // Display text
	Description string                 // Optional description
	Metadata    map[string]interface{} // Additional metadata
}

// ParameterDef defines a parameter that can be prompted for
type ParameterDef struct {
	Name        string                                  // Parameter name
	DisplayName string                                  // Human-readable name
	Required    bool                                    // Is this required?
	Type        ParameterType                           // Type of parameter
	Validator   func(interface{}) error                 // Validation function
	Fetcher     func(context.Context) ([]Option, error) // Function to fetch options from API
	Default     interface{}                             // Default value
	DependsOn   []string                                // Other params this depends on
	Description string                                  // Help text
}

// ParameterType defines the type of parameter input
type ParameterType int

const (
	TypeString ParameterType = iota
	TypePassword
	TypeSelect
	TypeMultiSelect
	TypeConfirm
	TypeInt
)

// PromptManager manages interactive prompts
type PromptManager struct {
	ctx         context.Context
	values      map[string]interface{}
	definitions []ParameterDef
	interactive bool
}

// NewPromptManager creates a new prompt manager
func NewPromptManager(ctx context.Context) *PromptManager {
	return &PromptManager{
		ctx:         ctx,
		values:      make(map[string]interface{}),
		interactive: CanRunInteractive(),
	}
}

// CanRunInteractive checks if we can run in interactive mode
func CanRunInteractive() bool {
	// Check if running in terminal
	if !term.IsTerminal(int(os.Stdin.Fd())) {
		return false // In pipe or script
	}

	// Check for CI environment
	if os.Getenv("CI") != "" {
		return false
	}

	// Check for explicit disable
	if os.Getenv("NHNCLOUD_NONINTERACTIVE") != "" {
		return false
	}

	return true
}

// IsInteractive returns whether interactive mode is enabled
func (pm *PromptManager) IsInteractive() bool {
	return pm.interactive
}

// SetDefinitions sets the parameter definitions
func (pm *PromptManager) SetDefinitions(defs []ParameterDef) {
	pm.definitions = defs
}

// SetProvidedValues sets values that were provided via flags
func (pm *PromptManager) SetProvidedValues(values map[string]interface{}) {
	for k, v := range values {
		if v != nil && v != "" && v != 0 && v != false {
			pm.values[k] = v
		}
	}
}

// GetMissingRequired returns list of missing required parameters
func (pm *PromptManager) GetMissingRequired() []string {
	var missing []string
	for _, def := range pm.definitions {
		if def.Required {
			if _, exists := pm.values[def.Name]; !exists {
				missing = append(missing, def.Name)
			}
		}
	}
	return missing
}

// ValidateProvidedValues validates all provided values
func (pm *PromptManager) ValidateProvidedValues() map[string]error {
	errors := make(map[string]error)
	for _, def := range pm.definitions {
		if value, exists := pm.values[def.Name]; exists {
			if def.Validator != nil {
				if err := def.Validator(value); err != nil {
					errors[def.Name] = err
				}
			}
		}
	}
	return errors
}

// PromptForMissing prompts for missing required parameters
func (pm *PromptManager) PromptForMissing() error {
	if !pm.interactive {
		missing := pm.GetMissingRequired()
		if len(missing) > 0 {
			return fmt.Errorf("missing required parameters: %s", strings.Join(missing, ", "))
		}
		return nil
	}

	missing := pm.GetMissingRequired()
	if len(missing) == 0 {
		return nil
	}

	fmt.Println("\nMissing required parameters. Entering interactive mode...")
	fmt.Println()

	// Show what's already provided
	for _, def := range pm.definitions {
		if value, exists := pm.values[def.Name]; exists && def.Required {
			fmt.Printf("✓ %s: %v\n", def.DisplayName, pm.formatValue(def, value))
		}
	}

	if len(pm.values) > 0 {
		fmt.Println()
	}

	// Prompt for missing
	for _, def := range pm.definitions {
		if def.Required {
			if _, exists := pm.values[def.Name]; !exists {
				value, err := pm.promptForParameter(def)
				if err != nil {
					return err
				}
				pm.values[def.Name] = value
			}
		}
	}

	return nil
}

// PromptForOptional asks if user wants to configure optional parameters
func (pm *PromptManager) PromptForOptional() error {
	if !pm.interactive {
		return nil
	}

	// Check if there are optional parameters
	hasOptional := false
	for _, def := range pm.definitions {
		if !def.Required {
			if _, exists := pm.values[def.Name]; !exists {
				hasOptional = true
				break
			}
		}
	}

	if !hasOptional {
		return nil
	}

	// Ask if user wants to configure optional parameters
	configure := false
	prompt := &survey.Confirm{
		Message: "Would you like to configure additional options?",
		Default: false,
	}
	if err := survey.AskOne(prompt, &configure); err != nil {
		return err
	}

	if !configure {
		return nil
	}

	// Prompt for optional parameters
	for _, def := range pm.definitions {
		if !def.Required {
			if _, exists := pm.values[def.Name]; !exists {
				// Check dependencies
				if !pm.dependenciesSatisfied(def) {
					continue
				}

				value, err := pm.promptForParameter(def)
				if err != nil {
					// Optional parameters can be skipped on error
					continue
				}
				pm.values[def.Name] = value
			}
		}
	}

	return nil
}

// promptForParameter prompts for a single parameter
func (pm *PromptManager) promptForParameter(def ParameterDef) (interface{}, error) {
	switch def.Type {
	case TypeString:
		return pm.promptString(def)
	case TypePassword:
		return pm.promptPassword(def)
	case TypeSelect:
		return pm.promptSelect(def)
	case TypeMultiSelect:
		return pm.promptMultiSelect(def)
	case TypeConfirm:
		return pm.promptConfirm(def)
	case TypeInt:
		return pm.promptInt(def)
	default:
		return pm.promptString(def)
	}
}

// promptString prompts for a string value
func (pm *PromptManager) promptString(def ParameterDef) (string, error) {
	prompt := &survey.Input{
		Message: def.DisplayName,
		Help:    def.Description,
	}

	if def.Default != nil {
		prompt.Default = def.Default.(string)
	}

	var result string
	opts := []survey.AskOpt{}

	if def.Validator != nil {
		opts = append(opts, survey.WithValidator(survey.ComposeValidators(
			survey.Required,
			func(val interface{}) error {
				return def.Validator(val)
			},
		)))
	}

	err := survey.AskOne(prompt, &result, opts...)
	return result, err
}

// promptPassword prompts for a password
func (pm *PromptManager) promptPassword(def ParameterDef) (string, error) {
	prompt := &survey.Password{
		Message: def.DisplayName,
		Help:    def.Description,
	}

	var result string
	opts := []survey.AskOpt{}

	if def.Validator != nil {
		opts = append(opts, survey.WithValidator(func(val interface{}) error {
			return def.Validator(val)
		}))
	}

	err := survey.AskOne(prompt, &result, opts...)
	return result, err
}

// promptSelect prompts for a single selection
func (pm *PromptManager) promptSelect(def ParameterDef) (string, error) {
	// Fetch options if fetcher is available
	var options []Option
	var err error

	if def.Fetcher != nil {
		fmt.Printf("Fetching %s options...\n", strings.ToLower(def.DisplayName))
		options, err = def.Fetcher(pm.ctx)
		if err != nil {
			fmt.Printf("⚠ Could not fetch options: %v\n", err)
			// Fall back to manual input
			return pm.promptString(def)
		}
	} else {
		return pm.promptString(def)
	}

	if len(options) == 0 {
		fmt.Printf("No %s available\n", strings.ToLower(def.DisplayName))
		return "", fmt.Errorf("no options available")
	}

	// Create display options
	displayOptions := make([]string, len(options))
	valueMap := make(map[string]string)

	for i, opt := range options {
		display := opt.Display
		if opt.Description != "" {
			display = fmt.Sprintf("%s - %s", opt.Display, opt.Description)
		}
		displayOptions[i] = display
		valueMap[display] = opt.Value
	}

	// Set default if available
	defaultIndex := 0
	if def.Default != nil {
		defaultValue := def.Default.(string)
		for i, opt := range options {
			if opt.Value == defaultValue {
				defaultIndex = i
				break
			}
		}
	}

	prompt := &survey.Select{
		Message: def.DisplayName,
		Options: displayOptions,
		Default: displayOptions[defaultIndex],
	}

	var selected string
	err = survey.AskOne(prompt, &selected)
	if err != nil {
		return "", err
	}

	return valueMap[selected], nil
}

// promptMultiSelect prompts for multiple selections
func (pm *PromptManager) promptMultiSelect(def ParameterDef) ([]string, error) {
	// Fetch options if fetcher is available
	var options []Option
	var err error

	if def.Fetcher != nil {
		fmt.Printf("Fetching %s options...\n", strings.ToLower(def.DisplayName))
		options, err = def.Fetcher(pm.ctx)
		if err != nil {
			fmt.Printf("⚠ Could not fetch options: %v\n", err)
			return []string{}, nil // Optional, so return empty
		}
	} else {
		return []string{}, nil
	}

	if len(options) == 0 {
		fmt.Printf("No %s available\n", strings.ToLower(def.DisplayName))
		return []string{}, nil
	}

	// Create display options
	displayOptions := make([]string, len(options))
	valueMap := make(map[string]string)

	for i, opt := range options {
		display := opt.Display
		if opt.Description != "" {
			display = fmt.Sprintf("%s - %s", opt.Display, opt.Description)
		}
		displayOptions[i] = display
		valueMap[display] = opt.Value
	}

	prompt := &survey.MultiSelect{
		Message: def.DisplayName,
		Options: displayOptions,
	}

	var selected []string
	err = survey.AskOne(prompt, &selected)
	if err != nil {
		return []string{}, err
	}

	// Convert display values back to actual values
	result := make([]string, len(selected))
	for i, s := range selected {
		result[i] = valueMap[s]
	}

	return result, nil
}

// promptConfirm prompts for a yes/no confirmation
func (pm *PromptManager) promptConfirm(def ParameterDef) (bool, error) {
	defaultValue := false
	if def.Default != nil {
		defaultValue = def.Default.(bool)
	}

	prompt := &survey.Confirm{
		Message: def.DisplayName,
		Default: defaultValue,
	}

	var result bool
	err := survey.AskOne(prompt, &result)
	return result, err
}

// promptInt prompts for an integer value
func (pm *PromptManager) promptInt(def ParameterDef) (int, error) {
	prompt := &survey.Input{
		Message: def.DisplayName,
		Help:    def.Description,
	}

	if def.Default != nil {
		prompt.Default = fmt.Sprintf("%d", def.Default.(int))
	}

	var result string
	err := survey.AskOne(prompt, &result)
	if err != nil {
		return 0, err
	}

	var intResult int
	_, err = fmt.Sscanf(result, "%d", &intResult)
	if err != nil {
		return 0, fmt.Errorf("invalid integer value")
	}

	return intResult, nil
}

// dependenciesSatisfied checks if dependencies are met
func (pm *PromptManager) dependenciesSatisfied(def ParameterDef) bool {
	if len(def.DependsOn) == 0 {
		return true
	}

	for _, dep := range def.DependsOn {
		if _, exists := pm.values[dep]; !exists {
			return false
		}
	}
	return true
}

// formatValue formats a value for display
func (pm *PromptManager) formatValue(def ParameterDef, value interface{}) string {
	if def.Type == TypePassword {
		return "********"
	}
	return fmt.Sprintf("%v", value)
}

// GetValues returns all collected values
func (pm *PromptManager) GetValues() map[string]interface{} {
	return pm.values
}

// SetValue sets a specific value
func (pm *PromptManager) SetValue(key string, value interface{}) {
	pm.values[key] = value
}

// CollectValues collects all required values through interactive prompts
func (pm *PromptManager) CollectValues() (map[string]interface{}, error) {
	// Prompt for missing required parameters
	if err := pm.PromptForMissing(); err != nil {
		return nil, err
	}

	// Prompt for optional parameters
	if err := pm.PromptForOptional(); err != nil {
		return nil, err
	}

	return pm.values, nil
}

// ShowSummary displays a summary of the configuration
func (pm *PromptManager) ShowSummary(title string) error {
	if !pm.interactive {
		return nil
	}

	fmt.Println("\n" + strings.Repeat("═", 50))
	fmt.Printf("📋 %s\n", title)
	fmt.Println(strings.Repeat("═", 50))

	for _, def := range pm.definitions {
		if value, exists := pm.values[def.Name]; exists {
			fmt.Printf("%-20s: %v\n", def.DisplayName, pm.formatValue(def, value))
		}
	}

	fmt.Println(strings.Repeat("─", 50))
	return nil
}

// ConfirmExecution asks for final confirmation
func (pm *PromptManager) ConfirmExecution(message string) (bool, error) {
	if !pm.interactive {
		return true, nil // Non-interactive mode always proceeds
	}

	prompt := &survey.Confirm{
		Message: message,
		Default: true,
	}

	var result bool
	err := survey.AskOne(prompt, &result)
	return result, err
}

// Common validators

// ValidateInstanceName validates database instance name
func ValidateInstanceName(val interface{}) error {
	str, ok := val.(string)
	if !ok {
		return fmt.Errorf("value must be a string")
	}

	if len(str) < 4 || len(str) > 50 {
		return fmt.Errorf("name must be between 4 and 50 characters")
	}

	if !regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9-]*$`).MatchString(str) {
		return fmt.Errorf("name must start with a letter and contain only letters, numbers, and hyphens")
	}

	return nil
}

// ValidatePassword validates database password
func ValidatePassword(val interface{}) error {
	str, ok := val.(string)
	if !ok {
		return fmt.Errorf("value must be a string")
	}

	if len(str) < 8 {
		return fmt.Errorf("password must be at least 8 characters")
	}

	// Check for at least one uppercase, lowercase, and number
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(str)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(str)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(str)

	if !hasUpper || !hasLower || !hasNumber {
		return fmt.Errorf("password must contain uppercase, lowercase, and numbers")
	}

	return nil
}

// ValidatePort validates port number
func ValidatePort(val interface{}) error {
	port, ok := val.(int)
	if !ok {
		// Try to convert string to int
		if str, ok := val.(string); ok {
			_, err := fmt.Sscanf(str, "%d", &port)
			if err != nil {
				return fmt.Errorf("invalid port number")
			}
		} else {
			return fmt.Errorf("value must be a number")
		}
	}

	if port < 1024 || port > 65535 {
		return fmt.Errorf("port must be between 1024 and 65535")
	}

	return nil
}

// ValidateUsername validates database username
func ValidateUsername(val interface{}) error {
	str, ok := val.(string)
	if !ok {
		return fmt.Errorf("value must be a string")
	}

	if len(str) < 2 || len(str) > 32 {
		return fmt.Errorf("username must be between 2 and 32 characters")
	}

	if strings.ToLower(str) == "root" {
		return fmt.Errorf("'root' is not allowed as username")
	}

	return nil
}

// ValidateStorageSize validates storage size
func ValidateStorageSize(val interface{}) error {
	var size int
	switch v := val.(type) {
	case int:
		size = v
	case string:
		_, err := fmt.Sscanf(v, "%d", &size)
		if err != nil {
			return fmt.Errorf("invalid storage size")
		}
	default:
		return fmt.Errorf("value must be a number")
	}

	if size < 20 || size > 6000 {
		return fmt.Errorf("storage size must be between 20 and 6000 GB")
	}

	if size%10 != 0 {
		return fmt.Errorf("storage size must be in increments of 10 GB")
	}

	return nil
}

// ValidateBackupPeriod validates backup retention period
func ValidateBackupPeriod(val interface{}) error {
	var period int
	switch v := val.(type) {
	case int:
		period = v
	case string:
		_, err := fmt.Sscanf(v, "%d", &period)
		if err != nil {
			return fmt.Errorf("invalid backup period")
		}
	default:
		return fmt.Errorf("value must be a number")
	}

	if period < 0 || period > 35 {
		return fmt.Errorf("backup period must be between 0 and 35 days")
	}

	return nil
}

// ValidateTimeFormat validates HH:MM time format
func ValidateTimeFormat(val interface{}) error {
	str, ok := val.(string)
	if !ok {
		return fmt.Errorf("value must be a string")
	}

	// Check HH:MM format
	var hour, minute int
	_, err := fmt.Sscanf(str, "%d:%d", &hour, &minute)
	if err != nil {
		return fmt.Errorf("time must be in HH:MM format")
	}

	if hour < 0 || hour > 23 {
		return fmt.Errorf("hour must be between 00 and 23")
	}

	if minute < 0 || minute > 59 {
		return fmt.Errorf("minute must be between 00 and 59")
	}

	return nil
}
