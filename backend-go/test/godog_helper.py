import re
import subprocess
import argparse
import os

def normalize_go_indentation(content):
	"""Ensure Go code uses tabs for indentation."""
	# Replace spaces at the beginning of lines with tabs
	# This regex matches spaces (4 or 2) at the start of lines and replaces with tabs
	content = re.sub(r'^( {4})', r'\t', content, flags=re.MULTILINE)
	content = re.sub(r'^( {2})', r'\t', content, flags=re.MULTILINE)
	# Handle nested indentation (2 levels)
	content = re.sub(r'^(\t {4})', r'\t\t', content, flags=re.MULTILINE)
	content = re.sub(r'^(\t {2})', r'\t\t', content, flags=re.MULTILINE)
	# Handle nested indentation (3 levels)
	content = re.sub(r'^(\t\t {4})', r'\t\t\t', content, flags=re.MULTILINE)
	content = re.sub(r'^(\t\t {2})', r'\t\t\t', content, flags=re.MULTILINE)
	return content

def clear_godog_test_file(filepath):
	"""Clean up a godog test file by removing undefined functions and clearing the initializer."""
	with open(filepath, 'r') as file:
		content = file.read()

	# Extract initializer name
	initializer_match = re.search(r'func \w+\(t \*testing\.T\) {\n\tsuite := godog\.TestSuite{\n\t\tScenarioInitializer: (\w+),', 
							   content, re.MULTILINE)
	if not initializer_match:
		raise ValueError("Could not find initializer in test file")
	
	initializer_name = initializer_match.group(1)

	# Extract undefined functions before removing them
	undef_func_names = re.findall(r'^func (\w+)\(.*?\) error {\n\treturn godog\.ErrPending\n}\n', content, re.MULTILINE)
	
	# Remove undefined functions
	content = re.sub(r'^func \w+\(.*?\) error {\n\treturn godog\.ErrPending\n}\n', '', content, flags=re.MULTILINE)

	# Clear InitializeScenario function
	pattern = rf'^func {re.escape(initializer_name)}\(\w+ \*godog\.ScenarioContext\) {{\n(.+\n)+?^}}$'
	content = re.sub(pattern, f'func {initializer_name}(ctx *godog.ScenarioContext) {{\n}}', content, flags=re.MULTILINE)

	# Normalize indentation to use tabs
	content = normalize_go_indentation(content)

	# Write updated content back to the file
	with open(filepath, 'w') as file:
		file.write(content)

	return undef_func_names

def get_new_functions(filepath):
	"""Run Go tests and extract function signatures and steps."""
	try:
		result = subprocess.run(['go', 'test', filepath, '-v'], capture_output=True, text=True)
		output = result.stdout if result.stdout else result.stderr

		# Find function signatures and bodies
		funcs = re.findall(r'(func (\w+)\(.*?\) error {\n\treturn godog\.ErrPending\n}\n)', output)
		sorted_funcs = sorted(funcs, key=lambda x: x[1])  # Sort by function name

		# Extract step definitions
		steps = re.findall(r'^\tctx.Step\(.+\)$', output, re.MULTILINE)
		steps.sort()

		return ([body for body, _ in sorted_funcs], steps)
	except Exception as e:
		print(f"Error running tests: {e}")
		return [], []

def extract_all_functions(content):
	"""Extract all functions from the file content."""
	# Get the main test function
	test_func_pattern = r'^(func (\w+)\(t \*testing\.T\)[\s\S]+?^})$'
	test_func_match = re.search(test_func_pattern, content, re.MULTILINE)
	if not test_func_match:
		raise ValueError("Test function not found")
	
	test_func = test_func_match.group(1)
	test_func_name = test_func_match.group(2)
	
	# Find the initializer function name
	init_name_match = re.search(r'ScenarioInitializer: (\w+),', content)
	if not init_name_match:
		raise ValueError("ScenarioInitializer not found")
	
	initializer_name = init_name_match.group(1)
	
	# Extract initializer function
	init_func_pattern = rf'^(func {re.escape(initializer_name)}\(\w+ \*godog\.ScenarioContext\)[\s\S]+?^}})$'
	init_func_match = re.search(init_func_pattern, content, re.MULTILINE)
	if not init_func_match:
		# If the initializer is empty, it might just be a simple function
		init_func_pattern = rf'^(func {re.escape(initializer_name)}\(\w+ \*godog\.ScenarioContext\) {{[\s\S]*?^}})$'
		init_func_match = re.search(init_func_pattern, content, re.MULTILINE)
		if not init_func_match:
			raise ValueError("Initializer function not found")
	
	init_func = init_func_match.group(1)
	
	# Extract all other functions (step implementations)
	# Remove test function and initializer from content to avoid matching them again
	remaining_content = content
	remaining_content = re.sub(re.escape(test_func), '', remaining_content)
	remaining_content = re.sub(re.escape(init_func), '', remaining_content)
	
	# Match step implementation functions
	step_func_pattern = r'^(func (\w+)\(.*?\) error {[\s\S]+?^})$'
	step_funcs = re.findall(step_func_pattern, remaining_content, re.MULTILINE)
	
	# Sort step functions by name
	sorted_step_funcs = sorted(step_funcs, key=lambda x: x[1])
	
	return test_func, init_func, [func_body for func_body, _ in sorted_step_funcs]

def write_functions_to_test_file(filepath, new_functions=None, steps=None):
	"""Update the test file with sorted functions and steps."""
	with open(filepath, 'r') as file:
		content = file.read()

	# Extract existing functions and sort them
	test_func, init_func, step_funcs = extract_all_functions(content)
	
	# Add new functions from test run if provided
	if new_functions:
		step_funcs.extend(new_functions)
		# Remove duplicates (keeping the first occurrence)
		seen = set()
		unique_funcs = []
		for func in step_funcs:
			func_name = re.search(r'func (\w+)\(', func).group(1)
			if func_name not in seen:
				seen.add(func_name)
				unique_funcs.append(func)
		step_funcs = sorted(unique_funcs, key=lambda x: re.search(r'func (\w+)\(', x).group(1))
	
	# Extract package and import statements
	package_match = re.search(r'^package .*$', content, re.MULTILINE)
	imports_match = re.search(r'^import \([\s\S]+?\)$', content, re.MULTILINE)
	
	if not package_match:
		raise ValueError("Package declaration not found")
	
	package_stmt = package_match.group(0)
	imports_stmt = imports_match.group(0) if imports_match else ""
	
	# Build new file content with proper ordering:
	# 1. Package and imports
	# 2. Test function (first)
	# 3. Step implementation functions (sorted)
	# 4. Initializer function (last)
	
	new_content = f"{package_stmt}\n\n"
	if imports_stmt:
		new_content += f"{imports_stmt}\n\n"
	
	# Add test function first (after package and imports)
	new_content += f"{test_func}\n\n"
	
	# Add all step functions next
	for func in step_funcs:
		new_content += f"{func}\n\n"
	
	# Update steps in initializer function if provided
	if steps:
		init_body_pattern = r'(func \w+\(\w+ \*godog\.ScenarioContext\) {)(\n.*?)?(^\})'
		init_match = re.search(init_body_pattern, init_func, re.DOTALL | re.MULTILINE)
		if init_match:
			before = init_match.group(1)
			after = init_match.group(3)
			steps_code = '\n\t' + '\n\t'.join(steps) if steps else ''
			init_func = before + steps_code + '\n' + after
	
	# Add initializer last
	new_content += f"{init_func}\n"
	
	# Ensure consistent tab indentation
	new_content = normalize_go_indentation(new_content)
	
	# Write the sorted file
	with open(filepath, 'w') as file:
		file.write(new_content)

def create_new_test_file(filepath, feature_name=None, feature_path=None):
	"""Create a new Godog test file with a basic structure."""
	if feature_name is None:
		# Extract a default test name from the filepath
		base_name = filepath.split("/")[-1].replace("_test.go", "").replace(".go", "")
		if base_name:
			feature_name = "".join(word.capitalize() for word in base_name.split("_"))
		else:
			feature_name = "Feature"
	
	# Convert feature_name to snake_case for directory paths
	snake_case_name = ''.join(['_'+c.lower() if c.isupper() else c for c in feature_name]).lstrip('_')
	
	# Determine feature path
	if feature_path is None:
		# Try to use feature-specific directory if it exists
		specific_path = f"../features/{snake_case_name}/"
		if os.path.exists(specific_path):
			feature_path = specific_path
		else:
			feature_path = "../features"
	
	# Create template content for the new test file
	content = f"""package test

import (
	"os"
	"testing"

	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
)

func Test{feature_name}(t *testing.T) {{
	suite := godog.TestSuite{{
		ScenarioInitializer: InitializeScenario{feature_name},
		Options: &godog.Options{{
			Format:   "pretty",
			Paths:    []string{{"{feature_path}"}},
			Output:   colors.Colored(os.Stdout),
			TestingT: t,
			Strict:   true,
		}},
	}}

	if suite.Run() != 0 {{
		t.Fatal("non-zero status returned, failed to run feature tests")
	}}
}}

func InitializeScenario{feature_name}(ctx *godog.ScenarioContext) {{
}}
"""
	
	# Write the new file with tab indentation
	with open(filepath, 'w') as file:
		file.write(content)  # The template already uses tabs
	
	print(f"Created new test file: {filepath}")
	print(f"Test suite name: Test{feature_name}")
	print(f"Feature path: {feature_path}")
	return True

def normalize_feature_name(name):
    """Normalize a feature name to CamelCase."""
    # Remove any non-alphanumeric characters and replace with spaces
    cleaned = re.sub(r'[^a-zA-Z0-9]', ' ', name)
    # Split by spaces and capitalize each word
    return ''.join(word.capitalize() for word in cleaned.split())

def get_snake_case(name):
    """Convert a CamelCase name to snake_case."""
    # Add underscore before capital letters and convert to lowercase
    return re.sub(r'(?<!^)(?=[A-Z])', '_', name).lower()

def ensure_feature_folder(feature_name):
    """Ensure the feature folder exists, create it if it doesn't."""
    snake_name = get_snake_case(feature_name)
    feature_folder = f"../features/{snake_name}"
    
    if not os.path.exists(feature_folder):
        os.makedirs(feature_folder)
        # Create a basic .feature file as a placeholder
        feature_file_path = f"{feature_folder}/{snake_name}.feature"
        with open(feature_file_path, 'w') as f:
            f.write(f"Feature: {feature_name}\n\n")
        print(f"Created feature folder and file: {feature_folder}")
        return feature_folder
    
    print(f"Feature folder already exists: {feature_folder}")
    return feature_folder

def validate_args(args):
    """Validate command line arguments."""
    # Check if feature_name is provided and not empty
    if not args.feature_name or args.feature_name.strip() == "":
        return False, "Feature name cannot be empty"
    
    # Check feature path if provided
    if args.feature_path and not os.path.exists(os.path.dirname(args.feature_path)):
        return False, f"Parent directory for feature path does not exist: {os.path.dirname(args.feature_path)}"
        
    return True, ""

def get_test_file_path(feature_name):
    """Generate the test file path based on the feature name."""
    snake_name = get_snake_case(feature_name)
    return f"./{snake_name}_test.go"

def main():
    # Set up argument parser with detailed description
    parser = argparse.ArgumentParser(
        description='Godog test helper: A tool to create and manage BDD test files for Go applications.',
        epilog="""
Examples:
  # Create a new test file and feature folder for "UserRegistration"
  python godog_helper.py -n UserRegistration
  
  # Update steps for existing "UserRegistration" feature
  python godog_helper.py UserRegistration
  
  # Create a new test with custom feature path
  python godog_helper.py -n UserLogin -p "../features/auth"
        """,
        formatter_class=argparse.RawDescriptionHelpFormatter
    )
    
    parser.add_argument('feature_name', 
                        help='Name of the feature to process (will be normalized to CamelCase)')
    parser.add_argument('-n', '--new', action='store_true',
                        help='Create a new test file and feature folder if needed')
    parser.add_argument('-p', '--feature-path', 
                        help='Override default feature path (default: "../features/feature_name")')
    args = parser.parse_args()
    
    # Validate arguments
    valid, error_message = validate_args(args)
    if not valid:
        print(f"Error: {error_message}")
        parser.print_help()
        return 1
    
    # Normalize the feature name
    feature_name = normalize_feature_name(args.feature_name)
    
    # Determine the test file path
    godog_test_path = get_test_file_path(feature_name)
    
    # Handle the new flag
    if args.new:
        # Ensure feature folder exists
        feature_folder = ensure_feature_folder(feature_name)
        feature_path = args.feature_path or f"../features/{get_snake_case(feature_name)}"
        
        # Create the test file
        create_new_test_file(godog_test_path, feature_name, feature_path)
        
        # If the feature folder already existed, update steps
        if os.path.exists(feature_path):
            print("Feature folder exists, updating steps...")
            try:
                functions, steps = get_new_functions(godog_test_path)
                write_functions_to_test_file(godog_test_path, functions, steps)
                print("Steps updated successfully!")
            except Exception as e:
                print(f"Failed to update steps: {e}")
        
        return
    
    # Process the existing godog test file
    if not os.path.exists(godog_test_path):
        print(f"Test file not found: {godog_test_path}")
        print(f"Use -n flag to create a new test file for {feature_name}")
        return
    
    print("Clearing undefined functions...")
    undef_function_names = clear_godog_test_file(godog_test_path)
    if undef_function_names:
        print('Cleared functions:')
        for name in undef_function_names:
            print(f"- {name}")
    
    print("Getting new function definitions...")
    functions, steps = get_new_functions(godog_test_path)
    
    print("Updating test file with sorted functions...")
    write_functions_to_test_file(godog_test_path, functions, steps)
    print("Done!")

if __name__ == "__main__":
    main()
