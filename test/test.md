# Test Mermaid File

This is a test file with various issues.

## Test 1: Newline literal

```mermaid
graph LR
    A[Hello\nWorld] --> B[Another node]
```

## Test 2: Unquoted text with special chars

```mermaid
graph LR
    A[Function(foo)] --> B[Name: Value]
    B --> C[Data {key: value}]
```

## Test 3: Duplicate node

```mermaid
graph TD
    A[First] --> B[Second]
    A[Duplicate A] --> C[Third]
```

## Test 4: Undefined class

```mermaid
graph LR
    classDef red fill:#ff0000;
    A[Red node] --> B[Blue node]
    class A red;
    class B blue;
```

## Test 5: Invalid style

```mermaid
graph LR
    A[Test] --> B[Test]
    style A filld:red;
    style B stroke-width:2px;
```

## Test 6: Isolated node

```mermaid
graph LR
    A[Connected] --> B[Connected]
    C[Isolated]
```

## Test 7: HTML tags

```mermaid
graph LR
    A[<div>Hello</div>] --> B[World]
```

