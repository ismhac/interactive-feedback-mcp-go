#!/usr/bin/env python3
"""
Single Popup Desktop GUI for Interactive Feedback MCP
Creates one unified Tkinter dialog with all information and input
"""

import sys
import os
import subprocess
import tempfile
import json
import tkinter as tk
from tkinter import ttk, scrolledtext
from pathlib import Path

class SinglePopupDesktopGUI:
    def __init__(self, project_directory, prompt):
        self.project_directory = project_directory
        self.prompt = prompt
        self.feedback = None
        self.root = None
        
    def show_notification(self, title, message):
        """Show desktop notification"""
        try:
            # Try notify-send (Linux)
            if os.system("which notify-send > /dev/null 2>&1") == 0:
                os.system(f'notify-send "{title}" "{message}"')
            # Try osascript (macOS)
            elif os.system("which osascript > /dev/null 2>&1") == 0:
                os.system(f'osascript -e \'display notification "{message}" with title "{title}"\'')
            # Try msg (Windows)
            elif os.system("which msg > /dev/null 2>&1") == 0:
                os.system(f'msg * "{title}: {message}"')
            else:
                print(f"{title}: {message}")
        except:
                print(f"{title}: {message}")
    
    def apply_dark_theme(self):
        """Apply dark theme styling to the application"""
        # Dark theme colors
        bg_color = "#2b2b2b"  # Dark gray background
        fg_color = "#ffffff"  # White text
        select_bg = "#404040"  # Selection background
        select_fg = "#ffffff"  # Selection text
        entry_bg = "#3c3c3c"  # Entry background
        button_bg = "#404040"  # Button background
        button_fg = "#ffffff"  # Button text
        button_active = "#505050"  # Button active state
        
        # Configure root window
        self.root.configure(bg=bg_color)
        
        # Configure ttk styles
        style = ttk.Style()
        style.theme_use('clam')  # Use clam theme as base for better customization
        
        # Configure Frame style
        style.configure('TFrame', background=bg_color)
        
        # Configure Label style
        style.configure('TLabel', 
                       background=bg_color, 
                       foreground=fg_color,
                       font=('Arial', 10))
        
        # Configure Button style with rounded corners
        style.configure('TButton',
                       background=button_bg,
                       foreground=button_fg,
                       borderwidth=0,
                       focuscolor='none',
                       relief='flat',
                       padding=(10, 5))
        
        style.map('TButton',
                 background=[('active', button_active),
                           ('pressed', button_active)])
        
        # Configure Text widget style (for ScrolledText)
        text_style = {
            'bg': entry_bg,
            'fg': fg_color,
            'selectbackground': select_bg,
            'selectforeground': select_fg,
            'insertbackground': fg_color,
            'font': ('Consolas', 9),
            'relief': 'flat',
            'borderwidth': 0,
            'highlightthickness': 0
        }
        
        # Store text style for later use
        self.text_style = text_style
    
    def create_rounded_widget(self, parent, widget_class, **kwargs):
        """Create a widget with rounded corners using a frame"""
        # Create a frame with rounded appearance
        frame = tk.Frame(parent, bg=self.text_style['bg'], relief='flat', bd=0)
        
        # Create the actual widget
        widget = widget_class(frame, **kwargs)
        widget.pack(fill='both', expand=True, padx=4, pady=4)
        
        return frame, widget
    
    def create_single_dialog(self):
        """Create a single unified Tkinter dialog with all information and input"""
        # Get conversation history from config
        conversation_text = self.get_conversation_history()
        
        # Create the main window
        self.root = tk.Tk()
        self.root.title("Interactive Feedback MCP")
        
        # Set maximized window by default (keeps title bar with close/minimize buttons)
        try:
            # Try maximized window first (keeps title bar)
            self.root.state('zoomed')
        except:
            try:
                # Linux fallback - maximized window
                self.root.attributes('-zoomed', True)
            except:
                try:
                    # Alternative method for maximized window
                    self.root.wm_state('zoomed')
                except:
                    # Final fallback - set large geometry
                    self.root.geometry("1200x800")
        
        self.root.resizable(True, True)
        
        # Apply dark theme
        self.apply_dark_theme()
        
        # Make window stay on top
        self.root.attributes('-topmost', True)
        self.root.after_idle(lambda: self.root.attributes('-topmost', False))
        
        # Create main frame
        main_frame = ttk.Frame(self.root, padding="10")
        main_frame.grid(row=0, column=0, sticky=(tk.W, tk.E, tk.N, tk.S))
        
        # Configure grid weights
        self.root.columnconfigure(0, weight=1)
        self.root.rowconfigure(0, weight=1)
        main_frame.columnconfigure(0, weight=1)
        main_frame.rowconfigure(1, weight=1)
        
        # Title
        title_label = ttk.Label(main_frame, text="Interactive Feedback MCP", 
                               font=('Arial', 16, 'bold'))
        title_label.grid(row=0, column=0, pady=(0, 10), sticky=tk.W)
        
        # Information display area (read-only) with rounded corners
        info_frame, info_text = self.create_rounded_widget(main_frame, scrolledtext.ScrolledText, 
                                                          height=15, width=80, 
                                                          wrap=tk.WORD, state=tk.DISABLED,
                                                          **self.text_style)
        info_frame.grid(row=1, column=0, pady=(0, 10), sticky=(tk.W, tk.E, tk.N, tk.S))
        
        # Insert information
        info_content = f"""Project: {self.project_directory}

{conversation_text}

Current Prompt: {self.prompt}

Please provide your feedback below:"""
        
        info_text.config(state=tk.NORMAL)
        info_text.insert(tk.END, info_content)
        info_text.config(state=tk.DISABLED)
        
        # Feedback input area
        feedback_frame = ttk.Frame(main_frame)
        feedback_frame.grid(row=2, column=0, pady=(0, 10), sticky=(tk.W, tk.E))
        feedback_frame.columnconfigure(0, weight=1)
        
        feedback_label = ttk.Label(feedback_frame, text="Your feedback:")
        feedback_label.grid(row=0, column=0, sticky=tk.W, pady=(0, 5))
        
        # Feedback entry with rounded corners
        feedback_entry_frame, self.feedback_entry = self.create_rounded_widget(feedback_frame, tk.Text, 
                                                                              height=3, wrap=tk.WORD,
                                                                              **self.text_style)
        feedback_entry_frame.grid(row=1, column=0, sticky=(tk.W, tk.E), pady=(0, 10))
        
        # Buttons frame
        buttons_frame = ttk.Frame(main_frame)
        buttons_frame.grid(row=3, column=0, sticky=(tk.W, tk.E))
        
        # Copy Conversation button
        copy_btn = ttk.Button(buttons_frame, text="Copy Conversation", 
                             command=self.copy_conversation)
        copy_btn.grid(row=0, column=0, padx=(0, 10))
        
        # Submit button
        submit_btn = ttk.Button(buttons_frame, text="Submit", 
                               command=self.submit_feedback)
        submit_btn.grid(row=0, column=1, padx=(0, 10))
        
        # Cancel button
        cancel_btn = ttk.Button(buttons_frame, text="Cancel", 
                               command=self.cancel_feedback)
        cancel_btn.grid(row=0, column=2)
        
        # Focus on feedback entry
        self.feedback_entry.focus()
        
        # Bind Escape key to toggle maximized window
        self.root.bind('<Escape>', self.toggle_fullscreen)
        self.root.bind('<F11>', self.toggle_fullscreen)
        
        # Start the GUI
        self.root.mainloop()
        
        return self.feedback if self.feedback is not None else ""
    
    def toggle_fullscreen(self, event=None):
        """Toggle between maximized and normal window mode"""
        try:
            # Check current window state
            current_state = self.root.state()
            if current_state == 'zoomed':
                self.root.state('normal')
            else:
                self.root.state('zoomed')
        except:
            try:
                # Linux fallback - toggle zoomed attribute
                current_state = self.root.attributes('-zoomed')
                self.root.attributes('-zoomed', not current_state)
            except:
                try:
                    # Alternative method
                    current_state = self.root.wm_state()
                    if current_state == 'zoomed':
                        self.root.wm_state('normal')
                    else:
                        self.root.wm_state('zoomed')
                except:
                    pass
    
    def copy_conversation(self):
        """Copy conversation to clipboard"""
        try:
            conversation_text = self.get_conversation_text_for_copy()
            if conversation_text and conversation_text != "No previous user request found.":
                success = self.copy_to_clipboard(conversation_text)
                if success:
                    self.show_notification("Success", "Conversation copied to clipboard!")
                else:
                    self.show_notification("Error", "Failed to copy to clipboard")
            else:
                self.show_notification("Info", "No conversation to copy")
        except Exception as e:
            self.show_notification("Error", f"Failed to copy conversation: {str(e)}")
    
    def submit_feedback(self):
        """Submit feedback and close dialog"""
        self.feedback = self.feedback_entry.get("1.0", tk.END).strip()
        self.root.quit()
        self.root.destroy()
    
    def cancel_feedback(self):
        """Cancel and close dialog without feedback"""
        self.feedback = ""
        self.root.quit()
        self.root.destroy()
    
    def get_conversation_history(self):
        """Get conversation history from config file"""
        try:
            config_file = os.path.join(self.project_directory, '.interactive-feedback-config.json')
            if os.path.exists(config_file):
                with open(config_file, 'r') as f:
                    config = json.load(f)
                    history = config.get('conversation_history', [])
                    
                    if len(history) >= 2:
                        # Get the last assistant and user messages
                        last_assistant = None
                        last_user = None
                        
                        for entry in reversed(history):
                            if entry.get('role') == 'user' and last_user is None:
                                last_user = entry.get('content', '')
                            elif entry.get('role') == 'assistant' and last_assistant is None:
                                last_assistant = entry.get('content', '')
                            
                            if last_user and last_assistant:
                                break
                        
                        if last_user and last_assistant:
                            return f"""Previous Conversation:
```
user: {last_user}
agent: {last_assistant}
```"""
                    elif len(history) >= 1:
                        # Get the last user message only
                        last_user = None
                        for entry in reversed(history):
                            if entry.get('role') == 'user':
                                last_user = entry.get('content', '')
                                break
                        
                        if last_user:
                            return f"""Previous User Request:
```
user: {last_user}
```"""
            
            return "Previous Conversation: No previous conversation found."
            
        except Exception as e:
            return "Previous Conversation: Error loading conversation history."
    
    def get_conversation_text_for_copy(self):
        """Get conversation text formatted for copying"""
        try:
            config_file = os.path.join(self.project_directory, '.interactive-feedback-config.json')
            if os.path.exists(config_file):
                with open(config_file, 'r') as f:
                    config = json.load(f)
                    history = config.get('conversation_history', [])
                    
                    if len(history) >= 2:
                        # Get the last assistant and user messages
                        last_assistant = None
                        last_user = None
                        
                        for entry in reversed(history):
                            if entry.get('role') == 'user' and last_user is None:
                                last_user = entry.get('content', '')
                            elif entry.get('role') == 'assistant' and last_assistant is None:
                                last_assistant = entry.get('content', '')
                            
                            if last_user and last_assistant:
                                break
                        
                        if last_user and last_assistant:
                            return f"user: {last_user}\nagent: {last_assistant}"
                    elif len(history) >= 1:
                        # Get the last user message only
                        last_user = None
                        for entry in reversed(history):
                            if entry.get('role') == 'user':
                                last_user = entry.get('content', '')
                                break
                        
                        if last_user:
                            return f"user: {last_user}"
            
            return "No previous conversation found."
            
        except Exception as e:
            return "Error loading conversation history."
    
    def copy_to_clipboard(self, text):
        """Copy text to clipboard"""
        try:
            # First try Tkinter clipboard (most reliable)
            if self.root:
                self.root.clipboard_clear()
                self.root.clipboard_append(text)
                self.root.update()  # Force update
                return True
            
            # Fallback to system commands
            # Try xclip (Linux)
            if os.system("which xclip > /dev/null 2>&1") == 0:
                result = subprocess.run(['xclip', '-selection', 'clipboard'], 
                                      input=text, text=True, capture_output=True)
                return result.returncode == 0
            
            # Try xsel (Linux alternative)
            elif os.system("which xsel > /dev/null 2>&1") == 0:
                result = subprocess.run(['xsel', '--clipboard', '--input'], 
                                      input=text, text=True, capture_output=True)
                return result.returncode == 0
            
            # Try pbcopy (macOS)
            elif os.system("which pbcopy > /dev/null 2>&1") == 0:
                result = subprocess.run(['pbcopy'], input=text, text=True, capture_output=True)
                return result.returncode == 0
            
            # Try clip (Windows)
            elif os.system("which clip > /dev/null 2>&1") == 0:
                result = subprocess.run(['clip'], input=text, text=True, capture_output=True)
                return result.returncode == 0
            
            return False
            
        except Exception as e:
            return False
    
    def get_terminal_input(self):
        """Fallback terminal input"""
        print(f"\n{'='*60}")
        print("Interactive Feedback MCP")
        print(f"{'='*60}")
        print(f"Project: {self.project_directory}")
        print(f"Prompt: {self.prompt}")
        print(f"{'='*60}")
        print("Please provide your feedback (or press Enter to skip):")
        
        try:
            feedback = input("Your feedback: ").strip()
            return feedback  # Return empty string if no input
        except (EOFError, KeyboardInterrupt):
            return ""  # Return empty string on interrupt

def main():
    if len(sys.argv) != 3:
        print("Usage: python3 desktop_gui_single.py <project_directory> <prompt>")
        sys.exit(1)
    
    project_directory = sys.argv[1]
    prompt = sys.argv[2]
    
    # Create GUI without system notification
    gui = SinglePopupDesktopGUI(project_directory, prompt)
    
    # Create and show dialog
    feedback = gui.create_single_dialog()
    
    # Output feedback
    print(feedback)

if __name__ == "__main__":
    main()