# Desktop Launcher for noise.sh

## Quick Setup

### Option 1: Automatic Setup (Recommended)
Run the setup script:
```bash
./scripts/create_desktop_launcher.sh
```

This will create a launcher on your Desktop that you can double-click to launch noise.sh.

### Option 2: Manual Setup
1. The launcher file `Launch noise.sh.command` is already in the project root
2. Copy it to your Desktop:
   ```bash
   cp "Launch noise.sh.command" ~/Desktop/
   ```
3. Make sure it's executable:
   ```bash
   chmod +x ~/Desktop/"Launch noise.sh.command"
   ```

## Usage

Simply double-click the `Launch noise.sh.command` file on your Desktop. The script will:
- Check if Go is installed
- Build the application if needed
- Launch noise.sh in Terminal

## Features

- **Auto-build**: Automatically builds the app if the binary doesn't exist
- **Error handling**: Shows helpful error messages if something goes wrong
- **Terminal integration**: Launches in Terminal so you can see output and interact with the app

## Troubleshooting

### "Go is not installed" error
Install Go 1.21+ from https://go.dev

### "Could not find noise.sh directory" error
Make sure the launcher file is in the project root directory, or update the script path.

### Permission denied
Make sure the file is executable:
```bash
chmod +x "Launch noise.sh.command"
```

### Build fails
Check that you're in the project directory and all dependencies are installed:
```bash
go mod download
```

