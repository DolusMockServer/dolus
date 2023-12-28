#!/bin/bash

GIT_REPO="git@github.com:DolusMockServer/dolus.git"
DESTINATION_DIR="$HOME/cue/github.com/DolusMockServer/dolus"

mkdir -p $DESTINATION_DIR

# Clone the repository to a temporary directory
temp_dir=$(mktemp -d)
git clone --depth 1 "$GIT_REPO" "$temp_dir" 

# Get all tags in the repository
tags=($(git -C "$temp_dir" tag))

# Sort the tags by semantic versioning
sorted_tags=($(printf '%s\n' "${tags[@]}" | sort -V))

if [ ${#sorted_tags[@]} -gt 0 ]; then
    latest_tag=${sorted_tags[-1]}
    git -C "$temp_dir" checkout "$latest_tag"
else
    default_branch=$(git remote show origin | grep "HEAD branch" | awk '{print $NF}')
    git -C "$temp_dir" checkout $default_branch
fi

cp -r "$temp_dir/cue-expectations" "$DESTINATION_DIR"

# Clean up temporary directory
rm -rf "$temp_dir"


echo "Dolus expectations cue pkg installed at $DESTINATION_DIR"
