# Enhanced Sequential Book Generator

A comprehensive book generation system with advanced features for resume, pause, logging, and chunked processing.

## 🚀 Key Features

### ✨ Enhanced Capabilities
- **Resume Functionality**: Automatically save state and resume from interruptions
- **Pause/Resume Control**: Real-time control during generation
- **Comprehensive Logging**: Detailed terminal logging with color coding and progress bars
- **Chunked Processing**: Process books in smaller batches for better resource management
- **Signal Handling**: Graceful handling of system signals (SIGINT, SIGUSR1, SIGUSR2)
- **Interactive Commands**: Real-time control via terminal commands
- **Auto-save**: Automatic state saving at regular intervals
- **Progress Tracking**: Detailed progress metrics and quality assessments

### 🎯 Problem Solutions
1. **Timeout Issues**: Chunked processing prevents long-running operations from timing out
2. **Resume Capability**: Never lose progress due to interruptions
3. **Clear Progress Visibility**: Real-time logging shows exactly what's happening
4. **Resource Management**: Better memory and CPU usage through controlled processing
5. **Graceful Control**: Clean pause/stop/resume without data corruption

## 🔧 Installation & Setup

### Build the Enhanced Generator
```bash
# Make build script executable
chmod +x build_enhanced_generator.sh

# Build the enhanced generator
./build_enhanced_generator.sh
```

### Alternative: Manual Build
```bash
go build -o enhanced_sequential_generator enhanced_sequential_generator.go progress_tracker.go outline_generator.go
```

## 🎮 Usage

### Start Generation
```bash
# Start fresh generation
./enhanced_sequential_generator

# Or use the management script
./manage_generator.sh start
```

### Interactive Commands (while running)
Type these commands into the generator console:
- `pause` - Pause generation after current operation
- `resume` - Resume paused generation  
- `stop` - Gracefully stop and save state
- `status` - Show current progress and metrics
- `save` - Manually save current state
- `help` - Show available commands

### Signal Control
```bash
# Get the process ID
PID=$(pgrep enhanced_sequential_generator)

# Pause generation
kill -USR1 $PID

# Resume generation  
kill -USR2 $PID

# Graceful stop
kill -INT $PID

# Force kill (not recommended)
kill -9 $PID
```

### Management Script
```bash
# All-in-one management
./manage_generator.sh <command>

# Available commands:
./manage_generator.sh build    # Build the generator
./manage_generator.sh start    # Start generation
./manage_generator.sh pause    # Pause via signal
./manage_generator.sh resume   # Resume via signal  
./manage_generator.sh stop     # Graceful stop
./manage_generator.sh kill     # Force kill
./manage_generator.sh status   # Check status
./manage_generator.sh logs     # Show recent logs
./manage_generator.sh clean    # Clean state files
./manage_generator.sh monitor  # Real-time monitoring
```

## 📊 Progress Tracking

### Real-time Logging
The enhanced generator provides comprehensive logging:
- **Color-coded messages**: Success (green), errors (red), info (blue), progress (yellow)
- **Progress bars**: Visual indicators of completion percentage
- **Timestamps**: All log entries include precise timestamps
- **Operation context**: Clear indication of current operation

### State Files
- `enhanced_generator_state_*.json` - Main state file for resumption
- `enhanced_generator_state_*.json.backup` - Automatic backup
- `output/books/logs/generation_log_*.txt` - Detailed text logs
- `progress_tracker_state.json` - Progress tracking data

### Progress Reports
- Individual book progress with quality metrics
- Phase-by-phase completion tracking  
- System-wide performance metrics
- Automatic final report generation

## 🔄 Resume Functionality

### Automatic Resume Detection
When starting the generator:
1. Checks for existing state files
2. Prompts user to resume or start fresh
3. Loads previous progress and continues from last checkpoint

### Resume Points
The system can resume from:
- Before outline generation
- After outline, before phases
- Any phase within a book
- Between books in a chunk
- Between chunks

### State Preservation
- Book queue and configuration
- Current book and phase numbers
- Completed books list
- Failed books list  
- Quality metrics and progress data
- Generation timestamps

## ⚙️ Configuration

### Default Settings
```go
Model:            "llama3.2"
TotalBooks:       10
ChunkSize:        3  // Books per chunk
AutoSaveInterval: 30 * time.Second
OutputDir:        "output/books"
```

### Customization
Edit the `NewEnhancedSequentialGenerator()` function to modify:
- AI model selection
- Number of books to generate
- Chunk size for processing
- Auto-save frequency
- Output directory location

## 🚨 Error Handling

### Graceful Failures
- Individual phase failures don't stop the entire process
- Book failures are logged and tracked
- Automatic retry capabilities (configurable)
- Detailed error logging with context

### Recovery Options
- Resume from last successful checkpoint
- Skip failed books and continue
- Retry failed operations with different parameters
- Manual intervention points for complex issues

## 📈 Performance Features

### Chunked Processing
- Processes books in configurable batches
- Prevents memory accumulation
- Enables progress checkpoints
- Reduces impact of individual failures

### Resource Monitoring
- CPU and memory usage tracking
- Network call monitoring
- Disk I/O measurement
- LLM token usage tracking

### Optimization
- Efficient state serialization
- Minimal memory footprint
- Optimized logging performance
- Smart caching strategies

## 🔍 Monitoring & Debugging

### Real-time Monitoring
```bash
./manage_generator.sh monitor
```
Shows:
- Current operation and progress
- Resource usage (CPU, memory)
- Recent log entries
- Process status

### Log Analysis
```bash
./manage_generator.sh logs
```
Displays recent log entries for troubleshooting.

### Status Checking
```bash
./manage_generator.sh status
```
Shows comprehensive status including:
- Process information
- State file existence
- Progress metrics
- System health

## 🎯 Book Flow & Logical Progression

### Maintained Context
The enhanced generator preserves logical flow by:
- Maintaining awareness of previously completed books
- Considering book themes and categories for coherent progression
- Ensuring phase dependencies are respected
- Preserving chapter ordering and thematic consistency

### Chapter Flow
Each book maintains internal consistency:
- Chapter 1 → 2 → 3 → 4 → 5 → 6 logical progression
- Consistent character development
- Thematic coherence across chapters
- Proper narrative arc development

### Quality Assurance
- Continuous quality metric tracking
- Content coherence validation
- Thematic consistency checks
- Progress quality assessment

## 📝 Example Session

```bash
# Start generation
./enhanced_sequential_generator

# Output will show:
╔═══════════════════════════════════════════════════════════════════════════════╗
║        📚 ENHANCED SEQUENTIAL BOOK GENERATION SYSTEM 📚                       ║
║    ✨ Features: Resume • Pause • Real-time Logging • Chunked Processing ✨    ║
╚═══════════════════════════════════════════════════════════════════════════════╝

[15:04:05] ℹ️  INFO: Configuration - Session ID: session_1642517045
[15:04:05] ℹ️  INFO: Configuration - Total books: 10
[15:04:05] 🔄 PROGRESS: Loading Themes (10.0%) [██░░░░░░░░░░░░░░░░░░] Reading configuration
[15:04:06] ✅ SUCCESS: Theme Selection - Selected 10 diverse themes
[15:04:06] 🔄 PROGRESS: Chunk Processing (0.0%) [░░░░░░░░░░░░░░░░░░░░] Starting chunk 1/4

# Interactive control:
pause    # Pauses generation
resume   # Resumes generation  
status   # Shows current status
stop     # Graceful stop
```

## 🆘 Troubleshooting

### Common Issues

**Build Failures**
```bash
# Ensure all dependencies are present
go mod tidy
go build -o enhanced_sequential_generator enhanced_sequential_generator.go progress_tracker.go outline_generator.go
```

**Permission Issues**
```bash
chmod +x enhanced_sequential_generator
chmod +x manage_generator.sh
chmod +x build_enhanced_generator.sh
```

**State Corruption**
```bash
# Clean corrupted state
./manage_generator.sh clean
# Start fresh
./manage_generator.sh start
```

**Ollama Connection Issues**
```bash
# Ensure Ollama is running
ollama serve

# Check model availability
ollama list
ollama pull llama3.2
```

### Support

For issues or feature requests:
1. Check the logs: `./manage_generator.sh logs`
2. Verify status: `./manage_generator.sh status`
3. Review state files in the output directory
4. Check Ollama connectivity and model availability

---

**Enhanced Sequential Book Generator** - Co-authored by animality.ai
*Making book generation resumable, manageable, and transparent.*