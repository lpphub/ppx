# Restructure New Command Templates Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Restructure the `ppx new` command to generate Go projects with standard `cmd/` and `internal/` directory layout.

**Architecture:** Rename and reorganize template files, update all import paths in templates, modify generator code to create new directory structure and map templates to new output paths.

**Tech Stack:** Go, text/template, embed.FS

---

## Task 1: Create New Template Directory Structure

**Files:**
- Create: `generator/templates/cmd/api/main.go.tmpl` (copy from existing main.go.tmpl)

**Step 1: Create cmd/api template directory**

Run: `mkdir -p generator/templates/cmd/api`

**Step 2: Copy main.go.tmpl to new location**

Run: `cp generator/templates/main.go.tmpl generator/templates/cmd/api/main.go.tmpl`

**Step 3: Verify copy succeeded**

Run: `ls -la generator/templates/cmd/api/`
Expected: `main.go.tmpl` file exists

---

## Task 2: Rename Module Templates (init.go → module.go)

**Files:**
- Create: `generator/templates/module/user/module.go.tmpl`
- Create: `generator/templates/module/auth/module.go.tmpl`
- Create: `generator/templates/module/post/module.go.tmpl`
- Delete: `generator/templates/module/user/init.go.tmpl`
- Delete: `generator/templates/module/auth/init.go.tmpl`
- Delete: `generator/templates/module/post/init.go.tmpl`

**Step 1: Rename user/init.go.tmpl**

Run: `mv generator/templates/module/user/init.go.tmpl generator/templates/module/user/module.go.tmpl`

**Step 2: Rename auth/init.go.tmpl**

Run: `mv generator/templates/module/auth/init.go.tmpl generator/templates/module/auth/module.go.tmpl`

**Step 3: Rename post/init.go.tmpl**

Run: `mv generator/templates/module/post/init.go.tmpl generator/templates/module/post/module.go.tmpl`

**Step 4: Verify renames**

Run: `ls generator/templates/module/*/module.go.tmpl`
Expected: Three files listed

---

## Task 3: Rename Repository Templates (repository.go → repo.go)

**Files:**
- Create: `generator/templates/module/user/repo.go.tmpl`
- Create: `generator/templates/module/post/repo.go.tmpl`
- Delete: `generator/templates/module/user/repository.go.tmpl`
- Delete: `generator/templates/module/post/repository.go.tmpl`

**Step 1: Rename user/repository.go.tmpl**

Run: `mv generator/templates/module/user/repository.go.tmpl generator/templates/module/user/repo.go.tmpl`

**Step 2: Rename post/repository.go.tmpl**

Run: `mv generator/templates/module/post/repository.go.tmpl generator/templates/module/post/repo.go.tmpl`

**Step 3: Verify renames**

Run: `ls generator/templates/module/*/repo.go.tmpl`
Expected: Two files listed

---

## Task 4: Create Platform Template Directory

**Files:**
- Create: `generator/templates/platform/` directory
- Create: `generator/templates/platform/db/` directory

**Step 1: Create platform directories**

Run: `mkdir -p generator/templates/platform/db`

**Step 2: Verify directories exist**

Run: `ls -la generator/templates/platform/`
Expected: `db/` directory exists

---

## Task 5: Move Infra Templates to Platform

**Files:**
- Move: `generator/templates/infra/*.tmpl` → `generator/templates/platform/`
- Rename: `database.go.tmpl` → `db/mysql.go.tmpl`

**Step 1: Move init.go.tmpl**

Run: `mv generator/templates/infra/init.go.tmpl generator/templates/platform/init.go.tmpl`

**Step 2: Move config.go.tmpl**

Run: `mv generator/templates/infra/config.go.tmpl generator/templates/platform/config.go.tmpl`

**Step 3: Move and rename database.go.tmpl**

Run: `mv generator/templates/infra/database.go.tmpl generator/templates/platform/db/mysql.go.tmpl`

**Step 4: Move jwt directory**

Run: `mv generator/templates/infra/jwt generator/templates/platform/`

**Step 5: Verify platform templates**

Run: `ls -la generator/templates/platform/ && ls -la generator/templates/platform/db/ && ls -la generator/templates/platform/jwt/`
Expected: All files in correct locations

**Step 6: Remove empty infra directory**

Run: `rmdir generator/templates/infra`
Expected: Directory removed (may fail if not empty, check manually)

---

## Task 6: Update Main Template Import Paths

**Files:**
- Modify: `generator/templates/cmd/api/main.go.tmpl`

**Step 1: Read current main.go.tmpl**

Run: `cat generator/templates/cmd/api/main.go.tmpl`

**Step 2: Update imports in main.go.tmpl**

Edit the file to change:
- `"{{.ModuleName}}/infra"` → `"{{.ModuleName}}/internal/platform"`
- `"{{.ModuleName}}/server"` → `"{{.ModuleName}}/internal/server"`

The main.go.tmpl should have updated import paths pointing to `internal/platform` and `internal/server`.

**Step 3: Verify changes**

Run: `grep -E "internal/platform|internal/server" generator/templates/cmd/api/main.go.tmpl`
Expected: Lines with new import paths

---

## Task 7: Update Platform Templates

**Files:**
- Modify: `generator/templates/platform/init.go.tmpl`
- Modify: `generator/templates/platform/config.go.tmpl`
- Modify: `generator/templates/platform/db/mysql.go.tmpl`
- Modify: `generator/templates/platform/jwt/jwt.go.tmpl`

**Step 1: Update init.go.tmpl imports**

Edit `generator/templates/platform/init.go.tmpl`:
- Change package from `infra` to `platform`
- Update all import paths containing `{{.ModuleName}}/module` to `{{.ModuleName}}/internal/modules`
- Update all import paths containing `{{.ModuleName}}/shared` to `{{.ModuleName}}/internal/shared`

**Step 2: Update config.go.tmpl**

Edit `generator/templates/platform/config.go.tmpl`:
- Change package from `infra` to `platform`

**Step 3: Update db/mysql.go.tmpl**

Edit `generator/templates/platform/db/mysql.go.tmpl`:
- Change package from `infra` to `db`
- Function names stay the same (or update to match new package)

**Step 4: Update jwt/jwt.go.tmpl**

Edit `generator/templates/platform/jwt/jwt.go.tmpl`:
- Package stays `jwt`
- Update any imports referencing `{{.ModuleName}}/infra` to `{{.ModuleName}}/internal/platform`

**Step 5: Verify all platform templates**

Run: `grep -r "package " generator/templates/platform/`
Expected: All show correct package names (platform, db, jwt)

---

## Task 8: Update Server Templates

**Files:**
- Move: `generator/templates/server/app.go.tmpl` (keep location, update imports)
- Modify: `generator/templates/server/core/module.go.tmpl`
- Modify: `generator/templates/server/core/registry.go.tmpl`
- Modify: `generator/templates/server/helper/helper.go.tmpl`
- Modify: `generator/templates/server/middleware/auth.go.tmpl`
- Modify: `generator/templates/server/middleware/cors.go.tmpl`

**Step 1: Update server templates package**

Server templates will need to move to `internal/server` in output. For now, update imports in templates:
- `{{.ModuleName}}/module/` → `{{.ModuleName}}/internal/modules/`
- `{{.ModuleName}}/infra` → `{{.ModuleName}}/internal/platform`
- `{{.ModuleName}}/shared` → `{{.ModuleName}}/internal/shared`

**Step 2: Update each server template**

Edit each file in `generator/templates/server/` to update import paths:
- `app.go.tmpl`
- `core/module.go.tmpl`
- `core/registry.go.tmpl`
- `helper/helper.go.tmpl`
- `middleware/auth.go.tmpl`
- `middleware/cors.go.tmpl`

**Step 3: Verify server template imports**

Run: `grep -r "internal/modules" generator/templates/server/`
Expected: All templates with module imports show new path

---

## Task 9: Update Shared Templates

**Files:**
- Modify: All files in `generator/templates/shared/`

**Step 1: Check shared templates for imports**

Run: `grep -r "{{.ModuleName}}" generator/templates/shared/`
Expected: Find any import references

**Step 2: Update any import references**

If any shared templates reference `{{.ModuleName}}/module`, update to `{{.ModuleName}}/internal/modules`.

---

## Task 10: Update Module Templates

**Files:**
- Modify: All files in `generator/templates/module/user/`
- Modify: All files in `generator/templates/module/auth/`
- Modify: All files in `generator/templates/module/post/`
- Modify: `generator/templates/module/contract/user.go.tmpl`

**Step 1: Update user module templates**

Edit each file in `generator/templates/module/user/`:
- `module.go.tmpl` - Update imports from `{{.ModuleName}}/infra` to `{{.ModuleName}}/internal/platform`, `{{.ModuleName}}/shared` to `{{.ModuleName}}/internal/shared`, `{{.ModuleName}}/module/contract` to `{{.ModuleName}}/internal/modules/contract`
- `handler.go.tmpl` - Update imports
- `service.go.tmpl` - Update imports
- `repo.go.tmpl` - Update imports, rename repository references to repo
- `model.go.tmpl` - Check for imports
- `dto.go.tmpl` - Check for imports

**Step 2: Update auth module templates**

Edit each file in `generator/templates/module/auth/`:
- `module.go.tmpl` - Update imports
- `handler.go.tmpl` - Update imports
- `service.go.tmpl` - Update imports
- `dto.go.tmpl` - Check for imports

**Step 3: Update post module templates**

Edit each file in `generator/templates/module/post/`:
- `module.go.tmpl` - Update imports
- `handler.go.tmpl` - Update imports
- `service.go.tmpl` - Update imports
- `repo.go.tmpl` - Update imports, rename repository references to repo
- `model.go.tmpl` - Check for imports
- `dto.go.tmpl` - Check for imports

**Step 4: Update contract template**

Edit `generator/templates/module/contract/user.go.tmpl`:
- Update package declaration and any imports

**Step 5: Verify all module imports**

Run: `grep -r "{{.ModuleName}}/module" generator/templates/module/`
Expected: No results (all updated to internal/modules)

Run: `grep -r "{{.ModuleName}}/infra" generator/templates/module/`
Expected: No results (all updated to internal/platform)

---

## Task 11: Update Generator project.go - createDirectories

**Files:**
- Modify: `generator/project.go`

**Step 1: Update createDirectories function**

Edit `generator/project.go` function `createDirectories`:

Replace:
```go
directories := []string{
    "config",
    "module/contract",
    "infra/jwt",
    "module/auth",
    "module/user",
    "module/post",
    "server/core",
    "server/helper",
    "server/middleware",
    "shared/consts",
    "shared/errs",
    "shared/pagination",
    "shared/strutils",
}
```

With:
```go
directories := []string{
    "cmd/api",
    "config",
    "internal/modules/contract",
    "internal/modules/auth",
    "internal/modules/user",
    "internal/modules/post",
    "internal/platform/db",
    "internal/platform/jwt",
    "internal/server/core",
    "internal/server/helper",
    "internal/server/middleware",
    "internal/shared/consts",
    "internal/shared/errs",
    "internal/shared/pagination",
    "internal/shared/strutils",
}
```

**Step 2: Verify changes**

Run: `grep -A20 "func createDirectories" generator/project.go`
Expected: New directory structure

---

## Task 12: Update Generator project.go - processTemplates

**Files:**
- Modify: `generator/project.go`

**Step 1: Update processTemplates mapping**

Edit `generator/project.go` function `processTemplates`:

Replace the `templates` map with:
```go
templates := map[string]string{
    "templates/cmd/api/main.go.tmpl":           "cmd/api/main.go",
    "templates/go.mod.tmpl":                    "go.mod",
    "templates/Makefile.tmpl":                  "Makefile",
    "templates/Dockerfile.tmpl":                "Dockerfile",
    "templates/gitignore.tmpl":                 ".gitignore",
    "templates/env.example.tmpl":               ".env.example",
    "templates/config/config.yml.tmpl":         "config/config.yml",
    "templates/module/contract/user.go.tmpl":   "internal/modules/contract/user.go",
    "templates/platform/init.go.tmpl":          "internal/platform/init.go",
    "templates/platform/config.go.tmpl":        "internal/platform/config.go",
    "templates/platform/db/mysql.go.tmpl":      "internal/platform/db/mysql.go",
    "templates/platform/jwt/jwt.go.tmpl":       "internal/platform/jwt/jwt.go",
    "templates/server/app.go.tmpl":             "internal/server/app.go",
    "templates/server/helper/helper.go.tmpl":   "internal/server/helper/helper.go",
    "templates/server/middleware/auth.go.tmpl": "internal/server/middleware/auth.go",
    "templates/server/middleware/cors.go.tmpl": "internal/server/middleware/cors.go",
    "templates/server/core/module.go.tmpl":     "internal/server/core/module.go",
    "templates/server/core/registry.go.tmpl":   "internal/server/core/registry.go",
    "templates/shared/consts/constants.go.tmpl":  "internal/shared/consts/constants.go",
    "templates/shared/errs/errors.go.tmpl":       "internal/shared/errs/errors.go",
    "templates/shared/pagination/cursor.go.tmpl": "internal/shared/pagination/cursor.go",
    "templates/shared/pagination/offset.go.tmpl": "internal/shared/pagination/offset.go",
    "templates/shared/strutils/string.go.tmpl":   "internal/shared/strutils/string.go",
    "templates/module/user/module.go.tmpl":        "internal/modules/user/module.go",
    "templates/module/user/model.go.tmpl":         "internal/modules/user/model.go",
    "templates/module/user/dto.go.tmpl":           "internal/modules/user/dto.go",
    "templates/module/user/handler.go.tmpl":       "internal/modules/user/handler.go",
    "templates/module/user/service.go.tmpl":       "internal/modules/user/service.go",
    "templates/module/user/repo.go.tmpl":          "internal/modules/user/repo.go",
    "templates/module/auth/module.go.tmpl":        "internal/modules/auth/module.go",
    "templates/module/auth/dto.go.tmpl":           "internal/modules/auth/dto.go",
    "templates/module/auth/handler.go.tmpl":       "internal/modules/auth/handler.go",
    "templates/module/auth/service.go.tmpl":       "internal/modules/auth/service.go",
    "templates/module/post/module.go.tmpl":        "internal/modules/post/module.go",
    "templates/module/post/model.go.tmpl":         "internal/modules/post/model.go",
    "templates/module/post/dto.go.tmpl":           "internal/modules/post/dto.go",
    "templates/module/post/handler.go.tmpl":       "internal/modules/post/handler.go",
    "templates/module/post/service.go.tmpl":       "internal/modules/post/service.go",
    "templates/module/post/repo.go.tmpl":          "internal/modules/post/repo.go",
}
```

**Step 2: Verify changes**

Run: `go build .`
Expected: No errors

---

## Task 13: Update Generator project.go - printSuccess

**Files:**
- Modify: `generator/project.go`

**Step 1: Update printSuccess function**

Edit `generator/project.go` function `printSuccess`:

Replace the directory structure output with:
```go
func printSuccess(projectName string) {
    color.Green("\n🎉 Project '%s' created successfully!", projectName)

    color.Cyan("\n📂 Generated directory structure:")
    fmt.Printf("   %s/\n", projectName)
    fmt.Printf("   ├── cmd/\n")
    fmt.Printf("   │   └── api/\n")
    fmt.Printf("   │       └── main.go\n")
    fmt.Printf("   ├── config/\n")
    fmt.Printf("   │   └── config.yml\n")
    fmt.Printf("   ├── internal/\n")
    fmt.Printf("   │   ├── modules/\n")
    fmt.Printf("   │   │   ├── contract/   # Interface definitions\n")
    fmt.Printf("   │   │   ├── auth/       # Authentication module\n")
    fmt.Printf("   │   │   ├── user/       # User module\n")
    fmt.Printf("   │   │   └── post/       # Demo CRUD module\n")
    fmt.Printf("   │   ├── platform/\n")
    fmt.Printf("   │   │   ├── db/\n")
    fmt.Printf("   │   │   └── jwt/\n")
    fmt.Printf("   │   ├── server/\n")
    fmt.Printf("   │   │   ├── core/\n")
    fmt.Printf("   │   │   ├── helper/\n")
    fmt.Printf("   │   │   └── middleware/\n")
    fmt.Printf("   │   └── shared/\n")
    fmt.Printf("   │       ├── consts/\n")
    fmt.Printf("   │       ├── errs/\n")
    fmt.Printf("   │       ├── pagination/\n")
    fmt.Printf("   │       └── strutils/\n")
    fmt.Printf("   ├── go.mod\n")
    fmt.Printf("   ├── Makefile\n")
    fmt.Printf("   └── Dockerfile\n")

    color.Cyan("\n📋 Next steps:")
    fmt.Printf("   1. cd %s\n", projectName)
    fmt.Printf("   2. Update config/config.yml with your database credentials\n")
    fmt.Printf("   3. cp .env.example .env && edit .env for local development\n")
    fmt.Printf("   4. go mod tidy\n")
    fmt.Printf("   5. go run ./cmd/api\n")

    color.Yellow("\n⚠ Don't forget:")
    fmt.Printf("   - Update config/config.yml (Database, Redis, JWT settings)\n")
    fmt.Printf("   - Default server port: 8080\n")

    color.Cyan("\n📚 Documentation: https://github.com/lpphub/ppx\n")
}
```

**Step 2: Verify changes**

Run: `go build .`
Expected: No errors

---

## Task 14: Update cmd/new.go Help Text

**Files:**
- Modify: `cmd/new.go`

**Step 1: Update Long description**

Edit `cmd/new.go` and update the `Long` field of `newCmd`:

```go
Long: `Create a new Go web project with modular architecture.

Examples:
  ppx new myapp
  ppx new myapp --module github.com/user/myapp

Generated Project Structure:
  myapp/
  ├── cmd/
  │   └── api/
  │       └── main.go
  ├── config/
  │   └── config.yml
  ├── internal/
  │   ├── modules/
  │   │   ├── contract/      # Contract/Interface definitions
  │   │   ├── auth/          # Authentication module
  │   │   ├── user/          # User module
  │   │   └── post/          # Demo post module (CRUD example)
  │   ├── platform/
  │   │   ├── db/
  │   │   └── jwt/
  │   ├── server/
  │   │   ├── app.go
  │   │   ├── core/
  │   │   ├── helper/
  │   │   └── middleware/
  │   └── shared/
  │       ├── consts/
  │       ├── errs/
  │       ├── pagination/
  │       └── strutils/
  ├── go.mod
  ├── Makefile
  └── Dockerfile`,
```

**Step 2: Verify changes**

Run: `go build .`
Expected: No errors

---

## Task 15: Update README.md

**Files:**
- Modify: `README.md`

**Step 1: Update Generated Project Structure section**

Edit `README.md` to reflect new structure:

```markdown
## Generated Project Structure

\`\`\`
myproject/
├── cmd/
│   └── api/
│       └── main.go           # Application entry point
├── config/
│   └── config.yml            # Configuration file
├── internal/
│   ├── modules/
│   │   ├── contract/         # Interface definitions
│   │   │   └── user.go
│   │   ├── auth/             # Authentication module
│   │   │   ├── module.go
│   │   │   ├── dto.go
│   │   │   ├── handler.go
│   │   │   └── service.go
│   │   ├── user/             # User module
│   │   │   ├── module.go
│   │   │   ├── model.go
│   │   │   ├── dto.go
│   │   │   ├── handler.go
│   │   │   ├── service.go
│   │   │   └── repo.go
│   │   └── post/             # Demo CRUD module
│   │       ├── module.go
│   │       ├── model.go
│   │       ├── dto.go
│   │       ├── handler.go
│   │       ├── service.go
│   │       └── repo.go
│   ├── platform/
│   │   ├── db/
│   │   │   └── mysql.go      # Database connection
│   │   ├── config.go         # Config loader
│   │   ├── init.go           # Infrastructure init
│   │   └── jwt/
│   │       └── jwt.go        # JWT utilities
│   ├── server/
│   │   ├── app.go            # HTTP server
│   │   ├── core/
│   │   │   ├── module.go
│   │   │   └── registry.go
│   │   ├── helper/
│   │   │   └── helper.go
│   │   └── middleware/
│   │       ├── auth.go       # JWT authentication
│   │       └── cors.go       # CORS middleware
│   └── shared/
│       ├── consts/           # Constants
│       ├── errs/             # Error definitions
│       ├── pagination/       # Pagination utilities
│       └── strutils/         # String utilities
├── go.mod
├── Makefile
└── Dockerfile
\`\`\`
```

**Step 2: Update Quick Start section**

Change `go run .` to `go run ./cmd/api`

**Step 3: Verify changes**

Run: `cat README.md`
Expected: Updated structure visible

---

## Task 16: Remove Old main.go.tmpl

**Files:**
- Delete: `generator/templates/main.go.tmpl`

**Step 1: Remove old template**

Run: `rm generator/templates/main.go.tmpl`

**Step 2: Verify removal**

Run: `ls generator/templates/main.go.tmpl`
Expected: "No such file or directory"

---

## Task 17: Build and Test

**Files:**
- None (testing only)

**Step 1: Build the ppx tool**

Run: `go build -o ppx .`
Expected: No errors

**Step 2: Create test project**

Run: `./ppx new testproject`
Expected: Project created successfully

**Step 3: Verify directory structure**

Run: `tree testproject -L 3`
Expected: Structure matches design with `cmd/api/` and `internal/`

**Step 4: Verify generated project builds**

Run: `cd testproject && go mod tidy && go build ./...`
Expected: No errors

**Step 5: Clean up test project**

Run: `rm -rf testproject`

---

## Task 18: Final Commit

**Step 1: Stage all changes**

Run: `git add -A`

**Step 2: Commit changes**

Run: `git commit -m "refactor: restructure new command templates to use cmd/ and internal/ layout

- Move main.go to cmd/api/main.go
- Move module/ to internal/modules/
- Move infra/ to internal/platform/
- Move server/ to internal/server/
- Move shared/ to internal/shared/
- Rename init.go to module.go in modules
- Rename repository.go to repo.go in modules
- Rename database.go to db/mysql.go
- Update all import paths in templates
- Update generator/project.go directory and template mappings
- Update README.md with new structure"`

**Step 3: Verify commit**

Run: `git log -1`
Expected: Commit message visible