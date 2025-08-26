#!/bin/bash

# Enhanced Sequential Book Generator Management Script

BINARY_NAME="enhanced_sequential_generator"

case "$1" in
    "build")
        echo "🔨 Building enhanced generator..."
        ./build_enhanced_generator.sh
        ;;
    
    "start")
        echo "🚀 Starting enhanced sequential generator..."
        if [ -f "$BINARY_NAME" ]; then
            ./"$BINARY_NAME"
        else
            echo "❌ Binary not found. Run './manage_generator.sh build' first."
        fi
        ;;
    
    "pause")
        echo "⏸️  Pausing generator..."
        PID=$(pgrep $BINARY_NAME)
        if [ -n "$PID" ]; then
            kill -USR1 $PID
            echo "📋 Pause signal sent to process $PID"
        else
            echo "❌ Generator not running"
        fi
        ;;
    
    "resume")
        echo "▶️  Resuming generator..."
        PID=$(pgrep $BINARY_NAME)
        if [ -n "$PID" ]; then
            kill -USR2 $PID
            echo "📋 Resume signal sent to process $PID"
        else
            echo "❌ Generator not running"
        fi
        ;;
    
    "stop")
        echo "🛑 Stopping generator gracefully..."
        PID=$(pgrep $BINARY_NAME)
        if [ -n "$PID" ]; then
            kill -INT $PID
            echo "📋 Stop signal sent to process $PID"
        else
            echo "❌ Generator not running"
        fi
        ;;
    
    "kill")
        echo "🚨 Force killing generator..."
        PID=$(pgrep $BINARY_NAME)
        if [ -n "$PID" ]; then
            kill -9 $PID
            echo "📋 Force kill signal sent to process $PID"
        else
            echo "❌ Generator not running"
        fi
        ;;
    
    "status")
        echo "📊 Checking generator status..."
        PID=$(pgrep $BINARY_NAME)
        if [ -n "$PID" ]; then
            echo "✅ Generator is running (PID: $PID)"
            echo "📋 Memory usage:"
            ps -p $PID -o pid,ppid,cmd,%mem,%cpu,etime
            
            # Check for state files
            if ls enhanced_generator_state_*.json 1> /dev/null 2>&1; then
                echo "💾 State files found:"
                ls -la enhanced_generator_state_*.json
            fi
            
            # Check for log files
            if [ -d "output/books/logs" ]; then
                echo "📝 Log files:"
                ls -la output/books/logs/
            fi
        else
            echo "❌ Generator not running"
            
            # Check for resumable state
            if ls enhanced_generator_state_*.json 1> /dev/null 2>&1; then
                echo "🔄 Resumable state files found:"
                ls -la enhanced_generator_state_*.json
            fi
        fi
        ;;
    
    "logs")
        echo "📝 Showing recent logs..."
        if [ -d "output/books/logs" ]; then
            LOG_FILE=$(ls -t output/books/logs/generation_log_*.txt | head -1)
            if [ -n "$LOG_FILE" ]; then
                echo "📄 Showing last 50 lines of $LOG_FILE:"
                echo "----------------------------------------"
                tail -50 "$LOG_FILE"
            else
                echo "❌ No log files found"
            fi
        else
            echo "❌ Log directory not found"
        fi
        ;;
    
    "clean")
        echo "🧹 Cleaning up state and log files..."
        rm -f enhanced_generator_state_*.json
        rm -f enhanced_generator_state_*.json.backup
        rm -rf output/books/logs/
        echo "✅ Cleanup complete"
        ;;
    
    "monitor")
        echo "📊 Monitoring generator in real-time..."
        echo "Press Ctrl+C to stop monitoring"
        while true; do
            clear
            echo "=== Enhanced Sequential Generator Monitor ==="
            echo "Time: $(date)"
            echo ""
            
            PID=$(pgrep $BINARY_NAME)
            if [ -n "$PID" ]; then
                echo "✅ Status: Running (PID: $PID)"
                echo "📊 Resource Usage:"
                ps -p $PID -o pid,ppid,cmd,%mem,%cpu,etime
                echo ""
                
                # Show latest log entries
                if [ -d "output/books/logs" ]; then
                    LOG_FILE=$(ls -t output/books/logs/generation_log_*.txt | head -1)
                    if [ -n "$LOG_FILE" ]; then
                        echo "📝 Latest Log Entries:"
                        echo "----------------------------------------"
                        tail -10 "$LOG_FILE"
                    fi
                fi
            else
                echo "❌ Status: Not Running"
            fi
            
            sleep 5
        done
        ;;
    
    "help"|*)
        echo "📋 Enhanced Sequential Book Generator Management"
        echo ""
        echo "Usage: $0 <command>"
        echo ""
        echo "Commands:"
        echo "  build   - Build the enhanced generator binary"
        echo "  start   - Start the generator"
        echo "  pause   - Pause current generation"
        echo "  resume  - Resume paused generation"
        echo "  stop    - Gracefully stop generation"
        echo "  kill    - Force kill the generator"
        echo "  status  - Show current status and state"
        echo "  logs    - Show recent log entries"
        echo "  clean   - Clean up state and log files"
        echo "  monitor - Real-time monitoring (Ctrl+C to exit)"
        echo "  help    - Show this help message"
        echo ""
        echo "📊 Interactive Commands (while generator is running):"
        echo "  Type into generator console:"
        echo "    pause   - Pause generation"
        echo "    resume  - Resume generation"
        echo "    stop    - Graceful stop"
        echo "    status  - Show status"
        echo "    save    - Manual save"
        echo "    help    - Show help"
        echo ""
        ;;
esac