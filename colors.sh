#!/usr/bin/env bash

source "${XDG_STATE_HOME:-$HOME/.local/state}/rong/colors.bash"

colors=(
  "$COLOR_0" "$COLOR_1" "$COLOR_2" "$COLOR_3"
  "$COLOR_4" "$COLOR_5" "$COLOR_6" "$COLOR_7"
  "$COLOR_8" "$COLOR_9" "$COLOR_A" "$COLOR_B"
  "$COLOR_C" "$COLOR_D" "$COLOR_E" "$COLOR_F"
)

for i in {0..15}; do
  hex="${colors[$i]}"
  # Convert hex (#RRGGBB) to decimal RGB
  r=$((16#${hex:1:2}))
  g=$((16#${hex:3:2}))
  b=$((16#${hex:5:2}))

  # Print block with color and index
  printf "\e[48;2;%d;%d;%dm  %2X  \e[0m" "$r" "$g" "$b" "$i"

  # Break line after 8 colors
  if (((i + 1) % 8 == 0)); then echo; fi
done
