# Minecraft Server Builder (MCSB)

*MCSB is a lightweight command-line utility for rapidly provisioning Minecraft servers.*

---

## Prerequisites

- **Go toolchain** — install from <https://go.dev/dl/>

---

## Installation

```bash
    git clone https://github.com/NickAwrist/mcsb
    cd mcsb
```

---

## Usage

```bash
    go run .
```

The interactive wizard walks you through:

1. Selecting the server framework (PaperMC, Vanilla, or Spigot)
2. Choosing the Minecraft version
3. Entering a server name
4. Specifying a port
5. Accepting Mojang’s EULA

After confirmation, MCSB downloads the required JAR, generates the configuration, and places a fully configured server folder on your desktop—ready to launch.
