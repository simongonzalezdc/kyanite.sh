# GitHub Repository Setup Instructions

## Steps to Create and Push to GitHub

1. **Create a new repository on GitHub:**
   - Go to https://github.com/new
   - Repository name: `voxforge` (lowercase)
   - Description: "Voice-to-song tool that transforms vocals into music"
   - Choose: Public or Private
   - Do NOT initialize with README, .gitignore, or license (we already have these)
   - Click "Create repository"

2. **Add the remote and push:**
   ```bash
   cd /Users/simongonzalezdecruz/Desktop/VoxForge
   git remote add origin https://github.com/YOUR_USERNAME/voxforge.git
   git branch -M main
   git push -u origin main
   ```

   Replace `YOUR_USERNAME` with your GitHub username.

3. **If you need to authenticate:**
   - Use a Personal Access Token (PAT) instead of password
   - Or use SSH: `git remote set-url origin git@github.com:YOUR_USERNAME/voxforge.git`

## Repository is Ready!

After pushing, your repository will be available at:
`https://github.com/YOUR_USERNAME/voxforge`

