# Sourcetool Protocol Buffers

Protocol Buffers definitions for Sourcetool.

## Setup

1. Install buf
```bash
brew install bufbuild/buf/buf
```

2. Update dependencies
```bash
make mod-update
```

3. Generate code
```bash
make generate
```

## Usage

### Go

Add this module as a dependency:

```bash
go get github.com/trysourcetool/sourcetool/proto
```

Then import the generated packages:

```go
import (
    commonv1 "github.com/trysourcetool/sourcetool/proto/go/common/v1"
    exceptionv1 "github.com/trysourcetool/sourcetool/proto/go/exception/v1"
    pagev1 "github.com/trysourcetool/sourcetool/proto/go/page/v1"
    websocketv1 "github.com/trysourcetool/sourcetool/proto/go/websocket/v1"
    widgetv1 "github.com/trysourcetool/sourcetool/proto/go/widget/v1"
)
```

### TypeScript/JavaScript

#### For Library Maintainers

To generate TypeScript/JavaScript code:

1. Install buf and its plugins
```bash
brew install bufbuild/buf/buf
```

2. Generate code
```bash
make generate
```

This will:
- Generate TypeScript/JavaScript code in the `ts` directory using the `buf.build/community/timostamm-protobuf-ts` plugin
- Create both `.js` and `.d.ts` files
- Include necessary runtime dependencies

3. Install dependencies and build the package
```bash
cd ts
npm install  # Install dependencies including TypeScript
npm run build  # Run tsc to compile TypeScript
```

4. Debug the package locally

a. Create a tarball to verify package contents:
```bash
cd ts
npm pack  # Creates a .tgz file containing what would be published
tar -tf sourcetool/proto-0.1.0.tgz  # List contents to verify included files
```
This helps verify:
- All necessary files are included (check files field in package.json)
- No unnecessary files are included
- Directory structure is correct

b. Test in another project using npm link:
```bash
# In the sourcetool/proto/ts directory
cd ts
npm link  # Creates a global symlink

# In your test project directory
npm link @trysourcetool/proto  # Links to your local package
```

c. Verify the imports work in your test project:
```typescript
import { Message } from '@trysourcetool/proto/websocket/v1/message';
// Try using the types and verify compilation works
// Test the runtime by creating and using message instances
```

5. Publish the package

The package is automatically published to npm when a new version tag is pushed to GitHub. To publish a new version:

1. Create and push a new version tag:
```bash
git tag v1.0.0  # Use appropriate version number
git push origin v1.0.0
```

2. The GitHub Actions workflow will:
   - Generate the TypeScript/JavaScript code
   - Build the package
   - Publish to npm automatically

Note: For publishing @trysourcetool scoped package:
- You need to be a member of the trysourcetool organization on npm
- You need to set the NPM_TOKEN secret in your GitHub repository settings

#### For Library Users

Add the package as a dependency:

```bash
npm install @trysourcetool/proto
# or
yarn add @trysourcetool/proto
```

The package includes:
- Generated TypeScript type definitions (.d.ts files)
- Generated JavaScript code (.js files)
- Runtime dependencies (@protobuf-ts/runtime)

Then import and use the generated code:

```typescript
import { Message } from '@trysourcetool/proto/websocket/v1/message';
import { Page } from '@trysourcetool/proto/page/v1/page';
import { Widget } from '@trysourcetool/proto/widget/v1/widget';
import { Exception } from '@trysourcetool/proto/exception/v1/exception';
import { Common } from '@trysourcetool/proto/common/v1/common';
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
