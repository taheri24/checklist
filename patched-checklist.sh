go run ./cmd/checklist 
{  echo "\`\`\`patch"; git --no-pager diff --no-color ; echo "\`\`\`";cat header.txt;cat selected.txt; } #| xclip -selection clipboard