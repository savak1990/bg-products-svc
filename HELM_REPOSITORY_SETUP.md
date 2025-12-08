# Setting Up GitHub Pages as a Helm Chart Repository

This guide explains how to publish your Helm chart to GitHub Pages, making it available as a public Helm repository that others can install from.

## Overview

GitHub Pages can serve as a free Helm chart repository by hosting an `index.yaml` file that Helm uses to discover and download charts. This is a popular approach for open-source projects.

---

## Prerequisites

- GitHub repository with a Helm chart (e.g., `bg-products-svc`)
- Git installed locally
- Helm CLI installed
- Write access to the repository

---

## Step 1: Bootstrap GitHub Pages Branch

GitHub Pages requires a dedicated branch (typically `gh-pages`) to serve static content. You need to create this branch once:

```bash
# Navigate to your repository root
cd /path/to/bg-products-svc

# Ensure you're on main branch with latest changes
git checkout main
git pull origin main

# Create and switch to a new orphan branch (no commit history)
git checkout --orphan gh-pages

# Remove all files from staging
git reset --hard

# Create initial files
echo "# Helm Chart Repository for bg-products-svc" > README.md
touch .nojekyll   # Prevents GitHub Pages from processing with Jekyll

# Commit the initial setup
git add README.md .nojekyll
git commit -m "Initialize gh-pages for Helm chart repository"

# Push the new branch to GitHub
git push -u origin gh-pages

# Switch back to main branch
git checkout main
```

### Why These Commands?

- **`git checkout --orphan gh-pages`**: Creates a new branch with no history (clean slate)
- **`git reset --hard`**: Removes all files from the working directory
- **`.nojekyll` file**: Tells GitHub Pages not to process the site with Jekyll (important for preserving `index.yaml` format)
- **`README.md`**: Optional documentation for the repository page

---

## Step 2: Enable GitHub Pages

1. Go to your GitHub repository in a browser
2. Navigate to **Settings** → **Pages**
3. Under **Source**, select:
   - **Branch**: `gh-pages`
   - **Folder**: `/ (root)`
4. Click **Save**

GitHub will provide a URL like: `https://savak1990.github.io/bg-products-svc/`

---

## Step 3: Package and Publish Your Helm Chart

### Option A: Manual Publishing

```bash
# From your repository root on main branch
cd deploy/helm

# Package the chart
helm package bg-products-svc
# This creates: bg-products-svc-0.0.1.tgz

# Generate or update the index
helm repo index . --url https://savak1990.github.io/bg-products-svc/

# Switch to gh-pages branch
git checkout gh-pages

# Copy the package and index to gh-pages
cp deploy/helm/bg-products-svc-*.tgz .
cp deploy/helm/index.yaml .

# Commit and push
git add *.tgz index.yaml
git commit -m "Publish bg-products-svc chart version 0.0.1"
git push origin gh-pages

# Switch back to main
git checkout main
```

### Option B: Automated with GitHub Actions (Recommended)

Create `.github/workflows/release-charts.yml`:

```yaml
name: Release Helm Charts

on:
  push:
    tags:
      - 'v*.*.*'  # Triggers on version tags like v0.0.1, v1.2.3

jobs:
  release:
    runs-on: ubuntu-latest
    permissions:
      contents: write  # Required to push to gh-pages
    steps:
      - name: Checkout code
        uses: actions/checkout@v4
        with:
          fetch-depth: 0

      - name: Configure Git
        run: |
          git config user.name "$GITHUB_ACTOR"
          git config user.email "$GITHUB_ACTOR@users.noreply.github.com"

      - name: Install Helm
        uses: azure/setup-helm@v3
        with:
          version: 'v3.13.0'

      - name: Run chart-releaser
        uses: helm/chart-releaser-action@v1.6.0
        env:
          CR_TOKEN: "${{ secrets.GITHUB_TOKEN }}"
        with:
          charts_dir: deploy/helm
          skip_existing: true
```

**How it works:**
- Triggered when you push a git tag (e.g., `git tag v0.0.1 && git push origin v0.0.1`)
- Automatically packages the chart
- Updates `index.yaml`
- Pushes to `gh-pages` branch
- Creates a GitHub Release

---

## Step 4: Using Your Helm Repository

Once published, anyone can use your Helm repository:

### Add Your Repository

```bash
# Add the repository
helm repo add bg-charts https://savak1990.github.io/bg-products-svc/

# Update repository index
helm repo update

# Verify it was added
helm repo list
```

### Search for Charts

```bash
# Search for your chart
helm search repo bg-products-svc

# Search with all versions
helm search repo bg-products-svc --versions
```

### Install from Repository

```bash
# Install latest version
helm install my-release bg-charts/bg-products-svc \
  --namespace bg-products \
  --create-namespace

# Install specific version
helm install my-release bg-charts/bg-products-svc \
  --version 0.0.1 \
  --namespace bg-products \
  --create-namespace
```

### Remove Repository

```bash
# Remove when no longer needed
helm repo remove bg-charts
```

---

## Step 5: Publishing Updates

### Manual Update

```bash
# Make changes to your chart
# Update version in Chart.yaml

# Package new version
cd deploy/helm
helm package bg-products-svc

# Update index
helm repo index . --url https://savak1990.github.io/bg-products-svc/ --merge index.yaml

# Switch to gh-pages and publish
git checkout gh-pages
cp deploy/helm/bg-products-svc-*.tgz .
cp deploy/helm/index.yaml .
git add *.tgz index.yaml
git commit -m "Release version X.Y.Z"
git push origin gh-pages
git checkout main
```

### Automated Update (with GitHub Actions)

Simply tag and push:

```bash
# Bump version in deploy/helm/bg-products-svc/Chart.yaml
# Commit changes
git add deploy/helm/bg-products-svc/Chart.yaml
git commit -m "Bump chart version to 0.0.2"
git push origin main

# Create and push tag
git tag v0.0.2
git push origin v0.0.2

# GitHub Actions will automatically publish
```

---

## Directory Structure

After setup, your repository will have:

```
bg-products-svc/
├── main branch (source code)
│   ├── cmd/
│   ├── internal/
│   ├── deploy/
│   │   └── helm/
│   │       └── bg-products-svc/
│   │           ├── Chart.yaml
│   │           ├── values.yaml
│   │           └── templates/
│   └── .github/
│       └── workflows/
│           └── release-charts.yml
│
└── gh-pages branch (Helm repository)
    ├── index.yaml                    # Helm repository index
    ├── bg-products-svc-0.0.1.tgz     # Chart package v0.0.1
    ├── bg-products-svc-0.0.2.tgz     # Chart package v0.0.2
    ├── README.md
    └── .nojekyll
```

---

## Alternative: GitHub Releases

Instead of GitHub Pages, you can use GitHub Releases:

```bash
# Package chart
helm package deploy/helm/bg-products-svc

# Create GitHub release and upload
gh release create v0.0.1 bg-products-svc-0.0.1.tgz \
  --title "Release v0.0.1" \
  --notes "Initial release"
```

Then users install from GitHub Releases:

```bash
helm install my-release \
  https://github.com/savak1990/bg-products-svc/releases/download/v0.0.1/bg-products-svc-0.0.1.tgz
```

---

## Troubleshooting

### GitHub Pages Not Serving Files

1. Check GitHub Pages is enabled in Settings → Pages
2. Verify `gh-pages` branch exists and has content
3. Wait 1-2 minutes for GitHub to build the site
4. Check the site URL is correct (username.github.io/repo-name)

### Index.yaml Not Found

```bash
# Verify index.yaml exists in gh-pages branch
git checkout gh-pages
ls -la index.yaml

# Check the URL in index.yaml matches your GitHub Pages URL
cat index.yaml
```

### Helm Repo Add Fails

```bash
# Test the URL directly
curl -L https://savak1990.github.io/bg-products-svc/index.yaml

# If 404, wait for GitHub Pages to deploy
# If connection refused, check repository is public
```

### Chart Not Updating

```bash
# Force update repository cache
helm repo update

# Remove and re-add repository
helm repo remove bg-charts
helm repo add bg-charts https://savak1990.github.io/bg-products-svc/
```

---

## Best Practices

1. **Version your charts**: Always increment version in `Chart.yaml` before releasing
2. **Use semantic versioning**: Follow `MAJOR.MINOR.PATCH` (e.g., `0.0.1`, `1.0.0`)
3. **Automate releases**: Use GitHub Actions for consistent releases
4. **Test before releasing**: Use `helm lint` and `helm template` to validate
5. **Document changes**: Maintain a CHANGELOG.md in your chart
6. **Sign your charts**: Use `helm package --sign` for production charts
7. **Keep index.yaml updated**: Always use `--merge` flag when updating index

---

## Security Considerations

- **Public access**: GitHub Pages repositories are public (even in private repos)
- **No authentication**: Anyone can download your charts
- **Consider private registries**: For proprietary charts, use:
  - ChartMuseum
  - Harbor
  - AWS ECR
  - Google Artifact Registry
  - Azure Container Registry

---

## Summary

### For First-Time Setup:
```bash
# 1. Create gh-pages branch
git checkout --orphan gh-pages
git reset --hard
echo "# Helm chart repo" > README.md
touch .nojekyll
git add README.md .nojekyll
git commit -m "Initialize gh-pages"
git push -u origin gh-pages
git checkout main

# 2. Enable GitHub Pages in repository settings

# 3. Add GitHub Actions workflow (optional but recommended)
```

### For Each Release:
```bash
# Manual:
helm package deploy/helm/bg-products-svc
helm repo index . --url https://USERNAME.github.io/REPO/ --merge index.yaml
# Copy to gh-pages and push

# Automated:
git tag v0.0.2
git push origin v0.0.2
# GitHub Actions handles the rest
```

### For Users:
```bash
helm repo add bg-charts https://USERNAME.github.io/REPO/
helm install my-release bg-charts/bg-products-svc
```

---

Your Helm chart repository is now ready for public consumption! 🎉
