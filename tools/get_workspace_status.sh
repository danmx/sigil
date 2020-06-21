#! /bin/sh

set -eu

appName="sigil"
appVersion="${VERSION:-0.5.3}"
gitCommit="${GIT_COMMIT:-$(git rev-parse HEAD)}"
gitBranch="${GIT_BRANCH:-$(git rev-parse --abbrev-ref HEAD)}"

IFS='.' read -r major minor patch << EOF
${appVersion}
EOF

cat <<EOF
STABLE_APPNAME ${appName}
STABLE_VERSION ${major}.${minor}.${patch}
STABLE_GIT_COMMIT ${gitCommit}
STABLE_GIT_BRANCH ${gitBranch}
STABLE_MAJOR_VERSION ${major}
STABLE_MINOR_VERSION ${major}.${minor}
EOF
