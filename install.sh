#!/bin/sh

# colors
reset="\033[0m"
bold="\033[1m"
color_green="\033[32m"
color_red="\033[31m"

# if you call install.sh --dry-run script will show all messages but without actual installation
# use it to ensure that script functions correctly
dry_run=false
if [ "$1" = "--dry-run" ]; then
    printf "${bold}${color_red}This is dry run, no actual actions will be performaed. Use this for testing installation script.${reset}\n"
    dry_run=true
fi

# path settings
workdir=$(pwd)
source_bin_dir="bin"
binary_file="pdfjuicer"
target_bin_dir="$HOME/bin/"

# path to compiled binary
source="$workdir/$source_bin_dir/$binary_file"

# create personal bin dir if not exists
echo "Creating directory $target_bin_dir for binary (if not exist)"
if ! $dry_run; then
  mkdir -p "$target_bin_dir"
fi

# copying to
echo "Copying $source to $target_bin_dir"
if ! $dry_run; then
  cp "$source" "$target_bin_dir"
fi

echo "Adding execute permission for user $(whoami)"
if ! $dry_run; then
  chmod u+x "$target_bin_dir$binary_file"
fi

printf "${bold}${color_green}Installation is completed.\n${reset}"
printf "Please add the following line to your shell configuration file\n"
printf "such as .bashrc, .bash_profile, .profile or .zshrc to include the new directory in your PATH:\n\n"
printf "${bold}export PATH=\"\$PATH:$target_bin_dir\"${reset}\n\n"
printf "Then reload terminal and run ${bold}$binary_file${reset} in terminal.\n"
printf "If installation was successful you will see info about app and help after calling ${bold}$binary_file${reset}\n"
