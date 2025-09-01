#!/bin/bash

SESSION="dev"

# Kill any old session with the same name
tmux kill-session -t $SESSION 2>/dev/null

# Create a new detached session with a named window
tmux new-session -d -s $SESSION -n editor -c ~/projects/myapp

# ───────────────────────────────
# Window 1: Editor
# ───────────────────────────────
tmux send-keys -t $SESSION:editor 'nvim .' C-m

# ───────────────────────────────
# Window 2: Server
# ───────────────────────────────
tmux new-window -t $SESSION:1 -n server -c ~/projects/myapp
tmux send-keys -t $SESSION:server 'export NODE_ENV=development' C-m
tmux send-keys -t $SESSION:server 'npm run dev' C-m

# ───────────────────────────────
# Window 3: Shell with 3 panes
# ───────────────────────────────
tmux new-window -t $SESSION:2 -n shell -c ~/projects/myapp

# Split into left and right
tmux split-window -h -t $SESSION:2 -c ~/projects/myapp
# Split bottom right into two (so now: left, top-right, bottom-right)
tmux split-window -v -t $SESSION:2.1 -c ~/projects/myapp

# Send commands into panes
tmux send-keys -t $SESSION:2.0 'git status' C-m
tmux send-keys -t $SESSION:2.1 'htop' C-m
tmux send-keys -t $SESSION:2.2 'tail -f logs/app.log' C-m

# Set a nice tiled layout
tmux select-layout -t $SESSION:2 tiled

# ───────────────────────────────
# Window 4: Database
# ───────────────────────────────
tmux new-window -t $SESSION:3 -n database -c ~/projects/myapp
tmux send-keys -t $SESSION:database 'psql -U postgres mydb' C-m

# ───────────────────────────────
# Set Environment Variable for the Session
# ───────────────────────────────
tmux set-environment -t $SESSION API_URL "http://localhost:3000"

# ───────────────────────────────
# Hooks Example: Run when new window is created
# ───────────────────────────────
tmux set-hook -t $SESSION window-created 'display-message "New window created in session #{session_name}"'

# ───────────────────────────────
# Status Bar Customization (per session)
# ───────────────────────────────
tmux set-option -t $SESSION status on
tmux set-option -t $SESSION status-bg black
tmux set-option -t $SESSION status-fg green
tmux set-option -t $SESSION status-left "[#S]"
tmux set-option -t $SESSION status-right "%Y-%m-%d %H:%M:%S"

# ───────────────────────────────
# Final: Attach to the session
# ───────────────────────────────
tmux attach-session -t $SESSION
