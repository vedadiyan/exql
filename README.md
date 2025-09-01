# EXQL - Expression Query Language

A powerful, extensible expression language for Go applications that supports dynamic evaluation of complex expressions with variables, functions, and built-in libraries.

## Features

- **Dynamic Expression Evaluation** - Parse and evaluate expressions at runtime
- **Rich Type System** - Support for numbers, strings, booleans, lists, and maps
- **Variable Resolution** - Access variables and nested object properties
- **Function Calls** - Built-in and custom function support
- **Comprehensive Operators** - Arithmetic, comparison, logical, and membership operators
- **Built-in Libraries** - Extensive library collection for common operations
- **Field and Index Access** - Navigate complex data structures with ease
- **Custom Contexts** - Flexible context system for variable and function injection

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/vedadiyan/exql"
    "github.com/vedadiyan/exql/lang"
)

func main() {
    // Create a context with built-in libraries
    ctx := exql.NewDefaultContext(exql.WithBuiltInLibrary())
    
    // Set some variables
    ctx.SetVariable("user", lang.MapValue{
        "name": lang.StringValue("Alice"),
        "age":  lang.NumberValue(30),
        "active": lang.BoolValue(true),
    })
    
    // Evaluate expressions
    result, err := exql.Eval("user.name", ctx)
    if err != nil {
        panic(err)
    }
    fmt.Println(result) // Output: Alice
    
    // Complex expression
    result, err = exql.Eval("user.active and user.age > 18", ctx)
    if err != nil {
        panic(err)
    }
    fmt.Println(result) // Output: true
}
```

## Syntax

### Literals

```javascript
42          // Numbers
3.14        // Floating point
'hello'     // Single-quoted strings
"world"     // Double-quoted strings
true        // Boolean true
false       // Boolean false
[1, 2, 3]   // Lists
```

### Variables and Access

```javascript
user                    // Variable access
user.name              // Field access
user['name']           // Index access with string
items[0]               // Index access with number
users.name             // Field access on lists (maps to all elements)
```

### Operators

#### Arithmetic
```javascript
3 + 5       // Addition
10 - 3      // Subtraction
4 * 6       // Multiplication
15 / 3      // Division
```

#### Comparison
```javascript
x == 5      // Equality
x != 5      // Inequality
x < 10      // Less than
x <= 10     // Less than or equal
x > 5       // Greater than
x >= 5      // Greater than or equal
```

#### Logical
```javascript
true and false    // Logical AND
true or false     // Logical OR
not true          // Logical NOT
```

#### Membership
```javascript
'apple' in fruits        // Contains check
'grape' not in fruits    // Does not contain
```

### Function Calls

```javascript
func(arg1, arg2)         // Function call
string.upper('hello')    // Library function call
util.coalesce(a, b, c)   // Multiple arguments
```

## Built-in Libraries

EXQL comes with comprehensive built-in libraries:

- **string** - String manipulation functions
- **util** - Utility functions (conditionals, type conversion, etc.)
- **time** - Date and time operations
- **json** - JSON parsing and manipulation
- **list** - List/array operations
- **map** - Map/object operations
- **url** - URL parsing and manipulation
- **http** - HTTP utilities
- **crypt** - Cryptographic functions
- **ip** - IP address utilities

### Library Usage Examples

```go
ctx := exql.NewDefaultContext(exql.WithBuiltInLibrary())

// String operations
result, _ := exql.Eval("string.upper('hello')", ctx)  // "HELLO"

// Utility functions
result, _ := exql.Eval("util.coalesce(null, 'default')", ctx)  // "default"

// Time operations
result, _ := exql.Eval("time.now()", ctx)  // Current timestamp
```

## Custom Functions

Add your own functions to the context:

```go
ctx := exql.NewDefaultContext()

// Add custom function
ctx.SetFunction("double", func(args []lang.Value) (lang.Value, error) {
    if len(args) != 1 {
        return nil, fmt.Errorf("double expects 1 argument")
    }
    
    num := toNumber(args[0])
    return lang.NumberValue(num * 2), nil
})

// Use the function
result, _ := exql.Eval("double(21)", ctx)  // 42
```

## Advanced Usage

### Parse and Evaluate Separately

For better performance when evaluating the same expression multiple times:

```go
// Parse once
ast, err := exql.Parse("user.age > threshold")
if err != nil {
    panic(err)
}

// Evaluate multiple times with different contexts
ctx1 := exql.NewDefaultContext()
ctx1.SetVariable("user", userObject1)
ctx1.SetVariable("threshold", lang.NumberValue(18))
result1, _ := ast.Evaluate(ctx1)

ctx2 := exql.NewDefaultContext()
ctx2.SetVariable("user", userObject2)
ctx2.SetVariable("threshold", lang.NumberValue(21))
result2, _ := ast.Evaluate(ctx2)
```

### Context Configuration

```go
// Empty context
ctx := exql.NewDefaultContext()

// With built-in libraries
ctx := exql.NewDefaultContext(exql.WithBuiltInLibrary())

// With custom functions
funcs := map[string]lang.Function{
    "myFunc": myFunction,
}
ctx := exql.NewDefaultContext(exql.WithFunctions(funcs))

// Combined
ctx := exql.NewDefaultContext(
    exql.WithBuiltInLibrary(),
    exql.WithFunctions(funcs),
)
```

### Real-World Example: User Authorization

```go
func checkAuthorization(user User, resource Resource, action string) (bool, error) {
    ctx := exql.NewDefaultContext(exql.WithBuiltInLibrary())
    
    // Set up context variables
    ctx.SetVariable("user", lang.MapValue{
        "id":     lang.NumberValue(float64(user.ID)),
        "role":   lang.StringValue(user.Role),
        "active": lang.BoolValue(user.Active),
        "permissions": convertToListValue(user.Permissions),
    })
    
    ctx.SetVariable("resource", lang.MapValue{
        "type":     lang.StringValue(resource.Type),
        "owner_id": lang.NumberValue(float64(resource.OwnerID)),
        "public":   lang.BoolValue(resource.Public),
    })
    
    ctx.SetVariable("action", lang.StringValue(action))
    
    // Define authorization rule
    rule := `
        user.active and (
            user.role == 'admin' or 
            user.id == resource.owner_id or
            (resource.public and action in user.permissions)
        )
    `
    
    result, err := exql.Eval(rule, ctx)
    if err != nil {
        return false, err
    }
    
    if boolResult, ok := result.(lang.BoolValue); ok {
        return bool(boolResult), nil
    }
    
    return false, nil
}
```

## Type System

EXQL supports the following types:

- `lang.NumberValue` - Floating point numbers
- `lang.StringValue` - UTF-8 strings
- `lang.BoolValue` - Boolean true/false
- `lang.ListValue` - Ordered collections
- `lang.MapValue` - Key-value maps
- `nil` - Null/undefined values

### Type Conversion

EXQL performs automatic type conversion in many contexts:

```go
// Numbers
ToNumber("42")     // 42
ToNumber(true)     // 1
ToNumber(false)    // 0

// Booleans
ToBool(0)          // false
ToBool(42)         // true
ToBool("")         // false
ToBool("hello")    // true
```

## Error Handling

EXQL provides detailed error information for both parsing and evaluation:

```go
result, err := exql.Eval("invalid + + syntax", ctx)
if err != nil {
    fmt.Printf("Error: %v\n", err)
    // Includes position information and context
}
```

## Performance Considerations

- Use `Parse` once and `Evaluate` multiple times for repeated expressions
- Built-in libraries are loaded on-demand
- Context creation is lightweight
- Expression compilation is cached internally

## Testing

The project includes comprehensive test suites:

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run benchmarks
go test -bench=. ./...

# Integration tests
go test ./exql
```

## Contributing

Contributions are welcome! Please ensure:

1. All tests pass
2. New features include tests
3. Documentation is updated
4. Code follows Go conventions

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.
