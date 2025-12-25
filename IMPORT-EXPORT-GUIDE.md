# WizardCore Content Creator - Import/Export Guide

## 🚀 Overview

The Import/Export feature allows content creators to:
- **Export** complete pathways with all modules, exercises, and test cases as JSON
- **Import** pathways from JSON files to quickly set up content
- **Share** content templates with other creators
- **Backup** your work before making major changes
- **Migrate** content between different WizardCore instances

## 📁 File Structure

Exported pathways include:
- Pathway metadata (title, description, level, duration, etc.)
- All modules with their metadata
- All exercises with full content and test cases
- Test cases with inputs, expected outputs, and scoring

## 🎯 How to Use

### 1. **Export a Pathway**
1. Go to the Creator Dashboard (`/creator/dashboard`)
2. Find the pathway you want to export
3. Click the **Export** button (📥 icon) on the pathway card
4. The system will:
   - Download a JSON file named `pathway_title_export.json`
   - Show a preview of the exported data
   - Option to copy JSON to clipboard

### 2. **Import a Pathway**
1. Go to the Creator Dashboard
2. Click the **Import Pathway** button in the top-right
3. Either:
   - Drag & drop a JSON file
   - Click to browse and select a file
4. The system will:
   - Validate the JSON structure
   - Show a preview of what will be imported
   - Let you import as a draft

### 3. **Use the Sample Template**
A sample pathway template is available at:
```
/public/templates/sample-pathway.json
```
Or download directly: `/templates/sample-pathway.json`

This template includes:
- A complete "Python for Cybersecurity" pathway
- 2 modules with 3 exercises
- Various test cases (visible and hidden)
- Different difficulty levels
- Example solutions and starter code

## 🔧 Technical Details

### Export Format
```json
{
  "pathway": {
    "title": "Pathway Title",
    "subtitle": "Optional subtitle",
    "description": "Detailed description...",
    "level": "Beginner|Intermediate|Advanced|Expert",
    "duration_weeks": 4,
    "color_gradient": "from-blue-500 to-purple-600",
    "icon": "🎓",
    "prerequisites": ["Basic knowledge"],
    "sort_order": 1,
    "modules": [
      {
        "title": "Module Title",
        "description": "Module description",
        "sort_order": 1,
        "estimated_hours": 10,
        "xp_reward": 100,
        "exercises": [
          {
            "title": "Exercise Title",
            "difficulty": "BEGINNER|INTERMEDIATE|ADVANCED",
            "points": 10,
            "time_limit_minutes": 30,
            "sort_order": 1,
            "objectives": ["Learn something"],
            "content": "Markdown content...",
            "description": "Exercise description",
            "constraints": ["No external libraries"],
            "hints": ["Use string methods"],
            "starter_code": "# Write your solution",
            "solution_code": "print('Hello')",
            "language_id": 71,  // Python 3.8.1
            "tags": ["python", "beginner"],
            "test_cases": [
              {
                "input": "\"test\"",
                "expected_output": "Hello test",
                "is_hidden": false,
                "points": 5,
                "sort_order": 1
              }
            ]
          }
        ]
      }
    ]
  },
  "metadata": {
    "exported_at": "2025-12-25T08:00:00Z",
    "version": "1.0",
    "creator_id": "uuid-here"
  }
}
```

### Supported Languages
- **71**: Python 3.8.1
- **50**: C GCC 9.2.0
- **54**: C++ GCC 9.2.0
- **62**: Java OpenJDK 13
- **63**: JavaScript Node 12
- **82**: SQL SQLite 3.27

### Validation Rules
When importing, the system checks:
- Required fields (title, level, duration)
- Valid language IDs
- At least one module per pathway
- At least one test case per exercise
- No duplicate sort orders within modules/exercises
- Valid difficulty levels (BEGINNER, INTERMEDIATE, ADVANCED)
- Valid status (draft, published)

## 💡 Use Cases

### For Content Creators
1. **Backup your work**: Export before making major changes
2. **Template library**: Create reusable exercise templates
3. **Collaboration**: Share pathways with other creators
4. **Version control**: Track changes by exporting versions

### For Administrators
1. **Onboarding**: Provide new creators with starter content
2. **Content standardization**: Distribute approved templates
3. **Migration**: Move content between environments
4. **Backup strategy**: Regular exports for disaster recovery

### For Educators
1. **Curriculum sharing**: Share complete learning paths
2. **Exercise banks**: Build libraries of coding challenges
3. **Adaptation**: Import and modify existing content
4. **Assessment**: Standardized test cases across classes

## 🛠️ API Endpoints

### Export a Pathway
```
GET /api/v1/content-creator/pathways/{id}/export
```
- Requires authentication
- User must own the pathway
- Returns complete pathway JSON

### Import a Pathway
```
POST /api/v1/content-creator/pathways/import
```
- Requires authentication
- User must be a content creator
- Accepts JSON in request body
- Returns created pathway with ID

## 🚨 Troubleshooting

### Common Issues

1. **"Invalid JSON file"**
   - Ensure file is valid JSON
   - Check for trailing commas
   - Validate with a JSON validator

2. **"Missing required fields"**
   - Check that pathway has title, level, duration
   - Ensure modules have titles
   - Ensure exercises have titles and language_id

3. **"Duplicate sort order"**
   - Each module must have unique sort_order
   - Each exercise within a module must have unique sort_order

4. **"No test cases"**
   - Every exercise must have at least one test case
   - Test cases array cannot be empty

### Testing Your JSON
1. Use the sample template as reference
2. Validate with: `python -m json.tool your-file.json`
3. Test import with small files first
4. Check browser console for detailed errors

## 📚 Best Practices

1. **Start with the sample template** - Use it as a reference
2. **Export before editing** - Always have a backup
3. **Use meaningful filenames** - Include pathway title and date
4. **Version your exports** - Add version numbers or dates
5. **Test imports** - Import your exports to verify they work
6. **Keep exercises modular** - Each exercise should be self-contained
7. **Use descriptive test cases** - Clear expected outputs help debugging
8. **Include hidden test cases** - For proper assessment

## 🔄 Migration Workflow

1. **Export from source system**
   ```bash
   # Use API or dashboard export
   GET /pathways/{id}/export
   ```

2. **Validate exported file**
   ```bash
   python -m json.tool exported.json > validated.json
   ```

3. **Import to target system**
   ```bash
   # Use API or dashboard import
   POST /pathways/import
   ```

4. **Verify imported content**
   - Check all modules and exercises
   - Test exercise solutions
   - Verify test cases work

## 🎨 Custom Templates

Create your own templates by:
1. Exporting a well-structured pathway
2. Modifying the JSON (change titles, content, etc.)
3. Saving as a template file
4. Sharing with other creators

Example template structure:
```
templates/
├── python-basics.json
├── web-security.json
├── cryptography.json
└── network-pentesting.json
```

## 📞 Support

If you encounter issues:
1. Check the validation errors in the import preview
2. Compare with the sample template
3. Ensure your JSON matches the expected format
4. Contact system administrator for API issues

---

**Happy Creating!** 🚀

With the import/export feature, you can now easily share, backup, and migrate your content creation work across different WizardCore instances.