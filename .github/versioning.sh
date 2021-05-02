#!/usr/bin/env bash
if test "$BASH" = "" || "$BASH" -uc "a=();true \"\${a[@]}\"" 2>/dev/null; then
    # Bash 4.4, Zsh
    set -euo pipefail
else
    # Bash 4.3 and older chokes on empty arrays with set -u.
    set -eo pipefail
fi
shopt -s nullglob
export DEBIAN_FRONTEND=noninteractive


# https://github.com/anordal/shellharden/blob/main/how_to_do_things_safely_in_bash.md
# http://wiki.bash-hackers.org/syntax/pe#use_a_default_value


# Usage:
# .github/versioning.sh v0.0.1

VERSION=${1:-versionNotSet}
if [ "$VERSION" == "versionNotSet"  ]; then
    printf "\n\n VERSION should not be empty\n"
    exit
fi


printf "\n creating git tag: ${VERSION} \n"
printf "\n with commit message, see Changelong: https://github.com/komuw/ote/blob/main/CHANGELOG.md \n" && \
  git tag -a "${VERSION}" -m "see Changelong: https://github.com/komuw/ote/blob/main/CHANGELOG.md"
printf "\n git push the tag::\n" && git push --all -u --follow-tags
