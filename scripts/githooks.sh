#!/bin/bash

# Setup local Git hooks.

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

###########################################################
# Pre-Push hook
###########################################################
echo "Creating Git 'pre-push' hook in .git/hooks/pre-push..."
cat > "$SCRIPT_DIR/../.git/hooks/pre-push" << EOF
#!/bin/sh

echo "Linting..."
make lint || exit 1

echo "Static analysis..."
make sast || exit 1

echo "Module checks..."
make modcheck || exit 1

echo "Tests and coverage..."
make test || exit 1

exit 0
EOF

chmod 755 .git/hooks/pre-push
