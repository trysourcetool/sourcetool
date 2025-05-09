# Sourcetool

**Backend logic seamlessly becomes reality.**

Sourcetool transforms your backend code into powerful internal tools. No frontend skills required.

[📚 Documentation](https://docs.trysourcetool.com) | [💬 Discord Community](https://discord.com/invite/K76agfQQKP)

![sourcetool_image](https://github.com/user-attachments/assets/7ab3ddeb-cb12-4153-8b26-974693c67866)

## 🌟 About Sourcetool

We develop Sourcetool, an open-source internal tool builder that handles frontend complexities automatically, allowing developers to focus on implementing business logic in backend code only.

### Backend-First Development
Focus on your business logic while we handle the UI. Build complete internal tools using only Go code. No frontend expertise required.

### Type-Safe & Flexible
Built with Go's type system for reliability. Create robust applications with type-safe APIs and seamless integration.

*Watch our demo video and see Sourcetool in action!*

https://github.com/user-attachments/assets/6c96ac38-8150-4d3d-a4ad-abab083cb77c

## 🚀 Get Started

> **Note:** While our cloud version is coming soon, you can start using Sourcetool today by deploying it in your own environment.

1. **Deploy Sourcetool**
   - Follow our [Deployment Guide](https://docs.trysourcetool.com/docs/getting-started/deployment) to set up Sourcetool in your environment
   - Use Docker for quick and easy deployment

2. **Install the Sourcetool SDK**
   ```bash
   go get github.com/trysourcetool/sourcetool-go
   ```

3. **Write your first internal tool**
   ```go
   package main

   import (
       "github.com/trysourcetool/sourcetool-go"
       "github.com/trysourcetool/sourcetool-go/textinput"
       "github.com/trysourcetool/sourcetool-go/table"
   )

   func listUsersPage(ui sourcetool.UIBuilder) error {
       ui.Markdown("## Users")

       // Search form
       name := ui.TextInput("Name", textinput.WithPlaceholder("Enter name to search"))
       
       // Fetch users from the database
       users, err := listUsers(name)
       if err != nil {
           return err
       }

       // Display users in a table
       ui.Table(users, table.WithHeader("Users List"))
       
       return nil
   }

   func main() {
       s := sourcetool.New(&sourcetool.Config{
           APIKey:   "your_api_key",
           Endpoint: "wss://your-sourcetool-instance"  // Your self-hosted Sourcetool endpoint
       })
       
       // Register pages
       s.Page("/users", "Users List", listUsersPage)
       
       if err := s.Listen(); err != nil {
           log.Fatal(err)
       }
   }
   ```

## 🚢 Deployment

Sourcetool can be easily deployed using Docker in your environment. We provide comprehensive deployment documentation covering:
- Infrastructure requirements (PostgreSQL, Redis)
- Docker image configuration
- Environment variables setup
- Production best practices

For detailed instructions, check out our [Deployment Guide](https://docs.trysourcetool.com/docs/getting-started/deployment).

## 🛠️ Local Development Setup

To get started with local development, follow these steps:

1. **Install Prerequisites**

   Make sure you have the following tools installed:
   - [Docker & Docker Compose](https://docs.docker.com/get-docker/)
   - [GNU core utilities](https://www.gnu.org/software/coreutils/) (`head`, `base64`, `sed`, `cat`, etc.)
   - (Optional) [Make](https://www.gnu.org/software/make/) for easier command usage

2. **Set Up Environment Variables**

   Run the setup script to generate your `.env` file and configure secrets:

   ```bash
   ./devtools/setup_local.sh
   ```

   - This script will check for required tools, generate secure keys, and interactively prompt you for Google OAuth and SMTP settings.
   - You can skip Google OAuth and SMTP setup during the script and edit `.env` later if needed.

3. **Start the Development Environment**

   Use Docker Compose to start all services:

   ```bash
   make start
   ```

   This will launch the backend, frontend, database, and other dependencies.

4. **Access the Application**

   - Application (Frontend & API): [http://localhost:3000](http://localhost:3000)

5. **Stopping Services**

   To stop all running services:

   ```bash
   make stop
   ```

## ❓ FAQ

### What is Sourcetool?
Sourcetool is an open-source internal tool builder that enables you to build full-featured internal tools without writing any frontend code. It handles all frontend complexities automatically, allowing you to focus on implementing business logic in your backend code.

### Do I need frontend skills to use Sourcetool?
No. As an internal tool builder, Sourcetool lets you create complete applications using only Go. The system automatically handles all UI rendering and interactions without requiring any frontend code.

### What types of applications can I build with Sourcetool?
Admin panels, dashboards, data management systems, monitoring tools, and any application where development speed is more important than custom UI/UX.

### Is Sourcetool secure?
Yes, Sourcetool is designed with security in mind. You deploy and run Sourcetool applications on your own infrastructure, keeping your code and sensitive data within your control.

### Is Sourcetool free to use?
Check out the [Sourcetool website](https://trysourcetool.com/) for pricing information.

### How does Sourcetool differ from Retool?
Retool uses a GUI-based drag-and-drop approach, while Sourcetool is **code-first**, making it a strong **Retool alternative**. This makes Sourcetool ideal for the AI era where code can be easily understood and modified by AI tools. With type-safe backend code, your applications are Git version-controllable and integrate seamlessly with AI-assisted development workflows.

## 📚 Resources

- [Documentation](https://docs.trysourcetool.com)
- [Discord Community](https://discord.com/invite/K76agfQQKP)
- [GitHub Repository](https://github.com/trysourcetool/sourcetool)
- [Security Policy](SECURITY.md)
- [Code of Conduct](CODE_OF_CONDUCT.md)

---

<div align="center">
Made with ❤️ by the Sourcetool Team
</div>
