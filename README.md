# ObsidianOps CLI (`obops`)

*ObsidianOps CLI (`obops`) is a lightweight command-line utility for rapidly provisioning Minecraft servers.*

Whether you're looking to quickly set up a server for yourself, friends, or for development, `obops` streamlines the process with an interactive wizard.

---

## Features

* Quickly set up Minecraft servers.
* Interactive CLI wizard to guide you through setup.
* Supports server types:
    * PaperMC (recommended for performance and plugins)
    * Vanilla (the original Minecraft experience)
* Automatic download of server JARs.
* Configuration of server properties (name, port).
* Automatic EULA acceptance (after user confirmation).
* Cross-platform (Windows, macOS, Linux).

---

## Installation

You can install `obops` using one of the package managers below, or by downloading directly from GitHub Releases.

### macOS & Linux (via Homebrew)

1.  **Add the Tap:**
    First, you need to add the `NickAwrist/obsidian-ops-cli` tap. This tells Homebrew where to find the `obops` formula.
    ```bash
    brew tap NickAwrist/obsidian-ops-cli https://github.com/NickAwrist/obsidian-ops-cli
    ```

2.  **Install `obops`:**
    Once the tap is added, you can install `obops`:
    ```bash
    brew install obops
    ```

    To upgrade in the future: `brew upgrade obops`

### Windows (via Chocolatey)

1.  **Prerequisite:** Ensure you have Chocolatey installed. If not, follow the instructions at [chocolatey.org/install](https://chocolatey.org/install).

2.  **Install `obops`:**
    Open PowerShell (preferably as Administrator) and run:
    ```powershell
    choco install obops
    ```

    To upgrade in the future: `choco upgrade obops`

### Manual Installation (All Platforms)

1.  Go to the [**GitHub Releases page**](https://github.com/NickAwrist/obsidian-ops-cli/releases).
2.  Download the archive (`.tar.gz` or `.zip`) for your operating system and architecture.
3.  Extract the `obops` (or `obops.exe` for Windows) executable.
4.  Move the executable to a directory that is included in your system's PATH environment variable (e.g., `/usr/local/bin` on macOS/Linux, or a custom directory on Windows).

---

## Usage

The primary command to start provisioning a new server is `obops create`.

```bash
obops create
```