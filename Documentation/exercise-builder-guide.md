# Exercise Builder with Judge0 Integration - User Guide

## Overview

The Exercise Builder is a comprehensive tool for content creators to build interactive coding exercises with live Judge0 testing. It provides a complete workflow from writing the exercise description to testing solutions against multiple test cases.

## Features

### ✅ Comprehensive Exercise Creation
- **Basic Information**: Title, description, difficulty, points, time limits
- **Rich Content**: Markdown-based exercise content
- **Learning Objectives**: Define what students will learn
- **Constraints**: Specify technical constraints
- **Progressive Hints**: Help students when they're stuck
- **Tagging**: Organize exercises with tags

### ✅ Live Code Editor (Monaco)
- **Syntax Highlighting**: Full Monaco editor with IntelliSense
- **Multiple Languages**: Python, C, C++, Java, JavaScript, SQL
- **Starter Code**: Provide initial code for students
- **Solution Code**: Your reference solution

### ✅ Integrated Judge0 Testing
- **Live Testing**: Test your solution against test cases in real-time
- **Single Test Execution**: Run individual test cases
- **Batch Testing**: Test all cases at once
- **Detailed Results**: See execution time, output, errors
- **Pass/Fail Indicators**: Clear visual feedback

### ✅ Advanced Test Case Manager
- **Visible/Hidden Tests**: Control what students see
- **Point Allocation**: Assign points per test case
- **Reorderable**: Drag to reorder test cases
- **Inline Testing**: Test each case individually
- **Rich Feedback**: See expected vs actual output

## How to Use

### Step 1: Navigate to Exercise Builder

```
/creator/exercises/new?module_id=YOUR_MODULE_ID
```

You must provide a `module_id` query parameter.

### Step 2: Fill in Exercise Details

#### Basic Information Tab
1. **Title**: Give your exercise a clear, descriptive title
   - Example: "Stack Buffer Overflow Exploitation"

2. **Description**: Short summary of the exercise
   - Example: "Learn to exploit a simple stack buffer overflow vulnerability"

3. **Difficulty**: Choose from Beginner, Intermediate, or Advanced

4. **Points**: Reward value (typically 100-1000)

5. **Time Limit**: How long students have (in minutes)

6. **Language**: Select the programming language (Python, C, C++, Java, JS, SQL)

#### Exercise Content
Write detailed instructions in Markdown:

```markdown
# Buffer Overflow Challenge

## Problem Description
You are given a vulnerable C program that uses `strcpy()` without bounds checking.
Your task is to write a Python exploit script that overwrites the return address.

## Examples
Input: (none)
Expected Output: "Exploit successful!"

## Tips
- The buffer is 64 bytes
- Little-endian architecture
- Use struct.pack() for addresses
```

#### Learning Objectives
Define what students will learn:
- "Understand stack memory layout"
- "Identify buffer overflow vulnerabilities"
- "Calculate offsets to return address"

#### Constraints
Technical limitations:
- "Must use Python 3"
- "No external libraries"
- "Input size <= 1024 bytes"

#### Hints (Progressive)
Students can unlock hints one at a time:
1. "Look for the strcpy() call"
2. "The return address is at buffer + 72"
3. "Use 0xdeadbeef as your test address"

#### Tags
Organize your exercise:
- buffer-overflow
- stack
- exploitation
- memory-corruption

### Step 3: Write Code

#### Starter Code
Provide initial code for students:

```python
#!/usr/bin/env python3
import struct

# TODO: Write your exploit here
payload = b""

print(payload.decode('latin-1'))
```

#### Solution Code
Your working solution:

```python
#!/usr/bin/env python3
import struct

# Overflow buffer (64 bytes) + saved EBP (8 bytes) = 72 bytes
padding = b"A" * 72

# Return address (0xdeadbeef in little-endian)
ret_addr = struct.pack("<Q", 0xdeadbeef)

payload = padding + ret_addr
print(payload.decode('latin-1'))
```

### Step 4: Create Test Cases

Click "Add Test Case" to create test cases:

#### Example Test Case 1 (Visible)
- **Input**: (leave empty for no input)
- **Expected Output**: 
  ```
  Exploit successful!
  ```
- **Points**: 100
- **Visibility**: Visible (students can see this)

#### Example Test Case 2 (Hidden)
- **Input**: (empty)
- **Expected Output**:
  ```
  Shell spawned!
  ```
- **Points**: 50
- **Visibility**: Hidden (students don't see this)

### Step 5: Test Your Solution

Click **"Test Solution"** to run your solution code against all test cases:

```
✓ All Tests Passed! (2/2)

Test Case #1: ✓ Passed
Test Case #2: ✓ Passed
```

Or test individual cases by clicking the play button on each test case.

### Step 6: Save

Click **"Save Exercise"** to save as a draft. The exercise will be:
- Saved with `status: draft`
- Visible only to you
- Ready for further editing

Later you can submit it for review to publish it to students.

## Judge0 Testing Features

### What Gets Tested
- **Source Code**: Your solution code from the Code tab
- **Language**: The selected language (e.g., Python 71)
- **Input**: The stdin from each test case
- **Expected Output**: Compared character-by-character (with trim)

### Test Results
For each test case, you'll see:
- ✓/✗ Pass/Fail indicator
- Execution time (in milliseconds)
- Actual output vs expected output
- Compilation errors (if any)
- Runtime errors (if any)

### Example Success Result
```
✓ Test Passed
Execution time: 23ms
```

### Example Failure Result
```
✗ Test Failed

Expected:
Hello, World!

Got:
Hello World

Error: Missing comma in output
```

## Best Practices

### 1. Write Clear Instructions
- Use Markdown formatting for readability
- Include examples with expected input/output
- Explain the problem context (why is this useful?)

### 2. Test Thoroughly
- Create at least 3-5 test cases
- Include edge cases (empty input, large input, special characters)
- Mix visible and hidden test cases
- Test your solution before saving!

### 3. Progressive Difficulty
- Start with a simple visible test case
- Add more complex hidden test cases
- Award more points for harder test cases

### 4. Provide Good Starter Code
- Give students a solid foundation
- Include helpful comments
- Show the expected structure

### 5. Use Hints Wisely
- Make hint #1 gentle (general direction)
- Make hint #2 more specific (where to look)
- Make hint #3 very specific (almost giving it away)

## API Integration

The Exercise Builder connects to these backend endpoints:

```
POST /api/v1/content-creator/exercises
```

**Payload**:
```json
{
  "module_id": "uuid",
  "title": "Exercise Title",
  "difficulty": "BEGINNER",
  "points": 100,
  "time_limit_minutes": 30,
  "objectives": ["Learn X", "Practice Y"],
  "content": "# Markdown content",
  "description": "Short description",
  "constraints": ["Constraint 1"],
  "hints": ["Hint 1", "Hint 2"],
  "starter_code": "# Starter",
  "solution_code": "# Solution",
  "language_id": 71,
  "tags": ["tag1", "tag2"],
  "status": "draft",
  "test_cases": [
    {
      "input": "",
      "expected_output": "Hello!",
      "is_hidden": false,
      "points": 100,
      "sort_order": 0
    }
  ]
}
```

## Judge0 Language IDs

| Language | ID | Monaco Language |
|----------|-----|----------------|
| Python 3.8.1 | 71 | python |
| C (GCC 9.2.0) | 50 | c |
| C++ (GCC 9.2.0) | 54 | cpp |
| Java (OpenJDK 13) | 62 | java |
| JavaScript (Node 12) | 63 | javascript |
| SQL (SQLite 3.27) | 82 | sql |

## Troubleshooting

### Test Always Fails
- Check for extra whitespace in expected output
- Verify newlines are included/excluded correctly
- Test the solution code manually in Judge0
- Check if the language ID matches your code

### Code Doesn't Execute
- Verify Judge0 is running (`http://localhost:2358/about`)
- Check for syntax errors in solution code
- Verify the language ID is correct
- Check Judge0 logs for errors

### Slow Testing
- Judge0 may take 1-3 seconds per test case
- Batch testing runs all cases sequentially
- Be patient with multiple test cases

### Save Fails
- Check that module_id is valid
- Verify you have content_creator role
- Ensure at least one test case exists
- Check network tab for error details

## Keyboard Shortcuts

In Monaco Editor:
- `Ctrl/Cmd + S`: Save (will trigger exercise save)
- `Ctrl/Cmd + F`: Find
- `Ctrl/Cmd + H`: Find and replace
- `Ctrl/Cmd + /`: Toggle comment
- `Alt + Up/Down`: Move line up/down

## Next Steps

After creating exercises, you can:
1. **Preview**: See how students will experience it
2. **Submit for Review**: Send to admin for approval
3. **Publish**: Make it available to students
4. **Analytics**: Track student performance

## Support

For issues or questions:
- Check the console for error messages
- Verify backend is running
- Test Judge0 connectivity
- Review the API documentation

---

**Happy Exercise Building!** 🎓💻
